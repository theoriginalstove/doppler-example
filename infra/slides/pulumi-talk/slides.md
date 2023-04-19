---
# try also 'default' to start simple
theme: seriph
# random image from a curated Unsplash collection by Anthony
# like them? see https://unsplash.com/collections/94734566/slidev
background: https://source.unsplash.com/collection/94734566/1920x1080
# apply any windi css classes to the current slide
class: 'text-center'
# https://sli.dev/custom/highlighters.html
highlighter: shiki
# show line numbers in code blocks
lineNumbers: false
# some information about the slides, markdown enabled
info: |
  ## Slidev Starter Template
  Presentation slides for developers.

  Learn more at [Sli.dev](https://sli.dev)
# persist drawings in exports and build
drawings:
  persist: false
# use UnoCSS
css: unocss
---

# Using Go for Infrastructure as Code

The Source Code:

 https://github.com/theoriginalstove/doppler-example

 Branch: pulumi-example


<!--
The last comment block of each slide will be treated as slide notes. It will be visible and editable in Presenter Mode along with the slide. [Read more in the docs](https://sli.dev/guide/syntax.html#notes)
-->

---

# Who am I?

Steven Turturo
- **Software Engineer**: @Finxact
- **Github**: TheOriginalStove
- **Background**: Worked in healthcare, transportation & logistics, ad-tech startups, and now banking.
    - Started with VBA for Excel & T-SQL back in 2013, and tried all sorts of programming languages like Ruby, Python, C#/.NET
    - Learned Go back in 2019, but working professionally with Go since early 2022.
- **Random Facts**:
    - Recently became addicted to finding OSS projects to contribute to.
    - Owner of a German Shedder(Shepherd)
    - Make too many references to movies and tv shows randomly.
    - Nickname and github username comes from the Bridesmaids movie




<style>
h1 {
  background-color: #2B90B6;
  background-image: linear-gradient(45deg, #4EC5D4 10%, #146b8c 20%);
  background-size: 100%;
  -webkit-background-clip: text;
  -moz-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-text-fill-color: transparent;
}
</style>

<!--
Here is another comment.
-->

---

# What is Infrastructure as Code? 

Programatically defining what our services need to run using code.

- YAML
- JSON
- HCL
- Go

#### Options
 1. Teraform (HCL)
 1. Ansible (Python)
 1. Pulumi (Go, JS/TS, C#, Python, Java, YAML)
 1. Crossplane (YAML)
 1. And others

 Great comparison video of a few IaC options - [https://www.youtube.com/watch?v=RaoKcJGchKM](https://www.youtube.com/watch?v=RaoKcJGchKM)

---
layout: image-right
image: https://source.unsplash.com/collection/94734566/1920x1080
---

# Pulumi basic concepts & definitions

- Project 
- Stacks
- Resources
- State/Backends
    - S3
    - Pulumi Cloud
    - Cloud Storage
- 

---

# Goal of this Walk through

- Deploy the Doppler Example App to GCP Cloud Run
- Define all the resources needed to secure the app
- Give it an ingress and IP address

