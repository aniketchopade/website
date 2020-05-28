---
title: strace with containers
date: 2019-11-30
for: Myself
---
# Strace
* strace is an essential tool for debugging anything on a Linux system. You should use strace anytime you want to understand what a process is doing. Reaching for strace as the first line of defense when debugging anything is a great way to quickly gather context about a problem.
# Get container details
There are number of ways you can get container's PID
1. get from htop. 
    * Use filter function and search with ENTRY level command in container. Get PID details
2. use control groups
    * Docker relies on control groups (cgroups). cgroups provide a way of partitioning off resources for groups of users and tasks. Complete explaination here- https://www.kernel.org/doc/Documentation/cgroup-v1/cgroups.txt
    * use systemd-cgls and systemd-cgtop command 
# Start trace
    * strace -e trace=\!futex,epoll_pwait -Tfp <pid>
    * T will print time
    * f will print all threads for process
    * p for mentioning PID