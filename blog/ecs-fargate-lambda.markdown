---
title: What to use and when? ECS, Fargate or  Lambda
date: 2019-12-12
for: Myself
---
# What to use and when? ECS, Fargate or  Lambda
## Lambda
   * Used for batch loads, mobile back-ends, stateless web services 
   * Low TPS applications, which are in inception
   * Applications with no requirement of persistent storage
   * Pros
       - Quick to time to market
       - No infrastructure patching (serverless advantages)
   * Cons
       - As an application show steady growth, Lambda's cost increases linearly with the number of requests.
       - No EBS volume to attach to the Lambda function. Disk space is limited to 512MB.
       - The code is not directly portable on-premises.
## ECS-Fargate
   * Used for any containerized applications.
   * Containers can be scaled up and down when the application's traffic pattern is well known. Thus cost can be curtailed.
   * Non-persistent storage can be attached (10 GB)
   * Pros- 
       - The code is directly portable on-premises or on other clouds
       - No infrastructure patching (serverless advantages)
   * Cons
       - Cost-effectiveness vanishes when traffic inflow is steady.
## ECS-EC2
   * Used for any containerized applications. 
   * Cost-effective when container CPU reservation, memory reservation can be kept steady.
   * Reserved instances further reduce the cost significantly 
   * Cons
       - CPU, memory utilization should be kept under close watch. Underutilization reduces cost-benefits.