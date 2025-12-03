# gz

gz is a tool for single commit, chained git branches.

## Overview

For years, I have followed a git workflow that consists of feature branches with a single commit where the vast majority of my commits are done via `git commit -a --amend --no-edit`. These feature branches can be chained where a sub-feature branch has a single commit for its local changes, then the single commit for the parent branch, then commits from `main`. One of the tricky parts is when I commit changes to the parent branch. I typically visit all the child branches, reset them to the new `HEAD` of the parent branch, and `cherry-pick` the local changes back into the branch. This is all a little tedious to do by hand. That's what `gz` is for, it provides automations for this workflow.

## Notes

### Commmands

`create` - creates a new branch
`checkout` - checks out an existing branch. I would like for this to only work for gz branches, but there is the problem that the root branch doesn't currently have metadata. Perhaps it should?
`commit` - if there are no commits to this branch, it appends a commit. If there is already a commit, it amends the commit. There is an option to append ... giving you multiple commits in a single feature branch. This is perhaps an anti-pattern but PR reviewers often complain about not being able to see the changes as individual commits.
`rebase` - this will make sure the chain is valid all the way to the root branch. this is intended to be used after doing a `gz commit` to a parent branch.
`sync` - this will walk up to the root branch, pull new changes from upstream, then rebase all parent branches back up to the current branch. This might be better as `gz rebase --sync` or `gz rebase --pull` or something like that.

### Operations

- Visualize the branch chain
- List all the children of a branch
- Navigate up the chain
- Rebase the chain
- Rebase when my parent has already been rebased
- Visualize the entire branch tree
- Search for a brancy by name or changed file
- Detecting that a branch has been merged into the parent branch
