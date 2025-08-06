#!/bin/bash

# Release script for glo
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' 

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

if [ $# -eq 0 ]; then
    print_error "Please provide a version number (e.g., v1.0.0)"
    exit 1
fi

VERSION=$1

if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Version must be in format vX.Y.Z (e.g., v1.0.0)"
    exit 1
fi

print_status "Starting release process for $VERSION"

CURRENT_BRANCH=$(git branch --show-current)
if [[ $CURRENT_BRANCH != "main" && $CURRENT_BRANCH != "release/"* ]]; then
    print_warning "You're not on main or a release branch. Current branch: $CURRENT_BRANCH"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

if [[ -n $(git status --porcelain) ]]; then
    print_error "Working directory is not clean. Please commit or stash changes."
    exit 1
fi

print_status "Running tests..."
if ! go test ./...; then
    print_error "Tests failed. Please fix before releasing."
    exit 1
fi

print_status "Building application..."
if ! go build -o glo .; then
    print_error "Build failed."
    exit 1
fi

print_status "Testing built binary..."
if ! ./glo --version; then
    print_error "Built binary doesn't work."
    exit 1
fi

print_status "Creating and pushing tag $VERSION..."
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION

print_status "Release $VERSION has been created and pushed!"
print_status "GitHub Actions will now build and publish the release."
print_status "Check the progress at: https://github.com/DinethDilhara/glo/actions"

print_warning "Next steps:"
echo "1. Wait for GitHub Actions to complete the release"
echo "2. Verify the release at: https://github.com/DinethDilhara/glo/releases"
echo "3. Test Homebrew installation: brew install dinethdhilhara/tap/glo"
