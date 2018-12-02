package linked_lists

import "errors"

type Post struct {
	body        string
	publishDate int64
	next        *Post
}

type Feed struct {
	length int
	start  *Post
	end    *Post
}

func (f *Feed) Append(newPost *Post) {
	if f.length == 0 {
		f.start = newPost
		f.end = newPost
	} else {
		lastPost := f.end
		lastPost.next = newPost
		f.end = newPost
	}
	f.length++
}

func (f *Feed) Remove(publishDate int64) {
	if f.length == 0 {
		panic(errors.New("feed is empty"))
	}
	var previousPost *Post
	currentPost := f.start

	for currentPost.publishDate != publishDate {
		if currentPost.next == nil {
			panic(errors.New("no such post"))
		}
		previousPost = currentPost
		currentPost = currentPost.next
	}
	previousPost.next = currentPost.next

	f.length--
}

func (f *Feed) Insert(newPost *Post) {
	if f.length == 0 {
		f.start = newPost
		f.end = newPost
	} else {
		var previousPost *Post
		currentPost := f.start

		for currentPost.publishDate < newPost.publishDate {
			previousPost = currentPost
			currentPost = currentPost.next
		}
		previousPost.next = newPost
		newPost.next = currentPost
	}
	f.length++
}
