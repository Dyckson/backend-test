-- Seed data para popular a tabela beer_styles com 3 cervejas
INSERT INTO beer_styles (name, temp_min, temp_max) VALUES 
    ('IPA', 4.0, 7.0),
    ('Lager', 2.0, 5.0),
    ('Stout', 6.0, 10.0)
ON CONFLICT (name) DO NOTHING;
