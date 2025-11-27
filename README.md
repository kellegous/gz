# gz

gz is a tool for single commit, chained git branches.

## Overview

For years, I have followed a git workflow that consists of feature branches with a single commit where the vast majority of my commits are done via `git commit -a --amend --no-edit`. These feature branches can be chained where a sub-feature branch has a single commit for its local changes, then the single commit for the parent branch, then commits from `main`. One of the tricky parts is when I commit changes to the parent branch. I typically visit all the child branches, reset them to the new `HEAD` of the parent branch, and `cherry-pick` the local changes back into the branch. This is all a little tedious to do by hand. That's what `gz` is for, it provides automations for this workflow.
