package netex

import "github.com/plugimt/transat-backend/models"

func (n *NetexService) GetLines() ([]models.Line, error) {
	lines := []models.Line{}

	rows, err := n.db.Query("SELECT id, name, transport_mode, public_code, colour, text_colour FROM NETEX_Line ORDER BY route_sort_order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var line models.Line
		err := rows.Scan(&line.ID, &line.Name, &line.TransportMode, &line.Number, &line.BackgroundColour, &line.ForegroundColour)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (n *NetexService) GetLine(id string) (*models.Line, error) {
	var line models.Line
	err := n.db.QueryRow("SELECT id, name, transport_mode, public_code, colour, text_colour FROM NETEX_Line WHERE id = $1", id).Scan(&line.ID, &line.Name, &line.TransportMode, &line.Number, &line.BackgroundColour, &line.ForegroundColour)
	if err != nil {
		return nil, err
	}
	return &line, nil
}
