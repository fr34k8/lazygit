package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickDuringRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick commits from the subcommits view during a rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowGraph = "never"
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("base").
			NewBranch("first-branch").
			NewBranch("second-branch").
			Checkout("first-branch").
			EmptyCommit("one").
			EmptyCommit("two").
			Checkout("second-branch").
			EmptyCommit("three").
			EmptyCommit("four").
			Checkout("first-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-branch"),
				Contains("second-branch"),
				Contains("master"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("base"),
			).
			// copy commit 'three'
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("1 commit copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI two").IsSelected(),
				Contains("CI one"),
				Contains("CI base"),
			).
			SelectNextItem().
			Press(keys.Universal.Edit).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick CI two"),
				Contains("--- Commits ---"),
				Contains("     CI one").IsSelected(),
				Contains("     CI base"),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 1 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				t.Views().Information().Content(DoesNotContain("commit copied"))
			}).
			Lines(
				Contains("--- Pending rebase todos ---"),
				Contains("pick CI two"),
				Contains("--- Commits ---"),
				Contains("     CI three"),
				Contains("     CI one").IsSelected(),
				Contains("     CI base"),
			).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("CI two"),
				Contains("CI three"),
				Contains("CI one").IsSelected(),
				Contains("CI base"),
			)
	},
})
