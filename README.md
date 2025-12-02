# zkret-santa

Started with a deterministic name selection system using 256 bits of entropy using the go standard library. Because
that clearly wasn't enough, I ensure the final binary is bit-for-bit reproducbile. But what if the toolchain is
compromised? No problem, use [Stage<sup>x</sup>](https://stagex.tools/) for full-source bootstrapped, reproducible and
deterministic builds from 181 bytes of machine code to go. To ensure reprodcubility and provide proof, build
container images multiple times and verify their reproducibilty. Using Github Actions also allows for providing attestions of the
build process meeting SLSA Build L2 compliance (L3 coming soon). Because ZKP is hot right now, allow for the generation
of a proof that the given `<seed> + <participant-file>` pair were provided to the binary.

Because this clearly isn't enough:

* [ ] Trigger workflow to build on Gitlab
* [ ] Trigger workflow to build on Codeberg
* [ ] Wait for all builds and validate hashes prior to publishing
* [ ] SLSA Build L3+ compliance
* [ ] Images on Dockerhub
* [ ] Images on Quay.io
* [ ] Images on Gitlab
* [ ] Cosign signatures from all cloud reproducers
* [ ] Web UI that gets updates over service workers and only install after validating attestations
* [ ] Web UI verifies the attestation of the server running the binary
