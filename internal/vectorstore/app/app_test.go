package app_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aria3ppp/rag-server/internal/vectorstore/domain"
	"github.com/aria3ppp/rag-server/pkg/wait"

	"github.com/gavv/httpexpect/v2"
	grpc_codes "google.golang.org/grpc/codes"
)

func Test_App_InsertTexts_then_SearchText(t *testing.T) {
	t.Parallel()

	grpcGatewayPort := setupVectorStoreApp(t, &wait.Opts{Interval: time.Second, MaxRetries: 30})
	baseURL := fmt.Sprintf("http://localhost:%d", grpcGatewayPort)

	e := httpexpect.Default(t, baseURL)

	// insert texts: validation failed
	e.POST("/api/v1/insert_texts").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().IsEqual(map[string]any{
		"code":    grpc_codes.InvalidArgument,
		"details": []any{},
		"message": (&domain.InsertTextsInput{Texts: []*domain.InsertTextsInputText{}}).Validate(context.Background()).Error(),
	})

	// insert texts: ok
	e.POST("/api/v1/insert_texts").
		WithJSON(map[string]any{
			"texts": []map[string]any{
				{
					"text": "The quick brown fox jumps over the lazy dog, showcasing agility, speed, and precision in every leap across the vibrant meadow.",
					"metadata": map[string]any{
						"animal_kind": "fox",
						"color":       "brown",
					},
				},
				{
					"text": "A graceful white swan glides across the shimmering lake, its movements a testament to elegance and beauty.",
					"metadata": map[string]any{
						"animal_kind": "swan",
						"color":       "white",
					},
				},
				{
					"text": "Beneath the golden sun, a playful orange tabby cat pounces on falling leaves scattered across the lawn.",
					"metadata": map[string]any{
						"animal_kind": "cat",
						"color":       "orange",
					},
				},
				{
					"text": "High in the sky, a soaring bald eagle with brown and white feathers scans the forest below for its prey.",
					"metadata": map[string]any{
						"animal_kind": "eagle",
						"color":       "brown and white",
					},
				},
				{
					"text": "In the shadow of the tall trees, a quiet gray wolf prowls silently, blending seamlessly into the wilderness.",
					"metadata": map[string]any{
						"animal_kind": "wolf",
						"color":       "gray",
					},
				},
				{
					"text": "A bright yellow canary chirps joyfully, its song echoing through the lush green garden filled with blooming flowers.",
					"metadata": map[string]any{
						"animal_kind": "canary",
						"color":       "yellow",
					},
				},
				{
					"text": "A sleek black panther moves stealthily through the dense jungle, its piercing yellow eyes glowing in the darkness.",
					"metadata": map[string]any{
						"animal_kind": "panther",
						"color":       "black",
					},
				},
				{
					"text": "An industrious red squirrel gathers acorns beneath a tall oak tree, its bushy tail twitching with excitement.",
					"metadata": map[string]any{
						"animal_kind": "squirrel",
						"color":       "red",
					},
				},
				{
					"text": "On the icy tundra, a majestic white polar bear roams, its thick fur gleaming under the pale arctic sun.",
					"metadata": map[string]any{
						"animal_kind": "polar bear",
						"color":       "white",
					},
				},
				{
					"text": "Under the blue expanse of the ocean, a curious green sea turtle swims gracefully among colorful coral reefs.",
					"metadata": map[string]any{
						"animal_kind": "sea turtle",
						"color":       "green",
					},
				},
				// {
				// 	"text": "In a futuristic city filled with towering skyscrapers and neon lights, a silver robot patrols the streets tirelessly, ensuring peace and harmony among its human creators, who rely on it for safety.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "robot",
				// 		"color":       "silver",
				// 	},
				// },
				// {
				// 	"text": "A loyal golden retriever waits patiently by the door of a cozy cottage, wagging its tail in joyful anticipation of a walk through the blooming meadows and the quiet forest trails beyond.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "dog",
				// 		"color":       "golden",
				// 	},
				// },
				// {
				// 	"text": "A talkative parrot with blue and green feathers perches on a tree branch, mimicking the sounds of the jungle.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "parrot",
				// 		"color":       "blue and green",
				// 	},
				// },
				// {
				// 	"text": "Beneath the ocean waves, a playful gray dolphin leaps gracefully through the air, thrilling nearby spectators.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "dolphin",
				// 		"color":       "gray",
				// 	},
				// },
				// {
				// 	"text": "A curious white rabbit hops playfully through a lush meadow filled with vibrant wildflowers, stopping occasionally to nibble on fresh grass while keeping a watchful eye for any signs of danger.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "rabbit",
				// 		"color":       "white",
				// 	},
				// },
				// {
				// 	"text": "A majestic black stallion gallops across the open plains under a wide blue sky, its mane flowing in the wind as it moves with a powerful grace that inspires awe in those lucky enough to witness it.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "horse",
				// 		"color":       "black",
				// 	},
				// },
				// {
				// 	"text": "Under the moonlight, a silent brown owl swoops through the forest, its eyes fixed on unsuspecting prey.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "owl",
				// 		"color":       "brown",
				// 	},
				// },
				// {
				// 	"text": "On the icy shores of Antarctica, a waddling black and white penguin dives into the freezing water with a splash, swimming effortlessly among floating icebergs in search of its next meal beneath the frigid surface.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "penguin",
				// 		"color":       "black and white",
				// 	},
				// },
				// {
				// 	"text": "In the vast savanna under a golden sunset, a towering gray elephant uses its dexterous trunk to pluck leaves from a tall tree, its large ears flapping gently as it moves with serene confidence.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "elephant",
				// 		"color":       "gray",
				// 	},
				// },
				// {
				// 	"text": "Camouflaged against the leaves, a green chameleon slowly extends its tongue to catch an unsuspecting insect.",
				// 	"metadata": map[string]any{
				// 		"animal_kind": "chameleon",
				// 		"color":       "green",
				// 	},
				// },
			},
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().IsEmpty()

	// search text: validation failed
	e.POST("/api/v1/search_text").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().IsEqual(map[string]any{
		"code":    grpc_codes.InvalidArgument,
		"details": []any{},
		"message": (&domain.SearchTextInput{}).Validate(context.Background()).Error(),
	})

	// search text: ok
	e.POST("/api/v1/search_text").
		WithJSON(map[string]any{
			"text": "The quick brown fox leaps through the forest, its sharp eyes scanning for prey.",
			// "text":   "A golden dog runs joyfully through a green field, its tail wagging in excitement.",
			"top_k":     10,
			"min_score": 0.5,
			"filter":    nil,
		}).
		Expect().
		Status(http.StatusOK).
		JSON().Object().IsEqual(map[string]any{
		"similar_texts": []map[string]any{
			{
				"metadata": map[string]any{
					"animal_kind": "fox",
					"color":       "brown",
				},
				"score": 0.72626674,
				"text":  "The quick brown fox jumps over the lazy dog, showcasing agility, speed, and precision in every leap across the vibrant meadow.",
			},
			{
				"metadata": map[string]any{
					"animal_kind": "eagle",
					"color":       "brown and white",
				},
				"score": 0.6991191,
				"text":  "High in the sky, a soaring bald eagle with brown and white feathers scans the forest below for its prey.",
			},
			{
				"metadata": map[string]any{
					"animal_kind": "wolf",
					"color":       "gray",
				},
				"score": 0.57536006,
				"text":  "In the shadow of the tall trees, a quiet gray wolf prowls silently, blending seamlessly into the wilderness.",
			},
			{
				"metadata": map[string]any{
					"animal_kind": "squirrel",
					"color":       "red",
				},
				"score": 0.5270842,
				"text":  "An industrious red squirrel gathers acorns beneath a tall oak tree, its bushy tail twitching with excitement.",
			},
			{
				"metadata": map[string]any{
					"animal_kind": "polar bear",
					"color":       "white",
				},
				"score": 0.5223348,
				"text":  "On the icy tundra, a majestic white polar bear roams, its thick fur gleaming under the pale arctic sun.",
			},
			{
				"metadata": map[string]any{
					"animal_kind": "panther",
					"color":       "black",
				},
				"score": 0.5124951,
				"text":  "A sleek black panther moves stealthily through the dense jungle, its piercing yellow eyes glowing in the darkness.",
			},
			// {
			// 	"metadata": map[string]any{
			// 		"animal_kind": "cat",
			// 		"color":       "orange",
			// 	},
			// 	"score": 0.4757305,
			// 	"text":  "Beneath the golden sun, a playful orange tabby cat pounces on falling leaves scattered across the lawn.",
			// },
			// {
			// 	"metadata": map[string]any{
			// 		"animal_kind": "canary",
			// 		"color":       "yellow",
			// 	},
			// 	"score": 0.4666258,
			// 	"text":  "A bright yellow canary chirps joyfully, its song echoing through the lush green garden filled with blooming flowers.",
			// },
			// {
			// 	"metadata": map[string]any{
			// 		"animal_kind": "swan",
			// 		"color":       "white",
			// 	},
			// 	"score": 0.45963997,
			// 	"text":  "A graceful white swan glides across the shimmering lake, its movements a testament to elegance and beauty.",
			// },
			// {
			// 	"metadata": map[string]any{
			// 		"animal_kind": "sea turtle",
			// 		"color":       "green",
			// 	},
			// 	"score": 0.41122904,
			// 	"text":  "Under the blue expanse of the ocean, a curious green sea turtle swims gracefully among colorful coral reefs.",
			// },
		},
	})

	var end = true
	_ = end
}
