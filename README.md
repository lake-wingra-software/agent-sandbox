Agent Sandbox

This agent should be able to complete the following workflow:

- A request is presented to LLM via a tagged comment in a pull request.
- The LLM uses a tool call to read a single file to reply to the request.
- The LLM uses a tool call to rewrite a file to reply to the request.
- The changes are pushed to the pull request branch.

NEXT:
- The LLM has the (limited) full conversation history in context.