CREATE TABLE "location" (
  "id" int(11) NOT NULL AUTO_INCREMENT,
  "name" varchar(255) DEFAULT NULL,
  "zip" varchar(10) DEFAULT NULL,
  "lat" float DEFAULT NULL,
  "lon" float DEFAULT NULL,
  "town" varchar(255) DEFAULT NULL,
  "state" varchar(255) DEFAULT NULL,
  PRIMARY KEY ("id")
)

CREATE TABLE "forecast" (
  "location_id" int(11) NOT NULL,
  "index" int(11) NOT NULL,
  "datedesc" varchar(255) DEFAULT NULL,
  "summary" varchar(1024) DEFAULT NULL,
  "forecast" varchar(4096) DEFAULT NULL,
  "high" int(11) DEFAULT NULL,
  "low" int(11) DEFAULT NULL,
  "icon" varchar(1024) DEFAULT NULL,
  "timestamp" bigint(20) DEFAULT NULL,
  PRIMARY KEY ("location_id","index")
);
