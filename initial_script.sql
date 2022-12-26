INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Juan', 'Perez', 12345678, 'jperez', 'jperez1234', 'jperez@mail.com', 'observed');

INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Maria', 'Dominguez', 87654321, 'mdominguez', 'mdominguez1234', 'mdominguez@mail.com', 'observer');

INSERT INTO SchoolBuses (license_plate, model, brand, school_bus_license) 
VALUES ('11AAA222', 'Master', 'Renault', '11222');

INSERT INTO ObservedUsers (user_id, privacy_key, company_name, school_bus_id)
VALUES (1, 'juan.perez.12345678', 'company school bus', 1);

INSERT INTO ObserverUsers (user_id)
VALUES (2);

INSERT INTO Addresses (street, number, floor, apartament, zipCode, city, state, country, latitude, longitude, observer_user_id) 
VALUES ('25 de Mayo', 2864, '1', 'A', '3000', 'Santa Fe', 'Santa Fe', 'Argentina', '-31.642672604529235', '-60.70456270200879', 2);

INSERT INTO ObservedUsersObserverUsers (observed_user_id, observer_user_id)
VALUES (1, 2);

INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id)
VALUES ('Pilar', 'Dominguez', 'La Salle', '08:00:00', '12:00:00', 2);
