package models

import (
	"github.com/aws/aws-sdk-go/service/glacier"
	"github.com/jrudio/go-plex-client"
)

type Upload struct {
	Metadata              plex.Metadata
	ArchiveCreationOutput glacier.ArchiveCreationOutput
}

type Movie struct {
	Metadata  plex.Metadata
	Unwatched int
}

type Show struct {
	Metadata  plex.Metadata
	Unwatched int
}
