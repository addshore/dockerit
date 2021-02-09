package cmd

import "testing"

func TestSum(t *testing.T) {
    matches := imageRefMatchesImageName("composer", "composer")
    if !matches{
       t.Errorf("ERROR: %s, want: %s.", "A", "B")
    }
}
