# Pull Request

## Description

<!-- Provide a clear and concise description of what this PR accomplishes -->

### Type of Change

<!-- Please check the type of change your PR introduces -->

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Performance improvement
- [ ] Code refactoring (no functional changes)
- [ ] Documentation update
- [ ] Test improvements
- [ ] CI/CD improvements
- [ ] Other (please describe):

## Related Issues

<!-- Link to any related issues using "Closes #123" or "Relates to #123" -->

- Closes #
- Relates to #

## Changes Made

<!-- List the main changes made in this PR -->

- 
- 
- 

## Testing

<!-- Describe how you tested your changes -->

### Test Coverage

- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed
- [ ] All existing tests pass

### Test Commands Run

```bash
# List the commands you ran to test your changes
make check
go test ./... -v -race
```

### Test Results

<!-- Paste relevant test output or describe test results -->

## Documentation

<!-- Check all that apply -->

- [ ] Code changes are documented with comments
- [ ] Public APIs are documented with godoc comments
- [ ] README.md updated (if applicable)
- [ ] CHANGELOG.md updated (if applicable)
- [ ] Documentation in docs/ updated (if applicable)

## Performance Impact

<!-- If applicable, describe the performance impact of your changes -->

- [ ] No performance impact
- [ ] Performance improved
- [ ] Performance impact exists but acceptable
- [ ] Performance benchmarks included

## Breaking Changes

<!-- If this introduces breaking changes, describe them -->

- [ ] No breaking changes
- [ ] Breaking changes documented below

### Breaking Change Details

<!-- Describe any breaking changes and migration steps -->

## Security Considerations

<!-- Check all that apply -->

- [ ] No security implications
- [ ] Security review requested
- [ ] Input validation added/updated
- [ ] Error handling doesn't expose sensitive information
- [ ] No hardcoded secrets or credentials

## Checklist

<!-- Please check all the boxes that apply -->

### Code Quality

- [ ] Code follows the project's style guidelines
- [ ] Self-review of code performed
- [ ] Code is properly formatted (`make fmt`)
- [ ] Linting passes (`make lint`)
- [ ] No unnecessary debug code or comments left in

### Testing and Validation

- [ ] All tests pass locally (`make test`)
- [ ] Test coverage is adequate
- [ ] Changes have been tested on multiple platforms (if applicable)
- [ ] Edge cases have been considered and tested

### Documentation and Communication

- [ ] Commit messages follow conventional commit format
- [ ] PR description clearly explains the changes
- [ ] Any new dependencies are justified and documented

### Final Checks

- [ ] Changes are backwards compatible (or breaking changes are documented)
- [ ] No merge conflicts
- [ ] Branch is up to date with target branch
- [ ] Ready for review

## Additional Notes

<!-- Add any additional notes, concerns, or context for reviewers -->

## Screenshots/Examples (if applicable)

<!-- Add screenshots, command examples, or output examples if relevant -->

```bash
# Example usage
viesquery --help
```

## Reviewer Notes

<!-- Any specific areas you'd like reviewers to focus on -->

- Please pay attention to:
- Questions for reviewers:
- Areas of concern:

---

**By submitting this PR, I confirm that:**

- [ ] I have read and agree to the project's contributing guidelines
- [ ] My code follows the project's code of conduct  
- [ ] I understand this contribution will be licensed under the project's MIT license
