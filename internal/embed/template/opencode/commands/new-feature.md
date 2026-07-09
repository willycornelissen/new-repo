---
description: Create a new feature end-to-end. Generates TDD, PRD, roadmap, plan, code (src/<feature>/), review, and documentation. Usage: /new-feature <feature-name>
---

You are creating a new feature called "$ARGUMENTS". Execute the full pipeline below. Use `$ARGUMENTS` as the feature name/slug for all file paths.

## Step 1 — TDD (Technical Design Document)

Load the technical-design-doc-creator skill. Create a Technical Design Document for the feature "$ARGUMENTS". At the end, save the complete TDD to `specification/$ARGUMENTS/tdd.md`.

## Step 2 — PRD (Product Requirements Document)

Load the tlc-spec-driven skill. Read `specification/$ARGUMENTS/tdd.md` and use the Specify phase to create a Product Requirements Document from it. At the end, save the complete PRD to `specification/$ARGUMENTS/prd.md`.

## Step 3 — Roadmap

Load the tlc-spec-driven skill. Read `specification/$ARGUMENTS/tdd.md` and `specification/$ARGUMENTS/prd.md`. Use the Design and Tasks phases to create a structured roadmap with prioritized features, task breakdown, and estimated dependencies. At the end, save the complete roadmap to `specification/$ARGUMENTS/roadmap.md`.

## Step 4 — Plan

Load the tlc-spec-driven skill. Read `specification/$ARGUMENTS/prd.md` and `specification/$ARGUMENTS/roadmap.md`. Use the Design and Tasks phases to create a detailed implementation plan with task breakdown, dependencies, and execution order. At the end, save the complete plan to `specification/$ARGUMENTS/plan.md`.

## Step 5 — Generate Code

Load the tlc-spec-driven skill. Read `specification/$ARGUMENTS/plan.md` and execute the Tasks and Implement phases to generate the planned code. Generate all code inside the `src/$ARGUMENTS/` directory following the plan's task breakdown, dependencies, and execution order. Write tests as specified in the plan.

## Step 6 — Code Review

Load the code-review-skill. Review all code in `src/$ARGUMENTS/` and evaluate its coverage and alignment against the requirements in `specification/$ARGUMENTS/prd.md` and the tasks in `specification/$ARGUMENTS/plan.md`. At the end, save the complete review report to `specification/$ARGUMENTS/review.md`.

## Step 7 — Documentation

Load the docs-writer skill. Review the code in `src/$ARGUMENTS/` and the specs in `specification/$ARGUMENTS/`, then create comprehensive documentation for the feature. At the end, save all documentation files inside the `documentation/$ARGUMENTS/` directory.

---

After all 7 steps are complete, summarize what was created with the file paths.
