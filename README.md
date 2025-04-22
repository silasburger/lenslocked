I created a webapp in Go with hand-rolled authentication.
It's deployed on a DigitalOcean VM using a containerized
setup with Caddy as the reverse proxy and PostgreSQL as the
database. Users can upload photos to the file system and
organize them into shareable galleries.