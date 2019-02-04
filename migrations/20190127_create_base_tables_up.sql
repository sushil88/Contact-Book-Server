CREATE TABLE "users" (
  "id" serial,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  "email" varchar(255) NOT NULL DEFAULT '',
  "first_name" varchar(255) NOT NULL,
  "middle_name" varchar(255) NOT NULL DEFAULT '',
  "last_name" varchar(255) NOT NULL DEFAULT '',
  "active" bool NOT NULL DEFAULT true,
  "password" varchar(255) NOT NULL DEFAULT '',
  CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX  "users_email_uidx" ON "users" USING btree("email");

CREATE TABLE "contacts" (
  "id" serial,
  "user_id" integer NOT NULL,
  "first_name" varchar(255) NOT NULL,
  "middle_name" varchar(255) NOT NULL default '',
  "last_name" varchar(255) NOT NULL default '',
  "email_address" varchar(255),
  "phone_number" varchar(255),
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  CONSTRAINT "contacts_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "contacts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE UNIQUE INDEX  "contacts_email_address_uidx" ON "contacts" USING btree("email_address");
