package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"time"

	"christine.website/internal/blog"
	"christine.website/internal/jsonfeed"
	"christine.website/internal/middleware"
	"github.com/gorilla/feeds"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	blackfriday "github.com/russross/blackfriday"
	"github.com/sebest/xff"
	"github.com/snabb/sitemap"
	"within.website/ln"
	"within.website/ln/ex"
	"within.website/ln/opname"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "29384"
	}

	ctx := ln.WithF(opname.With(context.Background(), "main"), ln.F{
		"port":    port,
		"git_rev": gitRev,
	})

	_ = prometheus.Register(prommod.NewCollector("christine"))

	s, err := Build()
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("Build"))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/.within/health", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "OK", http.StatusOK)
	})
	mux.Handle("/", s)

	ln.Log(ctx, ln.Action("http_listening"))
	ln.FatalErr(ctx, http.ListenAndServeTLS(":"+port,"/certs/server.crt", "/certs/server.key", mux))	
}

// Site is the parent object for https://chopade.xyz's backend.
type Site struct {
	Posts  blog.Posts
	Talks  blog.Posts
	Resume template.HTML

	rssFeed  *feeds.Feed
	jsonFeed *jsonfeed.Feed

	mux   *http.ServeMux
	xffmw *xff.XFF
}

var gitRev = os.Getenv("GIT_REV")

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := opname.With(r.Context(), "site.ServeHTTP")
	ctx = ln.WithF(ctx, ln.F{
		"user_agent": r.Header.Get("User-Agent"),
	})
	r = r.WithContext(ctx)
	if gitRev != "" {
		w.Header().Add("X-Git-Rev", gitRev)
	}

	middleware.RequestID(s.xffmw.Handler(ex.HTTPLog(s.mux))).ServeHTTP(w, r)
}

var arbDate = time.Date(2019, time.May, 20, 18, 0, 0, 0, time.UTC)

// Build creates a new Site instance or fails.
func Build() (*Site, error) {
	smi := sitemap.New()
	smi.Add(&sitemap.URL{
		Loc:        "https://chopade.xyz/resume",
		LastMod:    &arbDate,
		ChangeFreq: sitemap.Monthly,
	})

	smi.Add(&sitemap.URL{
		Loc:        "https://chopade.xyz/contact",
		LastMod:    &arbDate,
		ChangeFreq: sitemap.Monthly,
	})

	smi.Add(&sitemap.URL{
		Loc:        "https://chopade.xyz/",
		LastMod:    &arbDate,
		ChangeFreq: sitemap.Monthly,
	})

	smi.Add(&sitemap.URL{
		Loc:        "https://chopade.xyz/blog",
		LastMod:    &arbDate,
		ChangeFreq: sitemap.Weekly,
	})

	xffmw, err := xff.Default()
	if err != nil {
		return nil, err
	}

	s := &Site{
		rssFeed: &feeds.Feed{
			Title:       "Aniket's Blog",
			Link:        &feeds.Link{Href: "https://chopade.xyz/blog"},
			Description: "My blog posts and rants about various technology things.",
			Author:      &feeds.Author{Name: "Aniket Chopade", Email: "aniketchopade@gmail.com"},
			Created:     bootTime,
			Copyright:   "This work is not copyright of Aniket. My viewpoints are my own and not the view of any employer past, current or future.",
		},
		jsonFeed: &jsonfeed.Feed{
			Version:     jsonfeed.CurrentVersion,
			Title:       "Aniket's Blog",
			HomePageURL: "https://chopade.xyz/",
			FeedURL:     "https://chopade.xyz/blog.json",
			Description: "My blog posts and rants about various technology things.",
			UserComment: "This is a JSON feed of my blogposts. For more information read: https://jsonfeed.org/version/1",
			Icon:        icon,
			Favicon:     icon,
			Author: jsonfeed.Author{
				Name:   "Aniket Chopade",
				Avatar: icon,
			},
		},
		mux:   http.NewServeMux(),
		xffmw: xffmw,
	}

	posts, err := blog.LoadPosts("./blog/", "blog")
	if err != nil {
		return nil, err
	}
	s.Posts = posts

	talks, err := blog.LoadPosts("./blog", "talks")
	if err != nil {
		return nil, err
	}
	s.Talks = talks

	var everything blog.Posts
	everything = append(everything, posts...)
	// everything = append(everything, talks...)

	sort.Sort(sort.Reverse(everything))

	resumeData, err := ioutil.ReadFile("./static/resume/resume.md")
	if err != nil {
		return nil, err
	}

	s.Resume = template.HTML(blackfriday.Run(resumeData))

	for _, item := range everything {
		s.rssFeed.Items = append(s.rssFeed.Items, &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: "https://chopade.xyz/" + item.Link},
			Description: item.Summary,
			Created:     item.Date,
			Content:     string(item.BodyHTML),
		})

		s.jsonFeed.Items = append(s.jsonFeed.Items, jsonfeed.Item{
			ID:            "https://chopade.xyz/" + item.Link,
			URL:           "https://chopade.xyz/" + item.Link,
			Title:         item.Title,
			DatePublished: item.Date,
			ContentHTML:   string(item.BodyHTML),
		})

		smi.Add(&sitemap.URL{
			Loc:        "https://chopade.xyz/" + item.Link,
			LastMod:    &item.Date,
			ChangeFreq: sitemap.Monthly,
		})
	}

	// Add HTTP routes here
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			s.renderTemplatePage("error.html", "can't find "+r.URL.Path).ServeHTTP(w, r)
			return
		}

		s.renderTemplatePage("index.html", nil).ServeHTTP(w, r)
	})
	s.mux.Handle("/metrics", promhttp.Handler())
	s.mux.Handle("/resume", middleware.Metrics("resume", s.renderTemplatePage("resume.html", s.Resume)))
	s.mux.Handle("/blog", middleware.Metrics("blog", s.renderTemplatePage("blogindex.html", s.Posts)))
	s.mux.Handle("/talks", middleware.Metrics("talks", s.renderTemplatePage("talkindex.html", s.Talks)))
	s.mux.Handle("/contact", middleware.Metrics("contact", s.renderTemplatePage("contact.html", nil)))
	s.mux.Handle("/blog.rss", middleware.Metrics("blog.rss", http.HandlerFunc(s.createFeed)))
	s.mux.Handle("/blog.atom", middleware.Metrics("blog.atom", http.HandlerFunc(s.createAtom)))
	s.mux.Handle("/blog.json", middleware.Metrics("blog.json", http.HandlerFunc(s.createJSONFeed)))
	s.mux.Handle("/blog/", middleware.Metrics("blogpost", http.HandlerFunc(s.showPost)))
	s.mux.Handle("/talks/", middleware.Metrics("talks", http.HandlerFunc(s.showTalk)))
	s.mux.Handle("/css/", http.FileServer(http.Dir(".")))
	s.mux.Handle("/static/", http.FileServer(http.Dir(".")))
	s.mux.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/js/sw.js")
	})
	s.mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/robots.txt")
	})
	s.mux.Handle("/sitemap.xml", middleware.Metrics("sitemap", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		_, _ = smi.WriteTo(w)
	})))

	return s, nil
}

const icon = "https://chopade.xyz/static/img/avatar.png"
