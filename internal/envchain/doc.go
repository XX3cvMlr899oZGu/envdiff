// Package envchain provides a composable, ordered pipeline for applying
// sequences of env map transformations.
//
// A Chain is built by registering named StepFunc values. When Run is called,
// each step receives the output of the previous step. If any step returns an
// error the pipeline halts and the error is returned, annotated with the
// failing step name.
//
// Example:
//
//	chain := envchain.New().
//		Add("filter",  myFilterStep).
//		Add("resolve", myResolveStep).
//		Add("trim",    myTrimStep)
//
//	result, err := chain.Run(rawEnv)
//
// Steps are pure functions: the input map is copied before the first step so
// the original map is never mutated.
package envchain
