# stackpr - Stack Pull Requests CLI Tool

Stackpr is a command-line tool for managing stacked pull requests in Git. Stacked PRs are a workflow where changes are built on top of each other in a linear sequence, making it easier to review and merge complex changes incrementally.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/stackpr.git
cd stackpr

# Build the binary
go build -o stackpr

# Move it to a directory in your PATH
mv stackpr /usr/local/bin/
```

## Core Concepts

Stackpr manages a stack of Git branches where:

- Each branch builds on the previous one
- Changes flow from the base branch up through the stack
- Branches can be reordered, added, or removed while maintaining the stack

The tool keeps track of these relationships in a `.stackpr.yaml` config file at the root of your repository.

### Visual Representation of Stacked PRs

```
                   main
                     |
                     ↓
  PR #1: ------> feature-a
                     |
                     ↓
  PR #2: ------> feature-b
                     |
                     ↓
  PR #3: ------> feature-c
```

Each branch builds on the previous one, forming a linear stack of changes. When `main` is updated, changes flow down through the stack:

```
                  main (updated)
                     |
                     ↓
                 feature-a (updated)
                     |
                     ↓
                 feature-b (updated)
                     |
                     ↓
                 feature-c (updated)
```

And when PRs are merged, the stack adjusts automatically:

```
                    main
                     |     [PR #1 merged]
                     |     feature-a → main
                     ↓
                 feature-b (now based on main)
                     |
                     ↓
                 feature-c (updated)
```

## Quick Start

```bash
# Initialize stackpr in your repository
stackpr init

# Set your base branch (typically main or master)
stackpr config base main

# Create a new branch for your first change
git checkout main
stackpr new feature-a

# Make your changes and commit them
# ...

# Create another branch that builds on feature-a
git checkout feature-a
stackpr new feature-b

# Make more changes and commit them
# ...

# Sync changes through the stack (e.g., if main was updated)
stackpr sync

# Push all branches to remote
stackpr push
```

## Command Reference Table

| Command | Description | Flags | Example |
|---------|-------------|-------|---------|
| `init` | Initialize stackpr in the current repository | - | `stackpr init` |
| `config <key> [value]` | View or modify configuration values | - | `stackpr config base main` |
| `new <branch>` | Add a new branch to the stack | - | `stackpr new feature-name` |
| `list` | List all branches in the stack in order | - | `stackpr list` |
| `sync` | Sync all branches in the stack | - | `stackpr sync` |
| `push` | Push all branches to the remote repository | `--force`, `--debug` | `stackpr push --force` |
| `status` | Show the status of all branches in the stack | - | `stackpr status` |
| `merge <branch>` | Merge a branch into its parent | - | `stackpr merge feature-a` |
| `reorder <branch1> <branch2> ...` | Reorder branches in the stack | - | `stackpr reorder a b c` |
| `remove <branch>` | Remove a branch from the stack | - | `stackpr remove feature-b` |
| `help [command]` | Get help for any command | - | `stackpr help sync` |

## Command Reference

### `init`

Initialize stackpr in the current repository.

```bash
stackpr init
```

This creates a `.stackpr.yaml` configuration file with default settings.

### `config <key> [value]`

View or modify configuration values.

```bash
# View the current base branch
stackpr config base

# Set the base branch
stackpr config base main

# View the sync mode
stackpr config syncMode

# Change the sync mode
stackpr config syncMode merge  # Options: rebase, merge, reset
```

### `new <branch>`

Add a new branch to the stack.

```bash
# Create a branch based on the current branch
stackpr new feature-name

# The branch is created and added to the stack
```

### `list`

List all branches in the stack in order.

```bash
stackpr list
```

Example output:
```
Stacked branches:
1: feature-a
2: feature-b
3: feature-c
```

### `sync`

Sync all branches in the stack. This updates each branch with changes from its parent.

```bash
stackpr sync
```

The sync operation uses the configured sync mode (rebase, merge, or reset).

### `push [--force]`

Push all branches in the stack to the remote repository.

```bash
# Normal push
stackpr push

# Force push (use with caution)
stackpr push --force
```

### `status`

Show the status of all branches in the stack.

```bash
stackpr status
```

Example output:
```
feature-a: Implement user authentication (a1b2c3d)
feature-b: Add email validation (e4f5g6h)
feature-c: Implement password reset (i7j8k9l)
```

### `merge <branch>`

Merge a branch into its parent.

```bash
# Merge feature-a into the base branch
stackpr merge feature-a

# Merge feature-b into feature-a
stackpr merge feature-b
```

### `reorder <branch1> <branch2> ...`

Reorder branches in the stack.

```bash
# Reorder to put feature-c before feature-b
stackpr reorder feature-a feature-c feature-b
```

### `remove <branch>`

Remove a branch from the stack (doesn't delete the branch from Git).

```bash
stackpr remove feature-b
```

### `completion [bash|zsh|fish|powershell]`

Generate shell completion scripts.

```bash
stackpr completion zsh
```

### `help [command]`

Get help for any command.

```bash
stackpr help sync
```

## How It Works - Command Internals

This section explains in detail how each stackpr command works internally.

### `init`

The `init` command:
1. Creates a new `.stackpr.yaml` configuration file in the repository root
2. Sets default values for the following configuration options:
   - `remote`: Set to "origin" (the default remote name in Git)
   - `syncMode`: Set to "rebase" (the default sync strategy)
   - `base`: Initially empty, to be set separately
   - `branches`: Initially empty array

The configuration file uses YAML format and is structured as follows:
```yaml
remote: origin
base: main
branches:
  - branch1
  - branch2
syncMode: rebase
```

### `config <key> [value]`

The `config` command reads or modifies the `.stackpr.yaml` configuration file:

1. If only a key is provided (e.g., `stackpr config base`), it reads the value from the config file and displays it
2. If both key and value are provided (e.g., `stackpr config base main`), it updates the config file with the new value
3. For the `syncMode` option, it validates that the value is one of: "rebase", "merge", or "reset"

The config file is loaded at the start of each command execution and saved back to disk if modified.

### `new <branch>`

The `new` command adds a branch to the stack:

1. Checks if the branch already exists in Git:
   - If it exists, it outputs a message but continues
   - If not, it creates a new Git branch from the current HEAD
2. Updates the `.stackpr.yaml` configuration by appending the branch to the `branches` array (if not already present)
3. The branch is created based on the current checked-out branch, following the stacked structure

This command doesn't change the current Git branch - it only creates or adds the branch to the configuration.

### `list`

The `list` command:
1. Reads the `.stackpr.yaml` configuration file
2. Extracts the `branches` array
3. Outputs each branch with its position in the stack
4. If no branches are configured, it displays "No stacked branches"

This is a read-only command that doesn't modify the repository or config file.

### `sync`

The `sync` command is one of the most complex operations:

1. Reads the `.stackpr.yaml` configuration to get the base branch, list of branches, and sync mode
2. For each branch in the stack, starting from the first:
   - Determines the parent branch (base for the first branch, previous branch for others)
   - Checks out the branch
   - Performs the synchronization based on the configured `syncMode`:
     - **rebase**: Uses `git rebase <parent>` to reapply commits from the branch on top of the parent
     - **merge**: Uses `git merge <parent>` to merge changes from the parent
     - **reset**: Uses `git reset --hard <parent>` then applies saved changes

The sync operation ensures that changes from parent branches flow down to their child branches, maintaining the stacked relationship.

### `push`

The `push` command:
1. Reads the `.stackpr.yaml` configuration to get the remote name and branch list
2. For each branch in the stack:
   - Verifies the branch exists locally
   - Pushes the branch to the remote using `git push <remote> <branch>`
   - If `--force` flag is provided, uses force push (`git push --force`) which can overwrite remote history

This command requires that all branches in the stack already exist locally.

### `status`

The `status` command:
1. Reads the `.stackpr.yaml` configuration to get the list of branches
2. For each branch:
   - Gets the latest commit information (message and hash)
   - Displays the branch name, commit message, and abbreviated commit hash
3. If a branch doesn't exist, it shows a warning but continues with the other branches

This is a read-only command that provides a summary of the current state of the stack.

### `merge <branch>`

The `merge` command:
1. Reads the `.stackpr.yaml` configuration to determine the branch's parent:
   - If it's the first branch in the stack, the parent is the base branch
   - Otherwise, the parent is the previous branch in the stack
2. Checks out the parent branch
3. Merges the specified branch into the parent using `git merge <branch>`
4. If successful, it doesn't modify the stack configuration - the branch remains in the config

Note that after merging, you typically want to run `sync` to update dependent branches.

### `reorder <branch1> <branch2> ...`

The `reorder` command:
1. Verifies that all specified branches are in the current stack
2. Verifies that the number of branches provided matches the current stack size
3. Updates the `.stackpr.yaml` configuration with the new branch order
4. Does not modify Git branches or their relationships - only updates the configuration

After reordering, you typically want to run `sync` to rebuild the branch relationships according to the new order.

### `remove <branch>`

The `remove` command:
1. Checks if the specified branch is in the stack configuration
2. Removes the branch from the `branches` array in the `.stackpr.yaml` configuration
3. Does not delete the Git branch itself - only removes it from the stack configuration

After removing a branch, you might need to run `sync` to update the relationships between the remaining branches.

### `completion [bash|zsh|fish|powershell]`

The `completion` command:
1. Generates shell completion scripts for the specified shell
2. Outputs the script to stdout, which can be redirected to a file or sourced directly
3. Supports bash, zsh, fish, and powershell shell formats

This command helps users set up tab completion for stackpr commands in their preferred shell.

### `help [command]`

The `help` command:
1. With no arguments, displays a list of all available commands and general usage information
2. With a command argument, displays detailed help for that specific command
3. Shows the command syntax, available flags, and a brief description

## Common Workflows

### Creating a New Feature Stack

```bash
# Start from your main branch
git checkout main

# Initialize stackpr if not already done
stackpr init
stackpr config base main

# Create your first feature branch
stackpr new feature-base
# Make changes, commit

# Add the next part of the feature
stackpr new feature-ui
# Make changes, commit

# Add final part of the feature
stackpr new feature-tests
# Make changes, commit

# Push all the branches
stackpr push
```

### Handling Changes to Base Branch

```bash
# When main is updated (e.g., other PRs merged)
git checkout main
git pull

# Sync the changes through your stack
stackpr sync

# Resolve any conflicts if necessary

# Push the updated stack
stackpr push
```

### Merging the Stack Incrementally

```bash
# Merge the first PR
stackpr merge feature-base

# Update the stack after the PR is merged
git checkout main
git pull
stackpr sync

# Merge the next PR
stackpr merge feature-ui

# Continue until all PRs are merged
```

### Handling Merge Conflicts

```bash
# When conflicts occur during sync
stackpr sync
# If conflicts occur, the command will fail

# Manually resolve conflicts
git checkout <conflicting-branch>
# Fix conflicts
git add .
git commit -m "Resolve conflicts"

# Continue syncing
stackpr sync
```

### Advanced Conflict Resolution Strategies

Stackpr handles various types of merge conflicts that can occur in stacked PR workflows:

#### Direct Conflicts Between Base and Feature Branch

When your base branch (e.g., main) and a feature branch modify the same file:

```bash
# When main has changes that conflict with your feature branch
git checkout main
git pull  # Get latest changes
stackpr sync  # Will encounter conflicts

# Git will mark conflicts in the files
# Edit the files to resolve conflicts

git add <conflicted-files>
git commit -m "Resolve conflicts with main"

# Continue the sync to update all dependent branches
stackpr sync
```

#### Conflicts in Stacked Branches

When changes to a parent branch conflict with child branches:

```bash
# Scenario: feature-a was modified after feature-b and feature-c were created
# Sync will propagate changes but might cause conflicts

stackpr sync  # Attempt to sync the stack

# If conflicts occur in feature-b:
git checkout feature-b
# Resolve conflicts manually
git add .
git commit -m "Resolve conflicts from feature-a changes"

# If conflicts occur in feature-c:
git checkout feature-c
# Resolve conflicts manually
git add .
git commit -m "Resolve conflicts from feature-b changes"

# Final sync to ensure everything is updated
stackpr sync
```

#### Edit/Edit Conflicts

When the same lines are edited differently in two branches:

```bash
# This is the most common and complex type of conflict
# Example: both main and feature-a modified the same function

git checkout feature-a  # Where the conflict occurred
# Open the file and look for conflict markers:
# <<<<<<< HEAD
# Current branch version
# =======
# Incoming branch version
# >>>>>>> incoming-branch

# Edit the file to create a combined solution
git add <conflicted-file>
git commit -m "Resolve edit conflict"

# Resume sync
stackpr sync
```

#### Tips for Conflict Resolution

1. **Understand the conflict**: Before resolving, understand what changed in both branches
2. **Use visual diff tools**: Tools like VSCode, GitKraken, or `git mergetool` can help visualize conflicts
3. **Communicate with teammates**: If the conflict involves code from another developer, discuss the best resolution
4. **Test after resolving**: Always test your code after resolving conflicts to ensure it still works as expected
5. **Consider changing sync mode**: If rebase conflicts are too complex, try `stackpr config syncMode merge` temporarily

### Reordering and Reworking

```bash
# Reorder branches
stackpr reorder feature-base feature-tests feature-ui

# Remove a branch from the stack
stackpr remove feature-tests

# Add a new branch in the middle
git checkout feature-base
stackpr new feature-middleware
git checkout feature-middleware
# Make changes, commit

# Sync to update the stack
stackpr sync
```

## Advanced Use Cases

### Working with Multiple Stacks

You can work with multiple stacks in the same repository by switching between different base branches.

```bash
# Create a stack on main
stackpr config base main
stackpr new main-feature-a
stackpr new main-feature-b

# Create a different stack on a release branch
git checkout release/v1.0
stackpr config base release/v1.0
stackpr new hotfix-a
stackpr new hotfix-b
```

### Using Different Sync Strategies

```bash
# Use rebase (default - creates clean history)
stackpr config syncMode rebase

# Use merge (preserves merge commits)
stackpr config syncMode merge

# Use reset (forces branches to match - destructive)
stackpr config syncMode reset
```

### Working with GitHub/GitLab Flow

```bash
# Create PRs for each branch in the stack
# Work through the reviews for each

# When feature-a is approved and merged:
git checkout main
git pull
stackpr sync
# This updates feature-b to be based on the updated main
```

## Troubleshooting

### Fix Misconfigured Stack

If your stack configuration is out of sync with git:

```bash
# Manually edit .stackpr.yaml to match the desired stack
# Then sync to ensure everything is connected
stackpr sync
```

### Recover from Failed Rebase

If a rebase fails during sync:

```bash
# Abort the rebase
git rebase --abort

# Fix the stack configuration if needed
# Try a different sync mode
stackpr config syncMode merge
stackpr sync
```

## Tips and Best Practices

1. Keep branches focused on a single logical change to make review easier
2. Regularly sync your stack to avoid complex conflicts
3. Consider using `--force` with caution when pushing rebased branches
4. Create PRs in the order of the stack for easier review
5. Wait for each PR to be merged before merging dependent PRs

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
