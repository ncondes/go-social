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

const testPassword = "Qwerty123$"
const userRoleName = "user"

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

func Seed(r *repositories.Repositories, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(usersAmount, db)
	for _, user := range users {
		err := r.UserRepository.CreateUser(ctx, user)
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

	followers := generateFollowers(users)
	for _, follower := range followers {
		err := r.FollowerRepository.FollowUser(ctx, follower.UserID, follower.FollowerID)
		if err != nil {
			log.Printf("Error seeding followers: %v\n", err)
		}
	}
}

func generateUsers(amount int, db *sql.DB) []*domain.User {
	users := make([]*domain.User, amount)
	// get role id by name
	var roleID int64

	err := db.QueryRow("SELECT id FROM roles WHERE name = $1", userRoleName).Scan(&roleID)
	if err != nil {
		log.Printf("Error getting role: %v\n", err)
		return users
	}

	for i := 0; i < amount; i++ {
		firstName := seedFirstNames[rand.Intn(len(seedFirstNames))]
		lastName := seedLastNames[rand.Intn(len(seedLastNames))]
		username := strings.ToLower(string(firstName[0])+strings.TrimSpace(lastName)) + strconv.Itoa(i)
		email := username + "@" + seedEmailDomains[rand.Intn(len(seedEmailDomains))]

		user := &domain.User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			Email:     email,
			RoleID:    roleID,
		}

		if err := user.HashPassword(testPassword); err != nil {
			log.Printf("Error hashing password for user %s: %v\n", username, err)
		}

		users[i] = user
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

func generateFollowers(users []*domain.User) []*domain.Follower {
	followers := make([]*domain.Follower, 0)
	totalUsers := len(users)

	for _, user := range users {
		var numToFollow int
		roll := rand.Float64()

		switch {
		case roll < 0.25: // 25% are casual users (5-15% of total)
			min := int(float64(totalUsers) * 0.05)
			max := int(float64(totalUsers) * 0.15)
			numToFollow = min + rand.Intn(max-min+1)

		case roll < 0.60: // 35% are regular users (15-35% of total)
			min := int(float64(totalUsers) * 0.15)
			max := int(float64(totalUsers) * 0.35)
			numToFollow = min + rand.Intn(max-min+1)

		case roll < 0.85: // 25% are active users (35-50% of total)
			min := int(float64(totalUsers) * 0.35)
			max := int(float64(totalUsers) * 0.50)
			numToFollow = min + rand.Intn(max-min+1)

		default: // 15% are power users (50-60% of total)
			min := int(float64(totalUsers) * 0.50)
			max := int(float64(totalUsers) * 0.60)
			numToFollow = min + rand.Intn(max-min+1)
		}

		followed := make(map[int64]bool)

		for len(followed) < numToFollow {
			targetUser := users[rand.Intn(len(users))]

			if targetUser.ID == user.ID || followed[targetUser.ID] {
				continue // Skip if user already followed or self
			}

			followed[targetUser.ID] = true // Mark as followed

			followers = append(followers, &domain.Follower{
				UserID:     targetUser.ID,
				FollowerID: user.ID,
			})
		}
	}

	return followers
}
