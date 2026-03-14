package common

import (
	"fmt"
	"strings"
)

// --- 1. VOICE MODELING ---

// VoiceProfile defines a specific persona's settings.
type VoiceProfile struct {
	Name         string
	Pitch        float64 // Semitones (-20.0 to 20.0)
	SpeakingRate float64 // Speed (0.25 to 4.0)
	Label        string
}

// Line represents a single piece of dialogue in a conversation.
type Line struct {
	Speaker VoiceProfile
	Text    string
}

// Conversation manages a sequence of lines to be converted to SSML.
type Conversation struct {
	Lines []Line
}

// ToSSML generates the final XML string for Google Cloud TTS.
func (c *Conversation) ToSSML() string {
	var sb strings.Builder
	sb.WriteString("<speak>")
	for _, line := range c.Lines {
		// We wrap each voice in a prosody tag to apply the specific profile's pitch/rate
		sb.WriteString(fmt.Sprintf(
			`<voice name="%s"><prosody pitch="%.1fst" rate="%.2f">%s</prosody></voice>`,
			line.Speaker.Name,
			line.Speaker.Pitch,
			line.Speaker.SpeakingRate,
			line.Text,
		))
	}
	sb.WriteString("</speak>")
	return sb.String()
}

// --- 2. THE VOICE LIBRARY (30 PERSONAS) ---

const (
	// Base Engine Names
	S_Male   = "en-US-Studio-Q"
	S_Female = "en-US-Studio-O"
	N_Male   = "en-US-Neural2-D"
	N_Female = "en-US-Neural2-F"
	N_Alt_M  = "en-US-Neural2-J"
	N_Alt_F  = "en-US-Neural2-H"
)

var (
	// --- ANNOUNCERS (Formal & Constant) ---
	Announcer        = VoiceProfile{S_Male, -1.0, 0.95, "Formal Announcer"}
	NarratorFemale   = VoiceProfile{S_Female, 0.0, 1.0, "Smooth Narrator"}
	NewsAnchor       = VoiceProfile{N_Male, -2.0, 1.05, "News Anchor"}

	// --- CHILDREN & TODDLERS ---
	ToddlerBoy       = VoiceProfile{N_Male, 16.0, 1.0, "Toddler Boy"}
	ToddlerGirl      = VoiceProfile{N_Female, 18.0, 1.0, "Toddler Girl"}
	YoungChildBoy    = VoiceProfile{N_Male, 12.0, 1.1, "Small Boy"}
	YoungChildGirl   = VoiceProfile{N_Female, 13.0, 1.1, "Small Girl"}
	PreTeenBoy       = VoiceProfile{N_Male, 3.0, 1.05, "Pre-teen Boy"}
	PreTeenGirl      = VoiceProfile{N_Alt_F, 7.0, 1.05, "Pre-teen Girl"}

	// --- YOUNG ADULTS (Gen Z / College age) ---
	CollegeMale      = VoiceProfile{N_Alt_M, 1.0, 1.15, "College Male"}
	CollegeFemale    = VoiceProfile{N_Alt_F, 2.0, 1.15, "College Female"}
	HyperTeen        = VoiceProfile{N_Female, 4.0, 1.3, "Excited Teenager"}
	BroMale          = VoiceProfile{N_Male, -1.0, 1.1, "Casual Young Man"}

	// --- ADULTS (30s - 40s) ---
	BusinessMale     = VoiceProfile{S_Male, 0.0, 1.0, "Professional Man"}
	BusinessFemale   = VoiceProfile{S_Female, 0.0, 1.0, "Professional Woman"}
	DeepVoiceMan     = VoiceProfile{N_Male, -7.0, 0.9, "Deep Bass Man"}
	WarmMotherly     = VoiceProfile{N_Female, -1.0, 0.95, "Motherly Voice"}
	StrictTeacher    = VoiceProfile{N_Alt_F, -2.0, 0.85, "Strict Teacher"}

	// --- MIDDLE AGED (50s) ---
	MatureMale       = VoiceProfile{N_Alt_M, -3.0, 0.9, "Middle-aged Man"}
	MatureFemale     = VoiceProfile{N_Alt_F, -3.0, 0.9, "Middle-aged Woman"}
	GruffContractor  = VoiceProfile{N_Male, -5.0, 0.95, "Gruff Man"}

	// --- SENIORS (60s+) ---
	Grandpa          = VoiceProfile{N_Alt_M, -4.0, 0.8, "Grandfather"}
	Grandma          = VoiceProfile{N_Alt_F, -2.0, 0.8, "Grandmother"}
	ElderlyRaspy     = VoiceProfile{N_Male, -6.0, 0.75, "Elderly Raspy"}

	// --- EMOTIONAL / CHARACTER STATES ---
	Whisperer        = VoiceProfile{N_Female, 0.0, 0.8, "Whispering Voice"}
	FastTalker       = VoiceProfile{N_Alt_M, 2.0, 1.6, "Fast Talker"}
	SlowMover        = VoiceProfile{N_Male, -2.0, 0.6, "Very Slow Speaker"}
	MonotoneRobot    = VoiceProfile{N_Alt_F, -10.0, 1.0, "Robotic Tone"}
	HighPitchPanic   = VoiceProfile{N_Female, 10.0, 1.4, "Panicked Voice"}
	LowGrumble       = VoiceProfile{N_Male, -10.0, 0.8, "Grumbling Voice"}
)