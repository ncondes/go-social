package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/repositories"
)

const (
	usersAmount    = 50
	postsAmount    = 100
	commentsAmount = 200
)

var seedFirstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Barbara", "David", "Elizabeth", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen",
}

var seedLastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas",
	"Taylor", "Moore", "Jackson", "Martin",
}

var seedPostTitles = []string{
	"Getting Started with Go", "Understanding Microservices", "Best Practices for REST APIs",
	"Introduction to Docker", "Mastering Git Workflows", "Database Design Patterns",
	"Clean Code Principles", "Testing Strategies in Go", "Building Scalable Systems",
	"API Security Best Practices", "Concurrency in Go", "Error Handling Techniques",
	"Performance Optimization Tips", "Debugging Like a Pro", "Code Review Guidelines",
	"Refactoring Legacy Code", "Modern Web Development", "Cloud Architecture Basics",
	"DevOps Fundamentals", "Continuous Integration Pipelines",
}

var seedPostContents = []string{
	"This is a comprehensive guide to help you understand the fundamentals and best practices.",
	"In this post, we'll explore various techniques and patterns that can improve your workflow.",
	"Let me share some insights I've gained from years of experience in software development.",
	"Here are some practical examples and code snippets to illustrate the key concepts.",
	"This tutorial will walk you through step-by-step instructions with detailed explanations.",
	"I've compiled a list of resources and tools that have been incredibly helpful in my journey.",
	"Today we're diving deep into advanced topics that will take your skills to the next level.",
	"Learn how to avoid common pitfalls and mistakes that many developers encounter.",
	"This article covers everything you need to know to get started with this technology.",
	"Discover the latest trends and innovations that are shaping the future of development.",
	"A detailed analysis of different approaches and their trade-offs in real-world scenarios.",
	"Practical tips and tricks that you can immediately apply to your projects.",
	"Understanding the underlying principles will help you make better architectural decisions.",
	"This guide includes benchmarks, comparisons, and recommendations based on extensive testing.",
	"Explore the ecosystem and learn about the most popular libraries and frameworks.",
	"Real-world case studies demonstrating how these concepts are applied in production.",
	"A beginner-friendly introduction with clear examples and easy-to-follow instructions.",
	"Advanced techniques for optimizing performance and improving code quality.",
	"Common questions answered with detailed explanations and working code examples.",
	"Best practices from industry experts and lessons learned from large-scale projects.",
}

var seedPostTags = []string{
	"golang", "programming", "webdev", "tutorial", "backend", "api", "database",
	"docker", "kubernetes", "microservices", "testing", "security", "performance",
	"architecture", "devops", "cloud", "bestpractices", "coding", "software", "tech",
}

var seedCommentContents = []string{
	"Great post! Thanks for sharing this valuable information.",
	"This is exactly what I was looking for. Very helpful!",
	"Interesting perspective. I learned something new today.",
	"Could you elaborate more on this topic? I'd love to know more.",
	"Excellent explanation! Clear and concise.",
	"I have a different approach that also works well in my experience.",
	"Thanks for the detailed tutorial. It helped me solve my problem.",
	"This is a must-read for anyone working with this technology.",
	"Well written and easy to follow. Keep up the good work!",
	"I disagree with some points, but overall a solid article.",
	"Bookmarking this for future reference. Very useful!",
	"Can you provide more examples? That would be really helpful.",
	"I've been struggling with this, and your post clarified everything.",
	"Amazing content! Looking forward to more posts like this.",
	"This should be part of the official documentation. So clear!",
	"I tried this approach and it works perfectly. Thanks!",
	"Insightful post. You covered all the important aspects.",
	"Quick question: does this work with the latest version?",
	"Fantastic write-up! Sharing this with my team.",
	"This helped me understand the concept much better. Appreciate it!",
}

var seedEmailDomains = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "outlook.com", "apple.com",
}

func Flush(db *sql.DB) {
	ctx := context.Background()

	queries := []string{
		"TRUNCATE TABLE comments CASCADE;",
		"TRUNCATE TABLE posts CASCADE;",
		"TRUNCATE TABLE users CASCADE;",
	}

	for _, query := range queries {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			log.Printf("Error flushing db: %v\n", err)
		}
	}
}

func Seed(r *repositories.Repositories) {
	ctx := context.Background()

	users := generateUsers(usersAmount)
	for _, user := range users {
		err := r.UserRepository.Create(ctx, user)
		if err != nil {
			log.Printf("Error seeding users: %v\n", err)
		}
	}

	posts := generatePosts(postsAmount, users)
	for _, post := range posts {
		err := r.PostRepository.Create(ctx, post)
		if err != nil {
			log.Printf("Error seeding posts: %v\n", err)
		}
	}

	comments := generateComments(commentsAmount, posts, users)
	for _, comment := range comments {
		err := r.CommentRepository.Create(ctx, comment)
		if err != nil {
			log.Printf("Error seeding comments: %v\n", err)
		}
	}
}

func generateUsers(amount int) []*domain.User {
	users := make([]*domain.User, amount)

	for i := 0; i < amount; i++ {
		firstName := seedFirstNames[rand.Intn(len(seedFirstNames))]
		lastName := seedLastNames[rand.Intn(len(seedLastNames))]
		username := strings.ToLower(string(firstName[0])+strings.TrimSpace(lastName)) + strconv.Itoa(i)
		email := username + "@" + seedEmailDomains[rand.Intn(len(seedEmailDomains))]

		users[i] = &domain.User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			Email:     email,
			Password:  "Qwerty123$",
		}
	}

	return users
}

func generatePosts(amount int, users []*domain.User) []*domain.Post {
	posts := make([]*domain.Post, amount)

	for i := 0; i < amount; i++ {
		numTags := rand.Intn(6) // 0 to 5 tags
		tags := make([]string, 0, numTags)

		// Randomly select unique tags
		idxs := make(map[int]bool)
		for len(idxs) < numTags {
			idx := rand.Intn(len(seedPostTags))

			if _, exists := idxs[idx]; !exists {
				idxs[idx] = true
				tags = append(tags, seedPostTags[idx])
			}
		}

		userID := users[rand.Intn(len(users))].ID
		title := seedPostTitles[rand.Intn(len(seedPostTitles))]
		content := seedPostContents[rand.Intn(len(seedPostContents))]

		posts[i] = &domain.Post{
			Title:   title,
			Content: content,
			Tags:    tags,
			UserID:  userID,
		}
	}

	return posts
}

func generateComments(amount int, posts []*domain.Post, users []*domain.User) []*domain.Comment {
	comments := make([]*domain.Comment, amount)

	for i := 0; i < amount; i++ {
		userID := users[rand.Intn(len(users))].ID
		postID := posts[rand.Intn(len(posts))].ID
		content := seedCommentContents[rand.Intn(len(seedCommentContents))]

		comments[i] = &domain.Comment{
			Content: content,
			UserID:  userID,
			PostID:  postID,
		}
	}

	return comments
}
