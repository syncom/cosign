//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/syncom/cosign/pkg/cosign"
)

func Clean() *ffcli.Command {
	var (
		flagset = flag.NewFlagSet("cosign clean", flag.ExitOnError)
	)

	return &ffcli.Command{
		Name:       "clean",
		ShortUsage: "cosign clean <image uri>",
		ShortHelp:  "Remove all signatures from an image",
		FlagSet:    flagset,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) != 1 {
				return flag.ErrHelp
			}

			return CleanCmd(ctx, args[0])
		},
	}
}

func CleanCmd(_ context.Context, imageRef string) error {
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return err
	}
	// TODO: just return the descriptor directly if we have a digest reference.
	desc, err := remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	dstRef, err := cosign.DestinationRef(ref, desc)
	if err != nil {
		return err
	}

	signRef := dstRef.Context().Tag(cosign.Munge(desc.Descriptor))
	fmt.Println(signRef)

	fmt.Println("Deleting signature metadata...")

	err = remote.Delete(signRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}

	return nil
}
