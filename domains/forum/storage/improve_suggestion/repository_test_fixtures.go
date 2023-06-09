package improve_suggestion_storage

import (
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 8, 10, 0, 0, time.UTC)
)

var Fixtures = []interface{}{
	// Requests.
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		Title:     "Test",
		Content:   "Dummy content.",
	},
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(1001),
		CreatedAt: baseTime.Add(10 * time.Minute),
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		Title:     "Test",
		Content:   "Dummy content updated.",
	},
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(1002),
		CreatedAt: baseTime.Add(time.Minute),
		Source:    test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		Title:     "New Test",
		Content:   "Dummy content updated.",
	},
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(5000),
		CreatedAt: baseTime,
		Source:    test_utils.NumberUUID(5000),
		UserID:    test_utils.NumberUUID(300),
		Title:     "Lorem Ipsum",
		Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a.",
	},
	&improve_request_storage.Model{
		ID:        test_utils.NumberUUID(6000),
		CreatedAt: baseTime.Add(30 * time.Minute),
		Source:    test_utils.NumberUUID(6000),
		UserID:    test_utils.NumberUUID(200),
		Title:     "New title Updated.",
		Content:   "qwertyuiopasdfghjklzxcvbnm",
	},
	// Suggestions.
	&Model{ // 10, validated
		ID:        test_utils.NumberUUID(1000),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		Validated: true,
		UpVotes:   10,
		Core: Core{
			RequestID: test_utils.NumberUUID(1000),
			Title:     "Test 2",
			Content:   "Dummy content.",
		},
	},
	&Model{ // 5, validated
		ID:        test_utils.NumberUUID(1001),
		CreatedAt: baseTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(201),
		Validated: true,
		UpVotes:   7,
		DownVotes: 2,
		Core: Core{
			RequestID: test_utils.NumberUUID(1000),
			Title:     "Test",
			Content:   "Smart content.",
		},
	},
	&Model{ // 13
		ID:        test_utils.NumberUUID(1002),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		UpVotes:   21,
		DownVotes: 8,
		Core: Core{
			RequestID: test_utils.NumberUUID(1001),
			Title:     "Test",
			Content:   "Simple content.",
		},
	},
	&Model{ // 3, validated
		ID:        test_utils.NumberUUID(1003),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		Validated: true,
		UpVotes:   16,
		DownVotes: 13,
		Core: Core{
			RequestID: test_utils.NumberUUID(1001),
			Title:     "Test 3",
			Content:   "Simple content 3.",
		},
	},
	&Model{ // -4
		ID:        test_utils.NumberUUID(1004),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(201),
		UpVotes:   4,
		DownVotes: 8,
		Core: Core{
			RequestID: test_utils.NumberUUID(1001),
			Title:     "Test 4",
			Content:   "Simple content 4.",
		},
	},
	&Model{ // 8
		ID:        test_utils.NumberUUID(1005),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		UpVotes:   32,
		DownVotes: 24,
		Core: Core{
			RequestID: test_utils.NumberUUID(1000),
			Title:     "Test 5",
			Content:   "Simple content 5.",
		},
	},
	&Model{ // 9
		ID:        test_utils.NumberUUID(1006),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(200),
		UpVotes:   9,
		Core: Core{
			RequestID: test_utils.NumberUUID(1000),
			Title:     "Test 6",
			Content:   "Simple content 6.",
		},
	},
	&Model{ // 0, validated
		ID:        test_utils.NumberUUID(1007),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(1000),
		UserID:    test_utils.NumberUUID(201),
		Validated: true,
		Core: Core{
			RequestID: test_utils.NumberUUID(1001),
			Title:     "Test 7",
			Content:   "Simple content 7.",
		},
	},
	&Model{ // 3
		ID:        test_utils.NumberUUID(2000),
		CreatedAt: baseTime,
		UpdatedAt: &updateTime,
		SourceID:  test_utils.NumberUUID(5000),
		UserID:    test_utils.NumberUUID(201),
		UpVotes:   4,
		DownVotes: 1,
		Core: Core{
			RequestID: test_utils.NumberUUID(5000),
			Title:     "Ipsum Lorem",
			Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum tempor congue aliquam. Nam ullamcorper mi lectus, et dictum urna imperdiet a.",
		},
	},
}
