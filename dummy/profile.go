package main

import (
	"math/rand"

	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/icrowley/fake"
)

func newRandomProfile() *pb.Profile {
	name := "🤖" + fake.Company()

	vendor := true
	if rand.Intn(3) == 0 {
		vendor = false
		name = fake.FullName()
	}

	moderator := false
	modInfo := pb.Moderator{}
	if rand.Intn(8) == 0 {
		moderator = true
		modInfo.Description = fake.Sentences()
		modInfo.TermsAndConditions = fake.Paragraphs()
		modInfo.Languages = []string{"English"}

		langCount := rand.Intn(6)
		for i := 0; i < langCount; i++ {
			modInfo.Languages = append(modInfo.Languages, fake.Language())
		}

		modInfo.Fee = &pb.Moderator_Fee{
			Percentage: float32(rand.Intn(75)),
			FeeType:    pb.Moderator_Fee_PERCENTAGE,
		}
	}

	return &pb.Profile{
		Handle:           "@" + fake.UserName(),
		Name:             name,
		Location:         fake.City() + ", " + fake.Country(),
		About:            fake.Paragraphs(),
		ShortDescription: fake.Sentences(),
		ContactInfo: &pb.Profile_Contact{
			Website:     fake.Word() + ".example.com",
			Email:       fake.EmailAddress(),
			PhoneNumber: "555-5555",
		},

		Vendor:        vendor,
		Nsfw:          isNSFW(),
		Moderator:     moderator,
		ModeratorInfo: &modInfo,

		Stats: &pb.Profile_Stats{
			FollowerCount:  uint32(rand.Intn(9999)),
			FollowingCount: uint32(rand.Intn(9999)),
			AverageRating:  rand.Float32() * 5,
			RatingCount:    uint32(rand.Intn(9999)),

			// TODO
			// ListingsCount
		},

		Colors: &pb.Profile_Colors{
			Primary:       "#" + fake.HexColor(),
			Secondary:     "#" + fake.HexColor(),
			Text:          "#" + fake.HexColor(),
			Highlight:     "#" + fake.HexColor(),
			HighlightText: "#" + fake.HexColor(),
		},
	}
}
