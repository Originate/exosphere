#!/bin/bash

colorpint notice "This is NOTICE from CLI" # Pipe to stdout with magenta color.
colorpint info "This is INFO from CLI" # Pipe to stdout with green color.
colorpint debug "This is DEBUG from CLI" # Pipe to stdout with  color.
colorpint trace "This is TRACE from CLI" # Pipe to stdout with white color.
colorpint warn "This is WARN from CLI" # Pipe to stdout with yellow color.
colorpint error "This is ERROR from CLI" # Pipe to stderr with red color.
colorpint fatal "This is FATAL from CLI" # Pipe to stderr with bgRed color.

