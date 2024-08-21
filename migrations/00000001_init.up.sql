-- world_logo ==========================================================================================================
CREATE TABLE world_logo (
   id varchar (26) NOT NULL,
   name varchar(256) NOT NULL,
   logo_path text NOT NULL,
   src_key varchar(256) NOT NULL,
   created_at timestamp with time zone NOT NULL DEFAULT NOW(),
   updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX world_logo_id_uindex ON world_logo USING btree (id);
CREATE UNIQUE INDEX world_logo_src_key_uindex ON world_logo USING btree (src_key);
CREATE INDEX world_logo_name_index ON world_logo USING btree (name);
