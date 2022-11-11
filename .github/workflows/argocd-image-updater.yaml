name: build-every-branch

on:
  push:
    branches:
      - '*'

env:
  DOCKERHUB_USERNAME: lyzhang1999

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Set outputs
        id: vars
        run: echo "::set-output name=sha_short::$(git rev-parse --short HEAD)"
      - name: Extract branch name
        shell: bash
        run: echo "##[set-output name=branch;]$(echo ${GITHUB_REF#refs/heads/})"
        id: branch
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ env.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}
      - name: Build backend and push
        uses: docker/build-push-action@v3
        with:
          context: backend
          push: true
          platforms: linux/amd64,linux/arm64
          tags: ${{ env.DOCKERHUB_USERNAME }}/backend:${{ steps.branch.outputs.branch }}-${{ steps.vars.outputs.sha_short }}
      - name: Build frontend and push
        uses: docker/build-push-action@v3
        with:
          context: frontend
          push: true
          platforms: linux/amd64,linux/arm64
          tags: ${{ env.DOCKERHUB_USERNAME }}/frontend:${{ steps.branch.outputs.branch }}-${{ steps.vars.outputs.sha_short }}