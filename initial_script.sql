INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Juan', 'Perez', 12345678, 'jperez', 'jperez1234', 'jperez@mail.com', 'observed');

INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Raul', 'Gonzalez', 12345679, 'rgonzalez', 'rgonzalez1234', 'rgonzalez@mail.com', 'observed');

INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Maria', 'Dominguez', 87654321, 'mdominguez', 'mdominguez1234', 'mdominguez@mail.com', 'observer');

INSERT INTO Users (name, last_name, id_number, username, password, email, type)
VALUES ('Milagros', 'Paredes', 87654321, 'mparedes', 'mparedes1234', 'mparedes@mail.com', 'observer');

INSERT INTO ObservedUsers (user_id, privacy_key, company_name)
VALUES (1, 'juan.perez.12345678', 'company school bus');

INSERT INTO ObservedUsers (user_id, privacy_key, company_name)
VALUES (2, 'raul.gonzalez.12345679', 'company school bus');

INSERT INTO ObserverUsers (user_id)
VALUES (3);

INSERT INTO ObserverUsers (user_id)
VALUES (4);

INSERT INTO SchoolBuses (license_plate, model, brand, license, observed_user_id)
VALUES ('11AAA222', 'Master', 'Renault', '11222', 1);

INSERT INTO SchoolBuses (license_plate, model, brand, license, observed_user_id)
VALUES ('11AAA333', 'Fiat', 'Ducato', '11333', 2);

INSERT INTO Addresses (name, street, number, zipCode, city, state, country, latitude, longitude, observer_user_id)
VALUES ('La Salle', '25 de Mayo', 2569, '3000', 'Santa Fe', 'Santa Fe', 'Argentina', '-31.646020244103223',
        '-60.70579978666576', 3);

INSERT INTO Addresses (name, street, number, floor, apartment, zipCode, city, state, country, latitude, longitude,
                       observer_user_id)
VALUES ('Casa', '25 de Mayo', 2681, '1', 'A', '3000', 'Santa Fe', 'Santa Fe', 'Argentina', '-31.644603142894496',
        '-60.70545280200867', 3);

INSERT INTO Addresses (name, street, number, zipCode, city, state, country, latitude, longitude, observer_user_id)
VALUES ('Dante Alighieri', '25 de Mayo', 2569, '3000', 'Santa Fe', 'Santa Fe', 'Argentina', '-31.646020244103223',
        '-60.70579978666576', 4);

INSERT INTO Addresses (name, street, number, floor, apartment, zipCode, city, state, country, latitude, longitude,
                       observer_user_id)
VALUES ('Casa', '25 de Mayo', 2765, '2', 'B', '3000', 'Santa Fe', 'Santa Fe', 'Argentina', '-31.643698274006496',
        '-60.70517895968019', 4);

INSERT INTO ObservedUsersObserverUsers (observed_user_id, observer_user_id)
VALUES (1, 3);

INSERT INTO ObservedUsersObserverUsers (observed_user_id, observer_user_id)
VALUES (2, 4);

INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id)
VALUES ('Pilar', 'Dominguez', 'La Salle', '08:00:00', '12:00:00', 3);

INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id)
VALUES ('Pia', 'Dominguez', 'La Salle', '08:00:00', '12:00:00', 3);

INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id)
VALUES ('Catalina', 'Paredes', 'Dante Alighieri', '08:00:00', '12:00:00', 4);

INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id)
VALUES ('Augusta', 'Paredes', 'Dante Alighieri', '08:00:00', '12:00:00', 4);
