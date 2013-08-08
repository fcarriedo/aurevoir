aurevoir
========

The simple code review tool

Add message in commit:

`review francisco, joel` - possible list of reviewers. Possible `pre-commit` reject if reviewer not found.
`review francisco.carriedo`

Server hook parses the commit message and (if necessary) creates a review entry in the DB and send an email to the reviewer.

MVP Features:

  * See diffs in web via REST interface:
    * commits:
      * GET `svn/project/platform/commits` gets a paginated list of commits (commit time desc)
    * diffs :
      * GET `svn/project/platform/commits/{commitId}` gets all files within the commit, comments, reviewer, status, 
    * status:
      * PUT `svn/project/platform/commits/{commitId}/status` params: status [LGTM/reject]
    * comments:
      * POST `svn/project/platform/commits/{commitId}/comments` params: line, author, comment
      * GET `svn/project/platform/commits/{commitId}/comments/{commentId}`
      * PUT `svn/project/platform/commits/{commitId}/comments/{commentId}`
  * Comment on lines
  * Ability to LGTM or decline-improve workflow (w/ optional msg) - Notifies commiter
  * Ability to reassign reviewer
  * Supports git

Advanced:

  * Render Markdown in comments
  * Supports SVN

