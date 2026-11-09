package hetzner

import (
	"context"
	"strconv"

	"codedock/internal/modules/servers"
)

func (a *Adapter) Verify(ctx context.Context) error {
	_, err := a.client.ListLocations(ctx)
	return err
}

func (a *Adapter) Offers(ctx context.Context) (servers.ManagedOffers, error) {
	locations, err := a.client.ListLocations(ctx)
	if err != nil {
		return servers.ManagedOffers{}, err
	}
	serverTypes, err := a.client.ListServerTypes(ctx)
	if err != nil {
		return servers.ManagedOffers{}, err
	}
	return buildOffers(locations, serverTypes), nil
}

func buildOffers(locations, serverTypes []map[string]any) servers.ManagedOffers {
	reviewedRegions := map[string]bool{}
	for _, region := range Regions() {
		reviewedRegions[region] = true
	}
	result := servers.ManagedOffers{Provider: servers.ManagedProviderHetzner, Images: Images(), Locations: []servers.ManagedLocation{}, Offers: []servers.ManagedOffer{}}
	for _, location := range locations {
		name := stringify(location["name"])
		if !reviewedRegions[name] {
			continue
		}
		result.Locations = append(result.Locations, servers.ManagedLocation{
			Name:    name,
			City:    stringify(location["city"]),
			Country: stringify(location["country"]),
		})
	}
	live := map[string]map[string]any{}
	for _, serverType := range serverTypes {
		if serverType["deprecation"] != nil {
			continue
		}
		live[stringify(serverType["name"])] = serverType
	}
	for _, tier := range Tiers() {
		serverType, ok := live[tier.Name]
		if !ok {
			continue
		}
		prices := locationPrices(serverType, reviewedRegions)
		if len(prices) == 0 {
			continue
		}
		result.Offers = append(result.Offers, servers.ManagedOffer{Tier: tier, Prices: prices})
	}
	return result
}

func locationPrices(serverType map[string]any, allowed map[string]bool) map[string]float64 {
	prices := map[string]float64{}
	entries, _ := serverType["prices"].([]any)
	for _, entry := range entries {
		price, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		location := stringify(price["location"])
		if !allowed[location] {
			continue
		}
		monthly, _ := price["price_monthly"].(map[string]any)
		value, err := strconv.ParseFloat(stringify(monthly["gross"]), 64)
		if err != nil {
			continue
		}
		prices[location] = value
	}
	return prices
}
