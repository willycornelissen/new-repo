# AGENTS.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.


## 5. Architecture Principles

**10 Key Principles:**

1. Well-defined boundaries | 2. Composability | 3. Independence | 4. Individual scale | 5. Explicit communication
2. Replaceability | 7. Deployment independence | 8. State isolation ⚠️ | 9. Observability | 10. Fail independence

### Progressive Documentation Loading

**CRITICAL**: Only load documents relevant to your current task. Do NOT load all documentation at once.

#### Decision Tree: What to Read (Priority Order)

**Implementation tasks (writing code):**

- **Creating controllers, services, or repositories** → `.opencode/docs/coding-patterns.md`
  - Repository pattern, lean controllers, transaction management, entity naming, state isolation
- **Integrating external APIs, third-party services, observability** → `.opencode/docs/integration-patterns.md`
  - Client encapsulation, injection patterns, logging, metrics, circuit breakers, event systems

**Architecture/design tasks** → handled automatically by the `modular-architecture` skill:

- Creating modules, evaluating module boundaries, assessing compliance, maturity assessments

#### Quick Reference by Task Type

| Task Type                         | Primary Doc                    | Notes                        |
| --------------------------------- | ------------------------------ | ---------------------------- |
| New entity/migration              | `.opencode/docs/coding-patterns.md`      | Entity naming section        |
| New controller/service/repository | `.opencode/docs/coding-patterns.md`      | Full patterns doc            |
| External API integration          | `.opencode/docs/integration-patterns.md` | Client encapsulation section |
| Logging/metrics/circuit breakers  | `.opencode/docs/integration-patterns.md` | Resilience sections          |
| Create new module                 | `modular-architecture` skill   | Auto-triggered               |
| Evaluate module boundaries        | `modular-architecture` skill   | Auto-triggered               |
| Architecture compliance check     | `modular-architecture` skill   | Auto-triggered               |
| Maturity assessment               | `modular-architecture` skill   | Auto-triggered               |
| Co-location, depth, suffixes      | `modular-architecture` skill   | P11–P17 structural principles |



---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
