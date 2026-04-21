## OpenSpec + Teambuilder

When using OpenSpec, route to project personas if they exist in `.claude/agents/`:

- `/opsx:explore` or `openspec-explore` → use `analyst.md` as the thinking partner
- `/opsx:propose` or `openspec-propose` → use `architect.md` to drive artifact creation
- `/opsx:apply` or `openspec-apply-change` → run the following loop until the Reviewer reports no blocking issues:
  1. Infer the next pending task type from `tasks.md` and implement it using the appropriate persona:
     - Design/UX tasks → `designer.md`
     - Implementation tasks → `programmer.md` or `programmer-<variant>.md`
     - Testing tasks → `tester.md`
     - Proceed without a persona if the relevant one doesn't exist yet
  2. Mark the task complete in `tasks.md` (`- [ ]` → `- [x]`)
  3. Run the Reviewer (`reviewer.md`) against all changes made in this iteration
  4. If the Reviewer raises **blocking issues**: fix them using the appropriate persona, then return to step 3
  5. If the Reviewer raises **warnings only**: note them and continue to the next pending task (step 1)
  6. If all tasks in `tasks.md` are complete and the Reviewer has no blocking issues: stop and report a summary of what was done and any outstanding warnings
