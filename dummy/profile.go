package main

import (
	"math/rand"

	"github.com/OpenBazaar/openbazaar-go/pb"
	"github.com/icrowley/fake"
)

func newRandomProfile(randomImages chan (*randomImage)) *pb.Profile {
	name := "🤖" + fake.Company()

	vendor := true

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

	avatar := <-randomImages
	header := <-randomImages

	return &pb.Profile{
		Handle:           "@" + fake.UserName(),
		Name:             name,
		Location:         fake.City() + ", " + fake.Country(),
		About:            fake.Paragraphs(),
		Email:            fake.EmailAddress(),
		ShortDescription: "I am a robot.\n\n" + fake.Sentences(),

		Vendor:    vendor,
		Nsfw:      isNSFW(),
		Moderator: moderator,
		ModInfo:   &modInfo,

		FollowerCount:  uint32(rand.Intn(9999)),
		FollowingCount: uint32(rand.Intn(9999)),
		AvgRating:      uint32(rand.Intn(5)),
		NumRatings:     uint32(rand.Intn(9999)),

		PrimaryColor:       "#" + fake.HexColor(),
		SecondaryColor:     "#" + fake.HexColor(),
		TextColor:          "#" + fake.HexColor(),
		HighlightColor:     "#" + fake.HexColor(),
		HighlightTextColor: "#" + fake.HexColor(),

		AvatarHashes: &pb.Profile_Image{
			Tiny:     avatar.Tiny,
			Small:    avatar.Small,
			Medium:   avatar.Medium,
			Large:    avatar.Large,
			Original: avatar.Original,
		},

		HeaderHashes: &pb.Profile_Image{
			Tiny:     header.Tiny,
			Small:    header.Small,
			Medium:   header.Medium,
			Large:    header.Large,
			Original: header.Original,
		},
	}
}
