package conftypes

import "github.com/derision-test/go-mockgen/input/banana"

type UnifiedWatchable interface {
	Watchable
	UnifiedQuerier
}

type UnifiedQuerier interface {
	ServiceConnectionQuerier
	SiteConfigQuerier
}

type WatchableSiteConfig interface {
	SiteConfigQuerier
	Watchable
}

type ServiceConnectionQuerier interface {
	ServiceConnections() ServiceConnections
}

type SiteConfigQuerier interface {
	SiteConfig() banana.Banana
}

type Watchable interface {
	Watch(func())
}
