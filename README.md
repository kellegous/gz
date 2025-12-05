# gz

gz is a tool for single commit, chained git branches.

## Overview

For years, I have followed a git workflow that consists of feature branches with a single commit where the vast majority of my commits are done via `git commit -a --amend --no-edit`. These feature branches can be chained where a sub-feature branch has a single commit for its local changes, then the single commit for the parent branch, then commits from `main`. One of the tricky parts is when I commit changes to the parent branch. I typically visit all the child branches, reset them to the new `HEAD` of the parent branch, and `cherry-pick` the local changes back into the branch. This is all a little tedious to do by hand. That's what `gz` is for, it provides automations for this workflow.

## Notes

How does one resolve a branch once it is merged? Imagine:

```
main <- feature-1 <-  feature-2
```

`feature-1` is merged into `main` but it is merged via a squash merge. We wanted an automated way to detect the merge and end up
in this state.

```
main (containing feature-1 changes in squashed commit) <- feature-2
```

This kind of works if `feature-1` is a single commit because `git` detects the single commit as a cherry-pick. But if `feature-1` is multiple commits, a rebase will create a merge conflict. One way around this is to have an explicit command to re-parent ... or to have a command to assert that a branch is in the parent, so blindly reset the local branch.
