package security

import (
	"net/http"
	"sync"
	"time"
)

// Structure pour gérer les requêtes d'un client
type client struct {
	requests []time.Time // Historique des requêtes
}

// Map pour stocker les clients et un mutex pour la sécurité des threads
var clients = make(map[string]*client)
var mu sync.Mutex

// Paramètres de la limitation de débit
const maxRequests = 12              // Nombre maximum de requêtes autorisées
const timeWindow = 10 * time.Second // Intervalle de temps pour limiter les requêtes

// Fonction principale pour gérer les requêtes avec limitation
func RateLimitedHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // Identifier le client par son IP

		// Verrouiller l'accès à la map des clients
		mu.Lock()
		defer mu.Unlock()

		// Récupérer les informations du client ou créer un nouvel enregistrement
		c, exists := clients[ip]
		if !exists {
			c = &client{requests: []time.Time{}}
			clients[ip] = c
		}

		// Nettoyer les requêtes hors de la fenêtre de temps
		now := time.Now()
		validRequests := []time.Time{}
		for _, t := range c.requests {
			if now.Sub(t) <= timeWindow {
				validRequests = append(validRequests, t)
			}
		}
		c.requests = validRequests

		// Vérifier si la limite de requêtes est atteinte
		if len(c.requests) >= maxRequests {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Ajouter la nouvelle requête à l'historique
		c.requests = append(c.requests, now)

		next(w, r)
	}
}

// Nettoyage périodique des clients inactifs
func CleanupInactiveClients() {
	for {
		time.Sleep(1 * time.Minute) // Nettoyer toutes les minutes

		mu.Lock()
		for ip, c := range clients {
			// Si un client n'a pas envoyé de requêtes depuis un certain temps, on le supprime
			if len(c.requests) == 0 || time.Since(c.requests[len(c.requests)-1]) > 2*timeWindow {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}
