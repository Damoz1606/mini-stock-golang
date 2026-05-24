--
-- PostgreSQL database dump
--


-- Dumped from database version 15.17
-- Dumped by pg_dump version 15.17

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: atlas_schema_revisions; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA atlas_schema_revisions;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: atlas_schema_revisions; Type: TABLE; Schema: atlas_schema_revisions; Owner: -
--

CREATE TABLE atlas_schema_revisions.atlas_schema_revisions (
    version character varying NOT NULL,
    description character varying NOT NULL,
    type bigint DEFAULT 2 NOT NULL,
    applied bigint DEFAULT 0 NOT NULL,
    total bigint DEFAULT 0 NOT NULL,
    executed_at timestamp with time zone NOT NULL,
    execution_time bigint NOT NULL,
    error text,
    error_stmt text,
    hash character varying NOT NULL,
    partial_hashes jsonb,
    operator_version character varying NOT NULL
);


--
-- Name: attributes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.attributes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "categoryId" uuid NOT NULL,
    name text NOT NULL,
    description text,
    "createdAt" bigint NOT NULL,
    "updatedAt" bigint NOT NULL
);


--
-- Name: attributes_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.attributes_items (
    "categoryId" uuid NOT NULL,
    "attributeId" uuid NOT NULL,
    "itemId" uuid NOT NULL
);


--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    "createdAt" bigint NOT NULL,
    "updatedAt" bigint NOT NULL
);


--
-- Name: items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.items (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    "categoryId" uuid NOT NULL,
    name text NOT NULL,
    description text,
    "createdAt" bigint NOT NULL,
    "updatedAt" bigint NOT NULL
);


--
-- Name: v_attributes; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_attributes AS
 SELECT attributes.id,
    attributes."categoryId",
    attributes.name,
    attributes.description
   FROM public.attributes;


--
-- Name: v_categories; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_categories AS
 SELECT categories.id,
    categories.name
   FROM public.categories;


--
-- Name: v_items; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.v_items AS
 SELECT items.id,
    items."categoryId",
    items.name,
    items.description
   FROM public.items;


--
-- Name: atlas_schema_revisions atlas_schema_revisions_pkey; Type: CONSTRAINT; Schema: atlas_schema_revisions; Owner: -
--

ALTER TABLE ONLY atlas_schema_revisions.atlas_schema_revisions
    ADD CONSTRAINT atlas_schema_revisions_pkey PRIMARY KEY (version);


--
-- Name: attributes_items attributes_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes_items
    ADD CONSTRAINT attributes_items_pkey PRIMARY KEY ("categoryId", "attributeId", "itemId");


--
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: items items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);


--
-- Name: attributes uq_attributes_category_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT uq_attributes_category_name UNIQUE ("categoryId", name);


--
-- Name: attributes_items uq_attributes_items_attribute_item; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes_items
    ADD CONSTRAINT uq_attributes_items_attribute_item UNIQUE ("attributeId", "itemId");


--
-- Name: categories uq_categories_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT uq_categories_name UNIQUE (name);


--
-- Name: items uq_items_category_name; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT uq_items_category_name UNIQUE ("categoryId", name);


--
-- Name: attributes fk_attributes_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT fk_attributes_category FOREIGN KEY ("categoryId") REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: attributes_items fk_attributes_items_attribute; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes_items
    ADD CONSTRAINT fk_attributes_items_attribute FOREIGN KEY ("attributeId") REFERENCES public.attributes(id) ON DELETE CASCADE;


--
-- Name: attributes_items fk_attributes_items_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes_items
    ADD CONSTRAINT fk_attributes_items_category FOREIGN KEY ("categoryId") REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: attributes_items fk_attributes_items_item; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.attributes_items
    ADD CONSTRAINT fk_attributes_items_item FOREIGN KEY ("itemId") REFERENCES public.items(id) ON DELETE CASCADE;


--
-- Name: items fk_items_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.items
    ADD CONSTRAINT fk_items_category FOREIGN KEY ("categoryId") REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


