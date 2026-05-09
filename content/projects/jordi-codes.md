---
title: jordi.codes SSH Site
date: 2026-05-07
---

The site you're reading this from as you can tell it's being served over SSH.

I got inspired by this guys selling coffee over SSH and I thought that it would be a fun way to share my information, I find that it's beautiful the fact that that's not the usual way you'll find someone's info nowadays, and I wanted to do something different.

## Why?

So the first thing you might wonder, why would I do this?

Indeed it's a good way to explore tech, it gives you a different perspective on how IT really works, for people that were alive during the 90s that might be obvious, but for gen Z guys like me, it's often difficult to really think of something else than a website when you think about sharing content online on a non-limited way. 

It's also true that nowadays it's becoming more and more difficult to be in an anonymous way online, the internet is full of trackers, cookies, and things like that. This site is an exploration and a real alternative of how things could have been. It's certainly not as powerful, it lacks one fundamental capability, which is rendering images, but who cares, right?

## Security

The security concerns about this project might be pretty obvious, but indeed this is more secure than you think. Wish provides a secure interface for creating SSH servers, this is not a server that defaults to a shell and runs an script. This is a custom implementation of the SSH protocol, so sice it doesn't default to a shell, it doesn't have access to the system, so it's pretty secure. It's just piping I/O to a TUI application, every client has its own instance of the TUI application.

## Tech

I initially wanted to build this from scratch when I then realized the gorgeous work that the folks at Charmbracelet have been doing, so I decided to build this on top of their stack, which is built in Go and it's really easy to use, and it has a lot of cool features that I wanted to use for this project.

- Go
- [wish](https://github.com/charmbracelet/wish) — SSH server
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [glamour](https://github.com/charmbracelet/glamour) — Markdown rendering
- [lipgloss](https://github.com/charmbracelet/lipgloss) — Styling
