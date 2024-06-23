-- +goose Up
-- +goose StatementBegin
INSERT INTO stocks.stocks (id, available, reserved)
VALUES
  (1076963, 100, 10),
  (1148162, 150, 20),
  (1625903, 200, 30),
  (2618151, 50, 5),
  (2956315, 300, 25),
  (2958025, 75, 10),
  (3596599, 90, 15),
  (3618852, 110, 10),
  (4288068, 60, 8),
  (4465995, 130, 18),
  (4487693, 45, 6),
  (4669069, 85, 9),
  (4678287, 100, 20),
  (4678816, 120, 14),
  (4679011, 70, 7),
  (4687693, 95, 12),
  (4996014, 110, 15),
  (5097510, 140, 17),
  (5415913, 160, 19),
  (5647362, 180, 22);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE stocks.stocks;
-- +goose StatementEnd
