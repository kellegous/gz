package internal

import (
	"encoding/hex"
	"encoding/json"

	"github.com/kellegous/gz"
	"github.com/kellegous/poop"
)

type Branch struct {
	Name        string `json:"name"`
	Commits     []SHA  `json:"commits"`
	Parent      string `json:"parent"`
	Description string `json:"description"`
}

func (b *Branch) ToProto() *gz.Branch {
	commits := make([][]byte, 0, len(b.Commits))
	for _, commit := range b.Commits {
		commits = append(commits, commit)
	}

	return &gz.Branch{
		Name:        b.Name,
		Commits:     commits,
		Parent:      b.Parent,
		Description: b.Description,
	}
}

func BranchFromProto(proto *gz.Branch) *Branch {
	commits := make([]SHA, 0, len(proto.Commits))
	for _, commit := range proto.Commits {
		commits = append(commits, SHA(commit))
	}

	return &Branch{
		Name:        proto.Name,
		Commits:     commits,
		Parent:      proto.Parent,
		Description: proto.Description,
	}
}

type SHA []byte

func (s SHA) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(s))
}

func (s *SHA) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	var err error
	*s, err = hex.DecodeString(str)
	return poop.Chain(err)
}
