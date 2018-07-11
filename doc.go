// Package btsync is a library that can be used to synchronize tables, column families, qualifiers and
// retention policies for a Cloud BigTable instance.
//
// It works by reading the desired state of the BT tables from a set of configuration files and diffing that against the
// CBT instance. The result of this process is a set of tasks that need to be run in order to bring the CBT instance to
// the desired state.
package btsync
