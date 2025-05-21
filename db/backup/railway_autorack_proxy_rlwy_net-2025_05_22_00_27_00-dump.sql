--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8 (Debian 16.8-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Homebrew)

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
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: realcampus_privacy_level; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.realcampus_privacy_level AS ENUM (
    'PUBLIC',
    'PRIVATE'
);


ALTER TYPE public.realcampus_privacy_level OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: articles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles (
    id_articles integer NOT NULL,
    name character varying(100) NOT NULL,
    price numeric(15,2) NOT NULL,
    picture character varying(500) NOT NULL,
    description character varying(200),
    id_articles_types integer NOT NULL,
    seller character varying(100) NOT NULL,
    buyer character varying(100),
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    sell_date timestamp without time zone
);


ALTER TABLE public.articles OWNER TO postgres;

--
-- Name: articles_id_articles_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_id_articles_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.articles_id_articles_seq OWNER TO postgres;

--
-- Name: articles_id_articles_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_id_articles_seq OWNED BY public.articles.id_articles;


--
-- Name: articles_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.articles_types (
    id_articles_types integer NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.articles_types OWNER TO postgres;

--
-- Name: articles_types_id_articles_types_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.articles_types_id_articles_types_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.articles_types_id_articles_types_seq OWNER TO postgres;

--
-- Name: articles_types_id_articles_types_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.articles_types_id_articles_types_seq OWNED BY public.articles_types.id_articles_types;


--
-- Name: bassine; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.bassine (
    id_bassine integer NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    creator character varying(100) NOT NULL
);


ALTER TABLE public.bassine OWNER TO postgres;

--
-- Name: bassine_id_bassine_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.bassine_id_bassine_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.bassine_id_bassine_seq OWNER TO postgres;

--
-- Name: bassine_id_bassine_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.bassine_id_bassine_seq OWNED BY public.bassine.id_bassine;


--
-- Name: caps; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.caps (
    id_caps integer NOT NULL,
    game_code character varying(6) NOT NULL,
    number_of_players smallint NOT NULL,
    finished boolean,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    finish_date timestamp without time zone,
    creator character varying(100) NOT NULL
);


ALTER TABLE public.caps OWNER TO postgres;

--
-- Name: caps_id_caps_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.caps_id_caps_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.caps_id_caps_seq OWNER TO postgres;

--
-- Name: caps_id_caps_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.caps_id_caps_seq OWNED BY public.caps.id_caps;


--
-- Name: clubs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clubs (
    id_clubs integer NOT NULL,
    name character varying(50) NOT NULL,
    picture character varying(500),
    link character varying(100),
    subtitle character varying(150) NOT NULL,
    description character varying(500),
    location character varying(100) NOT NULL,
    practice_time time without time zone,
    practice_day character varying(20),
    id_clubs_types integer NOT NULL
);


ALTER TABLE public.clubs OWNER TO postgres;

--
-- Name: clubs_id_clubs_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clubs_id_clubs_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clubs_id_clubs_seq OWNER TO postgres;

--
-- Name: clubs_id_clubs_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clubs_id_clubs_seq OWNED BY public.clubs.id_clubs;


--
-- Name: clubs_members; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clubs_members (
    email character varying(100) NOT NULL,
    id_clubs integer NOT NULL
);


ALTER TABLE public.clubs_members OWNER TO postgres;

--
-- Name: clubs_respos; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clubs_respos (
    email character varying(100) NOT NULL,
    id_clubs integer NOT NULL
);


ALTER TABLE public.clubs_respos OWNER TO postgres;

--
-- Name: clubs_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.clubs_types (
    id_clubs_types integer NOT NULL,
    name character varying(50) NOT NULL
);


ALTER TABLE public.clubs_types OWNER TO postgres;

--
-- Name: clubs_types_id_clubs_types_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.clubs_types_id_clubs_types_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.clubs_types_id_clubs_types_seq OWNER TO postgres;

--
-- Name: clubs_types_id_clubs_types_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.clubs_types_id_clubs_types_seq OWNED BY public.clubs_types.id_clubs_types;


--
-- Name: request_statistics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.request_statistics (
    id integer NOT NULL,
    user_email character varying(100),
    endpoint character varying(200) NOT NULL,
    method character varying(10) NOT NULL,
    request_received timestamp without time zone NOT NULL,
    response_sent timestamp without time zone NOT NULL,
    status_code integer NOT NULL,
    duration_ms integer NOT NULL
);


ALTER TABLE public.request_statistics OWNER TO postgres;

--
-- Name: endpoint_statistics; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.endpoint_statistics AS
 SELECT endpoint,
    method,
    count(*) AS request_count,
    avg(duration_ms) AS avg_duration_ms,
    min(duration_ms) AS min_duration_ms,
    max(duration_ms) AS max_duration_ms,
    (((sum(
        CASE
            WHEN ((status_code >= 200) AND (status_code < 300)) THEN 1
            ELSE 0
        END))::double precision / (count(*))::double precision) * (100)::double precision) AS success_rate_percent,
    min(request_received) AS first_request,
    max(request_received) AS last_request,
    sum(
        CASE
            WHEN ((status_code >= 200) AND (status_code < 300)) THEN 1
            ELSE 0
        END) AS success_count,
    sum(
        CASE
            WHEN (status_code >= 400) THEN 1
            ELSE 0
        END) AS error_count
   FROM public.request_statistics
  GROUP BY endpoint, method
  ORDER BY endpoint, method;


ALTER VIEW public.endpoint_statistics OWNER TO postgres;

--
-- Name: events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events (
    id_events integer NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(200),
    link character varying(100),
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone,
    location character varying(100) NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    picture character varying(500),
    slots integer,
    price numeric(15,2)
);


ALTER TABLE public.events OWNER TO postgres;

--
-- Name: events_attendents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events_attendents (
    email character varying(100) NOT NULL,
    id_events integer NOT NULL
);


ALTER TABLE public.events_attendents OWNER TO postgres;

--
-- Name: events_clubs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events_clubs (
    id_clubs integer NOT NULL,
    id_events integer NOT NULL
);


ALTER TABLE public.events_clubs OWNER TO postgres;

--
-- Name: events_hosts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events_hosts (
    email character varying(100) NOT NULL,
    id_events integer NOT NULL
);


ALTER TABLE public.events_hosts OWNER TO postgres;

--
-- Name: events_id_events_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.events_id_events_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.events_id_events_seq OWNER TO postgres;

--
-- Name: events_id_events_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.events_id_events_seq OWNED BY public.events.id_events;


--
-- Name: files; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.files (
    id_files integer NOT NULL,
    name character varying(100) NOT NULL,
    path character varying(500) NOT NULL,
    email character varying(100) NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_access_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.files OWNER TO postgres;

--
-- Name: files_id_files_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.files_id_files_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.files_id_files_seq OWNER TO postgres;

--
-- Name: files_id_files_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.files_id_files_seq OWNED BY public.files.id_files;


--
-- Name: games; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.games (
    id_games integer NOT NULL,
    name character varying(50) NOT NULL,
    description character varying(200),
    rules character varying(1000),
    picture character varying(500)
);


ALTER TABLE public.games OWNER TO postgres;

--
-- Name: games_id_games_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.games_id_games_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.games_id_games_seq OWNER TO postgres;

--
-- Name: games_id_games_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.games_id_games_seq OWNED BY public.games.id_games;


--
-- Name: global_statistics; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.global_statistics AS
 SELECT count(*) AS total_request_count,
    avg(duration_ms) AS global_avg_duration_ms,
    min(duration_ms) AS global_min_duration_ms,
    max(duration_ms) AS global_max_duration_ms,
    (((sum(
        CASE
            WHEN ((status_code >= 200) AND (status_code < 300)) THEN 1
            ELSE 0
        END))::double precision / (count(*))::double precision) * (100)::double precision) AS global_success_rate_percent,
    min(request_received) AS first_request,
    max(request_received) AS last_request,
    sum(
        CASE
            WHEN ((status_code >= 200) AND (status_code < 300)) THEN 1
            ELSE 0
        END) AS success_count,
    sum(
        CASE
            WHEN (status_code >= 400) THEN 1
            ELSE 0
        END) AS error_count
   FROM public.request_statistics;


ALTER VIEW public.global_statistics OWNER TO postgres;

--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.goose_db_version OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

ALTER TABLE public.goose_db_version ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.goose_db_version_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: languages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.languages (
    id_languages integer NOT NULL,
    name character varying(50) NOT NULL,
    code character varying(10) NOT NULL
);


ALTER TABLE public.languages OWNER TO postgres;

--
-- Name: languages_id_languages_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.languages_id_languages_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.languages_id_languages_seq OWNER TO postgres;

--
-- Name: languages_id_languages_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.languages_id_languages_seq OWNED BY public.languages.id_languages;


--
-- Name: newf; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.newf (
    email character varying(100) NOT NULL,
    password character varying(200) NOT NULL,
    password_updated_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    verification_code character varying(6) DEFAULT lpad((floor((random() * (1000000)::double precision)))::text, 6, '0'::text),
    verification_code_expiration timestamp without time zone DEFAULT (CURRENT_TIMESTAMP + '00:10:00'::interval),
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    phone_number character varying(20),
    profile_picture character varying(500),
    notification_token character varying(50),
    graduation_year smallint,
    campus character varying(10),
    id_newf integer NOT NULL,
    language smallint DEFAULT 1 NOT NULL,
    CONSTRAINT newf_campus_check CHECK (((campus)::text = ANY ((ARRAY['NANTES'::character varying, 'BREST'::character varying, 'RENNES'::character varying])::text[])))
);


ALTER TABLE public.newf OWNER TO postgres;

--
-- Name: newf_id_newf_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.newf_id_newf_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.newf_id_newf_seq OWNER TO postgres;

--
-- Name: newf_id_newf_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.newf_id_newf_seq OWNED BY public.newf.id_newf;


--
-- Name: newf_roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.newf_roles (
    email character varying(100) NOT NULL,
    id_roles integer NOT NULL
);


ALTER TABLE public.newf_roles OWNER TO postgres;

--
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
    email character varying(100) NOT NULL,
    id_services integer NOT NULL
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- Name: players_bassine; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.players_bassine (
    email character varying(100) NOT NULL,
    id_games integer NOT NULL,
    id_bassine integer NOT NULL,
    score integer
);


ALTER TABLE public.players_bassine OWNER TO postgres;

--
-- Name: players_caps; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.players_caps (
    email character varying(100) NOT NULL,
    id_games integer NOT NULL,
    id_caps integer NOT NULL,
    score integer NOT NULL
);


ALTER TABLE public.players_caps OWNER TO postgres;

--
-- Name: realcampus_friendships; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.realcampus_friendships (
    id_friendship integer NOT NULL,
    user_id character varying(100) NOT NULL,
    friend_id character varying(100) NOT NULL,
    status character varying(20) DEFAULT 'PENDING'::character varying NOT NULL,
    request_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    acceptance_date timestamp without time zone,
    CONSTRAINT realcampus_friendships_check CHECK (((user_id)::text <> (friend_id)::text)),
    CONSTRAINT realcampus_friendships_status_check CHECK (((status)::text = ANY ((ARRAY['PENDING'::character varying, 'ACCEPTED'::character varying, 'REJECTED'::character varying])::text[])))
);


ALTER TABLE public.realcampus_friendships OWNER TO postgres;

--
-- Name: realcampus_friendships_id_friendship_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.realcampus_friendships_id_friendship_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.realcampus_friendships_id_friendship_seq OWNER TO postgres;

--
-- Name: realcampus_friendships_id_friendship_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.realcampus_friendships_id_friendship_seq OWNED BY public.realcampus_friendships.id_friendship;


--
-- Name: realcampus_post_reactions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.realcampus_post_reactions (
    id_reaction integer NOT NULL,
    id_post integer NOT NULL,
    id_file integer NOT NULL,
    author_email character varying(100) NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.realcampus_post_reactions OWNER TO postgres;

--
-- Name: realcampus_post_reactions_id_reaction_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.realcampus_post_reactions_id_reaction_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.realcampus_post_reactions_id_reaction_seq OWNER TO postgres;

--
-- Name: realcampus_post_reactions_id_reaction_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.realcampus_post_reactions_id_reaction_seq OWNED BY public.realcampus_post_reactions.id_reaction;


--
-- Name: realcampus_posts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.realcampus_posts (
    id_post integer NOT NULL,
    id_file_1 integer NOT NULL,
    id_file_2 integer NOT NULL,
    author_email character varying(100) NOT NULL,
    location character varying(200),
    privacy public.realcampus_privacy_level DEFAULT 'PRIVATE'::public.realcampus_privacy_level NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.realcampus_posts OWNER TO postgres;

--
-- Name: realcampus_posts_id_post_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.realcampus_posts_id_post_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.realcampus_posts_id_post_seq OWNER TO postgres;

--
-- Name: realcampus_posts_id_post_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.realcampus_posts_id_post_seq OWNED BY public.realcampus_posts.id_post;


--
-- Name: request_statistics_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.request_statistics_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.request_statistics_id_seq OWNER TO postgres;

--
-- Name: request_statistics_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.request_statistics_id_seq OWNED BY public.request_statistics.id;


--
-- Name: reservations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.reservations (
    email character varying(100) NOT NULL,
    id_rooms integer NOT NULL,
    start_date timestamp without time zone NOT NULL,
    end_date timestamp without time zone NOT NULL
);


ALTER TABLE public.reservations OWNER TO postgres;

--
-- Name: restaurant; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.restaurant (
    id_restaurant integer NOT NULL,
    articles character varying(1000) NOT NULL,
    updated_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    language smallint DEFAULT 1 NOT NULL
);


ALTER TABLE public.restaurant OWNER TO postgres;

--
-- Name: restaurant_id_restaurant_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.restaurant_id_restaurant_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.restaurant_id_restaurant_seq OWNER TO postgres;

--
-- Name: restaurant_id_restaurant_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.restaurant_id_restaurant_seq OWNED BY public.restaurant.id_restaurant;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.roles (
    id_roles integer NOT NULL,
    name character varying(50) NOT NULL,
    description character varying(100)
);


ALTER TABLE public.roles OWNER TO postgres;

--
-- Name: roles_id_roles_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.roles_id_roles_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.roles_id_roles_seq OWNER TO postgres;

--
-- Name: roles_id_roles_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.roles_id_roles_seq OWNED BY public.roles.id_roles;


--
-- Name: rooms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rooms (
    id_rooms integer NOT NULL,
    name character varying(100) NOT NULL,
    description character varying(200),
    picture character varying(500),
    location character varying(100) NOT NULL
);


ALTER TABLE public.rooms OWNER TO postgres;

--
-- Name: rooms_id_rooms_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.rooms_id_rooms_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.rooms_id_rooms_seq OWNER TO postgres;

--
-- Name: rooms_id_rooms_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.rooms_id_rooms_seq OWNED BY public.rooms.id_rooms;


--
-- Name: services; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.services (
    id_services integer NOT NULL,
    name character varying(50) NOT NULL,
    type character varying(50),
    id_category character varying(50)
);


ALTER TABLE public.services OWNER TO postgres;

--
-- Name: services_id_services_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.services_id_services_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.services_id_services_seq OWNER TO postgres;

--
-- Name: services_id_services_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.services_id_services_seq OWNED BY public.services.id_services;


--
-- Name: support; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.support (
    id_support integer NOT NULL,
    email character varying(100) NOT NULL,
    subject character varying(200) NOT NULL,
    message character varying(5000) NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status character varying(20) DEFAULT 'pending'::character varying NOT NULL,
    response character varying(5000)
);


ALTER TABLE public.support OWNER TO postgres;

--
-- Name: support_id_support_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.support_id_support_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.support_id_support_seq OWNER TO postgres;

--
-- Name: support_id_support_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.support_id_support_seq OWNED BY public.support.id_support;


--
-- Name: support_media; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.support_media (
    id_support_media integer NOT NULL,
    id_support integer NOT NULL,
    id_files integer NOT NULL
);


ALTER TABLE public.support_media OWNER TO postgres;

--
-- Name: support_media_id_support_media_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.support_media_id_support_media_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.support_media_id_support_media_seq OWNER TO postgres;

--
-- Name: support_media_id_support_media_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.support_media_id_support_media_seq OWNED BY public.support_media.id_support_media;


--
-- Name: top_users_statistics; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.top_users_statistics AS
 SELECT user_email AS email,
    count(*) AS request_count,
    avg(duration_ms) AS avg_duration_ms,
    min(duration_ms) AS min_duration_ms,
    max(duration_ms) AS max_duration_ms,
    ((100.0 * (count(
        CASE
            WHEN (status_code < 400) THEN 1
            ELSE NULL::integer
        END))::numeric) / (count(*))::numeric) AS success_rate_percent,
    min(request_received) AS first_request,
    max(request_received) AS last_request,
    count(
        CASE
            WHEN (status_code < 400) THEN 1
            ELSE NULL::integer
        END) AS success_count,
    count(
        CASE
            WHEN (status_code >= 400) THEN 1
            ELSE NULL::integer
        END) AS error_count
   FROM public.request_statistics
  WHERE (user_email IS NOT NULL)
  GROUP BY user_email
  ORDER BY (count(*)) DESC
 LIMIT 10;


ALTER VIEW public.top_users_statistics OWNER TO postgres;

--
-- Name: traq; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.traq (
    id_traq integer NOT NULL,
    name character varying(50) NOT NULL,
    disabled boolean DEFAULT false NOT NULL,
    limited boolean DEFAULT false NOT NULL,
    alcohol numeric(4,2) NOT NULL,
    out_of_stock boolean DEFAULT false NOT NULL,
    creation_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    picture character varying(500) NOT NULL,
    description character varying(200) NOT NULL,
    price numeric(10,2) NOT NULL,
    price_half numeric(10,2) NOT NULL,
    id_traq_types integer NOT NULL
);


ALTER TABLE public.traq OWNER TO postgres;

--
-- Name: traq_id_traq_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.traq_id_traq_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.traq_id_traq_seq OWNER TO postgres;

--
-- Name: traq_id_traq_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.traq_id_traq_seq OWNED BY public.traq.id_traq;


--
-- Name: traq_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.traq_types (
    id_traq_types integer NOT NULL,
    name character varying(100) NOT NULL
);


ALTER TABLE public.traq_types OWNER TO postgres;

--
-- Name: traq_types_id_traq_types_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.traq_types_id_traq_types_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.traq_types_id_traq_types_seq OWNER TO postgres;

--
-- Name: traq_types_id_traq_types_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.traq_types_id_traq_types_seq OWNED BY public.traq_types.id_traq_types;


--
-- Name: washing_machines; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.washing_machines (
    id_washing_machines integer NOT NULL,
    type character varying(50) NOT NULL,
    broken boolean NOT NULL
);


ALTER TABLE public.washing_machines OWNER TO postgres;

--
-- Name: washing_machines_id_washing_machines_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.washing_machines_id_washing_machines_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.washing_machines_id_washing_machines_seq OWNER TO postgres;

--
-- Name: washing_machines_id_washing_machines_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.washing_machines_id_washing_machines_seq OWNED BY public.washing_machines.id_washing_machines;


--
-- Name: washing_machines_reports; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.washing_machines_reports (
    email character varying(100) NOT NULL,
    id_washing_machines integer NOT NULL,
    report_date timestamp without time zone NOT NULL
);


ALTER TABLE public.washing_machines_reports OWNER TO postgres;

--
-- Name: articles id_articles; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles ALTER COLUMN id_articles SET DEFAULT nextval('public.articles_id_articles_seq'::regclass);


--
-- Name: articles_types id_articles_types; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles_types ALTER COLUMN id_articles_types SET DEFAULT nextval('public.articles_types_id_articles_types_seq'::regclass);


--
-- Name: bassine id_bassine; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bassine ALTER COLUMN id_bassine SET DEFAULT nextval('public.bassine_id_bassine_seq'::regclass);


--
-- Name: caps id_caps; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.caps ALTER COLUMN id_caps SET DEFAULT nextval('public.caps_id_caps_seq'::regclass);


--
-- Name: clubs id_clubs; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs ALTER COLUMN id_clubs SET DEFAULT nextval('public.clubs_id_clubs_seq'::regclass);


--
-- Name: clubs_types id_clubs_types; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_types ALTER COLUMN id_clubs_types SET DEFAULT nextval('public.clubs_types_id_clubs_types_seq'::regclass);


--
-- Name: events id_events; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events ALTER COLUMN id_events SET DEFAULT nextval('public.events_id_events_seq'::regclass);


--
-- Name: files id_files; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files ALTER COLUMN id_files SET DEFAULT nextval('public.files_id_files_seq'::regclass);


--
-- Name: games id_games; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.games ALTER COLUMN id_games SET DEFAULT nextval('public.games_id_games_seq'::regclass);


--
-- Name: languages id_languages; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.languages ALTER COLUMN id_languages SET DEFAULT nextval('public.languages_id_languages_seq'::regclass);


--
-- Name: newf id_newf; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf ALTER COLUMN id_newf SET DEFAULT nextval('public.newf_id_newf_seq'::regclass);


--
-- Name: realcampus_friendships id_friendship; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_friendships ALTER COLUMN id_friendship SET DEFAULT nextval('public.realcampus_friendships_id_friendship_seq'::regclass);


--
-- Name: realcampus_post_reactions id_reaction; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions ALTER COLUMN id_reaction SET DEFAULT nextval('public.realcampus_post_reactions_id_reaction_seq'::regclass);


--
-- Name: realcampus_posts id_post; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_posts ALTER COLUMN id_post SET DEFAULT nextval('public.realcampus_posts_id_post_seq'::regclass);


--
-- Name: request_statistics id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.request_statistics ALTER COLUMN id SET DEFAULT nextval('public.request_statistics_id_seq'::regclass);


--
-- Name: restaurant id_restaurant; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restaurant ALTER COLUMN id_restaurant SET DEFAULT nextval('public.restaurant_id_restaurant_seq'::regclass);


--
-- Name: roles id_roles; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles ALTER COLUMN id_roles SET DEFAULT nextval('public.roles_id_roles_seq'::regclass);


--
-- Name: rooms id_rooms; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms ALTER COLUMN id_rooms SET DEFAULT nextval('public.rooms_id_rooms_seq'::regclass);


--
-- Name: services id_services; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services ALTER COLUMN id_services SET DEFAULT nextval('public.services_id_services_seq'::regclass);


--
-- Name: support id_support; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support ALTER COLUMN id_support SET DEFAULT nextval('public.support_id_support_seq'::regclass);


--
-- Name: support_media id_support_media; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support_media ALTER COLUMN id_support_media SET DEFAULT nextval('public.support_media_id_support_media_seq'::regclass);


--
-- Name: traq id_traq; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq ALTER COLUMN id_traq SET DEFAULT nextval('public.traq_id_traq_seq'::regclass);


--
-- Name: traq_types id_traq_types; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq_types ALTER COLUMN id_traq_types SET DEFAULT nextval('public.traq_types_id_traq_types_seq'::regclass);


--
-- Name: washing_machines id_washing_machines; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.washing_machines ALTER COLUMN id_washing_machines SET DEFAULT nextval('public.washing_machines_id_washing_machines_seq'::regclass);


--
-- Data for Name: articles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.articles (id_articles, name, price, picture, description, id_articles_types, seller, buyer, creation_date, sell_date) FROM stdin;
\.


--
-- Data for Name: articles_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.articles_types (id_articles_types, name) FROM stdin;
\.


--
-- Data for Name: bassine; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.bassine (id_bassine, creation_date, creator) FROM stdin;
\.


--
-- Data for Name: caps; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.caps (id_caps, game_code, number_of_players, finished, creation_date, finish_date, creator) FROM stdin;
\.


--
-- Data for Name: clubs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clubs (id_clubs, name, picture, link, subtitle, description, location, practice_time, practice_day, id_clubs_types) FROM stdin;
\.


--
-- Data for Name: clubs_members; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clubs_members (email, id_clubs) FROM stdin;
\.


--
-- Data for Name: clubs_respos; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clubs_respos (email, id_clubs) FROM stdin;
\.


--
-- Data for Name: clubs_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.clubs_types (id_clubs_types, name) FROM stdin;
\.


--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events (id_events, name, description, link, start_date, end_date, location, creation_date, picture, slots, price) FROM stdin;
\.


--
-- Data for Name: events_attendents; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events_attendents (email, id_events) FROM stdin;
\.


--
-- Data for Name: events_clubs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events_clubs (id_clubs, id_events) FROM stdin;
\.


--
-- Data for Name: events_hosts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.events_hosts (email, id_events) FROM stdin;
\.


--
-- Data for Name: files; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.files (id_files, name, path, email, creation_date, last_access_date) FROM stdin;
40	b8aa4b8d-ce3c-4528-9936-71e3aee47ced.jpeg	/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	nathaniel.guitton@imt-atlantique.net	2025-05-20 09:29:57.732301	2025-05-20 09:48:07.034718
44	31714441-fe62-48c8-bf57-36dd24469c1f.jpeg	/data/31714441-fe62-48c8-bf57-36dd24469c1f_73c0a62bf1d69f33.jpeg	marina.carbone@imt-atlantique.net	2025-05-21 08:00:48.043164	2025-05-21 08:00:49.396175
17	0796447A-0228-434B-B77B-953A0CD1271C.jpg	/data/0796447A-0228-434B-B77B-953A0CD1271C.jpg	enzo.morvan@imt-atlantique.net	2025-04-08 09:31:02.125327	2025-04-12 13:54:20.464042
15	1EB9B061-EA4B-45AD-BBB8-02C8EE4B1559.jpg	/data/1EB9B061-EA4B-45AD-BBB8-02C8EE4B1559.jpg	maxime.bodin@imt-atlantique.net	2025-04-08 09:15:51.969382	2025-04-08 09:18:52.102616
22	E626A86F-4883-4C41-BBF2-FDDE3DAD7DFB.jpg	/data/E626A86F-4883-4C41-BBF2-FDDE3DAD7DFB.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:52:45.388442	2025-04-10 09:53:39.376158
16	1598C8BB-A979-4124-84C0-6C33E4F9F14A.jpg	/data/1598C8BB-A979-4124-84C0-6C33E4F9F14A.jpg	maxime.bodin@imt-atlantique.net	2025-04-08 09:15:52.335195	2025-04-09 20:24:15.738315
4	be4a8075-d4b6-47f6-bf6a-f5b7b22981f4.jpeg	/data/be4a8075-d4b6-47f6-bf6a-f5b7b22981f4.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-13 15:30:48.608299	2025-03-13 15:30:48.608299
5	99c5515b-0e13-4632-8c9d-37feb5255e42.jpeg	/data/99c5515b-0e13-4632-8c9d-37feb5255e42.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-13 15:45:48.488867	2025-03-13 15:56:29.058193
23	9ED041DB-E216-4EC3-AA78-96684B61EFE8.jpg	/data/9ED041DB-E216-4EC3-AA78-96684B61EFE8.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:52:45.787181	2025-04-10 09:53:39.504434
8	44AA0DC7-9545-441F-AB69-4B19F4E8F400.jpg	/data/44AA0DC7-9545-441F-AB69-4B19F4E8F400.jpg	enzo.morvan@imt-atlantique.net	2025-03-13 17:37:43.058615	2025-04-08 09:29:16.308907
6	b024cc57-df40-42fd-99e6-cf636b39b035.jpeg	/data/b024cc57-df40-42fd-99e6-cf636b39b035.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-13 15:59:05.95139	2025-03-13 16:02:46.391417
41	2F773180-2F5B-41EF-8DA2-1109D921BBE0.jpg	/data/2F773180-2F5B-41EF-8DA2-1109D921BBE0_4d1b0b7d141e7f9e.jpg	aurelien.pautet@imt-atlantique.net	2025-05-20 11:02:24.53706	2025-05-21 16:38:06.042802
14	FB362749-9650-4C11-8ECB-BD20AF5395B1.jpg	/data/FB362749-9650-4C11-8ECB-BD20AF5395B1.jpg	maxime.bodin@imt-atlantique.net	2025-04-08 09:09:51.565177	2025-04-09 20:40:09.435369
13	E4BD9E0D-A741-4EDB-9947-33F5C32D969F.jpg	/data/E4BD9E0D-A741-4EDB-9947-33F5C32D969F.jpg	maxime.bodin@imt-atlantique.net	2025-04-08 09:09:51.139295	2025-04-09 20:40:09.435362
18	B2B36083-4BB6-4494-B434-65B73E58FF6F.jpg	/data/B2B36083-4BB6-4494-B434-65B73E58FF6F.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:01:23.936334	2025-04-10 08:01:23.936334
19	FA3229BD-0807-4211-9CD4-2503D27A9DEE.jpg	/data/FA3229BD-0807-4211-9CD4-2503D27A9DEE.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:01:24.31601	2025-04-10 08:01:24.31601
30	2F92CB02-2D15-49ED-BAE0-9911EF10E81A.jpg	/data/2F92CB02-2D15-49ED-BAE0-9911EF10E81A.jpg	enzo.morvan@imt-atlantique.net	2025-04-12 14:23:26.801393	2025-04-15 20:21:46.434828
20	8237492B-3A16-42DA-BB01-AD39758E8A8C.jpg	/data/8237492B-3A16-42DA-BB01-AD39758E8A8C.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:43:59.868284	2025-04-10 08:43:59.868284
21	E9649E41-F440-4DD1-9284-5AAE1E8C2147.jpg	/data/E9649E41-F440-4DD1-9284-5AAE1E8C2147.jpg	maxime.bodin@imt-atlantique.net	2025-04-10 08:44:00.130186	2025-04-10 08:44:00.130186
7	c63f946c-4cf5-4826-926b-4a3614a9ecb8.jpeg	/data/c63f946c-4cf5-4826-926b-4a3614a9ecb8.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-13 16:14:14.330518	2025-03-14 18:02:03.780933
9	40e5edb7-5c02-4ac2-aacd-ebbcc9a7cc56.jpeg	/data/40e5edb7-5c02-4ac2-aacd-ebbcc9a7cc56.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-14 18:02:41.630997	2025-03-14 18:02:42.29493
10	a5a5ef71-2bc9-41b6-aaef-5804ffcce247.jpeg	/data/a5a5ef71-2bc9-41b6-aaef-5804ffcce247.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-14 18:04:55.035601	2025-03-14 18:04:55.947066
3	chavanel.png	/data/chavanel.png	yohann.chavanel@imt-atlantique.net	2025-03-13 13:28:10.46054	2025-04-01 09:02:02.052157
34	images.png	/data/images_d9b7046905ef1a64.png	test@imt-atlantique.net	2025-04-20 20:23:18.925665	2025-04-20 20:23:18.925665
11	572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	yohann.chavanel@imt-atlantique.net	2025-03-14 18:05:19.605384	2025-05-21 10:49:58.299177
25	6B6026C5-DD34-4753-A96E-F92BAAFC7D87.jpg	/data/6B6026C5-DD34-4753-A96E-F92BAAFC7D87.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:52:02.509828	2025-04-11 08:53:59.713766
24	D6589BB7-4969-4715-95E6-3AB83ABD61D4.jpg	/data/D6589BB7-4969-4715-95E6-3AB83ABD61D4.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:52:02.263508	2025-04-11 08:53:59.713755
27	2063872E-8E6D-426B-ACE1-436274264A1A.jpg	/data/2063872E-8E6D-426B-ACE1-436274264A1A_c28f7d559a4f5097.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:56:16.890622	2025-04-11 08:56:16.890622
29	3C466324-35C4-4720-8C16-53F643A10738.jpg	/data/3C466324-35C4-4720-8C16-53F643A10738_90af0dde225a405c.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:56:17.130798	2025-04-11 08:56:17.130798
28	3C466324-35C4-4720-8C16-53F643A10738.jpg	/data/3C466324-35C4-4720-8C16-53F643A10738.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:56:17.089347	2025-04-11 08:56:18.024875
26	2063872E-8E6D-426B-ACE1-436274264A1A.jpg	/data/2063872E-8E6D-426B-ACE1-436274264A1A.jpg	maxime.bodin@imt-atlantique.net	2025-04-11 08:56:16.830279	2025-04-11 08:56:18.024888
12	29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	lucie.delestre@imt-atlantique.net	2025-03-30 12:39:30.822238	2025-05-21 20:28:18.557522
45	8e13f62f-9789-481a-b93c-d4c2cef0ccec.jpeg	/data/8e13f62f-9789-481a-b93c-d4c2cef0ccec_db66b49d7de717b3.jpeg	pacome.cailleteau@imt-atlantique.net	2025-05-21 09:19:34.393129	2025-05-21 10:51:20.939712
33	New Project.png	/data/New_Project_943f029163e65ef3.png	test@imt-atlantique.net	2025-04-20 20:05:15.33312	2025-05-21 22:26:10.8526
31	8DA6336C-6D0D-472C-B5AB-25BE4078CBBF.jpg	/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	enzo.morvan@imt-atlantique.net	2025-04-15 20:22:34.09011	2025-05-21 11:34:29.522296
32	New Project.png	data/New_Project_8626bf35661a2046.png	yohann.chavanel@imt-atlantique.net	2025-04-20 20:03:46.065495	2025-04-20 20:03:46.065495
35	IMG_4222.JPG	/Users/yohann/Downloads/GitHub/Transat_2.0_Backend/data/IMG_4222_3d1dd73fcbf97f8b.JPG	yohann.chavanel@imt-atlantique.net	2025-05-20 08:52:38.862422	2025-05-20 08:52:38.862422
36	6a8937f3-fd98-486a-be2e-d597e9aa3576.jpeg	/data/6a8937f3-fd98-486a-be2e-d597e9aa3576_ea57c3fd8b37c04e.jpeg	test@imt-atlantique.net	2025-05-20 08:55:55.822822	2025-05-20 08:55:55.822822
37	297f8dfd-b5d0-4e66-91db-a6f33e8c33ee.jpeg	/data/297f8dfd-b5d0-4e66-91db-a6f33e8c33ee_18d87ab6543d751e.jpeg	test@imt-atlantique.net	2025-05-20 08:56:19.60368	2025-05-20 08:56:19.60368
38	b7f35d3b-3f52-479f-a613-1e16df857a29.jpeg	/data/b7f35d3b-3f52-479f-a613-1e16df857a29_69050dedd231faf0.jpeg	test@imt-atlantique.net	2025-05-20 09:05:06.142241	2025-05-20 09:05:06.142241
42	EDEB72E5-0F88-4926-8E9A-F2336C16ED38.jpg	/data/EDEB72E5-0F88-4926-8E9A-F2336C16ED38_53fa02e12cbd1f70.jpg	nathan.marie2@imt-atlantique.net	2025-05-20 17:44:09.007331	2025-05-21 06:49:59.561709
43	A7854B2E-474F-435D-B81B-BD24F41E477C.jpg	/data/A7854B2E-474F-435D-B81B-BD24F41E477C_f3b3b3e98f1fa56a.jpg	aurelien.moignet@imt-atlantique.net	2025-05-21 07:47:44.952546	2025-05-21 07:47:44.952546
39	467bdeed-c69f-4074-bb6e-765100a569f5.jpeg	/data/467bdeed-c69f-4074-bb6e-765100a569f5_85fe4e125190b7b4.jpeg	test@imt-atlantique.net	2025-05-20 09:17:37.29009	2025-05-21 16:45:32.076115
\.


--
-- Data for Name: games; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.games (id_games, name, description, rules, picture) FROM stdin;
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
1	0	t	2025-04-12 13:07:29.455033
2	2	t	2025-04-12 13:07:43.631487
3	1	t	2025-04-12 13:45:52.319726
4	3	t	2025-04-14 20:33:09.427138
5	4	t	2025-04-20 12:57:52.622221
6	5	t	2025-04-20 13:09:12.220478
7	6	t	2025-04-20 13:19:58.678629
8	7	t	2025-04-20 13:34:58.060986
10	8	t	2025-05-16 15:00:54.588412
\.


--
-- Data for Name: languages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.languages (id_languages, name, code) FROM stdin;
1	French	fr
2	English	en
3	Spanish	es
4	German	de
5	Italian	it
6	Portuguese	pt
7	Russian	ru
8	Chinese	zh
9	Korean	ko
\.


--
-- Data for Name: newf; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.newf (email, password, password_updated_date, verification_code, verification_code_expiration, creation_date, first_name, last_name, phone_number, profile_picture, notification_token, graduation_year, campus, id_newf, language) FROM stdin;
aurelien.pautet@imt-atlantique.net	$2a$10$8m9DAT2LtR7T.kfxan/YXOwmpkv77zZUXnop/6sDybqMVs.Xil3Iy	2025-05-20 09:00:54.734195	\N	\N	2025-05-20 09:00:54.734195	Aurelien 	Pautet	06 84 88 94 79	https://transat.destimt.fr/api/data/2F773180-2F5B-41EF-8DA2-1109D921BBE0_4d1b0b7d141e7f9e.jpg	ExponentPushToken[WD69TsLH76AOlCBrNkPVUm]	2027	\N	25	1
nathan.marie2@imt-atlantique.net	$2a$10$U0LrtMHZEFCOS7vMcUtXqOqaPDF5zbTPjg0nV0yuu6AL6iMGJfKi2	2025-05-19 17:24:03.150205	\N	\N	2025-05-19 17:24:03.150205	Nathan	MARIE2	\N	https://transat.destimt.fr/api/data/EDEB72E5-0F88-4926-8E9A-F2336C16ED38_53fa02e12cbd1f70.jpg	ExponentPushToken[sq7s90N_b9Ve6UbUPkBcJE]	\N	\N	22	6
enzo.morvan@imt-atlantique.net	$2a$10$h2cpPyZYnFThQ.Hg9TcXY.aGkWw6pa4sbJnIT7YvDwmoQGAUubH9W	2025-03-13 17:34:28.649816	278674	2025-02-16 10:57:47.890839	2025-02-16 10:47:47.890839	Enzo	MORVAN	0767847910	https://transat.destimt.fr/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	ExponentPushToken[WzPS3XDbFad9ud51erofkO]	2024	\N	2	1
matheo.vallee@imt-atlantique.net	$2a$10$nYZI9gAvcaE6EVR8yYY28eHIQZAYmns4gLK9RcTY.tPD74/NsixSq	2025-05-21 07:32:01.250094	\N	\N	2025-05-21 07:32:01.250094	Matheo	VALLEE	\N	\N	ExponentPushToken[okQKMyHnjw6ROuhOh7XPvW]	\N	\N	30	1
yohann.chavanel@imt-atlantique.net	$2a$10$FHP4uirIUcxb4kxoV4lNJudaFZH0K3H9DhJBGUT2CZqt5tQKGpPOK	2025-02-27 21:16:19.853675	467685	2025-02-27 21:26:06.557033	2025-02-16 01:00:39.965529	Yohann	CHAVANEL	0678139913	https://transat.destimt.fr/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	ExponentPushToken[1R_p_-JgSDC3uH8nM9T_PJ]	2027	NANTES	1	1
maxime.bodin@imt-atlantique.net	$2a$10$2P1/EgMIT2Fukpg4zTV.Be6NAM/kkB/CObANLLdAFjkeph7qQk8Y6	2025-04-07 12:29:58.114781	337459	2025-04-07 12:39:58.114781	2025-04-07 12:29:58.114781	Maxime	BODIN	\N	\N	ExponentPushToken[XTJt7PDZdErm90GthwCTVS]	\N	\N	14	1
aurelien.moignet@imt-atlantique.net	$2a$10$Ps1oWP0I2IEmeQz3nqdKlu0DK5d6aI3LBD0bfNVbf/YhWJ47eWbpa	2025-05-19 16:59:08.252438	\N	\N	2025-05-19 16:59:08.252438	aurelien	moignet	07 82 23 95 15	https://transat.destimt.fr/api/data/A7854B2E-474F-435D-B81B-BD24F41E477C_f3b3b3e98f1fa56a.jpg	ExponentPushToken[GhXAJwKnnvuGpL7QFXkCPJ]	2027	\N	21	2
marina.carbone@imt-atlantique.net	$2a$10$BIhqyj5QsQ64H7i0YQyu3ugwgrSGMLW1C78QXCd8aLb8nuftOA57.	2025-05-20 07:45:50.991456	\N	\N	2025-05-20 07:45:50.991456	Marina	Carbone	\N	https://transat.destimt.fr/api/data/31714441-fe62-48c8-bf57-36dd24469c1f_73c0a62bf1d69f33.jpeg	ExponentPushToken[0nIjjaJkI0W5vpEIP7gvQq]	2027	\N	24	1
lucie.delestre@imt-atlantique.net	$2a$10$Ni5bjApNQCzEMIwtsPVxKOkgYVYV7PSDM3.X3X/bW7f6nLcES8Gji	2025-03-04 08:39:24.903816	826689	2025-03-04 08:49:24.903816	2025-03-04 08:39:24.903816	Lucie	DELESTRE	\N	https://transat.destimt.fr/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	ExponentPushToken[S92U3LINeBCfl7bV_Z4DiM]	2027	\N	4	1
admin@imt-atlantique.net	$2a$10$wp.xlUCiqux42ZcjOmKbRORugmsgj7biEqVW9iJ24o9wb67TeTOlq	2025-04-20 13:22:30.974069	349104	2025-04-20 13:32:30.974069	2025-04-20 13:22:30.974069	Admin	User	\N	\N	\N	\N	\N	16	1
testtest@imt-atlantique.net	$2a$10$sGWry0s16Mj47EGwVSbTrejdYvg5iHLDt3QGKy8gObu9URT033XNe	2025-05-20 19:58:34.366796	702053	2025-05-20 20:08:34.366796	2025-05-20 19:58:34.366796	Testtest		\N	\N	\N	\N	\N	29	2
zephyr.dassouli@imt-atlantique.net	$2a$10$ditAJUsaZnWPVIWtEvbu5ubsfPsB0EFQ.OIAXDAgIPxXtdVoNvPj.	2025-05-16 14:36:37.884117	\N	\N	2025-05-16 14:36:37.884117	Zephyr	DASSOULI	\N	\N	ExponentPushToken[n3DTFePI3NHl20J28IH6rh]	\N	\N	19	2
test@imt-atlantique.net	$2a$10$LwCdvH3TrUd4HWSkqA.eC.3Ufmd0X.xUzhQYRmH9nr9xLlEasoP12	2025-05-03 14:54:47.643445	724412	2025-05-21 22:16:03.765507	2025-02-20 14:19:17.495832	Test	Account	066666666666	https://transat.destimt.fr/api/data/467bdeed-c69f-4074-bb6e-765100a569f5_85fe4e125190b7b4.jpeg	ExponentPushToken[Tk-9epDTEDy94pBnjZTpy3]	\N	\N	3	1
matis.byar@imt-atlantique.net	$2a$10$Vub3/z7bPPiNsR30rGyJE.ZU2rpgeuvSBI2YdbAMDKzV8chj0UJeO	2025-05-16 22:36:43.066785	\N	\N	2025-05-16 22:36:43.066785	Matis	BYAR	\N	\N	ExponentPushToken[wPcVntN0Ai2IqUO4rpywLQ]	\N	\N	20	2
ninon.basle-blin@imt-atlantique.net	$2a$10$qsGJph1xiZ3FHqBnL7zXf.wb0sa.PADvUtnQivk1tFCfQoniDkaX.	2025-05-20 11:01:15.79372	\N	\N	2025-05-20 11:01:15.79372	Ninon	BASLE-BLIN	\N	\N	ExponentPushToken[AMYgzmNcaZwaB266_RKfni]	\N	\N	26	1
nathaniel.guitton@imt-atlantique.net	$2a$10$nmczYSax3nKVXAQf68cnRemvCdVhskM832foXU0RzNt/86yyV./Km	2025-05-20 07:07:04.412714	\N	\N	2025-05-20 07:07:04.412714	Nathaniel	GUITTON	\N	https://transat.destimt.fr/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	ExponentPushToken[FqVLamJEAT7Qi-yFbNBhom]	\N	\N	23	4
pacome.cailleteau@imt-atlantique.net	$2a$10$bNCuG7XE7WrhME9iwlQ/MO/LhAo4QbTsf67S/lxk6GVg0SKmXHrWW	2025-05-21 09:15:54.885511	\N	\N	2025-05-21 09:15:54.885511	Pacôme	CAILLETEAU	0768050633	https://transat.destimt.fr/api/data/8e13f62f-9789-481a-b93c-d4c2cef0ccec_db66b49d7de717b3.jpeg	ExponentPushToken[ZsPivKG8Xo5avTV1-qUBzq]	2027	\N	31	1
\.


--
-- Data for Name: newf_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.newf_roles (email, id_roles) FROM stdin;
yohann.chavanel@imt-atlantique.net	2
enzo.morvan@imt-atlantique.net	2
test@imt-atlantique.net	2
lucie.delestre@imt-atlantique.net	2
maxime.bodin@imt-atlantique.net	2
admin@imt-atlantique.net	1
yohann.chavanel@imt-atlantique.net	1
zephyr.dassouli@imt-atlantique.net	2
matis.byar@imt-atlantique.net	2
aurelien.moignet@imt-atlantique.net	2
nathan.marie2@imt-atlantique.net	2
nathaniel.guitton@imt-atlantique.net	2
marina.carbone@imt-atlantique.net	2
aurelien.pautet@imt-atlantique.net	2
ninon.basle-blin@imt-atlantique.net	2
testtest@imt-atlantique.net	5
matheo.vallee@imt-atlantique.net	2
pacome.cailleteau@imt-atlantique.net	2
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (email, id_services) FROM stdin;
test@imt-atlantique.net	1
enzo.morvan@imt-atlantique.net	2
lucie.delestre@imt-atlantique.net	1
enzo.morvan@imt-atlantique.net	1
yohann.chavanel@imt-atlantique.net	1
yohann.chavanel@imt-atlantique.net	2
aurelien.moignet@imt-atlantique.net	2
aurelien.pautet@imt-atlantique.net	2
aurelien.pautet@imt-atlantique.net	3
maxime.bodin@imt-atlantique.net	1
matis.byar@imt-atlantique.net	1
nathaniel.guitton@imt-atlantique.net	2
maxime.bodin@imt-atlantique.net	2
admin@imt-atlantique.net	2
test@imt-atlantique.net	2
zephyr.dassouli@imt-atlantique.net	2
lucie.delestre@imt-atlantique.net	2
nathan.marie2@imt-atlantique.net	2
matis.byar@imt-atlantique.net	2
marina.carbone@imt-atlantique.net	2
nathaniel.guitton@imt-atlantique.net	3
enzo.morvan@imt-atlantique.net	3
aurelien.moignet@imt-atlantique.net	3
maxime.bodin@imt-atlantique.net	3
admin@imt-atlantique.net	3
test@imt-atlantique.net	3
zephyr.dassouli@imt-atlantique.net	3
lucie.delestre@imt-atlantique.net	3
nathan.marie2@imt-atlantique.net	3
matis.byar@imt-atlantique.net	3
marina.carbone@imt-atlantique.net	3
yohann.chavanel@imt-atlantique.net	3
ninon.basle-blin@imt-atlantique.net	2
ninon.basle-blin@imt-atlantique.net	3
testtest@imt-atlantique.net	2
testtest@imt-atlantique.net	3
matheo.vallee@imt-atlantique.net	1
matheo.vallee@imt-atlantique.net	2
matheo.vallee@imt-atlantique.net	3
pacome.cailleteau@imt-atlantique.net	1
pacome.cailleteau@imt-atlantique.net	3
aurelien.pautet@imt-atlantique.net	1
nathan.marie2@imt-atlantique.net	1
aurelien.moignet@imt-atlantique.net	1
marina.carbone@imt-atlantique.net	1
nathaniel.guitton@imt-atlantique.net	1
admin@imt-atlantique.net	1
testtest@imt-atlantique.net	1
zephyr.dassouli@imt-atlantique.net	1
ninon.basle-blin@imt-atlantique.net	1
pacome.cailleteau@imt-atlantique.net	2
\.


--
-- Data for Name: players_bassine; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.players_bassine (email, id_games, id_bassine, score) FROM stdin;
\.


--
-- Data for Name: players_caps; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.players_caps (email, id_games, id_caps, score) FROM stdin;
\.


--
-- Data for Name: realcampus_friendships; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.realcampus_friendships (id_friendship, user_id, friend_id, status, request_date, acceptance_date) FROM stdin;
\.


--
-- Data for Name: realcampus_post_reactions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.realcampus_post_reactions (id_reaction, id_post, id_file, author_email, creation_date) FROM stdin;
\.


--
-- Data for Name: realcampus_posts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.realcampus_posts (id_post, id_file_1, id_file_2, author_email, location, privacy, creation_date) FROM stdin;
5	14	13	maxime.bodin@imt-atlantique.net	New York	PUBLIC	2025-04-09 17:01:37.782167
6	22	23	maxime.bodin@imt-atlantique.net	TODO	PRIVATE	2025-04-10 08:52:45.988511
7	24	25	maxime.bodin@imt-atlantique.net	TODO	PRIVATE	2025-04-11 08:52:02.737996
8	27	29	maxime.bodin@imt-atlantique.net	TODO	PRIVATE	2025-04-11 08:56:17.333352
9	26	28	maxime.bodin@imt-atlantique.net	TODO	PRIVATE	2025-04-11 08:56:17.334369
\.


--
-- Data for Name: request_statistics; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.request_statistics (id, user_email, endpoint, method, request_received, response_sent, status_code, duration_ms) FROM stdin;
1224	\N	/status	GET	2025-04-20 20:47:26.414397	2025-04-20 20:47:26.414398	200	0
1225	\N	/status	GET	2025-04-20 20:47:26.589057	2025-04-20 20:47:26.589059	200	0
3	\N	/api/restaurant	GET	2025-04-20 15:10:17.332277	2025-04-20 15:10:17.51285	200	180
4	\N	/status	GET	2025-04-20 15:10:44.301629	2025-04-20 15:10:44.301738	200	0
5	\N	/api/auth/login	POST	2025-04-20 15:12:32.185371	2025-04-20 15:12:32.590573	200	405
6	test@imt-atlantique.net	/api/files	GET	2025-04-20 15:12:58.76685	2025-04-20 15:12:59.095477	200	328
7	test@imt-atlantique.net	/api/newf/	DELETE	2025-04-20 15:13:29.465327	2025-04-20 15:13:29.465474	200	0
8	\N	/api/restaurant	GET	2025-04-20 15:15:50.246765	2025-04-20 15:15:50.443729	200	196
9	\N	/api/auth/login	POST	2025-04-20 15:15:55.17964	2025-04-20 15:15:55.580923	401	401
10	\N	/api/auth/login	POST	2025-04-20 15:16:40.084594	2025-04-20 15:16:40.48246	401	397
11	\N	/api/statistics/global	GET	2025-04-20 15:19:59.740823	2025-04-20 15:19:59.740858	200	0
12	\N	/api/auth/login	POST	2025-04-20 15:20:06.210405	2025-04-20 15:20:06.593716	401	383
13	\N	/api/auth/login	POST	2025-04-20 15:20:54.115744	2025-04-20 15:20:54.512878	401	397
14	\N	/status	GET	2025-04-20 15:20:59.632016	2025-04-20 15:20:59.632031	200	0
15	\N	/status	GET	2025-04-20 15:22:07.779517	2025-04-20 15:22:07.779538	200	0
16	\N	/status	GET	2025-04-20 15:22:07.80114	2025-04-20 15:22:07.801142	200	0
17	\N	/status	GET	2025-04-20 15:22:07.794255	2025-04-20 15:22:07.794256	200	0
18	\N	/status	GET	2025-04-20 15:22:07.807767	2025-04-20 15:22:07.807768	200	0
19	\N	/status	GET	2025-04-20 15:22:07.786937	2025-04-20 15:22:07.786938	200	0
20	\N	/api/restaurant	GET	2025-04-20 15:22:12.421946	2025-04-20 15:22:12.628325	200	206
21	\N	/api/restaurant	GET	2025-04-20 15:22:12.646612	2025-04-20 15:22:12.646649	200	0
22	\N	/api/restaurant	GET	2025-04-20 15:22:12.637281	2025-04-20 15:22:12.637383	200	0
23	\N	/api/auth/register	POST	2025-04-20 15:22:18.54563	2025-04-20 15:22:20.36792	201	1822
24	\N	/api/auth/login	POST	2025-04-20 15:22:25.245985	2025-04-20 15:22:25.644072	401	398
25	\N	/api/auth/register	POST	2025-04-20 15:22:30.807539	2025-04-20 15:22:32.524201	201	1716
26	\N	/status	GET	2025-04-20 15:23:15.541253	2025-04-20 15:23:15.541273	200	0
27	\N	/status	GET	2025-04-20 15:23:15.548658	2025-04-20 15:23:15.548659	200	0
28	\N	/status	GET	2025-04-20 15:23:15.556661	2025-04-20 15:23:15.556663	200	0
29	\N	/status	GET	2025-04-20 15:23:15.564579	2025-04-20 15:23:15.56458	200	0
30	\N	/status	GET	2025-04-20 15:23:15.57214	2025-04-20 15:23:15.572141	200	0
31	\N	/api/restaurant	GET	2025-04-20 15:23:21.046785	2025-04-20 15:23:21.193921	200	147
32	\N	/api/restaurant	GET	2025-04-20 15:23:21.207717	2025-04-20 15:23:21.207754	200	0
33	\N	/api/restaurant	GET	2025-04-20 15:23:21.201246	2025-04-20 15:23:21.201288	200	0
34	\N	/api/statistics/global	GET	2025-04-20 15:23:30.825098	2025-04-20 15:23:30.988282	200	163
35	\N	/api/statistics/global	GET	2025-04-20 15:23:35.373293	2025-04-20 15:23:35.573287	200	199
36	\N	/api/statistics/endpoints	GET	2025-04-20 15:23:59.664095	2025-04-20 15:23:59.838371	200	174
37	test@imt-atlantique.net	/api/newf/	DELETE	2025-04-20 15:26:44.131819	2025-04-20 15:26:44.132276	200	0
38	\N	/api/auth/login	POST	2025-04-20 15:26:58.848927	2025-04-20 15:26:59.251323	200	402
39	yohann.chavanel@imt-atlantique.net	/api/newf/notification	GET	2025-04-20 15:27:51.873805	2025-04-20 15:27:51.873984	200	0
40	\N	/status	GET	2025-04-20 15:33:42.837926	2025-04-20 15:33:42.838032	200	0
41	\N	/nonexistent	GET	2025-04-20 15:33:48.5099	2025-04-20 15:33:48.509909	200	0
42	\N	/api/statistics/endpoints	GET	2025-04-20 15:33:54.065415	2025-04-20 15:33:54.251963	200	186
43	\N	/api/statistics/global	GET	2025-04-20 15:34:00.293535	2025-04-20 15:34:00.458041	200	164
44	\N	/nonexistent	GET	2025-04-20 15:35:02.015923	2025-04-20 15:35:02.016021	200	0
45	\N	/api/statistics/endpoints	GET	2025-04-20 15:35:12.083268	2025-04-20 15:35:12.344776	200	261
46	\N	/nonexistent-path	GET	2025-04-20 15:35:50.497413	2025-04-20 15:35:50.497534	500	0
47	\N	/status	GET	2025-04-20 15:35:55.531174	2025-04-20 15:35:55.531174	200	0
48	\N	/api/statistics/endpoints	GET	2025-04-20 15:36:01.090582	2025-04-20 15:36:01.268167	200	177
49	\N	/api/statistics/global	GET	2025-04-20 15:36:07.61493	2025-04-20 15:36:07.830113	200	215
50	yohann.chavanel@imt-atlantique.net	/api/newf/notificationbh	GET	2025-04-20 15:36:30.390509	2025-04-20 15:36:30.39078	500	0
51	yohann.chavanel@imt-atlantique.net	/api/newf/notificationbh$$	GET	2025-04-20 15:36:32.748639	2025-04-20 15:36:32.748772	500	0
52	\N	/api/auth/login	POST	2025-04-20 15:36:52.847684	2025-04-20 15:36:53.254651	401	406
53	\N	/.env	GET	2025-04-20 13:43:41.149826	2025-04-20 13:43:41.149853	200	0
54	\N	/status	GET	2025-04-20 14:26:50.034284	2025-04-20 14:26:50.03431	200	0
55	\N	/status	GET	2025-04-20 14:26:50.205908	2025-04-20 14:26:50.20591	200	0
56	\N	/api/statistics/global	GET	2025-04-20 14:26:50.358355	2025-04-20 14:26:50.360058	200	1
57	\N	/api/statistics/global	GET	2025-04-20 14:26:50.537577	2025-04-20 14:26:50.53843	200	0
58	\N	/api/statistics/endpoints	GET	2025-04-20 14:26:50.555373	2025-04-20 14:26:50.556968	200	1
59	\N	/api/statistics/endpoints	GET	2025-04-20 14:26:50.723885	2025-04-20 14:26:50.72604	200	2
60	\N	/status	GET	2025-04-20 14:27:03.938221	2025-04-20 14:27:03.938223	200	0
61	\N	/status	GET	2025-04-20 14:27:04.109436	2025-04-20 14:27:04.109438	200	0
62	\N	/api/statistics/global	GET	2025-04-20 14:27:04.171014	2025-04-20 14:27:04.178486	200	7
63	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:04.354707	2025-04-20 14:27:04.356031	200	1
64	\N	/api/statistics/global	GET	2025-04-20 14:27:04.354771	2025-04-20 14:27:04.37165	200	16
65	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:04.54567	2025-04-20 14:27:04.5475	200	1
66	\N	/status	GET	2025-04-20 14:27:08.298134	2025-04-20 14:27:08.298136	200	0
67	\N	/api/statistics/global	GET	2025-04-20 14:27:08.467494	2025-04-20 14:27:08.468703	200	1
68	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:08.651716	2025-04-20 14:27:08.65292	200	1
69	\N	/status	GET	2025-04-20 14:27:09.790595	2025-04-20 14:27:09.790596	200	0
70	\N	/api/statistics/global	GET	2025-04-20 14:27:09.959361	2025-04-20 14:27:09.960258	200	0
71	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:10.143692	2025-04-20 14:27:10.146633	200	2
72	\N	/status	GET	2025-04-20 14:27:11.08677	2025-04-20 14:27:11.086771	200	0
73	\N	/api/statistics/global	GET	2025-04-20 14:27:11.255554	2025-04-20 14:27:11.256605	200	1
74	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:11.430272	2025-04-20 14:27:11.431483	200	1
75	\N	/status	GET	2025-04-20 14:27:12.409906	2025-04-20 14:27:12.409908	200	0
76	\N	/api/statistics/global	GET	2025-04-20 14:27:12.57863	2025-04-20 14:27:12.579489	200	0
77	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:12.749235	2025-04-20 14:27:12.788702	200	39
78	\N	/status	GET	2025-04-20 14:27:30.626108	2025-04-20 14:27:30.626109	200	0
79	\N	/api/statistics/global	GET	2025-04-20 14:27:30.798592	2025-04-20 14:27:30.799495	200	0
80	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:30.973127	2025-04-20 14:27:30.974339	200	1
81	\N	/status	GET	2025-04-20 14:27:32.145827	2025-04-20 14:27:32.145829	200	0
82	\N	/api/statistics/global	GET	2025-04-20 14:27:32.317655	2025-04-20 14:27:32.318863	200	1
83	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:32.509297	2025-04-20 14:27:32.510512	200	1
84	\N	/status	GET	2025-04-20 14:27:33.911702	2025-04-20 14:27:33.911703	200	0
85	\N	/status	GET	2025-04-20 14:27:34.078739	2025-04-20 14:27:34.078741	200	0
86	\N	/api/statistics/global	GET	2025-04-20 14:27:34.258348	2025-04-20 14:27:34.287041	200	28
87	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:34.579395	2025-04-20 14:27:34.582763	200	3
88	\N	/status	GET	2025-04-20 14:27:35.515238	2025-04-20 14:27:35.515239	200	0
89	\N	/api/statistics/global	GET	2025-04-20 14:27:35.681783	2025-04-20 14:27:35.683245	200	1
90	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:35.858635	2025-04-20 14:27:35.860107	200	1
91	\N	/status	GET	2025-04-20 14:27:36.713592	2025-04-20 14:27:36.713594	200	0
92	\N	/api/statistics/global	GET	2025-04-20 14:27:36.88984	2025-04-20 14:27:36.892449	200	2
93	\N	/api/statistics/endpoints	GET	2025-04-20 14:27:37.071532	2025-04-20 14:27:37.084726	200	13
94	\N	/status	GET	2025-04-20 14:28:09.27255	2025-04-20 14:28:09.272551	200	0
95	\N	/status	GET	2025-04-20 14:28:09.471785	2025-04-20 14:28:09.471787	200	0
96	\N	/api/statistics/global	GET	2025-04-20 14:28:09.535071	2025-04-20 14:28:09.535944	200	0
97	\N	/api/statistics/global	GET	2025-04-20 14:28:09.702239	2025-04-20 14:28:09.70315	200	0
98	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:09.720092	2025-04-20 14:28:09.721185	200	1
99	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:09.889229	2025-04-20 14:28:09.893225	200	3
100	\N	/status	GET	2025-04-20 14:28:11.446513	2025-04-20 14:28:11.446515	200	0
1704	\N	/status	GET	2025-04-20 20:57:25.970128	2025-04-20 20:57:25.970129	200	0
101	\N	/api/statistics/global	GET	2025-04-20 14:28:11.613734	2025-04-20 14:28:11.614783	200	1
102	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:11.78786	2025-04-20 14:28:11.788845	200	0
103	\N	/status	GET	2025-04-20 14:28:14.142109	2025-04-20 14:28:14.142111	200	0
104	\N	/api/statistics/global	GET	2025-04-20 14:28:14.376352	2025-04-20 14:28:14.377483	200	1
105	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:14.568291	2025-04-20 14:28:14.573548	200	5
106	\N	/status	GET	2025-04-20 14:28:15.177732	2025-04-20 14:28:15.177734	200	0
107	\N	/api/statistics/global	GET	2025-04-20 14:28:15.36384	2025-04-20 14:28:15.373353	200	9
108	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:15.564826	2025-04-20 14:28:15.565837	200	1
109	\N	/status	GET	2025-04-20 14:28:16.078848	2025-04-20 14:28:16.07885	200	0
110	\N	/api/statistics/global	GET	2025-04-20 14:28:16.299644	2025-04-20 14:28:16.307874	200	8
111	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:16.522586	2025-04-20 14:28:16.523594	200	1
112	\N	/status	GET	2025-04-20 14:28:17.117536	2025-04-20 14:28:17.117538	200	0
113	\N	/api/statistics/global	GET	2025-04-20 14:28:17.3121	2025-04-20 14:28:17.314621	200	2
114	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:17.495691	2025-04-20 14:28:17.49687	200	1
115	\N	/status	GET	2025-04-20 14:28:18.26309	2025-04-20 14:28:18.263092	200	0
116	\N	/api/statistics/global	GET	2025-04-20 14:28:18.436573	2025-04-20 14:28:18.437373	200	0
117	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:18.606943	2025-04-20 14:28:18.608047	200	1
118	\N	/status	GET	2025-04-20 14:28:19.340724	2025-04-20 14:28:19.340726	200	0
119	\N	/api/statistics/global	GET	2025-04-20 14:28:19.50917	2025-04-20 14:28:19.509977	200	0
120	\N	/api/statistics/endpoints	GET	2025-04-20 14:28:19.699772	2025-04-20 14:28:19.700721	200	0
121	\N	/status	GET	2025-04-20 14:28:39.254945	2025-04-20 14:28:39.254947	200	0
122	\N	/status	GET	2025-04-20 14:29:04.490337	2025-04-20 14:29:04.490339	200	0
123	\N	/api/statistics/global	GET	2025-04-20 14:29:04.660177	2025-04-20 14:29:04.661121	200	0
124	\N	/api/statistics/endpoints	GET	2025-04-20 14:29:04.842598	2025-04-20 14:29:04.843763	200	1
125	\N	/status	GET	2025-04-20 14:29:09.256408	2025-04-20 14:29:09.25641	200	0
126	\N	/status	GET	2025-04-20 14:29:36.74423	2025-04-20 14:29:36.744232	200	0
127	\N	/status	GET	2025-04-20 14:29:36.914258	2025-04-20 14:29:36.91426	200	0
128	\N	/api/statistics/global	GET	2025-04-20 14:29:36.975027	2025-04-20 14:29:36.976015	200	0
129	\N	/api/statistics/global	GET	2025-04-20 14:29:37.147196	2025-04-20 14:29:37.148409	200	1
130	\N	/api/statistics/endpoints	GET	2025-04-20 14:29:37.160715	2025-04-20 14:29:37.161742	200	1
131	\N	/api/statistics/endpoints	GET	2025-04-20 14:29:37.328625	2025-04-20 14:29:37.336807	200	8
132	\N	/status	GET	2025-04-20 14:34:05.777612	2025-04-20 14:34:05.777639	200	0
133	\N	/status	GET	2025-04-20 14:34:05.946771	2025-04-20 14:34:05.946773	200	0
134	\N	/api/statistics/global	GET	2025-04-20 14:34:06.083665	2025-04-20 14:34:06.08465	200	0
135	\N	/api/statistics/global	GET	2025-04-20 14:34:06.251429	2025-04-20 14:34:06.252839	200	1
136	\N	/api/statistics/endpoints	GET	2025-04-20 14:34:06.269293	2025-04-20 14:34:06.270678	200	1
137	\N	/api/statistics/endpoints	GET	2025-04-20 14:34:06.437896	2025-04-20 14:34:06.438979	200	1
138	\N	/status	GET	2025-04-20 14:34:35.440491	2025-04-20 14:34:35.440493	200	0
139	\N	/status	GET	2025-04-20 14:35:05.714617	2025-04-20 14:35:05.71462	200	0
140	\N	/status	GET	2025-04-20 14:35:35.735015	2025-04-20 14:35:35.735017	200	0
141	\N	/status	GET	2025-04-20 14:36:05.416765	2025-04-20 14:36:05.416767	200	0
142	\N	/status	GET	2025-04-20 14:39:04.749544	2025-04-20 14:39:04.749563	200	0
143	\N	/api/statistics/global	GET	2025-04-20 14:39:04.917184	2025-04-20 14:39:04.918163	200	0
144	\N	/status	GET	2025-04-20 14:39:05.051234	2025-04-20 14:39:05.051238	200	0
145	\N	/api/statistics/endpoints	GET	2025-04-20 14:39:05.09616	2025-04-20 14:39:05.097387	200	1
146	\N	/api/statistics/global	GET	2025-04-20 14:39:05.22981	2025-04-20 14:39:05.233543	200	3
147	\N	/api/statistics/endpoints	GET	2025-04-20 14:39:05.436049	2025-04-20 14:39:05.437241	200	1
148	\N	/api/restaurant	GET	2025-04-20 14:41:46.650251	2025-04-20 14:41:47.590583	200	940
149	\N	/api/restaurant	GET	2025-04-20 14:41:47.757679	2025-04-20 14:41:47.75774	200	0
150	\N	/status	GET	2025-04-20 14:41:51.410824	2025-04-20 14:41:51.410826	200	0
151	\N	/status	GET	2025-04-20 14:41:51.576696	2025-04-20 14:41:51.576698	200	0
152	\N	/api/statistics/global	GET	2025-04-20 14:41:51.576692	2025-04-20 14:41:51.577962	200	1
153	\N	/api/statistics/global	GET	2025-04-20 14:41:51.749788	2025-04-20 14:41:51.750701	200	0
154	\N	/api/statistics/endpoints	GET	2025-04-20 14:41:51.76843	2025-04-20 14:41:51.769566	200	1
155	\N	/api/statistics/endpoints	GET	2025-04-20 14:41:51.945419	2025-04-20 14:41:51.948061	200	2
156	\N	/status	GET	2025-04-20 14:42:08.28338	2025-04-20 14:42:08.283381	200	0
157	\N	/status	GET	2025-04-20 14:42:08.470045	2025-04-20 14:42:08.470047	200	0
158	\N	/api/statistics/global	GET	2025-04-20 14:42:08.531755	2025-04-20 14:42:08.532608	200	0
159	\N	/api/statistics/global	GET	2025-04-20 14:42:08.707054	2025-04-20 14:42:08.708044	200	0
160	\N	/api/statistics/endpoints	GET	2025-04-20 14:42:08.725807	2025-04-20 14:42:08.727319	200	1
161	\N	/api/statistics/endpoints	GET	2025-04-20 14:42:08.894148	2025-04-20 14:42:08.895484	200	1
162	\N	/api/restaurant	GET	2025-04-20 14:42:09.933214	2025-04-20 14:42:09.933279	200	0
163	\N	/api/restaurant	GET	2025-04-20 14:42:10.108468	2025-04-20 14:42:10.10855	200	0
164	\N	/api/restaurant	GET	2025-04-20 14:43:32.865632	2025-04-20 14:43:32.865693	200	0
165	\N	/api/restaurant	GET	2025-04-20 14:43:33.032561	2025-04-20 14:43:33.032624	200	0
166	\N	/status	GET	2025-04-20 14:44:14.735936	2025-04-20 14:44:14.73594	200	0
167	\N	/status	GET	2025-04-20 14:44:14.902201	2025-04-20 14:44:14.902203	200	0
168	\N	/api/statistics/global	GET	2025-04-20 14:44:15.038793	2025-04-20 14:44:15.043613	200	4
169	\N	/api/statistics/global	GET	2025-04-20 14:44:15.212811	2025-04-20 14:44:15.214052	200	1
170	\N	/api/statistics/endpoints	GET	2025-04-20 14:44:15.212796	2025-04-20 14:44:15.214403	200	1
171	\N	/api/statistics/endpoints	GET	2025-04-20 14:44:15.381565	2025-04-20 14:44:15.382549	200	0
172	\N	/api/restaurant	GET	2025-04-20 14:44:28.67739	2025-04-20 14:44:28.677455	200	0
173	\N	/api/restaurant	GET	2025-04-20 14:44:28.843422	2025-04-20 14:44:28.843489	200	0
174	\N	/api/restaurant	GET	2025-04-20 14:44:36.305486	2025-04-20 14:44:36.305546	200	0
175	\N	/api/restaurant	GET	2025-04-20 14:47:03.731616	2025-04-20 14:47:03.731697	200	0
176	\N	/api/restaurant	GET	2025-04-20 14:50:40.459403	2025-04-20 14:50:40.459465	200	0
177	\N	/api/restaurant	GET	2025-04-20 14:53:57.618631	2025-04-20 14:53:57.618703	200	0
178	\N	/api/restaurant	GET	2025-04-20 14:54:11.424133	2025-04-20 14:54:11.424209	200	0
179	\N	/status	GET	2025-04-20 14:55:11.04786	2025-04-20 14:55:11.047879	200	0
180	\N	/status	GET	2025-04-20 14:55:11.219826	2025-04-20 14:55:11.219827	200	0
181	\N	/api/statistics/global	GET	2025-04-20 14:55:11.357053	2025-04-20 14:55:11.358475	200	1
182	\N	/api/statistics/global	GET	2025-04-20 14:55:11.525059	2025-04-20 14:55:11.525953	200	0
183	\N	/api/statistics/endpoints	GET	2025-04-20 14:55:11.542934	2025-04-20 14:55:11.544149	200	1
184	\N	/api/statistics/endpoints	GET	2025-04-20 14:55:11.711078	2025-04-20 14:55:11.712488	200	1
185	\N	/status	GET	2025-04-20 14:55:40.84957	2025-04-20 14:55:40.849574	200	0
186	\N	/status	GET	2025-04-20 14:55:44.644914	2025-04-20 14:55:44.644916	200	0
187	\N	/api/statistics/global	GET	2025-04-20 14:55:44.823845	2025-04-20 14:55:44.824966	200	1
188	\N	/api/statistics/endpoints	GET	2025-04-20 14:55:45.009309	2025-04-20 14:55:45.011642	200	2
189	\N	/api/restaurant	GET	2025-04-20 14:55:48.467563	2025-04-20 14:55:48.467615	200	0
190	\N	/api/restaurant	GET	2025-04-20 14:55:48.656526	2025-04-20 14:55:48.656578	200	0
191	\N	/status	GET	2025-04-20 14:55:51.991672	2025-04-20 14:55:51.991675	200	0
192	\N	/status	GET	2025-04-20 14:55:52.158665	2025-04-20 14:55:52.158667	200	0
193	\N	/api/statistics/global	GET	2025-04-20 14:55:52.316348	2025-04-20 14:55:52.325002	200	8
194	\N	/api/statistics/global	GET	2025-04-20 14:55:52.523686	2025-04-20 14:55:52.524491	200	0
195	\N	/api/statistics/endpoints	GET	2025-04-20 14:55:52.54188	2025-04-20 14:55:52.54331	200	1
196	\N	/api/statistics/endpoints	GET	2025-04-20 14:55:52.793207	2025-04-20 14:55:52.794411	200	1
197	\N	/status	GET	2025-04-20 14:59:09.546747	2025-04-20 14:59:09.546752	200	0
198	\N	/status	GET	2025-04-20 14:59:09.716048	2025-04-20 14:59:09.71605	200	0
199	\N	/api/statistics/global	GET	2025-04-20 14:59:09.852597	2025-04-20 14:59:09.85347	200	0
200	\N	/api/statistics/global	GET	2025-04-20 14:59:10.022136	2025-04-20 14:59:10.023388	200	1
201	\N	/api/statistics/endpoints	GET	2025-04-20 14:59:10.040922	2025-04-20 14:59:10.042226	200	1
202	\N	/api/statistics/endpoints	GET	2025-04-20 14:59:10.224056	2025-04-20 14:59:10.225562	200	1
203	\N	/status	GET	2025-04-20 15:02:28.560735	2025-04-20 15:02:28.560754	200	0
204	\N	/status	GET	2025-04-20 15:02:28.815627	2025-04-20 15:02:28.815629	200	0
205	\N	/api/statistics/global	GET	2025-04-20 15:02:28.980199	2025-04-20 15:02:28.982331	200	2
206	\N	/api/statistics/global	GET	2025-04-20 15:02:29.156741	2025-04-20 15:02:29.157784	200	1
207	\N	/api/statistics/endpoints	GET	2025-04-20 15:02:29.17462	2025-04-20 15:02:29.17734	200	2
208	\N	/api/statistics/endpoints	GET	2025-04-20 15:02:29.363451	2025-04-20 15:02:29.367579	200	4
209	\N	/status	GET	2025-04-20 15:02:58.714481	2025-04-20 15:02:58.714484	200	0
210	\N	/status	GET	2025-04-20 15:03:28.711714	2025-04-20 15:03:28.711717	200	0
211	\N	/status	GET	2025-04-20 15:03:58.712305	2025-04-20 15:03:58.712307	200	0
212	\N	/status	GET	2025-04-20 15:04:28.735374	2025-04-20 15:04:28.735394	200	0
213	\N	/status	GET	2025-04-20 15:05:41.598973	2025-04-20 15:05:41.598975	200	0
214	\N	/status	GET	2025-04-20 15:06:41.571452	2025-04-20 15:06:41.571454	200	0
215	\N	/status	GET	2025-04-20 15:07:41.57498	2025-04-20 15:07:41.574981	200	0
216	\N	/status	GET	2025-04-20 15:07:47.214572	2025-04-20 15:07:47.214592	200	0
217	\N	/api/statistics/global	GET	2025-04-20 15:07:47.386956	2025-04-20 15:07:47.387868	200	0
218	\N	/api/statistics/endpoints	GET	2025-04-20 15:07:47.555615	2025-04-20 15:07:47.557017	200	1
219	\N	/api/restaurant	GET	2025-04-20 15:09:04.952388	2025-04-20 15:09:04.952467	200	0
220	\N	/api/restaurant	GET	2025-04-20 15:09:05.117917	2025-04-20 15:09:05.117973	200	0
221	\N	/api/restaurant	GET	2025-04-20 15:10:46.309638	2025-04-20 15:10:46.309706	200	0
222	\N	/api/restaurant	GET	2025-04-20 15:10:46.477629	2025-04-20 15:10:46.477695	200	0
223	\N	/api/restaurant	GET	2025-04-20 15:25:03.580339	2025-04-20 15:25:03.58041	200	0
224	\N	/api/restaurant	GET	2025-04-20 15:25:03.751276	2025-04-20 15:25:03.751328	200	0
225	\N	/api/restaurant	GET	2025-04-20 15:26:01.05078	2025-04-20 15:26:01.050853	200	0
226	\N	/api/restaurant	GET	2025-04-20 15:26:01.273382	2025-04-20 15:26:01.273449	200	0
227	\N	/api/restaurant	GET	2025-04-20 15:29:58.795929	2025-04-20 15:29:58.796009	200	0
228	\N	/api/restaurant	GET	2025-04-20 15:29:58.96256	2025-04-20 15:29:58.962628	200	0
229	\N	/status	GET	2025-04-20 15:30:00.905079	2025-04-20 15:30:00.905085	200	0
230	\N	/status	GET	2025-04-20 15:30:01.072106	2025-04-20 15:30:01.072108	200	0
231	\N	/api/statistics/global	GET	2025-04-20 15:30:01.214659	2025-04-20 15:30:01.215719	200	1
232	\N	/api/statistics/global	GET	2025-04-20 15:30:01.384773	2025-04-20 15:30:01.385712	200	0
233	\N	/api/statistics/endpoints	GET	2025-04-20 15:30:01.402736	2025-04-20 15:30:01.403802	200	1
234	\N	/api/statistics/endpoints	GET	2025-04-20 15:30:01.569158	2025-04-20 15:30:01.570383	200	1
235	\N	/status	GET	2025-04-20 15:30:31.714066	2025-04-20 15:30:31.71409	200	0
236	\N	/status	GET	2025-04-20 15:31:01.741461	2025-04-20 15:31:01.741464	200	0
237	\N	/status	GET	2025-04-20 15:31:31.696806	2025-04-20 15:31:31.69681	200	0
238	\N	/status	GET	2025-04-20 15:32:01.735898	2025-04-20 15:32:01.735901	200	0
239	\N	/status	GET	2025-04-20 15:32:41.667867	2025-04-20 15:32:41.667869	200	0
240	\N	/status	GET	2025-04-20 15:33:41.598437	2025-04-20 15:33:41.598439	200	0
241	\N	/status	GET	2025-04-20 15:34:41.496363	2025-04-20 15:34:41.496364	200	0
242	\N	/status	GET	2025-04-20 15:35:41.494323	2025-04-20 15:35:41.494324	200	0
243	\N	/status	GET	2025-04-20 15:36:41.496416	2025-04-20 15:36:41.496417	200	0
244	\N	/status	GET	2025-04-20 15:37:01.493245	2025-04-20 15:37:01.493246	200	0
245	\N	/status	GET	2025-04-20 15:37:31.494561	2025-04-20 15:37:31.494563	200	0
246	\N	/status	GET	2025-04-20 15:37:34.380678	2025-04-20 15:37:34.380679	200	0
247	\N	/api/statistics/global	GET	2025-04-20 15:37:34.54874	2025-04-20 15:37:34.550273	200	1
248	\N	/api/statistics/endpoints	GET	2025-04-20 15:37:34.717866	2025-04-20 15:37:34.719304	200	1
249	\N	/status	GET	2025-04-20 15:38:07.610615	2025-04-20 15:38:07.610616	200	0
250	\N	/api/statistics/global	GET	2025-04-20 15:38:08.157715	2025-04-20 15:38:08.16288	200	5
251	\N	/api/statistics/endpoints	GET	2025-04-20 15:38:08.552015	2025-04-20 15:38:08.553179	200	1
252	\N	/status	GET	2025-04-20 15:38:18.76762	2025-04-20 15:38:18.767622	200	0
253	\N	/api/statistics/global	GET	2025-04-20 15:38:19.117055	2025-04-20 15:38:19.118071	200	1
254	\N	/api/statistics/endpoints	GET	2025-04-20 15:38:19.629895	2025-04-20 15:38:19.631206	200	1
255	\N	/status	GET	2025-04-20 15:38:37.417934	2025-04-20 15:38:37.417936	200	0
256	\N	/status	GET	2025-04-20 15:38:41.785015	2025-04-20 15:38:41.785034	200	0
257	\N	/status	GET	2025-04-20 15:39:41.591208	2025-04-20 15:39:41.59121	200	0
258	\N	/status	GET	2025-04-20 15:39:56.28256	2025-04-20 15:39:56.282563	200	0
259	\N	/api/statistics/global	GET	2025-04-20 15:39:56.458927	2025-04-20 15:39:56.507421	200	48
260	\N	/api/statistics/endpoints	GET	2025-04-20 15:39:56.681973	2025-04-20 15:39:56.683371	200	1
261	\N	/api/restaurant	GET	2025-04-20 15:40:10.97325	2025-04-20 15:40:10.973322	200	0
262	\N	/status	GET	2025-04-20 15:40:41.647464	2025-04-20 15:40:41.647466	200	0
263	\N	/api/restaurant	GET	2025-04-20 15:41:06.651253	2025-04-20 15:41:06.651329	200	0
264	\N	/status	GET	2025-04-20 15:41:22.271262	2025-04-20 15:41:22.271264	200	0
265	\N	/api/statistics/global	GET	2025-04-20 15:41:22.452185	2025-04-20 15:41:22.453036	200	0
266	\N	/api/statistics/endpoints	GET	2025-04-20 15:41:22.630203	2025-04-20 15:41:22.632963	200	2
267	\N	/status	GET	2025-04-20 15:41:41.748743	2025-04-20 15:41:41.748763	200	0
268	\N	/status	GET	2025-04-20 15:42:12.243499	2025-04-20 15:42:12.243501	200	0
269	\N	/status	GET	2025-04-20 15:42:31.864354	2025-04-20 15:42:31.864381	200	0
270	\N	/status	GET	2025-04-20 15:43:01.578361	2025-04-20 15:43:01.578363	200	0
271	\N	/status	GET	2025-04-20 15:43:39.644451	2025-04-20 15:43:39.644452	200	0
272	\N	/status	GET	2025-04-20 15:44:01.570751	2025-04-20 15:44:01.570752	200	0
273	\N	/status	GET	2025-04-20 15:44:31.852291	2025-04-20 15:44:31.852295	200	0
274	\N	/status	GET	2025-04-20 15:44:45.631807	2025-04-20 15:44:45.631808	200	0
275	\N	/status	GET	2025-04-20 15:44:45.799083	2025-04-20 15:44:45.799085	200	0
276	\N	/api/statistics/global	GET	2025-04-20 15:44:45.938974	2025-04-20 15:44:45.939998	200	1
277	\N	/api/statistics/global	GET	2025-04-20 15:44:46.108555	2025-04-20 15:44:46.11576	200	7
278	\N	/api/statistics/endpoints	GET	2025-04-20 15:44:46.108569	2025-04-20 15:44:46.117587	200	9
279	\N	/api/statistics/endpoints	GET	2025-04-20 15:44:46.299469	2025-04-20 15:44:46.310272	200	10
280	\N	/status	GET	2025-04-20 15:47:38.081811	2025-04-20 15:47:38.081836	200	0
281	\N	/status	GET	2025-04-20 15:47:40.193657	2025-04-20 15:47:40.193658	200	0
282	\N	/api/restaurant	GET	2025-04-20 15:47:40.410838	2025-04-20 15:47:40.41091	200	0
283	\N	/api/restaurant	GET	2025-04-20 15:48:10.793151	2025-04-20 15:48:10.793225	200	0
284	\N	/status	GET	2025-04-20 15:48:35.473309	2025-04-20 15:48:35.473332	200	0
285	\N	/api/statistics/global	GET	2025-04-20 15:48:35.707218	2025-04-20 15:48:35.708367	200	1
286	\N	/api/statistics/endpoints	GET	2025-04-20 15:48:35.91772	2025-04-20 15:48:35.919009	200	1
287	\N	/api/washingmachines	GET	2025-04-20 18:21:57.353408	2025-04-20 18:21:57.523609	200	170
288	\N	/api/washingmachines	GET	2025-04-20 18:23:15.259646	2025-04-20 18:23:15.361854	200	102
289	\N	/api/washingmachines	GET	2025-04-20 18:23:17.913426	2025-04-20 18:23:17.955634	200	42
290	\N	/api/washingmachines	GET	2025-04-20 18:23:19.09982	2025-04-20 18:23:19.151959	200	52
291	\N	/api/washingmachines	GET	2025-04-20 18:23:20.336847	2025-04-20 18:23:20.389692	200	52
292	\N	/api/washingmachines	GET	2025-04-20 18:24:04.471654	2025-04-20 18:24:04.574515	200	102
293	\N	/status	GET	2025-04-20 16:24:33.299853	2025-04-20 16:24:33.299872	200	0
294	\N	/api/statistics/global	GET	2025-04-20 16:24:33.467597	2025-04-20 16:24:33.468625	200	1
295	\N	/api/statistics/endpoints	GET	2025-04-20 16:24:33.656338	2025-04-20 16:24:33.657572	200	1
296	\N	/api/restaurant	GET	2025-04-20 16:24:55.129395	2025-04-20 16:24:55.129471	200	0
297	\N	/api/washingmachines	GET	2025-04-20 19:26:45.933274	2025-04-20 19:26:46.073069	200	139
298	\N	/api/washingmachines	GET	2025-04-20 19:28:10.024751	2025-04-20 19:28:10.149735	200	124
299	\N	/api/washingmachines	GET	2025-04-20 19:28:11.607352	2025-04-20 19:28:11.652496	200	45
300	\N	/api/washingmachines	GET	2025-04-20 19:28:12.570063	2025-04-20 19:28:12.613994	200	43
301	\N	/api/washingmachines	GET	2025-04-20 19:28:13.505716	2025-04-20 19:28:13.552936	200	47
302	\N	/api/washingmachines	GET	2025-04-20 19:28:14.438832	2025-04-20 19:28:14.476257	200	37
303	\N	/api/washingmachines	GET	2025-04-20 19:28:15.326072	2025-04-20 19:28:15.370594	200	44
304	\N	/api/washingmachines	GET	2025-04-20 19:28:16.289988	2025-04-20 19:28:16.334984	200	44
305	\N	/api/washingmachines	GET	2025-04-20 19:28:17.304036	2025-04-20 19:28:17.350035	200	45
306	\N	/api/washingmachines	GET	2025-04-20 19:28:18.40533	2025-04-20 19:28:18.452036	200	46
307	\N	/api/washingmachines	GET	2025-04-20 19:28:19.297048	2025-04-20 19:28:19.343472	200	46
308	\N	/api/restaurant	GET	2025-04-20 17:33:34.802442	2025-04-20 17:33:35.806189	200	1003
309	\N	/api/restaurant	GET	2025-04-20 17:33:35.970805	2025-04-20 17:33:35.970855	200	0
310	\N	/status	GET	2025-04-20 17:33:38.147705	2025-04-20 17:33:38.14771	200	0
311	\N	/status	GET	2025-04-20 17:33:38.309977	2025-04-20 17:33:38.309979	200	0
312	\N	/api/statistics/global	GET	2025-04-20 17:33:38.445842	2025-04-20 17:33:38.447382	200	1
313	\N	/api/statistics/global	GET	2025-04-20 17:33:38.607898	2025-04-20 17:33:38.608718	200	0
323	\N	/api/statistics/endpoints	GET	2025-04-20 17:34:17.333255	2025-04-20 17:34:17.335461	200	2
324	\N	/api/statistics/endpoints	GET	2025-04-20 17:34:17.502082	2025-04-20 17:34:17.503614	200	1
325	\N	/api/restaurant	GET	2025-04-20 17:34:20.38912	2025-04-20 17:34:20.389174	200	0
326	\N	/api/restaurant	GET	2025-04-20 17:34:20.551316	2025-04-20 17:34:20.551381	200	0
327	\N	/api/washingmachines	GET	2025-04-20 17:39:13.122658	2025-04-20 17:39:13.940554	200	817
328	\N	/api/washingmachines	GET	2025-04-20 17:40:08.321149	2025-04-20 17:40:09.092731	200	771
329	\N	/api/washingmachines	GET	2025-04-20 17:40:11.421657	2025-04-20 17:40:11.59509	200	173
330	\N	/api/washingmachines	GET	2025-04-20 17:40:12.47134	2025-04-20 17:40:12.64659	200	175
331	\N	/api/washingmachines	GET	2025-04-20 17:40:12.811697	2025-04-20 17:40:12.997195	200	185
332	\N	/api/washingmachines	GET	2025-04-20 17:40:13.334969	2025-04-20 17:40:13.520382	200	185
333	\N	/api/washingmachines	GET	2025-04-20 17:40:14.126574	2025-04-20 17:40:14.299952	200	173
334	\N	/api/washingmachines	GET	2025-04-20 17:40:15.364618	2025-04-20 17:40:15.539587	200	174
335	\N	/api/washingmachines	GET	2025-04-20 17:40:15.865337	2025-04-20 17:40:16.038434	200	173
336	\N	/api/washingmachines	GET	2025-04-20 17:41:12.773866	2025-04-20 17:41:13.533663	200	759
337	\N	/api/washingmachines	GET	2025-04-20 17:42:13.102756	2025-04-20 17:42:13.882128	200	779
338	\N	/api/restaurant	GET	2025-04-20 17:43:06.07214	2025-04-20 17:43:06.072199	200	0
339	\N	/api/restaurant	GET	2025-04-20 17:43:06.234799	2025-04-20 17:43:06.234859	200	0
340	\N	/api/washingmachines	GET	2025-04-20 17:43:10.26288	2025-04-20 17:43:10.90219	200	639
341	\N	/api/washingmachines	GET	2025-04-20 17:43:11.067244	2025-04-20 17:43:11.239405	200	172
342	\N	/status	GET	2025-04-20 17:43:16.214942	2025-04-20 17:43:16.214947	200	0
343	\N	/status	GET	2025-04-20 17:43:16.431134	2025-04-20 17:43:16.431136	200	0
344	\N	/api/statistics/global	GET	2025-04-20 17:43:16.565629	2025-04-20 17:43:16.566665	200	1
345	\N	/api/statistics/global	GET	2025-04-20 17:43:16.727478	2025-04-20 17:43:16.728527	200	1
1226	\N	/api/statistics/global	GET	2025-04-20 20:47:26.607391	2025-04-20 20:47:26.608775	200	1
1227	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:26.783702	2025-04-20 20:47:26.787434	200	3
1228	\N	/status	GET	2025-04-20 20:47:31.41908	2025-04-20 20:47:31.419081	200	0
1229	\N	/status	GET	2025-04-20 20:47:31.588358	2025-04-20 20:47:31.58836	200	0
1230	\N	/api/statistics/global	GET	2025-04-20 20:47:31.605221	2025-04-20 20:47:31.606841	200	1
1231	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:31.784706	2025-04-20 20:47:31.787918	200	3
1232	\N	/status	GET	2025-04-20 20:47:36.417268	2025-04-20 20:47:36.41727	200	0
1233	\N	/status	GET	2025-04-20 20:47:36.601062	2025-04-20 20:47:36.601063	200	0
1234	\N	/api/statistics/global	GET	2025-04-20 20:47:36.620315	2025-04-20 20:47:36.621618	200	1
1235	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:36.790328	2025-04-20 20:47:36.791972	200	1
1236	\N	/status	GET	2025-04-20 20:47:41.383652	2025-04-20 20:47:41.383654	200	0
1237	\N	/status	GET	2025-04-20 20:47:41.551594	2025-04-20 20:47:41.551596	200	0
1238	\N	/api/statistics/global	GET	2025-04-20 20:47:41.569329	2025-04-20 20:47:41.570307	200	0
1239	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:41.747566	2025-04-20 20:47:41.749375	200	1
1240	\N	/status	GET	2025-04-20 20:47:46.410276	2025-04-20 20:47:46.410278	200	0
1241	\N	/status	GET	2025-04-20 20:47:46.588611	2025-04-20 20:47:46.588612	200	0
1358	\N	/api/statistics/global	GET	2025-04-20 20:50:11.548448	2025-04-20 20:50:11.549935	200	1
1359	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:11.725065	2025-04-20 20:50:11.726963	200	1
1360	\N	/status	GET	2025-04-20 20:50:16.385982	2025-04-20 20:50:16.385983	200	0
1361	\N	/status	GET	2025-04-20 20:50:16.552498	2025-04-20 20:50:16.552499	200	0
1362	\N	/api/statistics/global	GET	2025-04-20 20:50:16.571301	2025-04-20 20:50:16.572425	200	1
1363	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:16.741506	2025-04-20 20:50:16.743902	200	2
1364	\N	/status	GET	2025-04-20 20:50:21.428138	2025-04-20 20:50:21.428139	200	0
1365	\N	/status	GET	2025-04-20 20:50:21.594965	2025-04-20 20:50:21.594967	200	0
1366	\N	/api/statistics/global	GET	2025-04-20 20:50:21.613574	2025-04-20 20:50:21.614856	200	1
1367	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:21.786864	2025-04-20 20:50:21.788578	200	1
1368	\N	/status	GET	2025-04-20 20:50:26.402022	2025-04-20 20:50:26.402024	200	0
1369	\N	/status	GET	2025-04-20 20:50:26.581576	2025-04-20 20:50:26.581578	200	0
1370	\N	/api/statistics/global	GET	2025-04-20 20:50:26.58781	2025-04-20 20:50:26.588749	200	0
1371	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:26.765205	2025-04-20 20:50:26.808096	200	42
1372	\N	/status	GET	2025-04-20 20:50:31.402443	2025-04-20 20:50:31.402444	200	0
1373	\N	/status	GET	2025-04-20 20:50:31.572923	2025-04-20 20:50:31.572925	200	0
1374	\N	/api/statistics/global	GET	2025-04-20 20:50:31.591016	2025-04-20 20:50:31.592378	200	1
1375	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:31.764544	2025-04-20 20:50:31.766301	200	1
1376	\N	/status	GET	2025-04-20 20:50:36.401825	2025-04-20 20:50:36.401826	200	0
1377	\N	/status	GET	2025-04-20 20:50:36.570337	2025-04-20 20:50:36.570339	200	0
1378	\N	/api/statistics/global	GET	2025-04-20 20:50:36.58881	2025-04-20 20:50:36.58986	200	1
1379	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:36.762139	2025-04-20 20:50:36.764537	200	2
1380	\N	/status	GET	2025-04-20 20:50:41.418222	2025-04-20 20:50:41.418224	200	0
1381	\N	/status	GET	2025-04-20 20:50:41.588687	2025-04-20 20:50:41.588688	200	0
1382	\N	/api/statistics/global	GET	2025-04-20 20:50:41.607064	2025-04-20 20:50:41.609131	200	2
1383	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:41.781548	2025-04-20 20:50:41.784097	200	2
1384	\N	/status	GET	2025-04-20 20:50:46.392221	2025-04-20 20:50:46.392223	200	0
1385	\N	/status	GET	2025-04-20 20:50:46.564801	2025-04-20 20:50:46.564803	200	0
1414	\N	/api/statistics/global	GET	2025-04-20 20:51:21.578925	2025-04-20 20:51:21.606107	200	27
1415	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:21.841474	2025-04-20 20:51:21.843354	200	1
1416	\N	/status	GET	2025-04-20 20:51:26.408588	2025-04-20 20:51:26.408589	200	0
1417	\N	/status	GET	2025-04-20 20:51:26.602674	2025-04-20 20:51:26.602676	200	0
1418	\N	/api/statistics/global	GET	2025-04-20 20:51:26.613577	2025-04-20 20:51:26.614988	200	1
1419	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:26.792991	2025-04-20 20:51:26.794852	200	1
1420	\N	/status	GET	2025-04-20 20:51:31.402485	2025-04-20 20:51:31.402486	200	0
1421	\N	/status	GET	2025-04-20 20:51:31.571237	2025-04-20 20:51:31.571239	200	0
1422	\N	/api/statistics/global	GET	2025-04-20 20:51:31.589355	2025-04-20 20:51:31.590801	200	1
1423	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:31.762426	2025-04-20 20:51:31.764186	200	1
1424	\N	/status	GET	2025-04-20 20:51:36.415358	2025-04-20 20:51:36.41536	200	0
1425	\N	/status	GET	2025-04-20 20:51:36.582789	2025-04-20 20:51:36.58279	200	0
1426	\N	/api/statistics/global	GET	2025-04-20 20:51:36.602664	2025-04-20 20:51:36.605243	200	2
1427	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:36.782328	2025-04-20 20:51:36.784438	200	2
1428	\N	/status	GET	2025-04-20 20:51:41.403839	2025-04-20 20:51:41.40384	200	0
1429	\N	/status	GET	2025-04-20 20:51:41.571629	2025-04-20 20:51:41.571631	200	0
1430	\N	/api/statistics/global	GET	2025-04-20 20:51:41.589663	2025-04-20 20:51:41.591021	200	1
1431	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:41.759576	2025-04-20 20:51:41.761286	200	1
1432	\N	/status	GET	2025-04-20 20:51:46.421731	2025-04-20 20:51:46.421733	200	0
1433	\N	/status	GET	2025-04-20 20:51:46.595759	2025-04-20 20:51:46.595761	200	0
1453	\N	/status	GET	2025-04-20 20:52:11.589086	2025-04-20 20:52:11.589087	200	0
1458	\N	/api/statistics/global	GET	2025-04-20 20:52:16.639721	2025-04-20 20:52:16.641036	200	1
1459	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:16.818361	2025-04-20 20:52:16.820321	200	1
1460	\N	/status	GET	2025-04-20 20:52:21.411644	2025-04-20 20:52:21.411645	200	0
1461	\N	/status	GET	2025-04-20 20:52:21.579437	2025-04-20 20:52:21.579438	200	0
1462	\N	/api/statistics/global	GET	2025-04-20 20:52:21.597384	2025-04-20 20:52:21.598832	200	1
1463	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:21.772649	2025-04-20 20:52:21.774503	200	1
1554	\N	/api/statistics/global	GET	2025-04-20 20:54:16.610369	2025-04-20 20:54:16.640747	200	30
1555	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:16.823549	2025-04-20 20:54:16.825648	200	2
1556	\N	/status	GET	2025-04-20 20:54:21.416373	2025-04-20 20:54:21.416374	200	0
1557	\N	/status	GET	2025-04-20 20:54:21.587443	2025-04-20 20:54:21.587444	200	0
314	\N	/api/statistics/endpoints	GET	2025-04-20 17:33:38.609327	2025-04-20 17:33:38.611353	200	2
315	\N	/api/statistics/endpoints	GET	2025-04-20 17:33:38.774763	2025-04-20 17:33:38.77597	200	1
316	\N	/status	GET	2025-04-20 17:33:45.685287	2025-04-20 17:33:45.685291	200	0
317	\N	/api/statistics/global	GET	2025-04-20 17:33:45.847167	2025-04-20 17:33:45.848262	200	1
318	\N	/api/statistics/endpoints	GET	2025-04-20 17:33:46.013734	2025-04-20 17:33:46.01553	200	1
319	\N	/status	GET	2025-04-20 17:34:16.812781	2025-04-20 17:34:16.812784	200	0
320	\N	/status	GET	2025-04-20 17:34:17.029555	2025-04-20 17:34:17.029556	200	0
321	\N	/api/statistics/global	GET	2025-04-20 17:34:17.16798	2025-04-20 17:34:17.168997	200	1
322	\N	/api/statistics/global	GET	2025-04-20 17:34:17.330938	2025-04-20 17:34:17.332155	200	1
346	\N	/api/statistics/endpoints	GET	2025-04-20 17:43:16.727491	2025-04-20 17:43:16.728952	200	1
347	\N	/api/statistics/endpoints	GET	2025-04-20 17:43:16.890468	2025-04-20 17:43:16.892071	200	1
348	\N	/api/washingmachines	GET	2025-04-20 17:43:26.146866	2025-04-20 17:43:26.90835	200	761
349	\N	/api/washingmachines	GET	2025-04-20 17:43:27.067358	2025-04-20 17:43:27.252923	200	185
350	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-20 17:44:18.891046	2025-04-20 17:44:18.893628	200	2
351	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-20 17:44:19.21849	2025-04-20 17:44:19.220745	200	2
352	\N	/api/restaurant	GET	2025-04-20 17:44:19.358772	2025-04-20 17:44:19.35884	200	0
353	\N	/api/washingmachines	GET	2025-04-20 17:46:16.459844	2025-04-20 17:46:17.230637	200	770
354	\N	/api/washingmachines	GET	2025-04-20 17:46:19.241083	2025-04-20 17:46:19.413289	200	172
355	\N	/api/washingmachines	GET	2025-04-20 17:46:20.622502	2025-04-20 17:46:20.806744	200	184
356	\N	/api/washingmachines	GET	2025-04-20 19:46:26.823003	2025-04-20 19:46:26.960287	200	137
357	\N	/api/washingmachines	GET	2025-04-20 19:46:29.814553	2025-04-20 19:46:29.861736	200	47
358	\N	/api/washingmachines	GET	2025-04-20 19:46:30.919182	2025-04-20 19:46:30.958333	200	39
359	\N	/api/washingmachines	GET	2025-04-20 19:46:32.173201	2025-04-20 19:46:32.218158	200	44
360	\N	/api/washingmachines	GET	2025-04-20 19:46:54.318767	2025-04-20 19:46:54.410777	200	92
361	\N	/api/washingmachines	GET	2025-04-20 17:46:58.410268	2025-04-20 17:46:59.10084	200	690
362	\N	/api/washingmachines	GET	2025-04-20 17:51:16.924723	2025-04-20 17:51:17.545548	200	620
363	\N	/api/washingmachines	GET	2025-04-20 17:54:05.156128	2025-04-20 17:54:05.837577	200	681
364	\N	/api/washingmachines	GET	2025-04-20 17:54:29.366778	2025-04-20 17:54:30.064609	200	697
365	\N	/api/washingmachines	GET	2025-04-20 17:55:55.834229	2025-04-20 17:55:56.472111	200	637
366	\N	/api/washingmachines	GET	2025-04-20 17:56:00.460654	2025-04-20 17:56:00.636144	200	175
367	\N	/api/washingmachines	GET	2025-04-20 17:56:02.388841	2025-04-20 17:56:02.573226	200	184
368	\N	/api/washingmachines	GET	2025-04-20 17:56:17.288213	2025-04-20 17:56:17.92138	200	633
369	\N	/api/washingmachines	GET	2025-04-20 17:56:19.336971	2025-04-20 17:56:19.511384	200	174
370	\N	/api/washingmachines	GET	2025-04-20 17:58:09.562028	2025-04-20 17:58:10.254258	200	692
371	\N	/api/washingmachines	GET	2025-04-20 17:58:11.510833	2025-04-20 17:58:11.685603	200	174
372	\N	/api/washingmachines	GET	2025-04-20 17:58:12.500612	2025-04-20 17:58:12.689128	200	188
373	\N	/api/washingmachines	GET	2025-04-20 17:59:29.402147	2025-04-20 17:59:30.093776	200	691
374	\N	/api/washingmachines	GET	2025-04-20 18:04:29.379646	2025-04-20 18:04:30.077525	200	697
375	\N	/api/washingmachines	GET	2025-04-20 20:08:43.777557	2025-04-20 20:08:43.964109	200	186
376	\N	/api/washingmachines	GET	2025-04-20 18:09:29.540931	2025-04-20 18:09:30.226104	200	685
377	\N	/api/washingmachines	GET	2025-04-20 18:10:15.867019	2025-04-20 18:10:16.562578	200	695
378	\N	/api/washingmachines	GET	2025-04-20 18:12:01.621483	2025-04-20 18:12:02.286503	200	665
379	\N	/api/restaurant	GET	2025-04-20 18:12:13.497789	2025-04-20 18:12:13.497892	200	0
380	\N	/status	GET	2025-04-20 18:12:19.099953	2025-04-20 18:12:19.099959	200	0
381	\N	/api/statistics/global	GET	2025-04-20 18:12:19.262246	2025-04-20 18:12:19.26408	200	1
382	\N	/api/statistics/endpoints	GET	2025-04-20 18:12:19.425102	2025-04-20 18:12:19.452585	200	27
383	\N	/status	GET	2025-04-20 18:12:49.233575	2025-04-20 18:12:49.233578	200	0
384	\N	/api/washingmachines	GET	2025-04-20 18:13:12.178721	2025-04-20 18:13:12.97407	200	795
385	\N	/api/washingmachines	GET	2025-04-20 18:18:12.315959	2025-04-20 18:18:13.023459	200	707
386	\N	/status	GET	2025-04-20 18:18:18.948104	2025-04-20 18:18:18.94814	200	0
387	\N	/api/statistics/global	GET	2025-04-20 18:18:19.144395	2025-04-20 18:18:19.145453	200	1
388	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:19.30967	2025-04-20 18:18:19.311042	200	1
389	\N	/status	GET	2025-04-20 18:18:41.538066	2025-04-20 18:18:41.53807	200	0
390	\N	/api/statistics/global	GET	2025-04-20 18:18:41.725133	2025-04-20 18:18:41.726081	200	0
391	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:41.886465	2025-04-20 18:18:41.887934	200	1
392	\N	/status	GET	2025-04-20 18:18:49.097182	2025-04-20 18:18:49.097184	200	0
393	\N	/status	GET	2025-04-20 18:18:51.947352	2025-04-20 18:18:51.947354	200	0
394	\N	/api/statistics/global	GET	2025-04-20 18:18:52.109426	2025-04-20 18:18:52.1112	200	1
395	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:52.280437	2025-04-20 18:18:52.281721	200	1
396	\N	/status	GET	2025-04-20 18:18:53.165715	2025-04-20 18:18:53.165717	200	0
397	\N	/api/statistics/global	GET	2025-04-20 18:18:53.341037	2025-04-20 18:18:53.346747	200	5
398	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:53.519436	2025-04-20 18:18:53.520551	200	1
399	\N	/status	GET	2025-04-20 18:18:54.190997	2025-04-20 18:18:54.190999	200	0
400	\N	/api/statistics/global	GET	2025-04-20 18:18:54.373434	2025-04-20 18:18:54.374449	200	1
401	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:54.545671	2025-04-20 18:18:54.546962	200	1
402	\N	/status	GET	2025-04-20 18:18:55.345083	2025-04-20 18:18:55.345085	200	0
403	\N	/api/statistics/global	GET	2025-04-20 18:18:55.509485	2025-04-20 18:18:55.510417	200	0
404	\N	/api/statistics/endpoints	GET	2025-04-20 18:18:55.675316	2025-04-20 18:18:55.676755	200	1
405	\N	/status	GET	2025-04-20 18:19:08.393897	2025-04-20 18:19:08.3939	200	0
406	\N	/api/statistics/global	GET	2025-04-20 18:19:08.557924	2025-04-20 18:19:08.559103	200	1
407	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:08.724906	2025-04-20 18:19:08.726155	200	1
408	\N	/status	GET	2025-04-20 18:19:10.597043	2025-04-20 18:19:10.597045	200	0
409	\N	/api/statistics/global	GET	2025-04-20 18:19:10.767687	2025-04-20 18:19:10.768705	200	1
410	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:10.93306	2025-04-20 18:19:10.934818	200	1
411	\N	/status	GET	2025-04-20 18:19:11.532553	2025-04-20 18:19:11.532553	200	0
412	\N	/api/statistics/global	GET	2025-04-20 18:19:11.697847	2025-04-20 18:19:11.699592	200	1
413	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:11.870285	2025-04-20 18:19:11.871512	200	1
414	\N	/status	GET	2025-04-20 18:19:12.546323	2025-04-20 18:19:12.546325	200	0
415	\N	/api/statistics/global	GET	2025-04-20 18:19:12.709589	2025-04-20 18:19:12.710534	200	0
416	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:12.88604	2025-04-20 18:19:12.887576	200	1
417	\N	/status	GET	2025-04-20 18:19:13.529924	2025-04-20 18:19:13.529925	200	0
418	\N	/api/statistics/global	GET	2025-04-20 18:19:13.693879	2025-04-20 18:19:13.694993	200	1
419	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:13.859811	2025-04-20 18:19:13.86121	200	1
420	\N	/status	GET	2025-04-20 18:19:14.490675	2025-04-20 18:19:14.490677	200	0
421	\N	/api/statistics/global	GET	2025-04-20 18:19:14.662804	2025-04-20 18:19:14.663714	200	0
422	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:14.831037	2025-04-20 18:19:14.83219	200	1
423	\N	/status	GET	2025-04-20 18:19:18.922986	2025-04-20 18:19:18.923008	200	0
424	\N	/status	GET	2025-04-20 18:19:20.467746	2025-04-20 18:19:20.467748	200	0
425	\N	/api/statistics/global	GET	2025-04-20 18:19:20.640176	2025-04-20 18:19:20.641171	200	0
426	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:20.803358	2025-04-20 18:19:20.804505	200	1
427	\N	/api/restaurant	GET	2025-04-20 18:19:22.339792	2025-04-20 18:19:22.339844	200	0
428	\N	/status	GET	2025-04-20 18:19:37.542047	2025-04-20 18:19:37.542049	200	0
429	\N	/api/statistics/global	GET	2025-04-20 18:19:37.709913	2025-04-20 18:19:37.71123	200	1
430	\N	/api/statistics/endpoints	GET	2025-04-20 18:19:37.872646	2025-04-20 18:19:37.874206	200	1
431	\N	/status	GET	2025-04-20 18:20:07.65865	2025-04-20 18:20:07.658674	200	0
432	\N	/status	GET	2025-04-20 18:20:17.599468	2025-04-20 18:20:17.599469	200	0
433	\N	/api/statistics/global	GET	2025-04-20 18:20:17.763335	2025-04-20 18:20:17.764393	200	1
434	\N	/api/statistics/endpoints	GET	2025-04-20 18:20:17.928594	2025-04-20 18:20:17.930201	200	1
435	\N	/status	GET	2025-04-20 18:20:37.52001	2025-04-20 18:20:37.520012	200	0
436	\N	/status	GET	2025-04-20 18:21:07.520246	2025-04-20 18:21:07.520247	200	0
437	\N	/status	GET	2025-04-20 18:21:37.546847	2025-04-20 18:21:37.546849	200	0
438	\N	/status	GET	2025-04-20 18:22:07.515868	2025-04-20 18:22:07.515868	200	0
439	\N	/status	GET	2025-04-20 18:22:37.522889	2025-04-20 18:22:37.52289	200	0
440	\N	/status	GET	2025-04-20 18:23:07.547929	2025-04-20 18:23:07.547941	200	0
441	\N	/status	GET	2025-04-20 18:23:37.520346	2025-04-20 18:23:37.520348	200	0
442	\N	/status	GET	2025-04-20 18:24:07.518914	2025-04-20 18:24:07.518916	200	0
443	\N	/status	GET	2025-04-20 18:24:37.570961	2025-04-20 18:24:37.570962	200	0
444	\N	/status	GET	2025-04-20 18:25:07.532526	2025-04-20 18:25:07.532528	200	0
445	\N	/status	GET	2025-04-20 18:25:37.660347	2025-04-20 18:25:37.660373	200	0
446	\N	/status	GET	2025-04-20 18:26:07.547442	2025-04-20 18:26:07.547443	200	0
447	\N	/status	GET	2025-04-20 18:26:37.788807	2025-04-20 18:26:37.788808	200	0
448	\N	/status	GET	2025-04-20 18:27:07.522834	2025-04-20 18:27:07.522837	200	0
449	\N	/status	GET	2025-04-20 18:27:30.397895	2025-04-20 18:27:30.397919	200	0
450	\N	/api/statistics/global	GET	2025-04-20 18:27:30.560809	2025-04-20 18:27:30.563464	200	2
451	\N	/api/statistics/endpoints	GET	2025-04-20 18:27:30.731609	2025-04-20 18:27:30.733126	200	1
452	\N	/status	GET	2025-04-20 18:27:35.890799	2025-04-20 18:27:35.8908	200	0
453	\N	/api/statistics/global	GET	2025-04-20 18:27:36.056919	2025-04-20 18:27:36.061194	200	4
454	\N	/api/statistics/endpoints	GET	2025-04-20 18:27:36.234276	2025-04-20 18:27:36.24326	200	8
455	\N	/statistics	GET	2025-04-20 20:27:45.456452	2025-04-20 20:27:45.456527	500	0
456	\N	/status	GET	2025-04-20 18:30:02.571084	2025-04-20 18:30:02.571122	200	0
457	\N	/api/statistics/global	GET	2025-04-20 18:30:02.736284	2025-04-20 18:30:02.73734	200	1
458	\N	/api/statistics/endpoints	GET	2025-04-20 18:30:02.898793	2025-04-20 18:30:02.902698	200	3
459	\N	/status	GET	2025-04-20 18:30:07.366031	2025-04-20 18:30:07.366033	200	0
460	\N	/api/restaurant	GET	2025-04-20 18:30:12.892398	2025-04-20 18:30:12.892457	200	0
461	\N	/status	GET	2025-04-20 18:30:19.165911	2025-04-20 18:30:19.165912	200	0
462	\N	/api/statistics/global	GET	2025-04-20 18:30:19.35606	2025-04-20 18:30:19.38885	200	32
463	\N	/api/statistics/endpoints	GET	2025-04-20 18:30:19.612462	2025-04-20 18:30:19.614013	200	1
464	\N	/status	GET	2025-04-20 18:30:24.158113	2025-04-20 18:30:24.158115	200	0
465	\N	/api/washingmachines	GET	2025-04-20 18:30:24.684353	2025-04-20 18:30:25.388161	200	703
466	\N	/api/restaurant	GET	2025-04-20 18:30:44.401222	2025-04-20 18:30:44.401289	200	0
467	\N	/status	GET	2025-04-20 18:31:44.151441	2025-04-20 18:31:44.151442	200	0
468	\N	/api/statistics/global	GET	2025-04-20 18:31:44.324838	2025-04-20 18:31:44.326005	200	1
469	\N	/api/statistics/endpoints	GET	2025-04-20 18:31:44.487484	2025-04-20 18:31:44.488787	200	1
470	\N	/status	GET	2025-04-20 18:31:48.870323	2025-04-20 18:31:48.870324	200	0
471	\N	/api/restaurant	GET	2025-04-20 18:31:49.734463	2025-04-20 18:31:49.734515	200	0
472	\N	/api/restaurant	GET	2025-04-20 18:35:30.965676	2025-04-20 18:35:30.965748	200	0
473	\N	/status	GET	2025-04-20 18:35:34.228252	2025-04-20 18:35:34.228257	200	0
474	\N	/api/statistics/global	GET	2025-04-20 18:35:34.401609	2025-04-20 18:35:34.402561	200	0
475	\N	/api/statistics/endpoints	GET	2025-04-20 18:35:34.563594	2025-04-20 18:35:34.584616	200	21
476	\N	/api/washingmachines	GET	2025-04-20 18:35:39.083276	2025-04-20 18:35:39.77295	200	689
477	\N	/api/washingmachines	GET	2025-04-20 18:35:56.603303	2025-04-20 18:35:57.305333	200	702
478	\N	/api/washingmachines	GET	2025-04-20 18:37:48.298283	2025-04-20 18:37:48.923174	200	624
479	\N	/api/restaurant	GET	2025-04-20 18:39:35.558285	2025-04-20 18:39:35.558362	200	0
480	\N	/api/restaurant	GET	2025-04-20 18:40:29.928902	2025-04-20 18:40:29.928971	200	0
481	\N	/status	GET	2025-04-20 18:40:35.366302	2025-04-20 18:40:35.366322	200	0
482	\N	/api/statistics/global	GET	2025-04-20 18:40:35.530242	2025-04-20 18:40:35.53141	200	1
483	\N	/api/statistics/endpoints	GET	2025-04-20 18:40:35.691942	2025-04-20 18:40:35.720947	200	29
484	\N	/api/restaurant	GET	2025-04-20 18:45:32.057876	2025-04-20 18:45:32.057943	200	0
485	\N	/api/restaurant	GET	2025-04-20 18:45:32.236385	2025-04-20 18:45:32.236436	200	0
486	\N	/api/restaurant	GET	2025-04-20 18:56:35.454502	2025-04-20 18:56:35.454587	200	0
487	\N	/api/restaurant	GET	2025-04-20 18:56:35.618269	2025-04-20 18:56:35.618316	200	0
488	\N	/status	GET	2025-04-20 18:56:37.910407	2025-04-20 18:56:37.910413	200	0
489	\N	/status	GET	2025-04-20 18:56:38.071923	2025-04-20 18:56:38.071924	200	0
490	\N	/api/statistics/global	GET	2025-04-20 18:56:38.208242	2025-04-20 18:56:38.209272	200	1
491	\N	/api/statistics/endpoints	GET	2025-04-20 18:56:38.381735	2025-04-20 18:56:38.382989	200	1
492	\N	/api/statistics/global	GET	2025-04-20 18:56:38.381737	2025-04-20 18:56:38.383229	200	1
493	\N	/api/statistics/endpoints	GET	2025-04-20 18:56:38.573635	2025-04-20 18:56:38.575544	200	1
494	\N	/status	GET	2025-04-20 18:56:42.902895	2025-04-20 18:56:42.902896	200	0
495	\N	/status	GET	2025-04-20 18:56:47.911195	2025-04-20 18:56:47.911196	200	0
496	\N	/status	GET	2025-04-20 18:56:52.91449	2025-04-20 18:56:52.914492	200	0
497	\N	/status	GET	2025-04-20 18:56:57.900206	2025-04-20 18:56:57.900207	200	0
498	\N	/status	GET	2025-04-20 18:57:02.903764	2025-04-20 18:57:02.903766	200	0
499	\N	/status	GET	2025-04-20 18:57:03.06591	2025-04-20 18:57:03.065912	200	0
500	\N	/api/statistics/global	GET	2025-04-20 18:57:03.235844	2025-04-20 18:57:03.237564	200	1
501	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:03.426066	2025-04-20 18:57:03.428596	200	2
502	\N	/status	GET	2025-04-20 18:57:04.542486	2025-04-20 18:57:04.542487	200	0
503	\N	/api/statistics/global	GET	2025-04-20 18:57:04.705751	2025-04-20 18:57:04.707126	200	1
504	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:04.868864	2025-04-20 18:57:04.870071	200	1
505	\N	/status	GET	2025-04-20 18:57:07.196197	2025-04-20 18:57:07.1962	200	0
506	\N	/api/statistics/global	GET	2025-04-20 18:57:07.379319	2025-04-20 18:57:07.380534	200	1
507	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:07.545412	2025-04-20 18:57:07.546578	200	1
508	\N	/status	GET	2025-04-20 18:57:07.900688	2025-04-20 18:57:07.900689	200	0
509	\N	/status	GET	2025-04-20 18:57:08.339242	2025-04-20 18:57:08.339243	200	0
510	\N	/api/statistics/global	GET	2025-04-20 18:57:08.502548	2025-04-20 18:57:08.503482	200	0
511	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:08.687021	2025-04-20 18:57:08.688639	200	1
512	\N	/status	GET	2025-04-20 18:57:09.482384	2025-04-20 18:57:09.482386	200	0
513	\N	/api/statistics/global	GET	2025-04-20 18:57:09.64623	2025-04-20 18:57:09.647462	200	1
514	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:09.813931	2025-04-20 18:57:09.815392	200	1
515	\N	/status	GET	2025-04-20 18:57:10.422294	2025-04-20 18:57:10.422295	200	0
516	\N	/api/statistics/global	GET	2025-04-20 18:57:10.586166	2025-04-20 18:57:10.58734	200	1
517	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:10.753431	2025-04-20 18:57:10.755208	200	1
518	\N	/status	GET	2025-04-20 18:57:11.320817	2025-04-20 18:57:11.320819	200	0
519	\N	/api/statistics/global	GET	2025-04-20 18:57:11.507282	2025-04-20 18:57:11.508153	200	0
520	\N	/api/statistics/endpoints	GET	2025-04-20 18:57:11.71751	2025-04-20 18:57:11.720401	200	2
521	\N	/status	GET	2025-04-20 18:57:12.900184	2025-04-20 18:57:12.900186	200	0
522	\N	/status	GET	2025-04-20 18:57:17.900062	2025-04-20 18:57:17.900064	200	0
523	\N	/status	GET	2025-04-20 18:57:22.904985	2025-04-20 18:57:22.904987	200	0
524	\N	/api/restaurant	GET	2025-04-20 18:57:26.895317	2025-04-20 18:57:26.89537	200	0
525	\N	/api/restaurant	GET	2025-04-20 18:57:27.409929	2025-04-20 18:57:27.409984	200	0
526	\N	/api/washingmachines	GET	2025-04-20 19:02:45.113289	2025-04-20 19:02:45.880411	200	767
527	\N	/api/washingmachines	GET	2025-04-20 19:02:46.04915	2025-04-20 19:02:46.233482	200	184
528	\N	/api/washingmachines	GET	2025-04-20 19:04:35.32782	2025-04-20 19:04:36.093246	200	765
529	\N	/api/restaurant	GET	2025-04-20 19:06:20.126739	2025-04-20 19:06:20.126813	200	0
530	\N	/api/restaurant	GET	2025-04-20 19:06:20.288123	2025-04-20 19:06:20.288175	200	0
531	\N	/status	GET	2025-04-20 19:08:07.979441	2025-04-20 19:08:07.979446	200	0
532	\N	/api/statistics/global	GET	2025-04-20 19:08:08.179298	2025-04-20 19:08:08.180355	200	1
533	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:08.379319	2025-04-20 19:08:08.38073	200	1
534	\N	/status	GET	2025-04-20 19:08:12.299332	2025-04-20 19:08:12.299333	200	0
535	\N	/api/statistics/global	GET	2025-04-20 19:08:12.53924	2025-04-20 19:08:12.54035	200	1
536	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:12.879274	2025-04-20 19:08:12.880585	200	1
537	\N	/status	GET	2025-04-20 19:08:14.318324	2025-04-20 19:08:14.318326	200	0
538	\N	/api/statistics/global	GET	2025-04-20 19:08:14.628422	2025-04-20 19:08:14.629377	200	0
539	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:14.839325	2025-04-20 19:08:14.840737	200	1
540	\N	/api/restaurant	GET	2025-04-20 19:08:23.74339	2025-04-20 19:08:23.743448	200	0
541	\N	/api/washingmachines	GET	2025-04-20 19:08:34.014895	2025-04-20 19:08:34.794999	200	780
542	\N	/api/washingmachines	GET	2025-04-20 19:08:34.954452	2025-04-20 19:08:35.141214	200	186
543	\N	/api/restaurant	GET	2025-04-20 19:08:52.258766	2025-04-20 19:08:52.258834	200	0
544	\N	/api/restaurant	GET	2025-04-20 19:08:52.42955	2025-04-20 19:08:52.429601	200	0
545	\N	/status	GET	2025-04-20 19:08:53.546907	2025-04-20 19:08:53.546911	200	0
546	\N	/status	GET	2025-04-20 19:08:53.787641	2025-04-20 19:08:53.787643	200	0
547	\N	/api/statistics/global	GET	2025-04-20 19:08:53.92198	2025-04-20 19:08:53.923322	200	1
548	\N	/api/statistics/global	GET	2025-04-20 19:08:54.082677	2025-04-20 19:08:54.083851	200	1
549	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:54.083911	2025-04-20 19:08:54.086088	200	2
550	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:54.248846	2025-04-20 19:08:54.250977	200	2
551	\N	/status	GET	2025-04-20 19:08:57.193484	2025-04-20 19:08:57.193486	200	0
552	\N	/api/statistics/global	GET	2025-04-20 19:08:57.38492	2025-04-20 19:08:57.386137	200	1
553	\N	/api/statistics/endpoints	GET	2025-04-20 19:08:57.553592	2025-04-20 19:08:57.555524	200	1
554	\N	/status	GET	2025-04-20 19:08:58.55366	2025-04-20 19:08:58.553661	200	0
555	\N	/api/restaurant	GET	2025-04-20 19:11:35.810739	2025-04-20 19:11:35.810814	200	0
556	\N	/api/restaurant	GET	2025-04-20 19:11:35.978418	2025-04-20 19:11:35.978481	200	0
557	\N	/status	GET	2025-04-20 19:12:02.50678	2025-04-20 19:12:02.506786	200	0
558	\N	/status	GET	2025-04-20 19:12:02.724909	2025-04-20 19:12:02.72491	200	0
559	\N	/api/statistics/global	GET	2025-04-20 19:12:02.860003	2025-04-20 19:12:02.861108	200	1
560	\N	/api/statistics/global	GET	2025-04-20 19:12:03.026869	2025-04-20 19:12:03.027951	200	1
561	\N	/api/statistics/endpoints	GET	2025-04-20 19:12:03.028079	2025-04-20 19:12:03.029851	200	1
562	\N	/api/statistics/endpoints	GET	2025-04-20 19:12:03.190215	2025-04-20 19:12:03.202284	200	12
563	\N	/status	GET	2025-04-20 19:12:07.514319	2025-04-20 19:12:07.514321	200	0
564	\N	/api/upload	POST	2025-04-20 22:03:29.799831	2025-04-20 22:03:29.801333	401	1
565	\N	/api/auth/login	POST	2025-04-20 22:03:40.299328	2025-04-20 22:03:40.702538	200	403
566	yohann.chavanel@imt-atlantique.net	/api/upload	POST	2025-04-20 22:03:45.79896	2025-04-20 22:03:46.144306	200	345
567	\N	/api/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:04:01.271625	2025-04-20 20:04:01.272979	404	1
568	\N	/favicon.ico	GET	2025-04-20 20:04:01.818799	2025-04-20 20:04:01.81881	500	0
569	\N	/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:04:09.153977	2025-04-20 20:04:09.153981	500	0
570	\N	/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:04:12.299919	2025-04-20 20:04:12.299923	500	0
571	\N	//api/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:04:18.067045	2025-04-20 20:04:18.06705	500	0
572	\N	/api/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:04:21.113614	2025-04-20 20:04:21.113653	404	0
573	yohann.chavanel@imt-atlantique.net	/api/files	GET	2025-04-20 22:04:26.000047	2025-04-20 22:04:26.335574	200	335
574	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-04-20 20:04:45.105676	2025-04-20 20:04:45.105875	200	0
575	yohann.chavanel@imt-atlantique.net	/api/all-files	GET	2025-04-20 22:04:52.57734	2025-04-20 22:04:52.578092	200	0
576	\N	/api/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:05:05.841512	2025-04-20 20:05:05.841575	404	0
577	test@imt-atlantique.net	/api/upload	POST	2025-04-20 20:05:15.324075	2025-04-20 20:05:15.335135	200	11
578	\N	/api/data/New_Project_8626bf35661a2046.png	GET	2025-04-20 20:05:19.661073	2025-04-20 20:05:19.661124	404	0
579	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:05:32.222701	2025-04-20 20:05:32.222803	200	0
580	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:05:41.165602	2025-04-20 20:05:41.165673	200	0
581	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:08:55.455997	2025-04-20 20:08:55.456084	304	0
582	\N	/api/restaurant	GET	2025-04-20 20:13:57.069807	2025-04-20 20:13:57.069886	200	0
583	\N	/api/restaurant	GET	2025-04-20 20:13:57.245416	2025-04-20 20:13:57.245474	200	0
584	\N	/api/restaurant	GET	2025-04-20 20:20:45.271621	2025-04-20 20:20:45.27171	200	0
585	\N	/api/restaurant	GET	2025-04-20 20:20:45.484155	2025-04-20 20:20:45.484221	200	0
586	\N	/api/washingmachines	GET	2025-04-20 20:21:09.304954	2025-04-20 20:21:10.21444	200	909
587	\N	/api/washingmachines	GET	2025-04-20 20:21:10.386884	2025-04-20 20:21:10.562751	200	175
588	\N	/api/restaurant	GET	2025-04-20 20:22:37.189288	2025-04-20 20:22:37.189357	200	0
589	\N	/api/restaurant	GET	2025-04-20 20:22:37.363328	2025-04-20 20:22:37.36338	200	0
590	test@imt-atlantique.net	/api/upload	POST	2025-04-20 20:23:18.924646	2025-04-20 20:23:18.927835	200	3
591	\N	/status	GET	2025-04-20 20:23:31.510987	2025-04-20 20:23:31.511001	200	0
592	\N	/status	GET	2025-04-20 20:23:31.716511	2025-04-20 20:23:31.716512	200	0
593	\N	/api/statistics/global	GET	2025-04-20 20:23:31.85355	2025-04-20 20:23:31.854733	200	1
594	\N	/api/statistics/global	GET	2025-04-20 20:23:32.020751	2025-04-20 20:23:32.048069	200	27
595	\N	/api/statistics/endpoints	GET	2025-04-20 20:23:32.038518	2025-04-20 20:23:32.048288	200	9
596	\N	/api/statistics/endpoints	GET	2025-04-20 20:23:32.215061	2025-04-20 20:23:32.216882	200	1
597	\N	/status	GET	2025-04-20 20:23:36.49354	2025-04-20 20:23:36.493544	200	0
598	\N	/status	GET	2025-04-20 20:23:41.495137	2025-04-20 20:23:41.495139	200	0
599	\N	/status	GET	2025-04-20 20:23:46.494532	2025-04-20 20:23:46.494534	200	0
600	\N	/status	GET	2025-04-20 20:23:51.505889	2025-04-20 20:23:51.50589	200	0
601	\N	/status	GET	2025-04-20 20:23:56.495753	2025-04-20 20:23:56.495754	200	0
602	\N	/status	GET	2025-04-20 20:24:01.495379	2025-04-20 20:24:01.495381	200	0
603	\N	/status	GET	2025-04-20 20:24:06.493433	2025-04-20 20:24:06.493435	200	0
604	\N	/status	GET	2025-04-20 20:24:11.495258	2025-04-20 20:24:11.49526	200	0
605	\N	/status	GET	2025-04-20 20:24:16.513857	2025-04-20 20:24:16.513859	200	0
606	\N	/status	GET	2025-04-20 20:24:21.49286	2025-04-20 20:24:21.492861	200	0
607	\N	/status	GET	2025-04-20 20:24:26.493199	2025-04-20 20:24:26.493201	200	0
608	\N	/status	GET	2025-04-20 20:24:31.498686	2025-04-20 20:24:31.498688	200	0
609	\N	/status	GET	2025-04-20 20:24:36.497549	2025-04-20 20:24:36.49755	200	0
610	\N	/status	GET	2025-04-20 20:24:41.495673	2025-04-20 20:24:41.495675	200	0
611	\N	/status	GET	2025-04-20 20:24:46.497139	2025-04-20 20:24:46.49714	200	0
612	\N	/status	GET	2025-04-20 20:24:51.495536	2025-04-20 20:24:51.495538	200	0
613	\N	/status	GET	2025-04-20 20:24:56.493153	2025-04-20 20:24:56.493154	200	0
614	\N	/status	GET	2025-04-20 20:25:01.496391	2025-04-20 20:25:01.496393	200	0
615	\N	/status	GET	2025-04-20 20:25:06.497217	2025-04-20 20:25:06.497218	200	0
616	\N	/status	GET	2025-04-20 20:25:11.497863	2025-04-20 20:25:11.497864	200	0
617	\N	/status	GET	2025-04-20 20:25:16.496469	2025-04-20 20:25:16.496471	200	0
618	\N	/status	GET	2025-04-20 20:25:21.504166	2025-04-20 20:25:21.504168	200	0
619	\N	/status	GET	2025-04-20 20:25:26.498253	2025-04-20 20:25:26.498254	200	0
620	\N	/status	GET	2025-04-20 20:25:31.495696	2025-04-20 20:25:31.495697	200	0
621	\N	/status	GET	2025-04-20 20:25:36.492924	2025-04-20 20:25:36.492926	200	0
622	\N	/status	GET	2025-04-20 20:25:41.524225	2025-04-20 20:25:41.524226	200	0
623	\N	/status	GET	2025-04-20 20:25:46.4925	2025-04-20 20:25:46.492501	200	0
624	\N	/status	GET	2025-04-20 20:25:51.511357	2025-04-20 20:25:51.511358	200	0
625	\N	/status	GET	2025-04-20 20:25:56.497341	2025-04-20 20:25:56.497343	200	0
626	\N	/status	GET	2025-04-20 20:26:01.500125	2025-04-20 20:26:01.500127	200	0
627	\N	/status	GET	2025-04-20 20:26:06.511272	2025-04-20 20:26:06.511273	200	0
628	\N	/status	GET	2025-04-20 20:26:11.496173	2025-04-20 20:26:11.496174	200	0
629	\N	/status	GET	2025-04-20 20:26:16.500155	2025-04-20 20:26:16.500157	200	0
630	\N	/status	GET	2025-04-20 20:26:21.49808	2025-04-20 20:26:21.498082	200	0
631	\N	/status	GET	2025-04-20 20:26:26.492366	2025-04-20 20:26:26.492368	200	0
632	\N	/status	GET	2025-04-20 20:26:31.500402	2025-04-20 20:26:31.500404	200	0
633	\N	/status	GET	2025-04-20 20:26:36.494393	2025-04-20 20:26:36.494395	200	0
634	\N	/status	GET	2025-04-20 20:26:41.493267	2025-04-20 20:26:41.493268	200	0
635	\N	/status	GET	2025-04-20 20:26:46.494222	2025-04-20 20:26:46.494223	200	0
636	\N	/status	GET	2025-04-20 20:26:51.499182	2025-04-20 20:26:51.499184	200	0
637	\N	/status	GET	2025-04-20 20:26:56.497821	2025-04-20 20:26:56.497822	200	0
638	\N	/status	GET	2025-04-20 20:27:01.496442	2025-04-20 20:27:01.496443	200	0
639	\N	/status	GET	2025-04-20 20:27:06.493808	2025-04-20 20:27:06.493809	200	0
640	\N	/status	GET	2025-04-20 20:27:11.493252	2025-04-20 20:27:11.493254	200	0
641	\N	/status	GET	2025-04-20 20:27:16.494384	2025-04-20 20:27:16.494385	200	0
642	\N	/status	GET	2025-04-20 20:27:21.495824	2025-04-20 20:27:21.495825	200	0
643	\N	/status	GET	2025-04-20 20:27:26.496518	2025-04-20 20:27:26.496519	200	0
644	\N	/status	GET	2025-04-20 20:27:31.496009	2025-04-20 20:27:31.49601	200	0
645	\N	/status	GET	2025-04-20 20:27:36.493438	2025-04-20 20:27:36.493439	200	0
646	\N	/status	GET	2025-04-20 20:27:41.497723	2025-04-20 20:27:41.497724	200	0
647	\N	/status	GET	2025-04-20 20:27:46.496413	2025-04-20 20:27:46.496415	200	0
648	\N	/status	GET	2025-04-20 20:27:51.500434	2025-04-20 20:27:51.500435	200	0
649	\N	/status	GET	2025-04-20 20:27:56.495223	2025-04-20 20:27:56.495225	200	0
650	\N	/status	GET	2025-04-20 20:28:01.496058	2025-04-20 20:28:01.49606	200	0
651	\N	/status	GET	2025-04-20 20:28:06.503551	2025-04-20 20:28:06.503552	200	0
652	\N	/status	GET	2025-04-20 20:28:11.494507	2025-04-20 20:28:11.494508	200	0
653	\N	/status	GET	2025-04-20 20:28:16.493465	2025-04-20 20:28:16.493467	200	0
654	\N	/status	GET	2025-04-20 20:28:21.497106	2025-04-20 20:28:21.497108	200	0
655	\N	/status	GET	2025-04-20 20:28:26.507658	2025-04-20 20:28:26.50766	200	0
656	\N	/status	GET	2025-04-20 20:28:31.50032	2025-04-20 20:28:31.500321	200	0
657	\N	/status	GET	2025-04-20 20:28:36.495507	2025-04-20 20:28:36.495508	200	0
658	\N	/status	GET	2025-04-20 20:28:41.497833	2025-04-20 20:28:41.497835	200	0
659	\N	/status	GET	2025-04-20 20:28:46.493829	2025-04-20 20:28:46.493831	200	0
660	\N	/status	GET	2025-04-20 20:28:51.496555	2025-04-20 20:28:51.496557	200	0
661	\N	/status	GET	2025-04-20 20:28:56.515337	2025-04-20 20:28:56.515338	200	0
662	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:28:56.910916	2025-04-20 20:28:56.911013	304	0
663	\N	/status	GET	2025-04-20 20:28:57.032144	2025-04-20 20:28:57.032146	200	0
664	\N	/api/statistics/global	GET	2025-04-20 20:28:57.206578	2025-04-20 20:28:57.208193	200	1
665	\N	/api/statistics/endpoints	GET	2025-04-20 20:28:57.375317	2025-04-20 20:28:57.376743	200	1
666	\N	/status	GET	2025-04-20 20:29:02.028732	2025-04-20 20:29:02.028733	200	0
667	\N	/status	GET	2025-04-20 20:29:07.026032	2025-04-20 20:29:07.026033	200	0
668	\N	/status	GET	2025-04-20 20:29:12.036756	2025-04-20 20:29:12.036757	200	0
669	\N	/status	GET	2025-04-20 20:29:17.028842	2025-04-20 20:29:17.028843	200	0
670	\N	/status	GET	2025-04-20 20:29:22.027656	2025-04-20 20:29:22.027657	200	0
671	\N	/status	GET	2025-04-20 20:29:27.030613	2025-04-20 20:29:27.030615	200	0
672	\N	/status	GET	2025-04-20 20:29:32.024408	2025-04-20 20:29:32.02441	200	0
673	\N	/status	GET	2025-04-20 20:29:37.024607	2025-04-20 20:29:37.024609	200	0
674	\N	/status	GET	2025-04-20 20:29:42.025974	2025-04-20 20:29:42.025975	200	0
675	\N	/status	GET	2025-04-20 20:29:47.027286	2025-04-20 20:29:47.027288	200	0
676	\N	/status	GET	2025-04-20 20:29:52.032454	2025-04-20 20:29:52.032455	200	0
677	\N	/status	GET	2025-04-20 20:29:57.028326	2025-04-20 20:29:57.028328	200	0
678	\N	/status	GET	2025-04-20 20:30:00.807393	2025-04-20 20:30:00.807394	200	0
679	\N	/api/statistics/global	GET	2025-04-20 20:30:00.993713	2025-04-20 20:30:00.994973	200	1
680	\N	/api/statistics/endpoints	GET	2025-04-20 20:30:01.162882	2025-04-20 20:30:01.164665	200	1
681	\N	/status	GET	2025-04-20 20:30:05.810626	2025-04-20 20:30:05.810628	200	0
682	\N	/status	GET	2025-04-20 20:30:10.817822	2025-04-20 20:30:10.817824	200	0
683	\N	/status	GET	2025-04-20 20:30:15.808504	2025-04-20 20:30:15.808506	200	0
684	\N	/status	GET	2025-04-20 20:30:20.498011	2025-04-20 20:30:20.498013	200	0
685	\N	/api/statistics/global	GET	2025-04-20 20:30:20.665831	2025-04-20 20:30:20.666932	200	1
686	\N	/status	GET	2025-04-20 20:30:20.80634	2025-04-20 20:30:20.806342	200	0
687	\N	/api/statistics/endpoints	GET	2025-04-20 20:30:20.832346	2025-04-20 20:30:20.838785	200	6
688	\N	/status	GET	2025-04-20 20:30:25.814365	2025-04-20 20:30:25.814367	200	0
689	\N	/status	GET	2025-04-20 20:30:30.814738	2025-04-20 20:30:30.81474	200	0
690	\N	/status	GET	2025-04-20 20:30:35.81069	2025-04-20 20:30:35.810692	200	0
691	\N	/status	GET	2025-04-20 20:30:40.810516	2025-04-20 20:30:40.810517	200	0
692	\N	/status	GET	2025-04-20 20:30:44.740778	2025-04-20 20:30:44.74078	200	0
693	\N	/api/statistics/global	GET	2025-04-20 20:30:44.984002	2025-04-20 20:30:44.985209	200	1
694	\N	/api/statistics/endpoints	GET	2025-04-20 20:30:45.154186	2025-04-20 20:30:45.161053	200	6
695	\N	/status	GET	2025-04-20 20:30:49.745638	2025-04-20 20:30:49.74564	200	0
696	\N	/status	GET	2025-04-20 20:30:54.744124	2025-04-20 20:30:54.744126	200	0
697	\N	/status	GET	2025-04-20 20:30:59.749727	2025-04-20 20:30:59.74973	200	0
698	\N	/status	GET	2025-04-20 20:31:04.743827	2025-04-20 20:31:04.743829	200	0
699	\N	/status	GET	2025-04-20 20:31:09.746413	2025-04-20 20:31:09.746414	200	0
700	\N	/status	GET	2025-04-20 20:31:14.741145	2025-04-20 20:31:14.741146	200	0
701	\N	/status	GET	2025-04-20 20:31:19.745272	2025-04-20 20:31:19.745273	200	0
702	\N	/status	GET	2025-04-20 20:31:24.744018	2025-04-20 20:31:24.74402	200	0
703	\N	/status	GET	2025-04-20 20:31:29.742183	2025-04-20 20:31:29.742185	200	0
704	\N	/status	GET	2025-04-20 20:31:34.748029	2025-04-20 20:31:34.748031	200	0
705	\N	/status	GET	2025-04-20 20:31:39.74286	2025-04-20 20:31:39.742862	200	0
706	\N	/status	GET	2025-04-20 20:31:44.745944	2025-04-20 20:31:44.745945	200	0
707	\N	/status	GET	2025-04-20 20:31:49.749621	2025-04-20 20:31:49.749623	200	0
708	\N	/status	GET	2025-04-20 20:31:54.74783	2025-04-20 20:31:54.747832	200	0
709	\N	/status	GET	2025-04-20 20:31:59.741219	2025-04-20 20:31:59.741221	200	0
710	\N	/status	GET	2025-04-20 20:32:04.754787	2025-04-20 20:32:04.754788	200	0
711	\N	/status	GET	2025-04-20 20:32:09.743596	2025-04-20 20:32:09.743597	200	0
712	\N	/status	GET	2025-04-20 20:32:14.745694	2025-04-20 20:32:14.745696	200	0
713	\N	/status	GET	2025-04-20 20:32:19.747401	2025-04-20 20:32:19.747402	200	0
714	\N	/status	GET	2025-04-20 20:32:24.743591	2025-04-20 20:32:24.743593	200	0
715	\N	/status	GET	2025-04-20 20:32:29.753919	2025-04-20 20:32:29.75392	200	0
716	\N	/status	GET	2025-04-20 20:32:34.748043	2025-04-20 20:32:34.748045	200	0
717	\N	/status	GET	2025-04-20 20:32:39.748292	2025-04-20 20:32:39.748294	200	0
718	\N	/status	GET	2025-04-20 20:32:44.745725	2025-04-20 20:32:44.745727	200	0
719	\N	/status	GET	2025-04-20 20:32:49.741221	2025-04-20 20:32:49.741222	200	0
720	\N	/status	GET	2025-04-20 20:32:54.750029	2025-04-20 20:32:54.750031	200	0
721	\N	/status	GET	2025-04-20 20:32:59.744756	2025-04-20 20:32:59.744758	200	0
722	\N	/status	GET	2025-04-20 20:33:04.741991	2025-04-20 20:33:04.741993	200	0
723	\N	/status	GET	2025-04-20 20:33:09.741906	2025-04-20 20:33:09.741908	200	0
724	\N	/status	GET	2025-04-20 20:33:14.745522	2025-04-20 20:33:14.745524	200	0
725	\N	/status	GET	2025-04-20 20:33:19.748183	2025-04-20 20:33:19.748184	200	0
726	\N	/status	GET	2025-04-20 20:33:24.74192	2025-04-20 20:33:24.741922	200	0
727	\N	/status	GET	2025-04-20 20:33:29.743237	2025-04-20 20:33:29.743239	200	0
728	\N	/status	GET	2025-04-20 20:33:34.744283	2025-04-20 20:33:34.744285	200	0
729	\N	/status	GET	2025-04-20 20:33:39.744925	2025-04-20 20:33:39.744926	200	0
730	\N	/status	GET	2025-04-20 20:33:44.744983	2025-04-20 20:33:44.744984	200	0
731	\N	/status	GET	2025-04-20 20:33:49.746837	2025-04-20 20:33:49.746838	200	0
732	\N	/status	GET	2025-04-20 20:33:54.744677	2025-04-20 20:33:54.744678	200	0
733	\N	/status	GET	2025-04-20 20:33:59.751527	2025-04-20 20:33:59.751529	200	0
734	\N	/status	GET	2025-04-20 20:34:04.74901	2025-04-20 20:34:04.749012	200	0
735	\N	/status	GET	2025-04-20 20:34:09.744651	2025-04-20 20:34:09.744653	200	0
736	\N	/status	GET	2025-04-20 20:34:14.754062	2025-04-20 20:34:14.754063	200	0
737	\N	/status	GET	2025-04-20 20:34:19.751067	2025-04-20 20:34:19.751069	200	0
738	\N	/status	GET	2025-04-20 20:34:24.746967	2025-04-20 20:34:24.746969	200	0
739	\N	/status	GET	2025-04-20 20:34:29.752311	2025-04-20 20:34:29.752312	200	0
740	\N	/status	GET	2025-04-20 20:34:34.747415	2025-04-20 20:34:34.747417	200	0
741	\N	/status	GET	2025-04-20 20:34:39.748064	2025-04-20 20:34:39.748065	200	0
742	\N	/status	GET	2025-04-20 20:34:44.745327	2025-04-20 20:34:44.745329	200	0
743	\N	/status	GET	2025-04-20 20:34:49.755863	2025-04-20 20:34:49.755864	200	0
744	\N	/status	GET	2025-04-20 20:34:54.7513	2025-04-20 20:34:54.751301	200	0
745	\N	/status	GET	2025-04-20 20:34:59.745457	2025-04-20 20:34:59.745458	200	0
746	\N	/status	GET	2025-04-20 20:35:04.742911	2025-04-20 20:35:04.742913	200	0
747	\N	/status	GET	2025-04-20 20:35:09.97956	2025-04-20 20:35:09.979587	200	0
748	\N	/status	GET	2025-04-20 20:35:14.994334	2025-04-20 20:35:14.994336	200	0
749	\N	/status	GET	2025-04-20 20:35:20.078453	2025-04-20 20:35:20.078456	200	0
750	\N	/status	GET	2025-04-20 20:35:24.902291	2025-04-20 20:35:24.902293	200	0
751	\N	/status	GET	2025-04-20 20:35:30.121395	2025-04-20 20:35:30.121398	200	0
752	\N	/status	GET	2025-04-20 20:35:34.980678	2025-04-20 20:35:34.980681	200	0
753	\N	/status	GET	2025-04-20 20:35:39.829586	2025-04-20 20:35:39.829587	200	0
754	\N	/status	GET	2025-04-20 20:35:45.091271	2025-04-20 20:35:45.091274	200	0
755	\N	/status	GET	2025-04-20 20:35:49.828022	2025-04-20 20:35:49.828024	200	0
756	\N	/status	GET	2025-04-20 20:35:54.841156	2025-04-20 20:35:54.841158	200	0
757	\N	/status	GET	2025-04-20 20:35:59.827198	2025-04-20 20:35:59.8272	200	0
758	\N	/status	GET	2025-04-20 20:36:04.951728	2025-04-20 20:36:04.951749	200	0
759	\N	/status	GET	2025-04-20 20:36:09.868339	2025-04-20 20:36:09.86834	200	0
760	\N	/status	GET	2025-04-20 20:36:14.835508	2025-04-20 20:36:14.83551	200	0
761	\N	/status	GET	2025-04-20 20:36:19.849186	2025-04-20 20:36:19.849188	200	0
762	\N	/status	GET	2025-04-20 20:36:24.827481	2025-04-20 20:36:24.827482	200	0
763	\N	/status	GET	2025-04-20 20:36:29.844882	2025-04-20 20:36:29.844883	200	0
764	\N	/status	GET	2025-04-20 20:36:34.863235	2025-04-20 20:36:34.863236	200	0
765	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:36:36.876143	2025-04-20 20:36:36.87625	304	0
766	\N	/status	GET	2025-04-20 20:36:37.134791	2025-04-20 20:36:37.134792	200	0
767	\N	/status	GET	2025-04-20 20:36:37.306526	2025-04-20 20:36:37.306527	200	0
768	\N	/api/statistics/global	GET	2025-04-20 20:36:37.4486	2025-04-20 20:36:37.44976	200	1
769	\N	/api/statistics/global	GET	2025-04-20 20:36:37.625802	2025-04-20 20:36:37.627098	200	1
770	\N	/api/statistics/endpoints	GET	2025-04-20 20:36:37.64419	2025-04-20 20:36:37.649977	200	5
771	\N	/api/statistics/endpoints	GET	2025-04-20 20:36:37.825011	2025-04-20 20:36:37.826975	200	1
772	\N	/status	GET	2025-04-20 20:36:42.146675	2025-04-20 20:36:42.146678	200	0
773	\N	/status	GET	2025-04-20 20:36:47.236168	2025-04-20 20:36:47.236169	200	0
774	\N	/status	GET	2025-04-20 20:36:52.138732	2025-04-20 20:36:52.138734	200	0
775	\N	/status	GET	2025-04-20 20:36:57.224839	2025-04-20 20:36:57.224841	200	0
776	\N	/status	GET	2025-04-20 20:37:02.140992	2025-04-20 20:37:02.140994	200	0
777	\N	/status	GET	2025-04-20 20:37:07.225102	2025-04-20 20:37:07.225103	200	0
778	\N	/status	GET	2025-04-20 20:37:12.162988	2025-04-20 20:37:12.16299	200	0
779	\N	/status	GET	2025-04-20 20:37:17.253626	2025-04-20 20:37:17.253627	200	0
780	\N	/status	GET	2025-04-20 20:37:22.141477	2025-04-20 20:37:22.141478	200	0
781	\N	/status	GET	2025-04-20 20:37:27.214854	2025-04-20 20:37:27.214855	200	0
782	\N	/status	GET	2025-04-20 20:37:32.155376	2025-04-20 20:37:32.155377	200	0
783	\N	/status	GET	2025-04-20 20:37:37.308022	2025-04-20 20:37:37.308023	200	0
784	\N	/status	GET	2025-04-20 20:37:42.140415	2025-04-20 20:37:42.140417	200	0
785	\N	/status	GET	2025-04-20 20:37:47.137018	2025-04-20 20:37:47.13702	200	0
786	\N	/status	GET	2025-04-20 20:37:52.256255	2025-04-20 20:37:52.256257	200	0
787	\N	/status	GET	2025-04-20 20:37:57.142879	2025-04-20 20:37:57.14288	200	0
788	\N	/status	GET	2025-04-20 20:38:02.231063	2025-04-20 20:38:02.231064	200	0
789	\N	/status	GET	2025-04-20 20:38:07.145838	2025-04-20 20:38:07.14584	200	0
790	\N	/status	GET	2025-04-20 20:38:12.232399	2025-04-20 20:38:12.2324	200	0
791	\N	/status	GET	2025-04-20 20:38:17.138851	2025-04-20 20:38:17.138853	200	0
792	\N	/status	GET	2025-04-20 20:38:22.225911	2025-04-20 20:38:22.225912	200	0
793	\N	/status	GET	2025-04-20 20:38:27.14092	2025-04-20 20:38:27.140922	200	0
794	\N	/status	GET	2025-04-20 20:38:32.144216	2025-04-20 20:38:32.144217	200	0
795	\N	/status	GET	2025-04-20 20:38:37.0418	2025-04-20 20:38:37.041802	200	0
796	\N	/api/statistics/global	GET	2025-04-20 20:38:37.216139	2025-04-20 20:38:37.217487	200	1
797	\N	/api/statistics/endpoints	GET	2025-04-20 20:38:37.397002	2025-04-20 20:38:37.398583	200	1
798	\N	/status	GET	2025-04-20 20:38:41.974782	2025-04-20 20:38:41.974783	200	0
799	\N	/status	GET	2025-04-20 20:38:42.142324	2025-04-20 20:38:42.142325	200	0
800	\N	/api/statistics/global	GET	2025-04-20 20:38:42.16006	2025-04-20 20:38:42.163441	200	3
801	\N	/api/statistics/endpoints	GET	2025-04-20 20:38:42.339574	2025-04-20 20:38:42.366248	200	26
802	\N	/status	GET	2025-04-20 20:38:46.963449	2025-04-20 20:38:46.96345	200	0
803	\N	/status	GET	2025-04-20 20:38:47.13213	2025-04-20 20:38:47.132132	200	0
804	\N	/api/statistics/global	GET	2025-04-20 20:38:47.150052	2025-04-20 20:38:47.168443	200	18
805	\N	/api/statistics/endpoints	GET	2025-04-20 20:38:47.474191	2025-04-20 20:38:47.476139	200	1
806	\N	/status	GET	2025-04-20 20:38:51.960613	2025-04-20 20:38:51.960615	200	0
807	\N	/status	GET	2025-04-20 20:38:52.406332	2025-04-20 20:38:52.406334	200	0
808	\N	/api/statistics/global	GET	2025-04-20 20:38:52.426749	2025-04-20 20:38:52.429673	200	2
809	\N	/api/statistics/endpoints	GET	2025-04-20 20:38:52.606356	2025-04-20 20:38:52.608282	200	1
810	\N	/status	GET	2025-04-20 20:38:56.970572	2025-04-20 20:38:56.970574	200	0
811	\N	/status	GET	2025-04-20 20:38:57.35722	2025-04-20 20:38:57.357222	200	0
812	\N	/api/statistics/global	GET	2025-04-20 20:38:57.375041	2025-04-20 20:38:57.376753	200	1
813	\N	/api/statistics/endpoints	GET	2025-04-20 20:38:57.557899	2025-04-20 20:38:57.590603	200	32
814	\N	/status	GET	2025-04-20 20:39:01.962321	2025-04-20 20:39:01.962322	200	0
815	\N	/status	GET	2025-04-20 20:39:02.412105	2025-04-20 20:39:02.412106	200	0
816	\N	/api/statistics/global	GET	2025-04-20 20:39:02.430034	2025-04-20 20:39:02.431388	200	1
817	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:02.605388	2025-04-20 20:39:02.607146	200	1
818	\N	/status	GET	2025-04-20 20:39:06.955069	2025-04-20 20:39:06.955071	200	0
819	\N	/status	GET	2025-04-20 20:39:07.409837	2025-04-20 20:39:07.409839	200	0
820	\N	/api/statistics/global	GET	2025-04-20 20:39:07.421926	2025-04-20 20:39:07.42308	200	1
821	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:07.614387	2025-04-20 20:39:07.615945	200	1
822	\N	/status	GET	2025-04-20 20:39:11.956146	2025-04-20 20:39:11.956148	200	0
823	\N	/status	GET	2025-04-20 20:39:12.350775	2025-04-20 20:39:12.350777	200	0
824	\N	/api/statistics/global	GET	2025-04-20 20:39:12.368564	2025-04-20 20:39:12.369618	200	1
825	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:12.552145	2025-04-20 20:39:12.554103	200	1
826	\N	/status	GET	2025-04-20 20:39:16.958662	2025-04-20 20:39:16.958663	200	0
827	\N	/status	GET	2025-04-20 20:39:17.255201	2025-04-20 20:39:17.255202	200	0
828	\N	/api/statistics/global	GET	2025-04-20 20:39:17.273016	2025-04-20 20:39:17.274088	200	1
829	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:17.567214	2025-04-20 20:39:17.569461	200	2
830	\N	/status	GET	2025-04-20 20:39:21.959741	2025-04-20 20:39:21.959742	200	0
831	\N	/status	GET	2025-04-20 20:39:22.251414	2025-04-20 20:39:22.251415	200	0
832	\N	/api/statistics/global	GET	2025-04-20 20:39:22.269323	2025-04-20 20:39:22.270897	200	1
833	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:22.581908	2025-04-20 20:39:22.585432	200	3
834	\N	/status	GET	2025-04-20 20:39:26.960358	2025-04-20 20:39:26.96036	200	0
835	\N	/status	GET	2025-04-20 20:39:27.133918	2025-04-20 20:39:27.13392	200	0
836	\N	/api/statistics/global	GET	2025-04-20 20:39:27.148077	2025-04-20 20:39:27.1496	200	1
837	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:27.323466	2025-04-20 20:39:27.327162	200	3
838	\N	/status	GET	2025-04-20 20:39:31.959467	2025-04-20 20:39:31.959468	200	0
839	\N	/status	GET	2025-04-20 20:39:32.127594	2025-04-20 20:39:32.127596	200	0
840	\N	/api/statistics/global	GET	2025-04-20 20:39:32.151931	2025-04-20 20:39:32.15856	200	6
841	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:32.338822	2025-04-20 20:39:32.346557	200	7
842	\N	/status	GET	2025-04-20 20:39:36.958192	2025-04-20 20:39:36.958194	200	0
843	\N	/status	GET	2025-04-20 20:39:37.136963	2025-04-20 20:39:37.136965	200	0
844	\N	/api/statistics/global	GET	2025-04-20 20:39:37.154877	2025-04-20 20:39:37.155925	200	1
845	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:37.330493	2025-04-20 20:39:37.333795	200	3
846	\N	/status	GET	2025-04-20 20:39:41.960305	2025-04-20 20:39:41.960306	200	0
847	\N	/status	GET	2025-04-20 20:39:42.139061	2025-04-20 20:39:42.139062	200	0
848	\N	/api/statistics/global	GET	2025-04-20 20:39:42.156822	2025-04-20 20:39:42.165095	200	8
849	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:42.348307	2025-04-20 20:39:42.361728	200	13
850	\N	/status	GET	2025-04-20 20:39:46.964497	2025-04-20 20:39:46.964499	200	0
851	\N	/status	GET	2025-04-20 20:39:47.13521	2025-04-20 20:39:47.135211	200	0
852	\N	/api/statistics/global	GET	2025-04-20 20:39:47.155045	2025-04-20 20:39:47.157233	200	2
853	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:47.343992	2025-04-20 20:39:47.347765	200	3
854	\N	/status	GET	2025-04-20 20:39:51.964571	2025-04-20 20:39:51.964573	200	0
855	\N	/status	GET	2025-04-20 20:39:52.137875	2025-04-20 20:39:52.137877	200	0
856	\N	/api/statistics/global	GET	2025-04-20 20:39:52.149447	2025-04-20 20:39:52.150475	200	1
857	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:52.325054	2025-04-20 20:39:52.336348	200	11
858	\N	/status	GET	2025-04-20 20:39:56.990245	2025-04-20 20:39:56.990246	200	0
859	\N	/status	GET	2025-04-20 20:39:57.156799	2025-04-20 20:39:57.1568	200	0
860	\N	/api/statistics/global	GET	2025-04-20 20:39:57.174542	2025-04-20 20:39:57.17541	200	0
861	\N	/api/statistics/endpoints	GET	2025-04-20 20:39:57.355412	2025-04-20 20:39:57.357333	200	1
862	\N	/status	GET	2025-04-20 20:40:01.960534	2025-04-20 20:40:01.960535	200	0
863	\N	/status	GET	2025-04-20 20:40:02.127264	2025-04-20 20:40:02.127266	200	0
864	\N	/api/statistics/global	GET	2025-04-20 20:40:02.145108	2025-04-20 20:40:02.147832	200	2
865	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:02.335231	2025-04-20 20:40:02.338759	200	3
866	\N	/status	GET	2025-04-20 20:40:06.959276	2025-04-20 20:40:06.959278	200	0
867	\N	/status	GET	2025-04-20 20:40:07.125296	2025-04-20 20:40:07.125297	200	0
868	\N	/api/statistics/global	GET	2025-04-20 20:40:07.143291	2025-04-20 20:40:07.144591	200	1
869	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:07.329081	2025-04-20 20:40:07.332873	200	3
870	\N	/status	GET	2025-04-20 20:40:11.960855	2025-04-20 20:40:11.960857	200	0
871	\N	/status	GET	2025-04-20 20:40:12.127049	2025-04-20 20:40:12.127054	200	0
872	\N	/api/statistics/global	GET	2025-04-20 20:40:12.144995	2025-04-20 20:40:12.146127	200	1
873	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:12.327406	2025-04-20 20:40:12.334264	200	6
874	\N	/status	GET	2025-04-20 20:40:16.961127	2025-04-20 20:40:16.961129	200	0
875	\N	/status	GET	2025-04-20 20:40:17.128113	2025-04-20 20:40:17.128114	200	0
876	\N	/api/statistics/global	GET	2025-04-20 20:40:17.14582	2025-04-20 20:40:17.148351	200	2
877	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:17.329261	2025-04-20 20:40:17.332786	200	3
878	\N	/status	GET	2025-04-20 20:40:21.961647	2025-04-20 20:40:21.961649	200	0
879	\N	/status	GET	2025-04-20 20:40:22.13017	2025-04-20 20:40:22.130172	200	0
880	\N	/api/statistics/global	GET	2025-04-20 20:40:22.147888	2025-04-20 20:40:22.148967	200	1
881	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:22.317733	2025-04-20 20:40:22.327632	200	9
882	\N	/status	GET	2025-04-20 20:40:26.960055	2025-04-20 20:40:26.960057	200	0
883	\N	/status	GET	2025-04-20 20:40:27.125906	2025-04-20 20:40:27.125907	200	0
884	\N	/api/statistics/global	GET	2025-04-20 20:40:27.143651	2025-04-20 20:40:27.144803	200	1
885	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:27.320692	2025-04-20 20:40:27.322484	200	1
886	\N	/status	GET	2025-04-20 20:40:31.96436	2025-04-20 20:40:31.964362	200	0
887	\N	/status	GET	2025-04-20 20:40:32.132409	2025-04-20 20:40:32.13241	200	0
903	\N	/status	GET	2025-04-20 20:40:52.130938	2025-04-20 20:40:52.130939	200	0
924	\N	/api/statistics/global	GET	2025-04-20 20:41:17.143511	2025-04-20 20:41:17.146291	200	2
925	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:17.336033	2025-04-20 20:41:17.349229	200	13
1242	\N	/api/statistics/global	GET	2025-04-20 20:47:46.606462	2025-04-20 20:47:46.632004	200	25
1243	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:46.810634	2025-04-20 20:47:46.812789	200	2
1244	\N	/status	GET	2025-04-20 20:47:51.413559	2025-04-20 20:47:51.41356	200	0
1245	\N	/status	GET	2025-04-20 20:47:51.581467	2025-04-20 20:47:51.581469	200	0
1246	\N	/api/statistics/global	GET	2025-04-20 20:47:51.599958	2025-04-20 20:47:51.60094	200	0
1247	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:51.787734	2025-04-20 20:47:51.78947	200	1
1248	\N	/status	GET	2025-04-20 20:47:56.478351	2025-04-20 20:47:56.478353	200	0
1249	\N	/status	GET	2025-04-20 20:47:56.66651	2025-04-20 20:47:56.666512	200	0
1250	\N	/api/statistics/global	GET	2025-04-20 20:47:56.684893	2025-04-20 20:47:56.686084	200	1
1251	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:56.864385	2025-04-20 20:47:56.866649	200	2
1252	\N	/status	GET	2025-04-20 20:48:01.417352	2025-04-20 20:48:01.417354	200	0
1253	\N	/status	GET	2025-04-20 20:48:01.589354	2025-04-20 20:48:01.589355	200	0
1254	\N	/api/statistics/global	GET	2025-04-20 20:48:01.607813	2025-04-20 20:48:01.60905	200	1
1255	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:01.788238	2025-04-20 20:48:01.790512	200	2
1256	\N	/status	GET	2025-04-20 20:48:06.434009	2025-04-20 20:48:06.434011	200	0
1257	\N	/status	GET	2025-04-20 20:48:06.600498	2025-04-20 20:48:06.600499	200	0
1258	\N	/api/statistics/global	GET	2025-04-20 20:48:06.618952	2025-04-20 20:48:06.619933	200	0
1259	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:06.790287	2025-04-20 20:48:06.792245	200	1
1260	\N	/status	GET	2025-04-20 20:48:11.429263	2025-04-20 20:48:11.429265	200	0
1261	\N	/status	GET	2025-04-20 20:48:11.598842	2025-04-20 20:48:11.598843	200	0
1262	\N	/api/statistics/global	GET	2025-04-20 20:48:11.616908	2025-04-20 20:48:11.618444	200	1
1263	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:11.79669	2025-04-20 20:48:11.798429	200	1
1264	\N	/status	GET	2025-04-20 20:48:16.394825	2025-04-20 20:48:16.394827	200	0
1265	\N	/status	GET	2025-04-20 20:48:16.564573	2025-04-20 20:48:16.564574	200	0
1266	\N	/api/statistics/global	GET	2025-04-20 20:48:16.584581	2025-04-20 20:48:16.586147	200	1
1267	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:16.760091	2025-04-20 20:48:16.76198	200	1
1268	\N	/status	GET	2025-04-20 20:48:21.392138	2025-04-20 20:48:21.39214	200	0
1269	\N	/status	GET	2025-04-20 20:48:21.558801	2025-04-20 20:48:21.558803	200	0
1270	\N	/api/statistics/global	GET	2025-04-20 20:48:21.576739	2025-04-20 20:48:21.578072	200	1
1271	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:21.747553	2025-04-20 20:48:21.765949	200	18
1464	\N	/status	GET	2025-04-20 20:52:26.401788	2025-04-20 20:52:26.401789	200	0
1465	\N	/status	GET	2025-04-20 20:52:26.570332	2025-04-20 20:52:26.570333	200	0
1466	\N	/api/statistics/global	GET	2025-04-20 20:52:26.588616	2025-04-20 20:52:26.590067	200	1
1467	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:26.771852	2025-04-20 20:52:26.773668	200	1
1468	\N	/status	GET	2025-04-20 20:52:31.402122	2025-04-20 20:52:31.402124	200	0
1469	\N	/status	GET	2025-04-20 20:52:31.569522	2025-04-20 20:52:31.569523	200	0
1470	\N	/api/statistics/global	GET	2025-04-20 20:52:31.587642	2025-04-20 20:52:31.588965	200	1
1471	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:31.764055	2025-04-20 20:52:31.767218	200	3
1472	\N	/status	GET	2025-04-20 20:52:36.401663	2025-04-20 20:52:36.401665	200	0
1473	\N	/status	GET	2025-04-20 20:52:36.571113	2025-04-20 20:52:36.571115	200	0
1474	\N	/api/statistics/global	GET	2025-04-20 20:52:36.589792	2025-04-20 20:52:36.591043	200	1
1475	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:36.768342	2025-04-20 20:52:36.770546	200	2
1476	\N	/status	GET	2025-04-20 20:52:41.386938	2025-04-20 20:52:41.38694	200	0
1477	\N	/status	GET	2025-04-20 20:52:41.554856	2025-04-20 20:52:41.554858	200	0
1478	\N	/api/statistics/global	GET	2025-04-20 20:52:41.573387	2025-04-20 20:52:41.628843	200	55
1479	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:41.802196	2025-04-20 20:52:41.804663	200	2
1480	\N	/status	GET	2025-04-20 20:52:46.403772	2025-04-20 20:52:46.403773	200	0
1481	\N	/status	GET	2025-04-20 20:52:46.573204	2025-04-20 20:52:46.573206	200	0
1482	\N	/api/statistics/global	GET	2025-04-20 20:52:46.592064	2025-04-20 20:52:46.593515	200	1
1483	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:46.769289	2025-04-20 20:52:46.771329	200	2
1484	\N	/status	GET	2025-04-20 20:52:51.405222	2025-04-20 20:52:51.405224	200	0
1485	\N	/status	GET	2025-04-20 20:52:51.572969	2025-04-20 20:52:51.572971	200	0
1486	\N	/api/statistics/global	GET	2025-04-20 20:52:51.590818	2025-04-20 20:52:51.592149	200	1
1487	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:51.768833	2025-04-20 20:52:51.770585	200	1
1488	\N	/status	GET	2025-04-20 20:52:56.397234	2025-04-20 20:52:56.397235	200	0
1489	\N	/status	GET	2025-04-20 20:52:56.566858	2025-04-20 20:52:56.566859	200	0
1490	\N	/api/statistics/global	GET	2025-04-20 20:52:56.583048	2025-04-20 20:52:56.584261	200	1
1491	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:56.755316	2025-04-20 20:52:56.75862	200	3
1492	\N	/status	GET	2025-04-20 20:53:01.420798	2025-04-20 20:53:01.420799	200	0
1493	\N	/status	GET	2025-04-20 20:53:01.594058	2025-04-20 20:53:01.594059	200	0
1494	\N	/api/statistics/global	GET	2025-04-20 20:53:01.610549	2025-04-20 20:53:01.611938	200	1
1495	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:01.78648	2025-04-20 20:53:01.788695	200	2
1496	\N	/status	GET	2025-04-20 20:53:06.412253	2025-04-20 20:53:06.412255	200	0
1497	\N	/status	GET	2025-04-20 20:53:06.589885	2025-04-20 20:53:06.589887	200	0
1498	\N	/api/statistics/global	GET	2025-04-20 20:53:06.597829	2025-04-20 20:53:06.599082	200	1
1499	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:06.77249	2025-04-20 20:53:06.775708	200	3
1500	\N	/status	GET	2025-04-20 20:53:11.427666	2025-04-20 20:53:11.427668	200	0
1501	\N	/status	GET	2025-04-20 20:53:11.602777	2025-04-20 20:53:11.602779	200	0
1502	\N	/api/statistics/global	GET	2025-04-20 20:53:11.615526	2025-04-20 20:53:11.616886	200	1
1503	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:11.795524	2025-04-20 20:53:11.797662	200	2
1504	\N	/status	GET	2025-04-20 20:53:16.413321	2025-04-20 20:53:16.413322	200	0
1505	\N	/status	GET	2025-04-20 20:53:16.589754	2025-04-20 20:53:16.589755	200	0
1506	\N	/api/statistics/global	GET	2025-04-20 20:53:16.60199	2025-04-20 20:53:16.603275	200	1
1507	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:16.775306	2025-04-20 20:53:16.777823	200	2
1508	\N	/status	GET	2025-04-20 20:53:21.405946	2025-04-20 20:53:21.405948	200	0
1509	\N	/status	GET	2025-04-20 20:53:21.573687	2025-04-20 20:53:21.573688	200	0
1510	\N	/api/statistics/global	GET	2025-04-20 20:53:21.591496	2025-04-20 20:53:21.592883	200	1
1511	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:21.763957	2025-04-20 20:53:21.765682	200	1
1512	\N	/status	GET	2025-04-20 20:53:26.410566	2025-04-20 20:53:26.410567	200	0
1513	\N	/status	GET	2025-04-20 20:53:26.583004	2025-04-20 20:53:26.583005	200	0
1514	\N	/api/statistics/global	GET	2025-04-20 20:53:26.605297	2025-04-20 20:53:26.606539	200	1
1515	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:26.786632	2025-04-20 20:53:26.788384	200	1
888	\N	/api/statistics/global	GET	2025-04-20 20:40:32.150259	2025-04-20 20:40:32.197265	200	47
889	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:32.435369	2025-04-20 20:40:32.437236	200	1
890	\N	/status	GET	2025-04-20 20:40:36.967801	2025-04-20 20:40:36.967802	200	0
891	\N	/status	GET	2025-04-20 20:40:37.140758	2025-04-20 20:40:37.140759	200	0
892	\N	/api/statistics/global	GET	2025-04-20 20:40:37.158544	2025-04-20 20:40:37.166563	200	8
893	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:37.360902	2025-04-20 20:40:37.362497	200	1
894	\N	/status	GET	2025-04-20 20:40:41.959861	2025-04-20 20:40:41.959863	200	0
895	\N	/status	GET	2025-04-20 20:40:42.152327	2025-04-20 20:40:42.152328	200	0
896	\N	/api/statistics/global	GET	2025-04-20 20:40:42.170126	2025-04-20 20:40:42.172219	200	2
897	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:42.344522	2025-04-20 20:40:42.346244	200	1
898	\N	/status	GET	2025-04-20 20:40:46.961561	2025-04-20 20:40:46.961563	200	0
899	\N	/status	GET	2025-04-20 20:40:47.13072	2025-04-20 20:40:47.130722	200	0
900	\N	/api/statistics/global	GET	2025-04-20 20:40:47.146896	2025-04-20 20:40:47.151258	200	4
901	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:47.335963	2025-04-20 20:40:47.34617	200	10
902	\N	/status	GET	2025-04-20 20:40:51.960582	2025-04-20 20:40:51.960584	200	0
904	\N	/api/statistics/global	GET	2025-04-20 20:40:52.130956	2025-04-20 20:40:52.13207	200	1
905	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:52.332513	2025-04-20 20:40:52.355343	200	22
906	\N	/status	GET	2025-04-20 20:40:56.966717	2025-04-20 20:40:56.966718	200	0
907	\N	/status	GET	2025-04-20 20:40:57.134133	2025-04-20 20:40:57.134134	200	0
908	\N	/api/statistics/global	GET	2025-04-20 20:40:57.152097	2025-04-20 20:40:57.154251	200	2
909	\N	/api/statistics/endpoints	GET	2025-04-20 20:40:57.333095	2025-04-20 20:40:57.335362	200	2
910	\N	/status	GET	2025-04-20 20:41:01.958532	2025-04-20 20:41:01.958534	200	0
911	\N	/status	GET	2025-04-20 20:41:02.126326	2025-04-20 20:41:02.126328	200	0
912	\N	/api/statistics/global	GET	2025-04-20 20:41:02.144024	2025-04-20 20:41:02.145737	200	1
913	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:02.333701	2025-04-20 20:41:02.344338	200	10
914	\N	/status	GET	2025-04-20 20:41:06.963406	2025-04-20 20:41:06.963407	200	0
915	\N	/status	GET	2025-04-20 20:41:07.130482	2025-04-20 20:41:07.130483	200	0
916	\N	/api/statistics/global	GET	2025-04-20 20:41:07.149061	2025-04-20 20:41:07.150686	200	1
917	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:07.353654	2025-04-20 20:41:07.35536	200	1
918	\N	/status	GET	2025-04-20 20:41:11.965002	2025-04-20 20:41:11.965004	200	0
919	\N	/status	GET	2025-04-20 20:41:12.132964	2025-04-20 20:41:12.132966	200	0
920	\N	/api/statistics/global	GET	2025-04-20 20:41:12.159491	2025-04-20 20:41:12.16598	200	6
921	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:12.365709	2025-04-20 20:41:12.367245	200	1
922	\N	/status	GET	2025-04-20 20:41:16.956134	2025-04-20 20:41:16.956135	200	0
923	\N	/status	GET	2025-04-20 20:41:17.125669	2025-04-20 20:41:17.125671	200	0
926	\N	/status	GET	2025-04-20 20:41:21.966904	2025-04-20 20:41:21.966906	200	0
927	\N	/status	GET	2025-04-20 20:41:22.140928	2025-04-20 20:41:22.140929	200	0
928	\N	/api/statistics/global	GET	2025-04-20 20:41:22.156571	2025-04-20 20:41:22.158945	200	2
929	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:22.337742	2025-04-20 20:41:22.38315	200	45
930	\N	/status	GET	2025-04-20 20:41:26.966153	2025-04-20 20:41:26.966155	200	0
931	\N	/status	GET	2025-04-20 20:41:27.183377	2025-04-20 20:41:27.183378	200	0
932	\N	/api/statistics/global	GET	2025-04-20 20:41:27.201193	2025-04-20 20:41:27.203931	200	2
933	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:27.387461	2025-04-20 20:41:27.427858	200	40
934	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:41:31.026471	2025-04-20 20:41:31.026565	304	0
935	\N	/status	GET	2025-04-20 20:41:31.318131	2025-04-20 20:41:31.318133	200	0
936	\N	/status	GET	2025-04-20 20:41:31.497485	2025-04-20 20:41:31.497486	200	0
937	\N	/api/statistics/global	GET	2025-04-20 20:41:31.515124	2025-04-20 20:41:31.516302	200	1
938	\N	/api/statistics/global	GET	2025-04-20 20:41:31.683971	2025-04-20 20:41:31.685056	200	1
939	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:31.701847	2025-04-20 20:41:31.704006	200	2
940	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:31.872071	2025-04-20 20:41:31.874082	200	2
941	\N	/status	GET	2025-04-20 20:41:36.41294	2025-04-20 20:41:36.412941	200	0
942	\N	/status	GET	2025-04-20 20:41:36.589124	2025-04-20 20:41:36.589125	200	0
943	\N	/api/statistics/global	GET	2025-04-20 20:41:36.730994	2025-04-20 20:41:36.732047	200	1
944	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:36.902045	2025-04-20 20:41:36.90575	200	3
945	\N	/status	GET	2025-04-20 20:41:38.79808	2025-04-20 20:41:38.798081	200	0
946	\N	/api/statistics/global	GET	2025-04-20 20:41:38.972732	2025-04-20 20:41:38.973947	200	1
947	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:39.171159	2025-04-20 20:41:39.172801	200	1
948	\N	/status	GET	2025-04-20 20:41:41.31931	2025-04-20 20:41:41.319312	200	0
949	\N	/status	GET	2025-04-20 20:41:41.516682	2025-04-20 20:41:41.516684	200	0
950	\N	/api/statistics/global	GET	2025-04-20 20:41:41.535267	2025-04-20 20:41:41.536285	200	1
951	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:41.703072	2025-04-20 20:41:41.70492	200	1
952	\N	/status	GET	2025-04-20 20:41:46.387433	2025-04-20 20:41:46.387434	200	0
953	\N	/status	GET	2025-04-20 20:41:46.555371	2025-04-20 20:41:46.555373	200	0
954	\N	/api/statistics/global	GET	2025-04-20 20:41:46.574018	2025-04-20 20:41:46.575058	200	1
955	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:46.752988	2025-04-20 20:41:46.754557	200	1
956	\N	/status	GET	2025-04-20 20:41:51.520925	2025-04-20 20:41:51.520945	200	0
957	\N	/status	GET	2025-04-20 20:41:51.69187	2025-04-20 20:41:51.691871	200	0
958	\N	/api/statistics/global	GET	2025-04-20 20:41:51.82848	2025-04-20 20:41:51.829724	200	1
959	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:52.010037	2025-04-20 20:41:52.012512	200	2
960	\N	/status	GET	2025-04-20 20:41:56.412558	2025-04-20 20:41:56.41256	200	0
961	\N	/status	GET	2025-04-20 20:41:56.592871	2025-04-20 20:41:56.592873	200	0
962	\N	/api/statistics/global	GET	2025-04-20 20:41:56.728814	2025-04-20 20:41:56.730175	200	1
963	\N	/api/statistics/endpoints	GET	2025-04-20 20:41:56.915588	2025-04-20 20:41:56.917009	200	1
964	\N	/status	GET	2025-04-20 20:42:01.398806	2025-04-20 20:42:01.398808	200	0
965	\N	/status	GET	2025-04-20 20:42:01.584562	2025-04-20 20:42:01.584564	200	0
966	\N	/api/statistics/global	GET	2025-04-20 20:42:01.584565	2025-04-20 20:42:01.585883	200	1
967	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:01.775267	2025-04-20 20:42:01.788783	200	13
968	\N	/status	GET	2025-04-20 20:42:06.396463	2025-04-20 20:42:06.396464	200	0
969	\N	/status	GET	2025-04-20 20:42:06.567702	2025-04-20 20:42:06.567703	200	0
970	\N	/api/statistics/global	GET	2025-04-20 20:42:06.704869	2025-04-20 20:42:06.705845	200	0
971	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:06.872736	2025-04-20 20:42:06.874608	200	1
972	\N	/status	GET	2025-04-20 20:42:11.421941	2025-04-20 20:42:11.421942	200	0
973	\N	/status	GET	2025-04-20 20:42:11.58929	2025-04-20 20:42:11.589292	200	0
974	\N	/api/statistics/global	GET	2025-04-20 20:42:11.606025	2025-04-20 20:42:11.607445	200	1
975	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:11.776114	2025-04-20 20:42:11.778524	200	2
976	\N	/status	GET	2025-04-20 20:42:16.429075	2025-04-20 20:42:16.429077	200	0
977	\N	/status	GET	2025-04-20 20:42:16.597507	2025-04-20 20:42:16.597508	200	0
978	\N	/api/statistics/global	GET	2025-04-20 20:42:16.612793	2025-04-20 20:42:16.613864	200	1
979	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:16.787326	2025-04-20 20:42:16.788785	200	1
980	\N	/status	GET	2025-04-20 20:42:21.433771	2025-04-20 20:42:21.433772	200	0
981	\N	/status	GET	2025-04-20 20:42:21.602595	2025-04-20 20:42:21.602597	200	0
982	\N	/api/statistics/global	GET	2025-04-20 20:42:21.620936	2025-04-20 20:42:21.623326	200	2
983	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:21.792433	2025-04-20 20:42:21.794519	200	2
984	\N	/status	GET	2025-04-20 20:42:26.56327	2025-04-20 20:42:26.563273	200	0
985	\N	/status	GET	2025-04-20 20:42:26.774072	2025-04-20 20:42:26.774074	200	0
986	\N	/api/statistics/global	GET	2025-04-20 20:42:26.912839	2025-04-20 20:42:26.913865	200	1
987	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:27.08545	2025-04-20 20:42:27.09155	200	6
988	\N	/status	GET	2025-04-20 20:42:31.398321	2025-04-20 20:42:31.398323	200	0
989	\N	/status	GET	2025-04-20 20:42:31.57259	2025-04-20 20:42:31.572592	200	0
990	\N	/api/statistics/global	GET	2025-04-20 20:42:31.590423	2025-04-20 20:42:31.591349	200	0
991	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:31.761847	2025-04-20 20:42:31.763565	200	1
992	\N	/status	GET	2025-04-20 20:42:36.405788	2025-04-20 20:42:36.405789	200	0
993	\N	/status	GET	2025-04-20 20:42:36.58514	2025-04-20 20:42:36.585141	200	0
994	\N	/api/statistics/global	GET	2025-04-20 20:42:36.593928	2025-04-20 20:42:36.595159	200	1
995	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:36.766345	2025-04-20 20:42:36.768193	200	1
996	\N	/status	GET	2025-04-20 20:42:41.385721	2025-04-20 20:42:41.385723	200	0
997	\N	/status	GET	2025-04-20 20:42:41.591189	2025-04-20 20:42:41.591191	200	0
998	\N	/api/statistics/global	GET	2025-04-20 20:42:41.609878	2025-04-20 20:42:41.611156	200	1
999	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:41.783991	2025-04-20 20:42:41.785773	200	1
1000	\N	/status	GET	2025-04-20 20:42:46.41591	2025-04-20 20:42:46.415915	200	0
1001	\N	/status	GET	2025-04-20 20:42:46.593186	2025-04-20 20:42:46.593187	200	0
1002	\N	/api/statistics/global	GET	2025-04-20 20:42:46.602469	2025-04-20 20:42:46.603454	200	0
1003	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:46.781532	2025-04-20 20:42:46.788879	200	7
1004	\N	/status	GET	2025-04-20 20:42:51.400761	2025-04-20 20:42:51.400763	200	0
1005	\N	/status	GET	2025-04-20 20:42:51.567536	2025-04-20 20:42:51.567538	200	0
1006	\N	/api/statistics/global	GET	2025-04-20 20:42:51.585889	2025-04-20 20:42:51.587166	200	1
1007	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:51.767586	2025-04-20 20:42:51.769558	200	1
1008	\N	/status	GET	2025-04-20 20:42:56.488084	2025-04-20 20:42:56.488089	200	0
1009	\N	/status	GET	2025-04-20 20:42:56.654835	2025-04-20 20:42:56.654836	200	0
1010	\N	/api/statistics/global	GET	2025-04-20 20:42:56.672832	2025-04-20 20:42:56.67419	200	1
1011	\N	/api/statistics/endpoints	GET	2025-04-20 20:42:56.864945	2025-04-20 20:42:56.866751	200	1
1012	\N	/status	GET	2025-04-20 20:43:01.415684	2025-04-20 20:43:01.415689	200	0
1013	\N	/status	GET	2025-04-20 20:43:01.582268	2025-04-20 20:43:01.582269	200	0
1014	\N	/api/statistics/global	GET	2025-04-20 20:43:01.600073	2025-04-20 20:43:01.601185	200	1
1015	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:01.777538	2025-04-20 20:43:01.780647	200	3
1016	\N	/status	GET	2025-04-20 20:43:06.38733	2025-04-20 20:43:06.387334	200	0
1017	\N	/status	GET	2025-04-20 20:43:06.5536	2025-04-20 20:43:06.553602	200	0
1018	\N	/api/statistics/global	GET	2025-04-20 20:43:06.571903	2025-04-20 20:43:06.573146	200	1
1019	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:06.760616	2025-04-20 20:43:06.762166	200	1
1020	\N	/status	GET	2025-04-20 20:43:11.381337	2025-04-20 20:43:11.381338	200	0
1021	\N	/status	GET	2025-04-20 20:43:11.560311	2025-04-20 20:43:11.560312	200	0
1022	\N	/api/statistics/global	GET	2025-04-20 20:43:11.695812	2025-04-20 20:43:11.697218	200	1
1023	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:11.865187	2025-04-20 20:43:11.86694	200	1
1024	\N	/status	GET	2025-04-20 20:43:16.412166	2025-04-20 20:43:16.412168	200	0
1026	\N	/api/statistics/global	GET	2025-04-20 20:43:16.593084	2025-04-20 20:43:16.59434	200	1
1065	\N	/status	GET	2025-04-20 20:44:06.557247	2025-04-20 20:44:06.557249	200	0
1069	\N	/status	GET	2025-04-20 20:44:11.570505	2025-04-20 20:44:11.570507	200	0
1272	\N	/status	GET	2025-04-20 20:48:26.419462	2025-04-20 20:48:26.419463	200	0
1273	\N	/status	GET	2025-04-20 20:48:26.588403	2025-04-20 20:48:26.588405	200	0
1274	\N	/api/statistics/global	GET	2025-04-20 20:48:26.606492	2025-04-20 20:48:26.607648	200	1
1275	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:26.789833	2025-04-20 20:48:26.792269	200	2
1276	\N	/status	GET	2025-04-20 20:48:31.400516	2025-04-20 20:48:31.400518	200	0
1277	\N	/status	GET	2025-04-20 20:48:31.566837	2025-04-20 20:48:31.566839	200	0
1278	\N	/api/statistics/global	GET	2025-04-20 20:48:31.585151	2025-04-20 20:48:31.586202	200	1
1279	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:31.758235	2025-04-20 20:48:31.760557	200	2
1280	\N	/status	GET	2025-04-20 20:48:36.419067	2025-04-20 20:48:36.419069	200	0
1281	\N	/status	GET	2025-04-20 20:48:36.58594	2025-04-20 20:48:36.585941	200	0
1282	\N	/api/statistics/global	GET	2025-04-20 20:48:36.604748	2025-04-20 20:48:36.605931	200	1
1283	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:36.794996	2025-04-20 20:48:36.797258	200	2
1284	\N	/status	GET	2025-04-20 20:48:41.421506	2025-04-20 20:48:41.421508	200	0
1285	\N	/status	GET	2025-04-20 20:48:41.661683	2025-04-20 20:48:41.661684	200	0
1286	\N	/api/statistics/global	GET	2025-04-20 20:48:41.678348	2025-04-20 20:48:41.679856	200	1
1287	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:41.875636	2025-04-20 20:48:41.877482	200	1
1288	\N	/status	GET	2025-04-20 20:48:46.400601	2025-04-20 20:48:46.400603	200	0
1289	\N	/status	GET	2025-04-20 20:48:46.579585	2025-04-20 20:48:46.579587	200	0
1290	\N	/api/statistics/global	GET	2025-04-20 20:48:46.58757	2025-04-20 20:48:46.588933	200	1
1291	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:46.763854	2025-04-20 20:48:46.765719	200	1
1292	\N	/status	GET	2025-04-20 20:48:51.440207	2025-04-20 20:48:51.440208	200	0
1293	\N	/status	GET	2025-04-20 20:48:51.606665	2025-04-20 20:48:51.606667	200	0
1294	\N	/api/statistics/global	GET	2025-04-20 20:48:51.624334	2025-04-20 20:48:51.625383	200	1
1295	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:51.795011	2025-04-20 20:48:51.797309	200	2
1296	\N	/status	GET	2025-04-20 20:48:56.40847	2025-04-20 20:48:56.408471	200	0
1297	\N	/status	GET	2025-04-20 20:48:56.581672	2025-04-20 20:48:56.581673	200	0
1298	\N	/api/statistics/global	GET	2025-04-20 20:48:56.598738	2025-04-20 20:48:56.600186	200	1
1299	\N	/api/statistics/endpoints	GET	2025-04-20 20:48:56.773147	2025-04-20 20:48:56.774937	200	1
1300	\N	/status	GET	2025-04-20 20:49:01.406623	2025-04-20 20:49:01.406624	200	0
1301	\N	/status	GET	2025-04-20 20:49:01.572162	2025-04-20 20:49:01.572164	200	0
1302	\N	/api/statistics/global	GET	2025-04-20 20:49:01.589804	2025-04-20 20:49:01.591119	200	1
1303	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:01.77095	2025-04-20 20:49:01.772606	200	1
1304	\N	/status	GET	2025-04-20 20:49:06.411389	2025-04-20 20:49:06.411391	200	0
1305	\N	/status	GET	2025-04-20 20:49:06.58593	2025-04-20 20:49:06.585931	200	0
1306	\N	/api/statistics/global	GET	2025-04-20 20:49:06.60425	2025-04-20 20:49:06.605609	200	1
1307	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:06.778074	2025-04-20 20:49:06.782025	200	3
1308	\N	/status	GET	2025-04-20 20:49:11.393147	2025-04-20 20:49:11.393149	200	0
1309	\N	/status	GET	2025-04-20 20:49:11.570072	2025-04-20 20:49:11.570075	200	0
1310	\N	/api/statistics/global	GET	2025-04-20 20:49:11.58838	2025-04-20 20:49:11.589726	200	1
1311	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:11.772861	2025-04-20 20:49:11.774513	200	1
1312	\N	/status	GET	2025-04-20 20:49:16.396248	2025-04-20 20:49:16.396249	200	0
1313	\N	/status	GET	2025-04-20 20:49:16.56924	2025-04-20 20:49:16.569242	200	0
1314	\N	/api/statistics/global	GET	2025-04-20 20:49:16.586923	2025-04-20 20:49:16.588796	200	1
1315	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:16.763195	2025-04-20 20:49:16.765402	200	2
1316	\N	/status	GET	2025-04-20 20:49:21.416383	2025-04-20 20:49:21.416385	200	0
1317	\N	/status	GET	2025-04-20 20:49:21.584364	2025-04-20 20:49:21.584366	200	0
1318	\N	/api/statistics/global	GET	2025-04-20 20:49:21.609437	2025-04-20 20:49:21.611048	200	1
1319	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:21.782539	2025-04-20 20:49:21.784898	200	2
1320	\N	/status	GET	2025-04-20 20:49:26.43199	2025-04-20 20:49:26.431992	200	0
1321	\N	/status	GET	2025-04-20 20:49:26.615917	2025-04-20 20:49:26.615919	200	0
1322	\N	/api/statistics/global	GET	2025-04-20 20:49:26.634072	2025-04-20 20:49:26.635165	200	1
1323	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:26.806143	2025-04-20 20:49:26.808609	200	2
1324	\N	/status	GET	2025-04-20 20:49:31.394299	2025-04-20 20:49:31.394301	200	0
1325	\N	/status	GET	2025-04-20 20:49:31.562626	2025-04-20 20:49:31.562628	200	0
1326	\N	/api/statistics/global	GET	2025-04-20 20:49:31.581208	2025-04-20 20:49:31.582833	200	1
1327	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:31.767412	2025-04-20 20:49:31.769445	200	2
1328	\N	/status	GET	2025-04-20 20:49:36.433639	2025-04-20 20:49:36.43364	200	0
1329	\N	/status	GET	2025-04-20 20:49:36.603991	2025-04-20 20:49:36.603993	200	0
1330	\N	/api/statistics/global	GET	2025-04-20 20:49:36.622616	2025-04-20 20:49:36.623647	200	1
1331	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:36.869568	2025-04-20 20:49:36.873837	200	4
1332	\N	/status	GET	2025-04-20 20:49:41.464752	2025-04-20 20:49:41.464754	200	0
1333	\N	/status	GET	2025-04-20 20:49:41.635086	2025-04-20 20:49:41.635092	200	0
1334	\N	/api/statistics/global	GET	2025-04-20 20:49:41.653518	2025-04-20 20:49:41.654537	200	1
1335	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:41.830372	2025-04-20 20:49:41.83282	200	2
1336	\N	/status	GET	2025-04-20 20:49:46.415651	2025-04-20 20:49:46.415653	200	0
1337	\N	/status	GET	2025-04-20 20:49:46.581164	2025-04-20 20:49:46.581165	200	0
1338	\N	/api/statistics/global	GET	2025-04-20 20:49:46.598915	2025-04-20 20:49:46.600195	200	1
1339	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:46.770546	2025-04-20 20:49:46.772505	200	1
1340	\N	/status	GET	2025-04-20 20:49:51.406397	2025-04-20 20:49:51.406398	200	0
1341	\N	/status	GET	2025-04-20 20:49:51.600394	2025-04-20 20:49:51.600395	200	0
1342	\N	/api/statistics/global	GET	2025-04-20 20:49:51.618893	2025-04-20 20:49:51.620325	200	1
1516	\N	/status	GET	2025-04-20 20:53:31.430347	2025-04-20 20:53:31.430348	200	0
1025	\N	/status	GET	2025-04-20 20:43:16.593149	2025-04-20 20:43:16.593151	200	0
1027	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:16.767609	2025-04-20 20:43:16.769555	200	1
1028	\N	/status	GET	2025-04-20 20:43:21.418387	2025-04-20 20:43:21.418389	200	0
1029	\N	/status	GET	2025-04-20 20:43:21.589642	2025-04-20 20:43:21.589643	200	0
1030	\N	/api/statistics/global	GET	2025-04-20 20:43:21.608304	2025-04-20 20:43:21.609389	200	1
1031	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:21.7803	2025-04-20 20:43:21.781803	200	1
1032	\N	/status	GET	2025-04-20 20:43:26.410287	2025-04-20 20:43:26.410292	200	0
1033	\N	/status	GET	2025-04-20 20:43:26.576841	2025-04-20 20:43:26.576843	200	0
1034	\N	/api/statistics/global	GET	2025-04-20 20:43:26.595159	2025-04-20 20:43:26.596328	200	1
1035	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:26.768523	2025-04-20 20:43:26.770181	200	1
1036	\N	/status	GET	2025-04-20 20:43:31.403124	2025-04-20 20:43:31.403126	200	0
1037	\N	/status	GET	2025-04-20 20:43:31.571461	2025-04-20 20:43:31.571463	200	0
1038	\N	/api/statistics/global	GET	2025-04-20 20:43:31.598121	2025-04-20 20:43:31.599395	200	1
1039	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:31.765123	2025-04-20 20:43:31.767271	200	2
1040	\N	/status	GET	2025-04-20 20:43:36.38441	2025-04-20 20:43:36.384412	200	0
1041	\N	/status	GET	2025-04-20 20:43:36.550484	2025-04-20 20:43:36.550485	200	0
1042	\N	/api/statistics/global	GET	2025-04-20 20:43:36.568325	2025-04-20 20:43:36.569631	200	1
1043	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:36.736822	2025-04-20 20:43:36.738856	200	2
1044	\N	/status	GET	2025-04-20 20:43:41.418896	2025-04-20 20:43:41.418897	200	0
1045	\N	/status	GET	2025-04-20 20:43:41.594335	2025-04-20 20:43:41.594336	200	0
1046	\N	/api/statistics/global	GET	2025-04-20 20:43:41.645114	2025-04-20 20:43:41.64647	200	1
1047	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:41.835773	2025-04-20 20:43:41.837416	200	1
1048	\N	/status	GET	2025-04-20 20:43:46.40319	2025-04-20 20:43:46.403192	200	0
1049	\N	/status	GET	2025-04-20 20:43:46.571935	2025-04-20 20:43:46.571937	200	0
1050	\N	/api/statistics/global	GET	2025-04-20 20:43:46.591418	2025-04-20 20:43:46.593279	200	1
1051	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:46.779068	2025-04-20 20:43:46.780689	200	1
1052	\N	/status	GET	2025-04-20 20:43:51.466436	2025-04-20 20:43:51.466437	200	0
1053	\N	/status	GET	2025-04-20 20:43:51.63722	2025-04-20 20:43:51.637222	200	0
1054	\N	/api/statistics/global	GET	2025-04-20 20:43:51.656014	2025-04-20 20:43:51.65707	200	1
1055	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:51.825087	2025-04-20 20:43:51.826703	200	1
1056	\N	/status	GET	2025-04-20 20:43:56.469043	2025-04-20 20:43:56.469054	200	0
1057	\N	/status	GET	2025-04-20 20:43:56.63768	2025-04-20 20:43:56.637682	200	0
1058	\N	/api/statistics/global	GET	2025-04-20 20:43:56.696539	2025-04-20 20:43:56.698271	200	1
1059	\N	/api/statistics/endpoints	GET	2025-04-20 20:43:56.947442	2025-04-20 20:43:56.949356	200	1
1060	\N	/status	GET	2025-04-20 20:44:01.412363	2025-04-20 20:44:01.412364	200	0
1061	\N	/status	GET	2025-04-20 20:44:01.580486	2025-04-20 20:44:01.580487	200	0
1062	\N	/api/statistics/global	GET	2025-04-20 20:44:01.599476	2025-04-20 20:44:01.600511	200	1
1063	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:01.772705	2025-04-20 20:44:01.774326	200	1
1064	\N	/status	GET	2025-04-20 20:44:06.389869	2025-04-20 20:44:06.38987	200	0
1066	\N	/api/statistics/global	GET	2025-04-20 20:44:06.557244	2025-04-20 20:44:06.558175	200	0
1067	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:06.729417	2025-04-20 20:44:06.731444	200	2
1068	\N	/status	GET	2025-04-20 20:44:11.394346	2025-04-20 20:44:11.394347	200	0
1070	\N	/api/statistics/global	GET	2025-04-20 20:44:11.570504	2025-04-20 20:44:11.571676	200	1
1071	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:11.755659	2025-04-20 20:44:11.758254	200	2
1072	\N	/status	GET	2025-04-20 20:44:16.393892	2025-04-20 20:44:16.393893	200	0
1073	\N	/status	GET	2025-04-20 20:44:16.56496	2025-04-20 20:44:16.564962	200	0
1074	\N	/api/statistics/global	GET	2025-04-20 20:44:16.582655	2025-04-20 20:44:16.583642	200	0
1075	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:16.753867	2025-04-20 20:44:16.755566	200	1
1076	\N	/status	GET	2025-04-20 20:44:21.392427	2025-04-20 20:44:21.392429	200	0
1077	\N	/status	GET	2025-04-20 20:44:21.558172	2025-04-20 20:44:21.558173	200	0
1078	\N	/api/statistics/global	GET	2025-04-20 20:44:21.575961	2025-04-20 20:44:21.577143	200	1
1079	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:21.745495	2025-04-20 20:44:21.747323	200	1
1080	\N	/status	GET	2025-04-20 20:44:26.385532	2025-04-20 20:44:26.385535	200	0
1081	\N	/status	GET	2025-04-20 20:44:26.551471	2025-04-20 20:44:26.551473	200	0
1082	\N	/api/statistics/global	GET	2025-04-20 20:44:26.574051	2025-04-20 20:44:26.575224	200	1
1083	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:26.741052	2025-04-20 20:44:26.743379	200	2
1084	\N	/status	GET	2025-04-20 20:44:31.394435	2025-04-20 20:44:31.394436	200	0
1085	\N	/status	GET	2025-04-20 20:44:31.560027	2025-04-20 20:44:31.560029	200	0
1086	\N	/api/statistics/global	GET	2025-04-20 20:44:31.577943	2025-04-20 20:44:31.579144	200	1
1087	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:31.754258	2025-04-20 20:44:31.755832	200	1
1088	\N	/status	GET	2025-04-20 20:44:36.382338	2025-04-20 20:44:36.38234	200	0
1089	\N	/status	GET	2025-04-20 20:44:36.548	2025-04-20 20:44:36.548001	200	0
1090	\N	/api/statistics/global	GET	2025-04-20 20:44:36.566314	2025-04-20 20:44:36.567312	200	0
1091	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:36.835549	2025-04-20 20:44:36.837166	200	1
1092	\N	/status	GET	2025-04-20 20:44:41.40581	2025-04-20 20:44:41.405811	200	0
1093	\N	/status	GET	2025-04-20 20:44:41.574055	2025-04-20 20:44:41.574057	200	0
1094	\N	/api/statistics/global	GET	2025-04-20 20:44:41.592448	2025-04-20 20:44:41.59347	200	1
1095	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:41.760666	2025-04-20 20:44:41.762457	200	1
1096	\N	/status	GET	2025-04-20 20:44:46.404449	2025-04-20 20:44:46.404451	200	0
1097	\N	/status	GET	2025-04-20 20:44:46.570523	2025-04-20 20:44:46.570525	200	0
1098	\N	/api/statistics/global	GET	2025-04-20 20:44:46.588815	2025-04-20 20:44:46.589941	200	1
1099	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:46.766966	2025-04-20 20:44:46.7698	200	2
1100	\N	/status	GET	2025-04-20 20:44:51.417046	2025-04-20 20:44:51.417048	200	0
1101	\N	/status	GET	2025-04-20 20:44:51.598166	2025-04-20 20:44:51.598167	200	0
1102	\N	/api/statistics/global	GET	2025-04-20 20:44:51.617644	2025-04-20 20:44:51.618719	200	1
1103	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:51.797757	2025-04-20 20:44:51.799491	200	1
1104	\N	/status	GET	2025-04-20 20:44:56.389392	2025-04-20 20:44:56.389393	200	0
1105	\N	/status	GET	2025-04-20 20:44:56.559167	2025-04-20 20:44:56.559169	200	0
1106	\N	/api/statistics/global	GET	2025-04-20 20:44:56.577254	2025-04-20 20:44:56.578322	200	1
1107	\N	/api/statistics/endpoints	GET	2025-04-20 20:44:56.750206	2025-04-20 20:44:56.751771	200	1
1108	\N	/status	GET	2025-04-20 20:45:01.41977	2025-04-20 20:45:01.419772	200	0
1109	\N	/status	GET	2025-04-20 20:45:01.588497	2025-04-20 20:45:01.588498	200	0
1110	\N	/api/statistics/global	GET	2025-04-20 20:45:01.606428	2025-04-20 20:45:01.607583	200	1
1111	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:01.778906	2025-04-20 20:45:01.780476	200	1
1112	\N	/status	GET	2025-04-20 20:45:06.393007	2025-04-20 20:45:06.393008	200	0
1113	\N	/status	GET	2025-04-20 20:45:06.559613	2025-04-20 20:45:06.559614	200	0
1114	\N	/api/statistics/global	GET	2025-04-20 20:45:06.57847	2025-04-20 20:45:06.579728	200	1
1115	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:06.758695	2025-04-20 20:45:06.76029	200	1
1116	\N	/status	GET	2025-04-20 20:45:11.414816	2025-04-20 20:45:11.414817	200	0
1117	\N	/status	GET	2025-04-20 20:45:11.59411	2025-04-20 20:45:11.594112	200	0
1118	\N	/api/statistics/global	GET	2025-04-20 20:45:11.594106	2025-04-20 20:45:11.595367	200	1
1119	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:11.773613	2025-04-20 20:45:11.775341	200	1
1120	\N	/status	GET	2025-04-20 20:45:16.387595	2025-04-20 20:45:16.387597	200	0
1121	\N	/status	GET	2025-04-20 20:45:16.554752	2025-04-20 20:45:16.554754	200	0
1122	\N	/api/statistics/global	GET	2025-04-20 20:45:16.573288	2025-04-20 20:45:16.57476	200	1
1123	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:16.753767	2025-04-20 20:45:16.756793	200	3
1124	\N	/status	GET	2025-04-20 20:45:21.404834	2025-04-20 20:45:21.404836	200	0
1125	\N	/status	GET	2025-04-20 20:45:21.571131	2025-04-20 20:45:21.571133	200	0
1126	\N	/api/statistics/global	GET	2025-04-20 20:45:21.589373	2025-04-20 20:45:21.591036	200	1
1127	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:21.771537	2025-04-20 20:45:21.773229	200	1
1128	\N	/status	GET	2025-04-20 20:45:26.404205	2025-04-20 20:45:26.404207	200	0
1129	\N	/status	GET	2025-04-20 20:45:26.573818	2025-04-20 20:45:26.573819	200	0
1130	\N	/api/statistics/global	GET	2025-04-20 20:45:26.591612	2025-04-20 20:45:26.592679	200	1
1131	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:26.770023	2025-04-20 20:45:26.771663	200	1
1132	\N	/status	GET	2025-04-20 20:45:31.395209	2025-04-20 20:45:31.395211	200	0
1133	\N	/status	GET	2025-04-20 20:45:31.562977	2025-04-20 20:45:31.562979	200	0
1134	\N	/api/statistics/global	GET	2025-04-20 20:45:31.581411	2025-04-20 20:45:31.582535	200	1
1135	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:31.767513	2025-04-20 20:45:31.769164	200	1
1136	\N	/status	GET	2025-04-20 20:45:36.419407	2025-04-20 20:45:36.419408	200	0
1137	\N	/status	GET	2025-04-20 20:45:36.586443	2025-04-20 20:45:36.586444	200	0
1138	\N	/api/statistics/global	GET	2025-04-20 20:45:36.604117	2025-04-20 20:45:36.605475	200	1
1139	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:36.782377	2025-04-20 20:45:36.784543	200	2
1140	\N	/status	GET	2025-04-20 20:45:41.409881	2025-04-20 20:45:41.409883	200	0
1141	\N	/status	GET	2025-04-20 20:45:41.588805	2025-04-20 20:45:41.588807	200	0
1142	\N	/api/statistics/global	GET	2025-04-20 20:45:41.607044	2025-04-20 20:45:41.608319	200	1
1143	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:41.776608	2025-04-20 20:45:41.779809	200	3
1144	\N	/status	GET	2025-04-20 20:45:46.426065	2025-04-20 20:45:46.426067	200	0
1145	\N	/status	GET	2025-04-20 20:45:46.594241	2025-04-20 20:45:46.594243	200	0
1146	\N	/api/statistics/global	GET	2025-04-20 20:45:46.612312	2025-04-20 20:45:46.613268	200	0
1147	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:46.784758	2025-04-20 20:45:46.78696	200	2
1148	\N	/status	GET	2025-04-20 20:45:51.391565	2025-04-20 20:45:51.391566	200	0
1149	\N	/status	GET	2025-04-20 20:45:51.560669	2025-04-20 20:45:51.56067	200	0
1150	\N	/api/statistics/global	GET	2025-04-20 20:45:51.579773	2025-04-20 20:45:51.583731	200	3
1151	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:51.761479	2025-04-20 20:45:51.763276	200	1
1152	\N	/status	GET	2025-04-20 20:45:56.419723	2025-04-20 20:45:56.419724	200	0
1153	\N	/status	GET	2025-04-20 20:45:56.599959	2025-04-20 20:45:56.599961	200	0
1154	\N	/api/statistics/global	GET	2025-04-20 20:45:56.609987	2025-04-20 20:45:56.611001	200	1
1155	\N	/api/statistics/endpoints	GET	2025-04-20 20:45:56.835787	2025-04-20 20:45:56.837386	200	1
1156	\N	/status	GET	2025-04-20 20:46:01.390507	2025-04-20 20:46:01.390509	200	0
1157	\N	/status	GET	2025-04-20 20:46:01.558165	2025-04-20 20:46:01.558166	200	0
1158	\N	/api/statistics/global	GET	2025-04-20 20:46:01.576691	2025-04-20 20:46:01.577835	200	1
1159	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:01.754053	2025-04-20 20:46:01.756743	200	2
1160	\N	/status	GET	2025-04-20 20:46:06.398647	2025-04-20 20:46:06.398649	200	0
1161	\N	/status	GET	2025-04-20 20:46:06.57312	2025-04-20 20:46:06.573122	200	0
1162	\N	/api/statistics/global	GET	2025-04-20 20:46:06.597259	2025-04-20 20:46:06.598977	200	1
1163	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:06.77327	2025-04-20 20:46:06.77508	200	1
1164	\N	/status	GET	2025-04-20 20:46:11.397029	2025-04-20 20:46:11.39703	200	0
1165	\N	/status	GET	2025-04-20 20:46:11.564642	2025-04-20 20:46:11.564643	200	0
1166	\N	/api/statistics/global	GET	2025-04-20 20:46:11.583001	2025-04-20 20:46:11.584001	200	0
1167	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:11.752216	2025-04-20 20:46:11.75424	200	2
1168	\N	/status	GET	2025-04-20 20:46:16.393944	2025-04-20 20:46:16.393945	200	0
1169	\N	/status	GET	2025-04-20 20:46:16.560955	2025-04-20 20:46:16.560956	200	0
1170	\N	/api/statistics/global	GET	2025-04-20 20:46:16.579471	2025-04-20 20:46:16.580699	200	1
1171	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:16.760477	2025-04-20 20:46:16.76223	200	1
1172	\N	/status	GET	2025-04-20 20:46:21.40359	2025-04-20 20:46:21.403591	200	0
1173	\N	/status	GET	2025-04-20 20:46:21.574409	2025-04-20 20:46:21.574411	200	0
1174	\N	/api/statistics/global	GET	2025-04-20 20:46:21.588303	2025-04-20 20:46:21.591276	200	2
1175	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:21.768068	2025-04-20 20:46:21.769875	200	1
1176	\N	/status	GET	2025-04-20 20:46:26.398739	2025-04-20 20:46:26.39874	200	0
1177	\N	/status	GET	2025-04-20 20:46:26.569456	2025-04-20 20:46:26.569458	200	0
1178	\N	/api/statistics/global	GET	2025-04-20 20:46:26.587723	2025-04-20 20:46:26.588709	200	0
1179	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:26.767836	2025-04-20 20:46:26.769671	200	1
1180	\N	/status	GET	2025-04-20 20:46:31.416798	2025-04-20 20:46:31.4168	200	0
1181	\N	/status	GET	2025-04-20 20:46:31.587638	2025-04-20 20:46:31.58764	200	0
1182	\N	/api/statistics/global	GET	2025-04-20 20:46:31.606054	2025-04-20 20:46:31.607275	200	1
1183	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:31.774858	2025-04-20 20:46:31.776429	200	1
1184	\N	/status	GET	2025-04-20 20:46:36.386743	2025-04-20 20:46:36.386745	200	0
1185	\N	/status	GET	2025-04-20 20:46:36.565954	2025-04-20 20:46:36.565956	200	0
1186	\N	/api/statistics/global	GET	2025-04-20 20:46:36.583917	2025-04-20 20:46:36.585288	200	1
1187	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:36.767752	2025-04-20 20:46:36.769409	200	1
1188	\N	/status	GET	2025-04-20 20:46:41.383212	2025-04-20 20:46:41.383214	200	0
1189	\N	/status	GET	2025-04-20 20:46:41.550753	2025-04-20 20:46:41.550755	200	0
1190	\N	/api/statistics/global	GET	2025-04-20 20:46:41.574226	2025-04-20 20:46:41.575264	200	1
1191	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:41.758101	2025-04-20 20:46:41.759848	200	1
1192	\N	/status	GET	2025-04-20 20:46:46.387794	2025-04-20 20:46:46.387795	200	0
1193	\N	/status	GET	2025-04-20 20:46:46.56993	2025-04-20 20:46:46.569932	200	0
1194	\N	/api/statistics/global	GET	2025-04-20 20:46:46.577122	2025-04-20 20:46:46.578067	200	0
1195	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:46.751569	2025-04-20 20:46:46.753699	200	2
1196	\N	/status	GET	2025-04-20 20:46:51.380017	2025-04-20 20:46:51.380018	200	0
1198	\N	/api/statistics/global	GET	2025-04-20 20:46:51.553559	2025-04-20 20:46:51.554982	200	1
1199	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:51.732147	2025-04-20 20:46:51.736906	200	4
1200	\N	/status	GET	2025-04-20 20:46:56.416459	2025-04-20 20:46:56.41646	200	0
1201	\N	/status	GET	2025-04-20 20:46:56.586657	2025-04-20 20:46:56.586658	200	0
1202	\N	/api/statistics/global	GET	2025-04-20 20:46:56.604509	2025-04-20 20:46:56.605963	200	1
1203	\N	/api/statistics/endpoints	GET	2025-04-20 20:46:56.773833	2025-04-20 20:46:56.77572	200	1
1204	\N	/status	GET	2025-04-20 20:47:01.417092	2025-04-20 20:47:01.417094	200	0
1205	\N	/status	GET	2025-04-20 20:47:01.584398	2025-04-20 20:47:01.5844	200	0
1206	\N	/api/statistics/global	GET	2025-04-20 20:47:01.602173	2025-04-20 20:47:01.603659	200	1
1207	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:01.777939	2025-04-20 20:47:01.781989	200	4
1208	\N	/status	GET	2025-04-20 20:47:06.39981	2025-04-20 20:47:06.399812	200	0
1209	\N	/status	GET	2025-04-20 20:47:06.565998	2025-04-20 20:47:06.565999	200	0
1210	\N	/api/statistics/global	GET	2025-04-20 20:47:06.583852	2025-04-20 20:47:06.585369	200	1
1211	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:06.759054	2025-04-20 20:47:06.760707	200	1
1212	\N	/status	GET	2025-04-20 20:47:11.41399	2025-04-20 20:47:11.413991	200	0
1213	\N	/status	GET	2025-04-20 20:47:11.632949	2025-04-20 20:47:11.632951	200	0
1214	\N	/api/statistics/global	GET	2025-04-20 20:47:11.650488	2025-04-20 20:47:11.651564	200	1
1215	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:11.823672	2025-04-20 20:47:11.82587	200	2
1216	\N	/status	GET	2025-04-20 20:47:16.378287	2025-04-20 20:47:16.378289	200	0
1217	\N	/status	GET	2025-04-20 20:47:16.543783	2025-04-20 20:47:16.543784	200	0
1218	\N	/api/statistics/global	GET	2025-04-20 20:47:16.562359	2025-04-20 20:47:16.563871	200	1
1219	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:16.739994	2025-04-20 20:47:16.742216	200	2
1220	\N	/status	GET	2025-04-20 20:47:21.412639	2025-04-20 20:47:21.41264	200	0
1221	\N	/status	GET	2025-04-20 20:47:21.618233	2025-04-20 20:47:21.618234	200	0
1222	\N	/api/statistics/global	GET	2025-04-20 20:47:21.63679	2025-04-20 20:47:21.63792	200	1
1223	\N	/api/statistics/endpoints	GET	2025-04-20 20:47:21.809536	2025-04-20 20:47:21.811687	200	2
1343	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:51.790191	2025-04-20 20:49:51.793615	200	3
1344	\N	/status	GET	2025-04-20 20:49:56.425912	2025-04-20 20:49:56.425913	200	0
1345	\N	/status	GET	2025-04-20 20:49:56.592813	2025-04-20 20:49:56.592813	200	0
1346	\N	/api/statistics/global	GET	2025-04-20 20:49:56.611741	2025-04-20 20:49:56.612806	200	1
1347	\N	/api/statistics/endpoints	GET	2025-04-20 20:49:56.780054	2025-04-20 20:49:56.782323	200	2
1348	\N	/status	GET	2025-04-20 20:50:01.414105	2025-04-20 20:50:01.414106	200	0
1349	\N	/status	GET	2025-04-20 20:50:01.583322	2025-04-20 20:50:01.583323	200	0
1350	\N	/api/statistics/global	GET	2025-04-20 20:50:01.601483	2025-04-20 20:50:01.602531	200	1
1351	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:01.771332	2025-04-20 20:50:01.773669	200	2
1352	\N	/status	GET	2025-04-20 20:50:06.391252	2025-04-20 20:50:06.391253	200	0
1353	\N	/status	GET	2025-04-20 20:50:06.561009	2025-04-20 20:50:06.56101	200	0
1354	\N	/api/statistics/global	GET	2025-04-20 20:50:06.579099	2025-04-20 20:50:06.580551	200	1
1355	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:06.752324	2025-04-20 20:50:06.755884	200	3
1197	\N	/status	GET	2025-04-20 20:46:51.553622	2025-04-20 20:46:51.553624	200	0
1356	\N	/status	GET	2025-04-20 20:50:11.379	2025-04-20 20:50:11.379001	200	0
1357	\N	/status	GET	2025-04-20 20:50:11.54835	2025-04-20 20:50:11.548351	200	0
1386	\N	/api/statistics/global	GET	2025-04-20 20:50:46.564892	2025-04-20 20:50:46.565994	200	1
1387	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:46.745518	2025-04-20 20:50:46.747859	200	2
1388	\N	/status	GET	2025-04-20 20:50:51.399816	2025-04-20 20:50:51.399818	200	0
1389	\N	/status	GET	2025-04-20 20:50:51.572873	2025-04-20 20:50:51.572874	200	0
1390	\N	/api/statistics/global	GET	2025-04-20 20:50:51.590615	2025-04-20 20:50:51.592087	200	1
1391	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:51.761707	2025-04-20 20:50:51.763532	200	1
1392	\N	/status	GET	2025-04-20 20:50:56.401347	2025-04-20 20:50:56.401349	200	0
1393	\N	/status	GET	2025-04-20 20:50:56.568724	2025-04-20 20:50:56.568725	200	0
1394	\N	/api/statistics/global	GET	2025-04-20 20:50:56.587444	2025-04-20 20:50:56.588705	200	1
1395	\N	/api/statistics/endpoints	GET	2025-04-20 20:50:56.769769	2025-04-20 20:50:56.771789	200	2
1396	\N	/status	GET	2025-04-20 20:51:01.38631	2025-04-20 20:51:01.386311	200	0
1397	\N	/status	GET	2025-04-20 20:51:01.586871	2025-04-20 20:51:01.586873	200	0
1398	\N	/api/statistics/global	GET	2025-04-20 20:51:01.60697	2025-04-20 20:51:01.608409	200	1
1399	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:01.793275	2025-04-20 20:51:01.795478	200	2
1400	\N	/status	GET	2025-04-20 20:51:06.398575	2025-04-20 20:51:06.398576	200	0
1401	\N	/status	GET	2025-04-20 20:51:06.567652	2025-04-20 20:51:06.567653	200	0
1402	\N	/api/statistics/global	GET	2025-04-20 20:51:06.586206	2025-04-20 20:51:06.587814	200	1
1403	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:06.765383	2025-04-20 20:51:06.76743	200	2
1404	\N	/status	GET	2025-04-20 20:51:11.408752	2025-04-20 20:51:11.408753	200	0
1405	\N	/status	GET	2025-04-20 20:51:11.579163	2025-04-20 20:51:11.579164	200	0
1406	\N	/api/statistics/global	GET	2025-04-20 20:51:11.597745	2025-04-20 20:51:11.599356	200	1
1407	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:11.779838	2025-04-20 20:51:11.781878	200	2
1408	\N	/status	GET	2025-04-20 20:51:16.40465	2025-04-20 20:51:16.404652	200	0
1409	\N	/status	GET	2025-04-20 20:51:16.571604	2025-04-20 20:51:16.571605	200	0
1410	\N	/api/statistics/global	GET	2025-04-20 20:51:16.590437	2025-04-20 20:51:16.591897	200	1
1411	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:16.762585	2025-04-20 20:51:16.764791	200	2
1412	\N	/status	GET	2025-04-20 20:51:21.392911	2025-04-20 20:51:21.392913	200	0
1413	\N	/status	GET	2025-04-20 20:51:21.560867	2025-04-20 20:51:21.560868	200	0
1434	\N	/api/statistics/global	GET	2025-04-20 20:51:46.595762	2025-04-20 20:51:46.596927	200	1
1435	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:46.77409	2025-04-20 20:51:46.776739	200	2
1436	\N	/status	GET	2025-04-20 20:51:51.423809	2025-04-20 20:51:51.42381	200	0
1437	\N	/status	GET	2025-04-20 20:51:51.594274	2025-04-20 20:51:51.594275	200	0
1438	\N	/api/statistics/global	GET	2025-04-20 20:51:51.611785	2025-04-20 20:51:51.613195	200	1
1439	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:51.796057	2025-04-20 20:51:51.807555	200	11
1440	\N	/status	GET	2025-04-20 20:51:56.39059	2025-04-20 20:51:56.390592	200	0
1441	\N	/status	GET	2025-04-20 20:51:56.567087	2025-04-20 20:51:56.567089	200	0
1442	\N	/api/statistics/global	GET	2025-04-20 20:51:56.585378	2025-04-20 20:51:56.586577	200	1
1443	\N	/api/statistics/endpoints	GET	2025-04-20 20:51:56.764418	2025-04-20 20:51:56.766696	200	2
1444	\N	/status	GET	2025-04-20 20:52:01.463005	2025-04-20 20:52:01.463007	200	0
1445	\N	/status	GET	2025-04-20 20:52:01.631218	2025-04-20 20:52:01.63122	200	0
1446	\N	/api/statistics/global	GET	2025-04-20 20:52:01.650287	2025-04-20 20:52:01.651608	200	1
1447	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:01.890692	2025-04-20 20:52:01.893557	200	2
1448	\N	/status	GET	2025-04-20 20:52:06.406455	2025-04-20 20:52:06.406456	200	0
1449	\N	/status	GET	2025-04-20 20:52:06.574689	2025-04-20 20:52:06.574691	200	0
1450	\N	/api/statistics/global	GET	2025-04-20 20:52:06.5931	2025-04-20 20:52:06.594179	200	1
1451	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:06.766815	2025-04-20 20:52:06.768751	200	1
1452	\N	/status	GET	2025-04-20 20:52:11.410734	2025-04-20 20:52:11.410736	200	0
1454	\N	/api/statistics/global	GET	2025-04-20 20:52:11.589086	2025-04-20 20:52:11.590297	200	1
1455	\N	/api/statistics/endpoints	GET	2025-04-20 20:52:11.777946	2025-04-20 20:52:11.77976	200	1
1456	\N	/status	GET	2025-04-20 20:52:16.397214	2025-04-20 20:52:16.397216	200	0
1457	\N	/status	GET	2025-04-20 20:52:16.639617	2025-04-20 20:52:16.639619	200	0
1517	\N	/status	GET	2025-04-20 20:53:31.608521	2025-04-20 20:53:31.608522	200	0
1518	\N	/api/statistics/global	GET	2025-04-20 20:53:31.62676	2025-04-20 20:53:31.628059	200	1
1519	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:31.801785	2025-04-20 20:53:31.803919	200	2
1520	\N	/status	GET	2025-04-20 20:53:36.41919	2025-04-20 20:53:36.419191	200	0
1521	\N	/status	GET	2025-04-20 20:53:36.58525	2025-04-20 20:53:36.585251	200	0
1522	\N	/api/statistics/global	GET	2025-04-20 20:53:36.61637	2025-04-20 20:53:36.617822	200	1
1523	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:36.7955	2025-04-20 20:53:36.797898	200	2
1524	\N	/status	GET	2025-04-20 20:53:41.406991	2025-04-20 20:53:41.406993	200	0
1525	\N	/status	GET	2025-04-20 20:53:41.588518	2025-04-20 20:53:41.58852	200	0
1526	\N	/api/statistics/global	GET	2025-04-20 20:53:41.607014	2025-04-20 20:53:41.60825	200	1
1527	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:41.780801	2025-04-20 20:53:41.782694	200	1
1528	\N	/status	GET	2025-04-20 20:53:46.39861	2025-04-20 20:53:46.398612	200	0
1529	\N	/status	GET	2025-04-20 20:53:46.567015	2025-04-20 20:53:46.567016	200	0
1530	\N	/api/statistics/global	GET	2025-04-20 20:53:46.585276	2025-04-20 20:53:46.586565	200	1
1531	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:46.768242	2025-04-20 20:53:46.771023	200	2
1532	\N	/status	GET	2025-04-20 20:53:51.38218	2025-04-20 20:53:51.382181	200	0
1533	\N	/status	GET	2025-04-20 20:53:51.550378	2025-04-20 20:53:51.550379	200	0
1534	\N	/api/statistics/global	GET	2025-04-20 20:53:51.575785	2025-04-20 20:53:51.576916	200	1
1535	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:51.747458	2025-04-20 20:53:51.750257	200	2
1536	\N	/status	GET	2025-04-20 20:53:56.40706	2025-04-20 20:53:56.407062	200	0
1537	\N	/status	GET	2025-04-20 20:53:56.573867	2025-04-20 20:53:56.573869	200	0
1538	\N	/api/statistics/global	GET	2025-04-20 20:53:56.594592	2025-04-20 20:53:56.595971	200	1
1539	\N	/api/statistics/endpoints	GET	2025-04-20 20:53:56.784346	2025-04-20 20:53:56.786066	200	1
1540	\N	/status	GET	2025-04-20 20:54:01.395986	2025-04-20 20:54:01.395988	200	0
1541	\N	/status	GET	2025-04-20 20:54:01.56306	2025-04-20 20:54:01.563061	200	0
1542	\N	/api/statistics/global	GET	2025-04-20 20:54:01.591241	2025-04-20 20:54:01.59248	200	1
1543	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:01.761977	2025-04-20 20:54:01.765025	200	3
1544	\N	/status	GET	2025-04-20 20:54:06.392794	2025-04-20 20:54:06.392795	200	0
1545	\N	/status	GET	2025-04-20 20:54:06.560988	2025-04-20 20:54:06.560989	200	0
1546	\N	/api/statistics/global	GET	2025-04-20 20:54:06.579586	2025-04-20 20:54:06.580927	200	1
1547	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:06.748481	2025-04-20 20:54:06.750298	200	1
1548	\N	/status	GET	2025-04-20 20:54:11.383393	2025-04-20 20:54:11.383394	200	0
1549	\N	/status	GET	2025-04-20 20:54:11.549923	2025-04-20 20:54:11.549925	200	0
1550	\N	/api/statistics/global	GET	2025-04-20 20:54:11.614253	2025-04-20 20:54:11.615998	200	1
1551	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:11.785507	2025-04-20 20:54:11.787559	200	2
1552	\N	/status	GET	2025-04-20 20:54:16.419263	2025-04-20 20:54:16.419265	200	0
1553	\N	/status	GET	2025-04-20 20:54:16.591843	2025-04-20 20:54:16.591844	200	0
1558	\N	/api/statistics/global	GET	2025-04-20 20:54:21.605817	2025-04-20 20:54:21.607262	200	1
1559	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:21.785055	2025-04-20 20:54:21.792624	200	7
1560	\N	/status	GET	2025-04-20 20:54:26.405662	2025-04-20 20:54:26.405664	200	0
1561	\N	/status	GET	2025-04-20 20:54:26.572989	2025-04-20 20:54:26.572991	200	0
1562	\N	/api/statistics/global	GET	2025-04-20 20:54:26.590851	2025-04-20 20:54:26.591915	200	1
1563	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:26.771298	2025-04-20 20:54:26.773329	200	2
1564	\N	/status	GET	2025-04-20 20:54:31.461039	2025-04-20 20:54:31.46104	200	0
1565	\N	/status	GET	2025-04-20 20:54:31.627889	2025-04-20 20:54:31.627891	200	0
1566	\N	/api/statistics/global	GET	2025-04-20 20:54:31.645622	2025-04-20 20:54:31.647005	200	1
1567	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:31.821615	2025-04-20 20:54:31.82436	200	2
1568	\N	/status	GET	2025-04-20 20:54:36.39332	2025-04-20 20:54:36.393322	200	0
1569	\N	/status	GET	2025-04-20 20:54:36.559763	2025-04-20 20:54:36.559765	200	0
1570	\N	/api/statistics/global	GET	2025-04-20 20:54:36.577878	2025-04-20 20:54:36.579209	200	1
1705	\N	/status	GET	2025-04-20 20:57:26.140442	2025-04-20 20:57:26.140443	200	0
1571	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:36.75842	2025-04-20 20:54:36.760118	200	1
1572	\N	/status	GET	2025-04-20 20:54:41.414522	2025-04-20 20:54:41.414524	200	0
1573	\N	/status	GET	2025-04-20 20:54:41.583193	2025-04-20 20:54:41.583194	200	0
1574	\N	/api/statistics/global	GET	2025-04-20 20:54:41.602119	2025-04-20 20:54:41.603692	200	1
1575	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:41.775264	2025-04-20 20:54:41.777791	200	2
1576	\N	/status	GET	2025-04-20 20:54:46.418633	2025-04-20 20:54:46.418635	200	0
1577	\N	/status	GET	2025-04-20 20:54:46.586922	2025-04-20 20:54:46.586924	200	0
1578	\N	/api/statistics/global	GET	2025-04-20 20:54:46.605498	2025-04-20 20:54:46.606804	200	1
1579	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:46.77744	2025-04-20 20:54:46.779422	200	1
1580	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 20:54:50.632381	2025-04-20 20:54:50.632464	304	0
1581	\N	/status	GET	2025-04-20 20:54:50.998994	2025-04-20 20:54:50.998996	200	0
1582	\N	/api/statistics/global	GET	2025-04-20 20:54:51.17185	2025-04-20 20:54:51.17346	200	1
1583	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:51.344295	2025-04-20 20:54:51.346577	200	2
1584	\N	/status	GET	2025-04-20 20:54:55.973782	2025-04-20 20:54:55.973783	200	0
1585	\N	/status	GET	2025-04-20 20:54:56.182039	2025-04-20 20:54:56.18204	200	0
1586	\N	/api/statistics/global	GET	2025-04-20 20:54:56.200919	2025-04-20 20:54:56.202839	200	1
1587	\N	/api/statistics/endpoints	GET	2025-04-20 20:54:56.380206	2025-04-20 20:54:56.40935	200	29
1588	\N	/status	GET	2025-04-20 20:55:00.974745	2025-04-20 20:55:00.974746	200	0
1589	\N	/status	GET	2025-04-20 20:55:01.170878	2025-04-20 20:55:01.170879	200	0
1590	\N	/api/statistics/global	GET	2025-04-20 20:55:01.19037	2025-04-20 20:55:01.192094	200	1
1591	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:01.365536	2025-04-20 20:55:01.370555	200	5
1592	\N	/status	GET	2025-04-20 20:55:05.98827	2025-04-20 20:55:05.988272	200	0
1593	\N	/status	GET	2025-04-20 20:55:06.190997	2025-04-20 20:55:06.190998	200	0
1594	\N	/api/statistics/global	GET	2025-04-20 20:55:06.208873	2025-04-20 20:55:06.210776	200	1
1595	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:06.398777	2025-04-20 20:55:06.400941	200	2
1596	\N	/status	GET	2025-04-20 20:55:10.968382	2025-04-20 20:55:10.968384	200	0
1597	\N	/status	GET	2025-04-20 20:55:11.13783	2025-04-20 20:55:11.137832	200	0
1598	\N	/api/statistics/global	GET	2025-04-20 20:55:11.165397	2025-04-20 20:55:11.167286	200	1
1599	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:11.342646	2025-04-20 20:55:11.347884	200	5
1600	\N	/status	GET	2025-04-20 20:55:15.998783	2025-04-20 20:55:15.998784	200	0
1601	\N	/status	GET	2025-04-20 20:55:16.175763	2025-04-20 20:55:16.175765	200	0
1602	\N	/api/statistics/global	GET	2025-04-20 20:55:16.194269	2025-04-20 20:55:16.196715	200	2
1603	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:16.391612	2025-04-20 20:55:16.416123	200	24
1604	\N	/status	GET	2025-04-20 20:55:20.974474	2025-04-20 20:55:20.974476	200	0
1605	\N	/status	GET	2025-04-20 20:55:21.157326	2025-04-20 20:55:21.157327	200	0
1610	\N	/api/statistics/global	GET	2025-04-20 20:55:26.177454	2025-04-20 20:55:26.182063	200	4
1611	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:26.379142	2025-04-20 20:55:26.381233	200	2
1612	\N	/status	GET	2025-04-20 20:55:30.986637	2025-04-20 20:55:30.986639	200	0
1613	\N	/status	GET	2025-04-20 20:55:31.155878	2025-04-20 20:55:31.15588	200	0
1614	\N	/api/statistics/global	GET	2025-04-20 20:55:31.174332	2025-04-20 20:55:31.177304	200	2
1615	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:31.350035	2025-04-20 20:55:31.352583	200	2
1616	\N	/status	GET	2025-04-20 20:55:35.973943	2025-04-20 20:55:35.973945	200	0
1617	\N	/status	GET	2025-04-20 20:55:36.142932	2025-04-20 20:55:36.142933	200	0
1618	\N	/api/statistics/global	GET	2025-04-20 20:55:36.161258	2025-04-20 20:55:36.162474	200	1
1619	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:36.346869	2025-04-20 20:55:36.359917	200	13
1620	\N	/status	GET	2025-04-20 20:55:40.982655	2025-04-20 20:55:40.982657	200	0
1621	\N	/status	GET	2025-04-20 20:55:41.149942	2025-04-20 20:55:41.149944	200	0
1645	\N	/status	GET	2025-04-20 20:56:11.174847	2025-04-20 20:56:11.174849	200	0
2765	\N	/status	GET	2025-04-20 23:55:29.049896	2025-04-20 23:55:29.049897	200	0
2766	\N	/status	GET	2025-04-20 23:55:29.254291	2025-04-20 23:55:29.254293	200	0
2767	\N	/api/statistics/global	GET	2025-04-20 23:55:29.273055	2025-04-20 23:55:29.275042	200	1
2768	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:29.48323	2025-04-20 23:55:29.486109	200	2
2769	\N	/status	GET	2025-04-20 23:55:34.024779	2025-04-20 23:55:34.024781	200	0
2770	\N	/status	GET	2025-04-20 23:55:34.201791	2025-04-20 23:55:34.201792	200	0
2771	\N	/api/statistics/global	GET	2025-04-20 23:55:34.211459	2025-04-20 23:55:34.213905	200	2
2772	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:34.413023	2025-04-20 23:55:34.416528	200	3
2773	\N	/status	GET	2025-04-20 23:55:39.048549	2025-04-20 23:55:39.04855	200	0
2774	\N	/status	GET	2025-04-20 23:55:39.221196	2025-04-20 23:55:39.221197	200	0
2775	\N	/api/statistics/global	GET	2025-04-20 23:55:39.239643	2025-04-20 23:55:39.248621	200	8
2776	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:39.450645	2025-04-20 23:55:39.453401	200	2
2777	\N	/status	GET	2025-04-20 23:55:44.069409	2025-04-20 23:55:44.069411	200	0
2778	\N	/status	GET	2025-04-20 23:55:44.238895	2025-04-20 23:55:44.238896	200	0
2779	\N	/api/statistics/global	GET	2025-04-20 23:55:44.262368	2025-04-20 23:55:44.278067	200	15
2780	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:44.475171	2025-04-20 23:55:44.488518	200	13
2781	\N	/status	GET	2025-04-20 23:55:49.041375	2025-04-20 23:55:49.04138	200	0
2782	\N	/status	GET	2025-04-20 23:55:49.218828	2025-04-20 23:55:49.218829	200	0
2783	\N	/api/statistics/global	GET	2025-04-20 23:55:49.293433	2025-04-20 23:55:49.327629	200	34
2784	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:49.506087	2025-04-20 23:55:49.509785	200	3
2785	\N	/status	GET	2025-04-20 23:55:54.039185	2025-04-20 23:55:54.039187	200	0
2786	\N	/status	GET	2025-04-20 23:55:54.220772	2025-04-20 23:55:54.220773	200	0
2787	\N	/api/statistics/global	GET	2025-04-20 23:55:54.238709	2025-04-20 23:55:54.241808	200	3
2788	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:54.448016	2025-04-20 23:55:54.452821	200	4
2789	\N	/status	GET	2025-04-20 23:55:59.045751	2025-04-20 23:55:59.045752	200	0
2790	\N	/status	GET	2025-04-20 23:55:59.21788	2025-04-20 23:55:59.217881	200	0
2811	\N	/api/statistics/global	GET	2025-04-20 23:56:24.224033	2025-04-20 23:56:24.233155	200	9
2812	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:24.457433	2025-04-20 23:56:24.459593	200	2
2813	\N	/status	GET	2025-04-20 23:56:29.045427	2025-04-20 23:56:29.045428	200	0
2814	\N	/status	GET	2025-04-20 23:56:29.215179	2025-04-20 23:56:29.215181	200	0
2815	\N	/api/statistics/global	GET	2025-04-20 23:56:29.236912	2025-04-20 23:56:29.238479	200	1
2816	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:29.426618	2025-04-20 23:56:29.431171	200	4
2817	\N	/status	GET	2025-04-20 23:56:34.026418	2025-04-20 23:56:34.02642	200	0
2818	\N	/status	GET	2025-04-20 23:56:34.199973	2025-04-20 23:56:34.199975	200	0
2819	\N	/api/statistics/global	GET	2025-04-20 23:56:34.218034	2025-04-20 23:56:34.221191	200	3
2820	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:34.410084	2025-04-20 23:56:34.412616	200	2
2821	\N	/status	GET	2025-04-20 23:56:39.025467	2025-04-20 23:56:39.025468	200	0
2822	\N	/status	GET	2025-04-20 23:56:39.194781	2025-04-20 23:56:39.194783	200	0
2823	\N	/api/statistics/global	GET	2025-04-20 23:56:39.212645	2025-04-20 23:56:39.215264	200	2
2824	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:39.408284	2025-04-20 23:56:39.410944	200	2
2825	\N	/status	GET	2025-04-20 23:56:44.035667	2025-04-20 23:56:44.035668	200	0
2826	\N	/status	GET	2025-04-20 23:56:44.205071	2025-04-20 23:56:44.205073	200	0
2827	\N	/api/statistics/global	GET	2025-04-20 23:56:44.222979	2025-04-20 23:56:44.224549	200	1
2828	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:44.416546	2025-04-20 23:56:44.419393	200	2
2829	\N	/status	GET	2025-04-20 23:56:49.044872	2025-04-20 23:56:49.044873	200	0
2830	\N	/status	GET	2025-04-20 23:56:49.215614	2025-04-20 23:56:49.215616	200	0
2831	\N	/api/statistics/global	GET	2025-04-20 23:56:49.231387	2025-04-20 23:56:49.233176	200	1
2832	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:49.431436	2025-04-20 23:56:49.433924	200	2
2833	\N	/status	GET	2025-04-20 23:56:54.052166	2025-04-20 23:56:54.052167	200	0
2834	\N	/status	GET	2025-04-20 23:56:54.221065	2025-04-20 23:56:54.221066	200	0
2921	\N	/status	GET	2025-04-20 23:58:43.121392	2025-04-20 23:58:43.121394	200	0
2926	\N	/api/statistics/global	GET	2025-04-20 23:58:48.196734	2025-04-20 23:58:48.20878	200	12
2927	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:48.421488	2025-04-20 23:58:48.424213	200	2
2928	\N	/status	GET	2025-04-20 23:58:52.945758	2025-04-20 23:58:52.94576	200	0
2929	\N	/status	GET	2025-04-20 23:58:53.113294	2025-04-20 23:58:53.113296	200	0
2930	\N	/api/statistics/global	GET	2025-04-20 23:58:53.132023	2025-04-20 23:58:53.13364	200	1
1606	\N	/api/statistics/global	GET	2025-04-20 20:55:21.173669	2025-04-20 20:55:21.178748	200	5
1607	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:21.375465	2025-04-20 20:55:21.377559	200	2
1608	\N	/status	GET	2025-04-20 20:55:25.982146	2025-04-20 20:55:25.982148	200	0
1609	\N	/status	GET	2025-04-20 20:55:26.168018	2025-04-20 20:55:26.16802	200	0
1622	\N	/api/statistics/global	GET	2025-04-20 20:55:41.16863	2025-04-20 20:55:41.192458	200	23
1623	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:41.46198	2025-04-20 20:55:41.463954	200	1
1624	\N	/status	GET	2025-04-20 20:55:45.965034	2025-04-20 20:55:45.965036	200	0
1625	\N	/status	GET	2025-04-20 20:55:46.19115	2025-04-20 20:55:46.191151	200	0
1626	\N	/api/statistics/global	GET	2025-04-20 20:55:46.203853	2025-04-20 20:55:46.205279	200	1
1627	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:46.38144	2025-04-20 20:55:46.383389	200	1
1628	\N	/status	GET	2025-04-20 20:55:50.973111	2025-04-20 20:55:50.973112	200	0
1629	\N	/status	GET	2025-04-20 20:55:51.139664	2025-04-20 20:55:51.139666	200	0
1630	\N	/api/statistics/global	GET	2025-04-20 20:55:51.157617	2025-04-20 20:55:51.158926	200	1
1631	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:51.336976	2025-04-20 20:55:51.339903	200	2
1632	\N	/status	GET	2025-04-20 20:55:55.965118	2025-04-20 20:55:55.96512	200	0
1633	\N	/status	GET	2025-04-20 20:55:56.132689	2025-04-20 20:55:56.132689	200	0
1634	\N	/api/statistics/global	GET	2025-04-20 20:55:56.151199	2025-04-20 20:55:56.153418	200	2
1635	\N	/api/statistics/endpoints	GET	2025-04-20 20:55:56.322506	2025-04-20 20:55:56.326217	200	3
1636	\N	/status	GET	2025-04-20 20:56:00.98471	2025-04-20 20:56:00.984712	200	0
1637	\N	/status	GET	2025-04-20 20:56:01.152902	2025-04-20 20:56:01.152904	200	0
1638	\N	/api/statistics/global	GET	2025-04-20 20:56:01.180334	2025-04-20 20:56:01.182188	200	1
1639	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:01.354462	2025-04-20 20:56:01.357143	200	2
1640	\N	/status	GET	2025-04-20 20:56:05.962939	2025-04-20 20:56:05.962941	200	0
1641	\N	/status	GET	2025-04-20 20:56:06.128464	2025-04-20 20:56:06.128465	200	0
1642	\N	/api/statistics/global	GET	2025-04-20 20:56:06.146652	2025-04-20 20:56:06.148206	200	1
1643	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:06.318094	2025-04-20 20:56:06.322591	200	4
1644	\N	/status	GET	2025-04-20 20:56:10.9958	2025-04-20 20:56:10.995802	200	0
1646	\N	/api/statistics/global	GET	2025-04-20 20:56:11.174844	2025-04-20 20:56:11.177825	200	2
1647	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:11.366635	2025-04-20 20:56:11.368996	200	2
1648	\N	/status	GET	2025-04-20 20:56:16.267781	2025-04-20 20:56:16.267782	200	0
1649	\N	/status	GET	2025-04-20 20:56:16.456228	2025-04-20 20:56:16.45623	200	0
1650	\N	/api/statistics/global	GET	2025-04-20 20:56:16.475604	2025-04-20 20:56:16.476865	200	1
1651	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:16.647728	2025-04-20 20:56:16.649583	200	1
1652	\N	/status	GET	2025-04-20 20:56:20.972005	2025-04-20 20:56:20.972007	200	0
1653	\N	/status	GET	2025-04-20 20:56:21.138306	2025-04-20 20:56:21.138307	200	0
1654	\N	/api/statistics/global	GET	2025-04-20 20:56:21.156157	2025-04-20 20:56:21.158026	200	1
1655	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:21.32896	2025-04-20 20:56:21.339481	200	10
1656	\N	/status	GET	2025-04-20 20:56:25.97779	2025-04-20 20:56:25.977792	200	0
1657	\N	/status	GET	2025-04-20 20:56:26.145148	2025-04-20 20:56:26.14515	200	0
1658	\N	/api/statistics/global	GET	2025-04-20 20:56:26.163053	2025-04-20 20:56:26.167244	200	4
1659	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:26.374971	2025-04-20 20:56:26.3772	200	2
1660	\N	/status	GET	2025-04-20 20:56:30.9846	2025-04-20 20:56:30.984602	200	0
1661	\N	/status	GET	2025-04-20 20:56:31.154121	2025-04-20 20:56:31.154122	200	0
1662	\N	/api/statistics/global	GET	2025-04-20 20:56:31.17169	2025-04-20 20:56:31.17422	200	2
1663	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:31.352324	2025-04-20 20:56:31.355511	200	3
1664	\N	/status	GET	2025-04-20 20:56:35.966245	2025-04-20 20:56:35.966246	200	0
1665	\N	/status	GET	2025-04-20 20:56:36.149871	2025-04-20 20:56:36.149872	200	0
1666	\N	/api/statistics/global	GET	2025-04-20 20:56:36.16847	2025-04-20 20:56:36.169781	200	1
1667	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:36.340561	2025-04-20 20:56:36.343284	200	2
1668	\N	/status	GET	2025-04-20 20:56:40.979909	2025-04-20 20:56:40.97991	200	0
1669	\N	/status	GET	2025-04-20 20:56:41.149435	2025-04-20 20:56:41.149436	200	0
1670	\N	/api/statistics/global	GET	2025-04-20 20:56:41.176764	2025-04-20 20:56:41.184497	200	7
1671	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:41.373467	2025-04-20 20:56:41.376749	200	3
1672	\N	/status	GET	2025-04-20 20:56:46.005583	2025-04-20 20:56:46.005584	200	0
1673	\N	/status	GET	2025-04-20 20:56:46.209822	2025-04-20 20:56:46.209824	200	0
1674	\N	/api/statistics/global	GET	2025-04-20 20:56:46.227589	2025-04-20 20:56:46.265727	200	38
1675	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:46.487565	2025-04-20 20:56:46.490404	200	2
1676	\N	/status	GET	2025-04-20 20:56:50.972818	2025-04-20 20:56:50.97282	200	0
1677	\N	/status	GET	2025-04-20 20:56:51.141723	2025-04-20 20:56:51.141725	200	0
1678	\N	/api/statistics/global	GET	2025-04-20 20:56:51.160676	2025-04-20 20:56:51.161933	200	1
1679	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:51.360882	2025-04-20 20:56:51.363749	200	2
1680	\N	/status	GET	2025-04-20 20:56:55.978727	2025-04-20 20:56:55.978729	200	0
1681	\N	/status	GET	2025-04-20 20:56:56.190349	2025-04-20 20:56:56.19035	200	0
1682	\N	/api/statistics/global	GET	2025-04-20 20:56:56.209025	2025-04-20 20:56:56.216372	200	7
1683	\N	/api/statistics/endpoints	GET	2025-04-20 20:56:56.418494	2025-04-20 20:56:56.42049	200	1
1684	\N	/status	GET	2025-04-20 20:57:00.988317	2025-04-20 20:57:00.988318	200	0
1685	\N	/status	GET	2025-04-20 20:57:01.161296	2025-04-20 20:57:01.161298	200	0
1686	\N	/api/statistics/global	GET	2025-04-20 20:57:01.179303	2025-04-20 20:57:01.181088	200	1
1687	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:01.360975	2025-04-20 20:57:01.363347	200	2
1688	\N	/status	GET	2025-04-20 20:57:05.983465	2025-04-20 20:57:05.983467	200	0
1689	\N	/status	GET	2025-04-20 20:57:06.151233	2025-04-20 20:57:06.151235	200	0
1690	\N	/api/statistics/global	GET	2025-04-20 20:57:06.16903	2025-04-20 20:57:06.171528	200	2
1691	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:06.342771	2025-04-20 20:57:06.356094	200	13
1692	\N	/status	GET	2025-04-20 20:57:11.006462	2025-04-20 20:57:11.006463	200	0
1693	\N	/status	GET	2025-04-20 20:57:11.17605	2025-04-20 20:57:11.176051	200	0
1694	\N	/api/statistics/global	GET	2025-04-20 20:57:11.195178	2025-04-20 20:57:11.204587	200	9
1695	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:11.403637	2025-04-20 20:57:11.405652	200	2
1696	\N	/status	GET	2025-04-20 20:57:16.001596	2025-04-20 20:57:16.001598	200	0
1697	\N	/status	GET	2025-04-20 20:57:16.171423	2025-04-20 20:57:16.171425	200	0
1698	\N	/api/statistics/global	GET	2025-04-20 20:57:16.191339	2025-04-20 20:57:16.195384	200	4
1699	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:16.381616	2025-04-20 20:57:16.383888	200	2
1700	\N	/status	GET	2025-04-20 20:57:20.998422	2025-04-20 20:57:20.998423	200	0
1701	\N	/status	GET	2025-04-20 20:57:21.171196	2025-04-20 20:57:21.171198	200	0
1702	\N	/api/statistics/global	GET	2025-04-20 20:57:21.189214	2025-04-20 20:57:21.190894	200	1
1703	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:21.361829	2025-04-20 20:57:21.363958	200	2
2791	\N	/api/statistics/global	GET	2025-04-20 23:55:59.235493	2025-04-20 23:55:59.247124	200	11
2792	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:59.446673	2025-04-20 23:55:59.44921	200	2
2793	\N	/status	GET	2025-04-20 23:56:04.043551	2025-04-20 23:56:04.043553	200	0
2794	\N	/status	GET	2025-04-20 23:56:04.212705	2025-04-20 23:56:04.212707	200	0
2795	\N	/api/statistics/global	GET	2025-04-20 23:56:04.231914	2025-04-20 23:56:04.237765	200	5
2796	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:04.43452	2025-04-20 23:56:04.438394	200	3
2797	\N	/status	GET	2025-04-20 23:56:09.043619	2025-04-20 23:56:09.043621	200	0
2798	\N	/status	GET	2025-04-20 23:56:09.215917	2025-04-20 23:56:09.215918	200	0
2799	\N	/api/statistics/global	GET	2025-04-20 23:56:09.233721	2025-04-20 23:56:09.242276	200	8
2800	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:09.442376	2025-04-20 23:56:09.445057	200	2
2801	\N	/status	GET	2025-04-20 23:56:14.048299	2025-04-20 23:56:14.048301	200	0
2802	\N	/status	GET	2025-04-20 23:56:14.256479	2025-04-20 23:56:14.25648	200	0
2803	\N	/api/statistics/global	GET	2025-04-20 23:56:14.275332	2025-04-20 23:56:14.285548	200	10
2804	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:14.509312	2025-04-20 23:56:14.512563	200	3
2805	\N	/status	GET	2025-04-20 23:56:19.056875	2025-04-20 23:56:19.056876	200	0
2806	\N	/status	GET	2025-04-20 23:56:19.225392	2025-04-20 23:56:19.225393	200	0
2807	\N	/api/statistics/global	GET	2025-04-20 23:56:19.243833	2025-04-20 23:56:19.247741	200	3
2808	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:19.455755	2025-04-20 23:56:19.458266	200	2
2809	\N	/status	GET	2025-04-20 23:56:24.035266	2025-04-20 23:56:24.035268	200	0
1706	\N	/api/statistics/global	GET	2025-04-20 20:57:26.158788	2025-04-20 20:57:26.16177	200	2
1707	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:26.331554	2025-04-20 20:57:26.333687	200	2
1708	\N	/status	GET	2025-04-20 20:57:30.977207	2025-04-20 20:57:30.977209	200	0
1709	\N	/status	GET	2025-04-20 20:57:31.149537	2025-04-20 20:57:31.149539	200	0
1710	\N	/api/statistics/global	GET	2025-04-20 20:57:31.168175	2025-04-20 20:57:31.173126	200	4
1711	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:31.363271	2025-04-20 20:57:31.365171	200	1
1712	\N	/status	GET	2025-04-20 20:57:35.990855	2025-04-20 20:57:35.990857	200	0
1713	\N	/status	GET	2025-04-20 20:57:36.189101	2025-04-20 20:57:36.189102	200	0
1714	\N	/api/statistics/global	GET	2025-04-20 20:57:36.207406	2025-04-20 20:57:36.219001	200	11
1715	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:36.408399	2025-04-20 20:57:36.41026	200	1
1716	\N	/status	GET	2025-04-20 20:57:40.985995	2025-04-20 20:57:40.985996	200	0
1717	\N	/status	GET	2025-04-20 20:57:41.156981	2025-04-20 20:57:41.156982	200	0
1718	\N	/api/statistics/global	GET	2025-04-20 20:57:41.171232	2025-04-20 20:57:41.172567	200	1
1719	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:41.346968	2025-04-20 20:57:41.349191	200	2
1720	\N	/status	GET	2025-04-20 20:57:45.965563	2025-04-20 20:57:45.965564	200	0
1721	\N	/status	GET	2025-04-20 20:57:46.132742	2025-04-20 20:57:46.132744	200	0
1722	\N	/api/statistics/global	GET	2025-04-20 20:57:46.150218	2025-04-20 20:57:46.152164	200	1
1723	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:46.321829	2025-04-20 20:57:46.324242	200	2
1724	\N	/status	GET	2025-04-20 20:57:51.012406	2025-04-20 20:57:51.012408	200	0
1725	\N	/status	GET	2025-04-20 20:57:51.183242	2025-04-20 20:57:51.183244	200	0
1726	\N	/api/statistics/global	GET	2025-04-20 20:57:51.210081	2025-04-20 20:57:51.211899	200	1
1727	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:51.452685	2025-04-20 20:57:51.45531	200	2
1728	\N	/status	GET	2025-04-20 20:57:55.993941	2025-04-20 20:57:55.993943	200	0
1729	\N	/status	GET	2025-04-20 20:57:56.162427	2025-04-20 20:57:56.162428	200	0
1730	\N	/api/statistics/global	GET	2025-04-20 20:57:56.162427	2025-04-20 20:57:56.164165	200	1
1731	\N	/api/statistics/endpoints	GET	2025-04-20 20:57:56.345701	2025-04-20 20:57:56.348394	200	2
1732	\N	/status	GET	2025-04-20 20:58:01.003276	2025-04-20 20:58:01.003278	200	0
1733	\N	/status	GET	2025-04-20 20:58:01.173867	2025-04-20 20:58:01.173868	200	0
1734	\N	/api/statistics/global	GET	2025-04-20 20:58:01.191637	2025-04-20 20:58:01.19366	200	2
1735	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:01.364918	2025-04-20 20:58:01.36891	200	3
1736	\N	/status	GET	2025-04-20 20:58:06.010963	2025-04-20 20:58:06.010964	200	0
1737	\N	/status	GET	2025-04-20 20:58:06.194739	2025-04-20 20:58:06.194741	200	0
1738	\N	/api/statistics/global	GET	2025-04-20 20:58:06.19941	2025-04-20 20:58:06.201291	200	1
1739	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:06.381853	2025-04-20 20:58:06.384095	200	2
1740	\N	/status	GET	2025-04-20 20:58:10.979951	2025-04-20 20:58:10.979952	200	0
1741	\N	/status	GET	2025-04-20 20:58:11.151039	2025-04-20 20:58:11.151041	200	0
1742	\N	/api/statistics/global	GET	2025-04-20 20:58:11.170048	2025-04-20 20:58:11.173899	200	3
1743	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:11.469796	2025-04-20 20:58:11.472102	200	2
1744	\N	/status	GET	2025-04-20 20:58:15.980545	2025-04-20 20:58:15.980547	200	0
1745	\N	/status	GET	2025-04-20 20:58:16.228819	2025-04-20 20:58:16.228821	200	0
1746	\N	/api/statistics/global	GET	2025-04-20 20:58:16.246383	2025-04-20 20:58:16.270462	200	24
1747	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:16.648142	2025-04-20 20:58:16.693228	200	45
1748	\N	/status	GET	2025-04-20 20:58:20.987203	2025-04-20 20:58:20.987204	200	0
1749	\N	/status	GET	2025-04-20 20:58:21.161412	2025-04-20 20:58:21.161414	200	0
1750	\N	/api/statistics/global	GET	2025-04-20 20:58:21.179537	2025-04-20 20:58:21.186871	200	7
1751	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:21.378379	2025-04-20 20:58:21.380568	200	2
1752	\N	/status	GET	2025-04-20 20:58:25.997208	2025-04-20 20:58:25.997209	200	0
1753	\N	/status	GET	2025-04-20 20:58:26.173854	2025-04-20 20:58:26.173855	200	0
1754	\N	/api/statistics/global	GET	2025-04-20 20:58:26.173853	2025-04-20 20:58:26.175195	200	1
1755	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:26.364511	2025-04-20 20:58:26.36688	200	2
1756	\N	/status	GET	2025-04-20 20:58:30.981051	2025-04-20 20:58:30.981052	200	0
1757	\N	/status	GET	2025-04-20 20:58:31.147791	2025-04-20 20:58:31.147793	200	0
1758	\N	/api/statistics/global	GET	2025-04-20 20:58:31.165647	2025-04-20 20:58:31.167197	200	1
1759	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:31.342964	2025-04-20 20:58:31.344921	200	1
1760	\N	/status	GET	2025-04-20 20:58:35.980395	2025-04-20 20:58:35.980396	200	0
1761	\N	/status	GET	2025-04-20 20:58:36.152085	2025-04-20 20:58:36.152086	200	0
1762	\N	/api/statistics/global	GET	2025-04-20 20:58:36.169746	2025-04-20 20:58:36.178303	200	8
1763	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:36.362108	2025-04-20 20:58:36.37167	200	9
1764	\N	/status	GET	2025-04-20 20:58:40.986592	2025-04-20 20:58:40.986593	200	0
1765	\N	/status	GET	2025-04-20 20:58:41.156261	2025-04-20 20:58:41.156263	200	0
1766	\N	/api/statistics/global	GET	2025-04-20 20:58:41.174064	2025-04-20 20:58:41.180247	200	6
1767	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:41.367272	2025-04-20 20:58:41.369199	200	1
1768	\N	/status	GET	2025-04-20 20:58:45.992781	2025-04-20 20:58:45.992782	200	0
1769	\N	/status	GET	2025-04-20 20:58:46.173093	2025-04-20 20:58:46.173095	200	0
1770	\N	/api/statistics/global	GET	2025-04-20 20:58:46.173093	2025-04-20 20:58:46.184204	200	11
1771	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:46.370662	2025-04-20 20:58:46.37275	200	2
1772	\N	/status	GET	2025-04-20 20:58:50.996618	2025-04-20 20:58:50.99662	200	0
1773	\N	/status	GET	2025-04-20 20:58:51.167878	2025-04-20 20:58:51.16788	200	0
1774	\N	/api/statistics/global	GET	2025-04-20 20:58:51.1879	2025-04-20 20:58:51.195168	200	7
1775	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:51.398959	2025-04-20 20:58:51.40088	200	1
1776	\N	/status	GET	2025-04-20 20:58:56.046598	2025-04-20 20:58:56.0466	200	0
1777	\N	/status	GET	2025-04-20 20:58:56.214528	2025-04-20 20:58:56.214529	200	0
1778	\N	/api/statistics/global	GET	2025-04-20 20:58:56.232532	2025-04-20 20:58:56.236712	200	4
1779	\N	/api/statistics/endpoints	GET	2025-04-20 20:58:56.406285	2025-04-20 20:58:56.408833	200	2
1780	\N	/status	GET	2025-04-20 20:59:00.983616	2025-04-20 20:59:00.983617	200	0
1781	\N	/status	GET	2025-04-20 20:59:01.173654	2025-04-20 20:59:01.173656	200	0
1782	\N	/api/statistics/global	GET	2025-04-20 20:59:01.173583	2025-04-20 20:59:01.178937	200	5
1783	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:01.381438	2025-04-20 20:59:01.383371	200	1
1784	\N	/status	GET	2025-04-20 20:59:05.998135	2025-04-20 20:59:05.998137	200	0
1785	\N	/status	GET	2025-04-20 20:59:06.170867	2025-04-20 20:59:06.170869	200	0
1786	\N	/api/statistics/global	GET	2025-04-20 20:59:06.204265	2025-04-20 20:59:06.207088	200	2
1787	\N	/status	GET	2025-04-20 20:59:06.76953	2025-04-20 20:59:06.769531	200	0
1788	\N	/api/statistics/global	GET	2025-04-20 20:59:06.968348	2025-04-20 20:59:06.969574	200	1
1789	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:07.13697	2025-04-20 20:59:07.181252	200	44
1790	\N	/status	GET	2025-04-20 20:59:11.779101	2025-04-20 20:59:11.779103	200	0
1791	\N	/status	GET	2025-04-20 20:59:11.953613	2025-04-20 20:59:11.953615	200	0
1792	\N	/api/statistics/global	GET	2025-04-20 20:59:11.98636	2025-04-20 20:59:11.987659	200	1
1793	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:12.192901	2025-04-20 20:59:12.195638	200	2
1794	\N	/status	GET	2025-04-20 20:59:16.725326	2025-04-20 20:59:16.725327	200	0
1795	\N	/status	GET	2025-04-20 20:59:16.902639	2025-04-20 20:59:16.902641	200	0
1796	\N	/api/statistics/global	GET	2025-04-20 20:59:16.908667	2025-04-20 20:59:16.909968	200	1
1797	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:17.096882	2025-04-20 20:59:17.099204	200	2
1798	\N	/status	GET	2025-04-20 20:59:21.738342	2025-04-20 20:59:21.738344	200	0
1799	\N	/status	GET	2025-04-20 20:59:21.906887	2025-04-20 20:59:21.906889	200	0
1800	\N	/api/statistics/global	GET	2025-04-20 20:59:21.922656	2025-04-20 20:59:21.924058	200	1
1801	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:22.09571	2025-04-20 20:59:22.098224	200	2
1802	\N	/status	GET	2025-04-20 20:59:26.695775	2025-04-20 20:59:26.695777	200	0
1803	\N	/status	GET	2025-04-20 20:59:26.865885	2025-04-20 20:59:26.865887	200	0
1804	\N	/api/statistics/global	GET	2025-04-20 20:59:27.01395	2025-04-20 20:59:27.016652	200	2
1805	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:27.190175	2025-04-20 20:59:27.195976	200	5
1806	\N	/status	GET	2025-04-20 20:59:31.802682	2025-04-20 20:59:31.802684	200	0
1807	\N	/status	GET	2025-04-20 20:59:31.97146	2025-04-20 20:59:31.971462	200	0
1808	\N	/api/statistics/global	GET	2025-04-20 20:59:31.990084	2025-04-20 20:59:31.991353	200	1
1809	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:32.159265	2025-04-20 20:59:32.161944	200	2
1810	\N	/status	GET	2025-04-20 20:59:36.741706	2025-04-20 20:59:36.741707	200	0
1811	\N	/status	GET	2025-04-20 20:59:36.913131	2025-04-20 20:59:36.913133	200	0
1812	\N	/api/statistics/global	GET	2025-04-20 20:59:36.931711	2025-04-20 20:59:36.932919	200	1
1813	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:37.105762	2025-04-20 20:59:37.108276	200	2
1814	\N	/status	GET	2025-04-20 20:59:41.736645	2025-04-20 20:59:41.736646	200	0
1815	\N	/status	GET	2025-04-20 20:59:41.904433	2025-04-20 20:59:41.904434	200	0
1816	\N	/api/statistics/global	GET	2025-04-20 20:59:41.922729	2025-04-20 20:59:41.92446	200	1
1817	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:42.094426	2025-04-20 20:59:42.101674	200	7
1818	\N	/status	GET	2025-04-20 20:59:46.717511	2025-04-20 20:59:46.717525	200	0
1819	\N	/status	GET	2025-04-20 20:59:46.883856	2025-04-20 20:59:46.883858	200	0
1820	\N	/api/statistics/global	GET	2025-04-20 20:59:46.906613	2025-04-20 20:59:46.908183	200	1
1821	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:47.082952	2025-04-20 20:59:47.090664	200	7
1822	\N	/status	GET	2025-04-20 20:59:51.701354	2025-04-20 20:59:51.701355	200	0
1823	\N	/status	GET	2025-04-20 20:59:51.87925	2025-04-20 20:59:51.879251	200	0
1839	\N	/status	GET	2025-04-20 21:00:11.904396	2025-04-20 21:00:11.904398	200	0
1871	\N	/status	GET	2025-04-20 21:00:51.87958	2025-04-20 21:00:51.879582	200	0
1924	\N	/api/statistics/global	GET	2025-04-20 21:01:56.904277	2025-04-20 21:01:56.90563	200	1
1925	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:57.083983	2025-04-20 21:01:57.086283	200	2
1926	\N	/status	GET	2025-04-20 21:02:01.727534	2025-04-20 21:02:01.727535	200	0
1927	\N	/status	GET	2025-04-20 21:02:01.895068	2025-04-20 21:02:01.895069	200	0
1928	\N	/api/statistics/global	GET	2025-04-20 21:02:01.913747	2025-04-20 21:02:01.915846	200	2
1929	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:02.09615	2025-04-20 21:02:02.098351	200	2
1930	\N	/status	GET	2025-04-20 21:02:06.728551	2025-04-20 21:02:06.728552	200	0
1931	\N	/status	GET	2025-04-20 21:02:06.895186	2025-04-20 21:02:06.895187	200	0
1932	\N	/api/statistics/global	GET	2025-04-20 21:02:06.91328	2025-04-20 21:02:06.914844	200	1
1933	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:07.10428	2025-04-20 21:02:07.107148	200	2
1934	\N	/status	GET	2025-04-20 21:02:11.738937	2025-04-20 21:02:11.738938	200	0
1935	\N	/status	GET	2025-04-20 21:02:11.918146	2025-04-20 21:02:11.918147	200	0
1936	\N	/api/statistics/global	GET	2025-04-20 21:02:11.936116	2025-04-20 21:02:11.93798	200	1
1937	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:12.10845	2025-04-20 21:02:12.15741	200	48
1938	\N	/status	GET	2025-04-20 21:02:16.714089	2025-04-20 21:02:16.714091	200	0
1939	\N	/status	GET	2025-04-20 21:02:16.881811	2025-04-20 21:02:16.881813	200	0
1940	\N	/api/statistics/global	GET	2025-04-20 21:02:16.902374	2025-04-20 21:02:16.904309	200	1
1941	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:17.094153	2025-04-20 21:02:17.098982	200	4
1942	\N	/status	GET	2025-04-20 21:02:21.702844	2025-04-20 21:02:21.702846	200	0
1943	\N	/status	GET	2025-04-20 21:02:21.868014	2025-04-20 21:02:21.868015	200	0
1944	\N	/api/statistics/global	GET	2025-04-20 21:02:21.886506	2025-04-20 21:02:21.887632	200	1
1945	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:22.0598	2025-04-20 21:02:22.062838	200	3
2810	\N	/status	GET	2025-04-20 23:56:24.205609	2025-04-20 23:56:24.205611	200	0
2835	\N	/api/statistics/global	GET	2025-04-20 23:56:54.238843	2025-04-20 23:56:54.243632	200	4
2836	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:54.453175	2025-04-20 23:56:54.456578	200	3
2837	\N	/status	GET	2025-04-20 23:56:59.054566	2025-04-20 23:56:59.054567	200	0
2838	\N	/status	GET	2025-04-20 23:56:59.231165	2025-04-20 23:56:59.231166	200	0
2839	\N	/api/statistics/global	GET	2025-04-20 23:56:59.24961	2025-04-20 23:56:59.254097	200	4
2840	\N	/api/statistics/endpoints	GET	2025-04-20 23:56:59.45034	2025-04-20 23:56:59.454151	200	3
2841	\N	/status	GET	2025-04-20 23:57:02.937011	2025-04-20 23:57:02.937012	200	0
2842	\N	/api/statistics/global	GET	2025-04-20 23:57:03.105041	2025-04-20 23:57:03.111887	200	6
2843	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:03.295276	2025-04-20 23:57:03.307515	200	12
2844	\N	/status	GET	2025-04-20 23:57:08.044152	2025-04-20 23:57:08.044161	200	0
2845	\N	/status	GET	2025-04-20 23:57:08.219441	2025-04-20 23:57:08.219442	200	0
2846	\N	/api/statistics/global	GET	2025-04-20 23:57:08.237157	2025-04-20 23:57:08.239067	200	1
2847	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:08.424333	2025-04-20 23:57:08.427004	200	2
2848	\N	/status	GET	2025-04-20 23:57:12.942108	2025-04-20 23:57:12.942109	200	0
2849	\N	/status	GET	2025-04-20 23:57:13.110447	2025-04-20 23:57:13.110448	200	0
2850	\N	/api/statistics/global	GET	2025-04-20 23:57:13.139176	2025-04-20 23:57:13.149064	200	9
2851	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:13.419911	2025-04-20 23:57:13.422386	200	2
2852	\N	/status	GET	2025-04-20 23:57:17.943863	2025-04-20 23:57:17.943864	200	0
2853	\N	/status	GET	2025-04-20 23:57:18.1175	2025-04-20 23:57:18.117502	200	0
2854	\N	/api/statistics/global	GET	2025-04-20 23:57:18.133255	2025-04-20 23:57:18.138288	200	5
2855	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:18.338266	2025-04-20 23:57:18.340731	200	2
2856	\N	/status	GET	2025-04-20 23:57:23.035899	2025-04-20 23:57:23.035901	200	0
2857	\N	/status	GET	2025-04-20 23:57:23.220664	2025-04-20 23:57:23.220666	200	0
2858	\N	/api/statistics/global	GET	2025-04-20 23:57:23.239265	2025-04-20 23:57:23.251376	200	12
2859	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:23.439207	2025-04-20 23:57:23.443963	200	4
2931	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:53.318585	2025-04-20 23:58:53.321227	200	2
2932	\N	/status	GET	2025-04-20 23:58:58.020927	2025-04-20 23:58:58.020929	200	0
2933	\N	/status	GET	2025-04-20 23:58:58.190772	2025-04-20 23:58:58.190774	200	0
2934	\N	/api/statistics/global	GET	2025-04-20 23:58:58.208701	2025-04-20 23:58:58.211336	200	2
2935	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:58.385671	2025-04-20 23:58:58.388269	200	2
2936	\N	/status	GET	2025-04-20 23:59:03.039551	2025-04-20 23:59:03.039553	200	0
2937	\N	/status	GET	2025-04-20 23:59:03.214912	2025-04-20 23:59:03.214914	200	0
2938	\N	/api/statistics/global	GET	2025-04-20 23:59:03.233263	2025-04-20 23:59:03.239627	200	6
2939	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:03.437753	2025-04-20 23:59:03.440136	200	2
2940	\N	/status	GET	2025-04-20 23:59:07.943457	2025-04-20 23:59:07.943459	200	0
2941	\N	/status	GET	2025-04-20 23:59:08.111976	2025-04-20 23:59:08.111978	200	0
2942	\N	/api/statistics/global	GET	2025-04-20 23:59:08.129703	2025-04-20 23:59:08.134942	200	5
2943	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:08.333433	2025-04-20 23:59:08.346626	200	13
2944	\N	/status	GET	2025-04-20 23:59:13.041693	2025-04-20 23:59:13.041695	200	0
2945	\N	/status	GET	2025-04-20 23:59:13.213613	2025-04-20 23:59:13.213614	200	0
2946	\N	/api/statistics/global	GET	2025-04-20 23:59:13.232086	2025-04-20 23:59:13.234034	200	1
2947	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:13.43764	2025-04-20 23:59:13.440977	200	3
2948	\N	/status	GET	2025-04-20 23:59:17.941283	2025-04-20 23:59:17.941285	200	0
2949	\N	/status	GET	2025-04-20 23:59:18.243382	2025-04-20 23:59:18.243383	200	0
2950	\N	/api/statistics/global	GET	2025-04-20 23:59:18.261831	2025-04-20 23:59:18.263599	200	1
2951	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:18.449811	2025-04-20 23:59:18.453244	200	3
2952	\N	/status	GET	2025-04-20 23:59:22.939875	2025-04-20 23:59:22.939877	200	0
2953	\N	/status	GET	2025-04-20 23:59:23.107782	2025-04-20 23:59:23.107784	200	0
2954	\N	/api/statistics/global	GET	2025-04-20 23:59:23.126332	2025-04-20 23:59:23.131326	200	4
2955	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:23.320061	2025-04-20 23:59:23.348414	200	28
2956	\N	/status	GET	2025-04-20 23:59:28.018438	2025-04-20 23:59:28.01844	200	0
2957	\N	/status	GET	2025-04-20 23:59:28.187833	2025-04-20 23:59:28.187835	200	0
2958	\N	/api/statistics/global	GET	2025-04-20 23:59:28.20591	2025-04-20 23:59:28.207789	200	1
2959	\N	/api/statistics/endpoints	GET	2025-04-20 23:59:28.390108	2025-04-20 23:59:28.392568	200	2
2960	\N	/api/washingmachines	GET	2025-04-20 23:59:32.98733	2025-04-20 23:59:33.697633	200	710
2961	\N	/api/washingmachines	GET	2025-04-20 23:59:33.864812	2025-04-20 23:59:34.042761	200	177
2962	\N	/status	GET	2025-04-21 00:01:26.299044	2025-04-21 00:01:26.299046	200	0
2963	\N	/status	GET	2025-04-21 00:01:26.467905	2025-04-21 00:01:26.467906	200	0
2964	\N	/api/statistics/global	GET	2025-04-21 00:01:26.486643	2025-04-21 00:01:26.48857	200	1
2965	\N	/api/statistics/global	GET	2025-04-21 00:01:26.65749	2025-04-21 00:01:26.659626	200	2
2966	\N	/api/statistics/endpoints	GET	2025-04-21 00:01:26.676002	2025-04-21 00:01:26.679524	200	3
2967	\N	/api/statistics/endpoints	GET	2025-04-21 00:01:26.847575	2025-04-21 00:01:26.851254	200	3
2968	\N	/api/restaurant	GET	2025-04-21 00:01:34.161532	2025-04-21 00:01:35.090249	200	928
2973	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 00:15:53.979149	2025-04-21 00:15:53.979239	304	0
1824	\N	/api/statistics/global	GET	2025-04-20 20:59:51.897446	2025-04-20 20:59:51.907284	200	9
1825	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:52.085501	2025-04-20 20:59:52.092613	200	7
1826	\N	/status	GET	2025-04-20 20:59:56.733527	2025-04-20 20:59:56.733529	200	0
1827	\N	/status	GET	2025-04-20 20:59:56.908815	2025-04-20 20:59:56.908817	200	0
1828	\N	/api/statistics/global	GET	2025-04-20 20:59:56.934516	2025-04-20 20:59:56.935621	200	1
1829	\N	/api/statistics/endpoints	GET	2025-04-20 20:59:57.104181	2025-04-20 20:59:57.106888	200	2
1830	\N	/status	GET	2025-04-20 21:00:01.7197	2025-04-20 21:00:01.719701	200	0
1831	\N	/status	GET	2025-04-20 21:00:01.885554	2025-04-20 21:00:01.885555	200	0
1832	\N	/api/statistics/global	GET	2025-04-20 21:00:01.905707	2025-04-20 21:00:01.907433	200	1
1833	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:02.09618	2025-04-20 21:00:02.103282	200	7
1834	\N	/status	GET	2025-04-20 21:00:06.701546	2025-04-20 21:00:06.701548	200	0
1835	\N	/status	GET	2025-04-20 21:00:06.872864	2025-04-20 21:00:06.872865	200	0
1836	\N	/api/statistics/global	GET	2025-04-20 21:00:06.891993	2025-04-20 21:00:06.893305	200	1
1837	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:07.109783	2025-04-20 21:00:07.119034	200	9
1838	\N	/status	GET	2025-04-20 21:00:11.729811	2025-04-20 21:00:11.729813	200	0
1840	\N	/api/statistics/global	GET	2025-04-20 21:00:11.904404	2025-04-20 21:00:11.905475	200	1
1841	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:12.083575	2025-04-20 21:00:12.093174	200	9
1842	\N	/status	GET	2025-04-20 21:00:16.726072	2025-04-20 21:00:16.726073	200	0
1843	\N	/status	GET	2025-04-20 21:00:16.894097	2025-04-20 21:00:16.894099	200	0
1844	\N	/api/statistics/global	GET	2025-04-20 21:00:16.912295	2025-04-20 21:00:16.913494	200	1
1845	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:17.086149	2025-04-20 21:00:17.0887	200	2
1846	\N	/status	GET	2025-04-20 21:00:21.703026	2025-04-20 21:00:21.703028	200	0
1847	\N	/status	GET	2025-04-20 21:00:21.871403	2025-04-20 21:00:21.871405	200	0
1848	\N	/api/statistics/global	GET	2025-04-20 21:00:21.900364	2025-04-20 21:00:21.902128	200	1
1849	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:22.079825	2025-04-20 21:00:22.089355	200	9
1850	\N	/status	GET	2025-04-20 21:00:26.718646	2025-04-20 21:00:26.718647	200	0
1851	\N	/status	GET	2025-04-20 21:00:26.883492	2025-04-20 21:00:26.883493	200	0
1852	\N	/api/statistics/global	GET	2025-04-20 21:00:26.901763	2025-04-20 21:00:26.903349	200	1
1853	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:27.077913	2025-04-20 21:00:27.081306	200	3
1854	\N	/status	GET	2025-04-20 21:00:31.7139	2025-04-20 21:00:31.713901	200	0
1855	\N	/status	GET	2025-04-20 21:00:31.88433	2025-04-20 21:00:31.884332	200	0
1856	\N	/api/statistics/global	GET	2025-04-20 21:00:31.902093	2025-04-20 21:00:31.903934	200	1
1857	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:32.078542	2025-04-20 21:00:32.084269	200	5
1858	\N	/status	GET	2025-04-20 21:00:36.707946	2025-04-20 21:00:36.707946	200	0
1859	\N	/status	GET	2025-04-20 21:00:36.874499	2025-04-20 21:00:36.874501	200	0
1860	\N	/api/statistics/global	GET	2025-04-20 21:00:36.894927	2025-04-20 21:00:36.896779	200	1
1861	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:37.077286	2025-04-20 21:00:37.079747	200	2
1862	\N	/status	GET	2025-04-20 21:00:41.724873	2025-04-20 21:00:41.724875	200	0
1863	\N	/status	GET	2025-04-20 21:00:41.893878	2025-04-20 21:00:41.893879	200	0
1864	\N	/api/statistics/global	GET	2025-04-20 21:00:41.912195	2025-04-20 21:00:41.913865	200	1
1865	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:42.083841	2025-04-20 21:00:42.088272	200	4
1866	\N	/status	GET	2025-04-20 21:00:46.733722	2025-04-20 21:00:46.733723	200	0
1867	\N	/status	GET	2025-04-20 21:00:46.955448	2025-04-20 21:00:46.95545	200	0
1868	\N	/api/statistics/global	GET	2025-04-20 21:00:46.97337	2025-04-20 21:00:46.974832	200	1
1869	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:47.150372	2025-04-20 21:00:47.153823	200	3
1870	\N	/status	GET	2025-04-20 21:00:51.710143	2025-04-20 21:00:51.710144	200	0
1872	\N	/api/statistics/global	GET	2025-04-20 21:00:51.87958	2025-04-20 21:00:51.883848	200	4
1873	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:52.06292	2025-04-20 21:00:52.071829	200	8
1874	\N	/status	GET	2025-04-20 21:00:56.735045	2025-04-20 21:00:56.735046	200	0
1875	\N	/status	GET	2025-04-20 21:00:56.902414	2025-04-20 21:00:56.902415	200	0
1876	\N	/api/statistics/global	GET	2025-04-20 21:00:56.920901	2025-04-20 21:00:56.922098	200	1
1877	\N	/api/statistics/endpoints	GET	2025-04-20 21:00:57.090332	2025-04-20 21:00:57.093459	200	3
1878	\N	/status	GET	2025-04-20 21:01:01.728851	2025-04-20 21:01:01.728853	200	0
1879	\N	/status	GET	2025-04-20 21:01:01.933868	2025-04-20 21:01:01.93387	200	0
1880	\N	/api/statistics/global	GET	2025-04-20 21:01:01.952195	2025-04-20 21:01:01.953467	200	1
1881	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:02.130674	2025-04-20 21:01:02.137907	200	7
1882	\N	/status	GET	2025-04-20 21:01:06.752333	2025-04-20 21:01:06.752334	200	0
1883	\N	/status	GET	2025-04-20 21:01:06.919212	2025-04-20 21:01:06.919213	200	0
1884	\N	/api/statistics/global	GET	2025-04-20 21:01:06.937769	2025-04-20 21:01:06.939139	200	1
1885	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:07.111178	2025-04-20 21:01:07.116967	200	5
1886	\N	/status	GET	2025-04-20 21:01:11.697519	2025-04-20 21:01:11.697521	200	0
1887	\N	/status	GET	2025-04-20 21:01:11.865745	2025-04-20 21:01:11.865747	200	0
1888	\N	/api/statistics/global	GET	2025-04-20 21:01:11.885305	2025-04-20 21:01:11.886488	200	1
1889	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:12.057355	2025-04-20 21:01:12.062731	200	5
1890	\N	/status	GET	2025-04-20 21:01:16.723879	2025-04-20 21:01:16.72388	200	0
1891	\N	/status	GET	2025-04-20 21:01:16.891051	2025-04-20 21:01:16.891052	200	0
1892	\N	/api/statistics/global	GET	2025-04-20 21:01:16.909033	2025-04-20 21:01:16.910923	200	1
1893	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:17.088606	2025-04-20 21:01:17.091613	200	3
1894	\N	/status	GET	2025-04-20 21:01:21.730025	2025-04-20 21:01:21.730026	200	0
1895	\N	/status	GET	2025-04-20 21:01:21.897514	2025-04-20 21:01:21.897516	200	0
1896	\N	/api/statistics/global	GET	2025-04-20 21:01:21.916155	2025-04-20 21:01:21.917548	200	1
1897	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:22.087113	2025-04-20 21:01:22.094248	200	7
1898	\N	/status	GET	2025-04-20 21:01:26.729065	2025-04-20 21:01:26.729066	200	0
1899	\N	/status	GET	2025-04-20 21:01:26.895942	2025-04-20 21:01:26.895943	200	0
1900	\N	/api/statistics/global	GET	2025-04-20 21:01:26.914334	2025-04-20 21:01:26.915548	200	1
1901	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:27.092062	2025-04-20 21:01:27.095093	200	3
1902	\N	/status	GET	2025-04-20 21:01:31.717473	2025-04-20 21:01:31.717474	200	0
1903	\N	/status	GET	2025-04-20 21:01:31.884793	2025-04-20 21:01:31.884794	200	0
1904	\N	/api/statistics/global	GET	2025-04-20 21:01:31.912578	2025-04-20 21:01:31.91389	200	1
1905	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:32.086103	2025-04-20 21:01:32.090179	200	4
1906	\N	/status	GET	2025-04-20 21:01:36.696787	2025-04-20 21:01:36.696788	200	0
1907	\N	/status	GET	2025-04-20 21:01:36.916281	2025-04-20 21:01:36.916283	200	0
1908	\N	/api/statistics/global	GET	2025-04-20 21:01:36.938404	2025-04-20 21:01:36.939848	200	1
1909	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:37.122374	2025-04-20 21:01:37.144946	200	22
1910	\N	/status	GET	2025-04-20 21:01:41.707381	2025-04-20 21:01:41.707383	200	0
1911	\N	/status	GET	2025-04-20 21:01:41.872022	2025-04-20 21:01:41.872023	200	0
1912	\N	/api/statistics/global	GET	2025-04-20 21:01:41.889784	2025-04-20 21:01:41.891228	200	1
1913	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:42.058637	2025-04-20 21:01:42.06075	200	2
1914	\N	/status	GET	2025-04-20 21:01:46.735517	2025-04-20 21:01:46.735519	200	0
1915	\N	/status	GET	2025-04-20 21:01:46.903075	2025-04-20 21:01:46.903076	200	0
1916	\N	/api/statistics/global	GET	2025-04-20 21:01:46.92154	2025-04-20 21:01:46.92346	200	1
1917	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:47.09488	2025-04-20 21:01:47.096991	200	2
1918	\N	/status	GET	2025-04-20 21:01:51.715768	2025-04-20 21:01:51.71577	200	0
1919	\N	/status	GET	2025-04-20 21:01:51.885985	2025-04-20 21:01:51.885986	200	0
1920	\N	/api/statistics/global	GET	2025-04-20 21:01:51.90431	2025-04-20 21:01:51.906151	200	1
1921	\N	/api/statistics/endpoints	GET	2025-04-20 21:01:52.093468	2025-04-20 21:01:52.095963	200	2
1922	\N	/status	GET	2025-04-20 21:01:56.714364	2025-04-20 21:01:56.714366	200	0
1923	\N	/status	GET	2025-04-20 21:01:56.90415	2025-04-20 21:01:56.904152	200	0
1946	\N	/status	GET	2025-04-20 21:02:26.727447	2025-04-20 21:02:26.727449	200	0
1947	\N	/status	GET	2025-04-20 21:02:26.895584	2025-04-20 21:02:26.895585	200	0
1948	\N	/api/statistics/global	GET	2025-04-20 21:02:26.917421	2025-04-20 21:02:26.91925	200	1
1949	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:27.099794	2025-04-20 21:02:27.109044	200	9
1950	\N	/status	GET	2025-04-20 21:02:31.727076	2025-04-20 21:02:31.727077	200	0
1951	\N	/status	GET	2025-04-20 21:02:31.893389	2025-04-20 21:02:31.893391	200	0
1952	\N	/api/statistics/global	GET	2025-04-20 21:02:31.911968	2025-04-20 21:02:31.913809	200	1
1953	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:32.08895	2025-04-20 21:02:32.093861	200	4
1954	\N	/status	GET	2025-04-20 21:02:36.7261	2025-04-20 21:02:36.726102	200	0
1955	\N	/status	GET	2025-04-20 21:02:36.892297	2025-04-20 21:02:36.892298	200	0
1956	\N	/api/statistics/global	GET	2025-04-20 21:02:36.910022	2025-04-20 21:02:36.911888	200	1
1957	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:37.079003	2025-04-20 21:02:37.081464	200	2
1958	\N	/status	GET	2025-04-20 21:02:41.717184	2025-04-20 21:02:41.717185	200	0
1959	\N	/status	GET	2025-04-20 21:02:41.903787	2025-04-20 21:02:41.903788	200	0
1960	\N	/api/statistics/global	GET	2025-04-20 21:02:41.912767	2025-04-20 21:02:41.914186	200	1
1961	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:42.092085	2025-04-20 21:02:42.099774	200	7
1962	\N	/status	GET	2025-04-20 21:02:46.720857	2025-04-20 21:02:46.720858	200	0
1963	\N	/status	GET	2025-04-20 21:02:46.887802	2025-04-20 21:02:46.887804	200	0
1964	\N	/api/statistics/global	GET	2025-04-20 21:02:46.906378	2025-04-20 21:02:46.907825	200	1
1965	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:47.07927	2025-04-20 21:02:47.086503	200	7
1966	\N	/status	GET	2025-04-20 21:02:51.729664	2025-04-20 21:02:51.729666	200	0
1967	\N	/status	GET	2025-04-20 21:02:51.895624	2025-04-20 21:02:51.895626	200	0
1968	\N	/api/statistics/global	GET	2025-04-20 21:02:51.914696	2025-04-20 21:02:51.916429	200	1
1969	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:52.127392	2025-04-20 21:02:52.132485	200	5
1970	\N	/status	GET	2025-04-20 21:02:56.742284	2025-04-20 21:02:56.742286	200	0
1971	\N	/status	GET	2025-04-20 21:02:56.911894	2025-04-20 21:02:56.911895	200	0
1972	\N	/api/statistics/global	GET	2025-04-20 21:02:56.930401	2025-04-20 21:02:56.931953	200	1
1973	\N	/api/statistics/endpoints	GET	2025-04-20 21:02:57.137199	2025-04-20 21:02:57.140183	200	2
1974	\N	/status	GET	2025-04-20 21:03:01.729472	2025-04-20 21:03:01.729474	200	0
1975	\N	/status	GET	2025-04-20 21:03:01.903954	2025-04-20 21:03:01.903955	200	0
1976	\N	/api/statistics/global	GET	2025-04-20 21:03:01.921707	2025-04-20 21:03:01.925382	200	3
1977	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:02.094095	2025-04-20 21:03:02.09803	200	3
1978	\N	/status	GET	2025-04-20 21:03:06.73114	2025-04-20 21:03:06.731142	200	0
1979	\N	/status	GET	2025-04-20 21:03:07.052884	2025-04-20 21:03:07.052886	200	0
1980	\N	/api/statistics/global	GET	2025-04-20 21:03:07.069958	2025-04-20 21:03:07.073012	200	3
1981	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:07.244048	2025-04-20 21:03:07.250141	200	6
1982	\N	/status	GET	2025-04-20 21:03:11.728125	2025-04-20 21:03:11.728126	200	0
1983	\N	/status	GET	2025-04-20 21:03:11.895608	2025-04-20 21:03:11.89561	200	0
1984	\N	/api/statistics/global	GET	2025-04-20 21:03:11.913988	2025-04-20 21:03:11.948068	200	34
1985	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:12.148505	2025-04-20 21:03:12.151101	200	2
1986	\N	/status	GET	2025-04-20 21:03:16.722316	2025-04-20 21:03:16.722317	200	0
1987	\N	/status	GET	2025-04-20 21:03:16.892937	2025-04-20 21:03:16.892938	200	0
1988	\N	/api/statistics/global	GET	2025-04-20 21:03:16.892859	2025-04-20 21:03:16.896906	200	4
1989	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:17.073684	2025-04-20 21:03:17.080391	200	6
1990	\N	/status	GET	2025-04-20 21:03:21.714545	2025-04-20 21:03:21.714547	200	0
1991	\N	/status	GET	2025-04-20 21:03:21.882178	2025-04-20 21:03:21.88218	200	0
1992	\N	/api/statistics/global	GET	2025-04-20 21:03:21.900119	2025-04-20 21:03:21.901459	200	1
1993	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:22.072532	2025-04-20 21:03:22.075547	200	3
1994	\N	/status	GET	2025-04-20 21:03:26.709344	2025-04-20 21:03:26.709345	200	0
1995	\N	/status	GET	2025-04-20 21:03:26.89038	2025-04-20 21:03:26.890382	200	0
1996	\N	/api/statistics/global	GET	2025-04-20 21:03:26.910818	2025-04-20 21:03:26.912351	200	1
1997	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:27.082027	2025-04-20 21:03:27.084513	200	2
1998	\N	/status	GET	2025-04-20 21:03:31.71491	2025-04-20 21:03:31.714911	200	0
1999	\N	/status	GET	2025-04-20 21:03:31.894743	2025-04-20 21:03:31.894745	200	0
2000	\N	/api/statistics/global	GET	2025-04-20 21:03:31.912807	2025-04-20 21:03:31.914082	200	1
2001	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:32.088854	2025-04-20 21:03:32.093248	200	4
2002	\N	/status	GET	2025-04-20 21:03:36.729477	2025-04-20 21:03:36.729479	200	0
2003	\N	/status	GET	2025-04-20 21:03:36.897073	2025-04-20 21:03:36.897075	200	0
2004	\N	/api/statistics/global	GET	2025-04-20 21:03:36.915976	2025-04-20 21:03:36.917431	200	1
2005	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:37.085896	2025-04-20 21:03:37.088284	200	2
2006	\N	/status	GET	2025-04-20 21:03:41.735528	2025-04-20 21:03:41.735529	200	0
2007	\N	/status	GET	2025-04-20 21:03:41.906911	2025-04-20 21:03:41.906913	200	0
2008	\N	/api/statistics/global	GET	2025-04-20 21:03:41.922966	2025-04-20 21:03:41.924706	200	1
2009	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:42.095411	2025-04-20 21:03:42.100082	200	4
2010	\N	/status	GET	2025-04-20 21:03:46.718357	2025-04-20 21:03:46.718359	200	0
2011	\N	/status	GET	2025-04-20 21:03:46.889433	2025-04-20 21:03:46.889434	200	0
2012	\N	/api/statistics/global	GET	2025-04-20 21:03:46.904662	2025-04-20 21:03:46.906458	200	1
2013	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:47.083251	2025-04-20 21:03:47.088106	200	4
2014	\N	/status	GET	2025-04-20 21:03:51.718827	2025-04-20 21:03:51.718829	200	0
2015	\N	/status	GET	2025-04-20 21:03:51.886229	2025-04-20 21:03:51.88623	200	0
2016	\N	/api/statistics/global	GET	2025-04-20 21:03:51.905647	2025-04-20 21:03:51.906942	200	1
2017	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:52.080342	2025-04-20 21:03:52.082916	200	2
2018	\N	/status	GET	2025-04-20 21:03:56.726134	2025-04-20 21:03:56.726137	200	0
2019	\N	/status	GET	2025-04-20 21:03:56.89479	2025-04-20 21:03:56.894792	200	0
2020	\N	/api/statistics/global	GET	2025-04-20 21:03:56.913251	2025-04-20 21:03:56.91496	200	1
2021	\N	/api/statistics/endpoints	GET	2025-04-20 21:03:57.090455	2025-04-20 21:03:57.094952	200	4
2022	\N	/status	GET	2025-04-20 21:04:01.719857	2025-04-20 21:04:01.719858	200	0
2023	\N	/status	GET	2025-04-20 21:04:01.893976	2025-04-20 21:04:01.893978	200	0
2024	\N	/api/statistics/global	GET	2025-04-20 21:04:01.893899	2025-04-20 21:04:01.895569	200	1
2025	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:02.076052	2025-04-20 21:04:02.080385	200	4
2026	\N	/status	GET	2025-04-20 21:04:06.726933	2025-04-20 21:04:06.726935	200	0
2027	\N	/status	GET	2025-04-20 21:04:06.894723	2025-04-20 21:04:06.894725	200	0
2028	\N	/api/statistics/global	GET	2025-04-20 21:04:06.912971	2025-04-20 21:04:06.914475	200	1
2029	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:07.084807	2025-04-20 21:04:07.087185	200	2
2030	\N	/status	GET	2025-04-20 21:04:11.721798	2025-04-20 21:04:11.7218	200	0
2031	\N	/status	GET	2025-04-20 21:04:11.889313	2025-04-20 21:04:11.889314	200	0
2032	\N	/api/statistics/global	GET	2025-04-20 21:04:11.909147	2025-04-20 21:04:11.910404	200	1
2033	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:12.087635	2025-04-20 21:04:12.090356	200	2
2034	\N	/status	GET	2025-04-20 21:04:16.732642	2025-04-20 21:04:16.732644	200	0
2035	\N	/status	GET	2025-04-20 21:04:16.901352	2025-04-20 21:04:16.901354	200	0
2036	\N	/api/statistics/global	GET	2025-04-20 21:04:16.918973	2025-04-20 21:04:16.920194	200	1
2037	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:17.09874	2025-04-20 21:04:17.1453	200	46
2038	\N	/status	GET	2025-04-20 21:04:21.705117	2025-04-20 21:04:21.705119	200	0
2039	\N	/status	GET	2025-04-20 21:04:21.883202	2025-04-20 21:04:21.883204	200	0
2040	\N	/api/statistics/global	GET	2025-04-20 21:04:21.905229	2025-04-20 21:04:21.908151	200	2
2041	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:22.081535	2025-04-20 21:04:22.088125	200	6
2042	\N	/status	GET	2025-04-20 21:04:26.724582	2025-04-20 21:04:26.724583	200	0
2043	\N	/status	GET	2025-04-20 21:04:26.909941	2025-04-20 21:04:26.909942	200	0
2044	\N	/api/statistics/global	GET	2025-04-20 21:04:26.929092	2025-04-20 21:04:26.930397	200	1
2045	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:27.106281	2025-04-20 21:04:27.108809	200	2
2046	\N	/status	GET	2025-04-20 21:04:31.715135	2025-04-20 21:04:31.715138	200	0
2047	\N	/status	GET	2025-04-20 21:04:31.882892	2025-04-20 21:04:31.882893	200	0
2048	\N	/api/statistics/global	GET	2025-04-20 21:04:31.900714	2025-04-20 21:04:31.901926	200	1
2049	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:32.073642	2025-04-20 21:04:32.076043	200	2
2050	\N	/status	GET	2025-04-20 21:04:36.841347	2025-04-20 21:04:36.841349	200	0
2051	\N	/status	GET	2025-04-20 21:04:37.014025	2025-04-20 21:04:37.014026	200	0
2052	\N	/api/statistics/global	GET	2025-04-20 21:04:37.028314	2025-04-20 21:04:37.03052	200	2
2053	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:37.212418	2025-04-20 21:04:37.215466	200	3
2054	\N	/status	GET	2025-04-20 21:04:41.694631	2025-04-20 21:04:41.694632	200	0
2055	\N	/status	GET	2025-04-20 21:04:41.861477	2025-04-20 21:04:41.861479	200	0
2056	\N	/api/statistics/global	GET	2025-04-20 21:04:41.879653	2025-04-20 21:04:41.880863	200	1
2057	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:42.060456	2025-04-20 21:04:42.064856	200	4
2058	\N	/status	GET	2025-04-20 21:04:46.732259	2025-04-20 21:04:46.73226	200	0
2059	\N	/status	GET	2025-04-20 21:04:46.900469	2025-04-20 21:04:46.900471	200	0
2060	\N	/api/statistics/global	GET	2025-04-20 21:04:46.927011	2025-04-20 21:04:46.928679	200	1
2061	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:47.097273	2025-04-20 21:04:47.10043	200	3
2062	\N	/status	GET	2025-04-20 21:04:51.715637	2025-04-20 21:04:51.715638	200	0
2063	\N	/status	GET	2025-04-20 21:04:51.884353	2025-04-20 21:04:51.884355	200	0
2064	\N	/api/statistics/global	GET	2025-04-20 21:04:51.902814	2025-04-20 21:04:51.904299	200	1
2065	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:52.07539	2025-04-20 21:04:52.131244	200	55
2066	\N	/status	GET	2025-04-20 21:04:56.733142	2025-04-20 21:04:56.733143	200	0
2067	\N	/status	GET	2025-04-20 21:04:56.901339	2025-04-20 21:04:56.90134	200	0
2068	\N	/api/statistics/global	GET	2025-04-20 21:04:56.919277	2025-04-20 21:04:56.920462	200	1
2069	\N	/api/statistics/endpoints	GET	2025-04-20 21:04:57.094521	2025-04-20 21:04:57.097296	200	2
2070	\N	/status	GET	2025-04-20 21:05:01.707062	2025-04-20 21:05:01.707063	200	0
2071	\N	/status	GET	2025-04-20 21:05:01.885954	2025-04-20 21:05:01.885955	200	0
2072	\N	/api/statistics/global	GET	2025-04-20 21:05:01.90415	2025-04-20 21:05:01.905653	200	1
2073	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:02.076024	2025-04-20 21:05:02.078029	200	2
2074	\N	/status	GET	2025-04-20 21:05:06.734687	2025-04-20 21:05:06.734689	200	0
2075	\N	/status	GET	2025-04-20 21:05:06.903699	2025-04-20 21:05:06.903701	200	0
2076	\N	/api/statistics/global	GET	2025-04-20 21:05:06.922101	2025-04-20 21:05:06.924511	200	2
2077	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:07.095593	2025-04-20 21:05:07.109258	200	13
2078	\N	/status	GET	2025-04-20 21:05:11.727262	2025-04-20 21:05:11.727264	200	0
2079	\N	/status	GET	2025-04-20 21:05:11.901505	2025-04-20 21:05:11.901506	200	0
2080	\N	/api/statistics/global	GET	2025-04-20 21:05:11.913597	2025-04-20 21:05:11.915346	200	1
2081	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:12.093684	2025-04-20 21:05:12.102691	200	9
2082	\N	/status	GET	2025-04-20 21:05:16.742661	2025-04-20 21:05:16.742662	200	0
2083	\N	/status	GET	2025-04-20 21:05:16.910778	2025-04-20 21:05:16.91078	200	0
2084	\N	/api/statistics/global	GET	2025-04-20 21:05:16.928579	2025-04-20 21:05:16.93025	200	1
2085	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:17.103731	2025-04-20 21:05:17.1066	200	2
2086	\N	/status	GET	2025-04-20 21:05:21.810448	2025-04-20 21:05:21.810449	200	0
2087	\N	/status	GET	2025-04-20 21:05:21.977533	2025-04-20 21:05:21.977535	200	0
2088	\N	/api/statistics/global	GET	2025-04-20 21:05:21.998088	2025-04-20 21:05:21.999572	200	1
2089	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:22.182068	2025-04-20 21:05:22.184548	200	2
2090	\N	/status	GET	2025-04-20 21:05:26.860437	2025-04-20 21:05:26.860438	200	0
2091	\N	/status	GET	2025-04-20 21:05:27.02753	2025-04-20 21:05:27.027531	200	0
2092	\N	/api/statistics/global	GET	2025-04-20 21:05:27.047574	2025-04-20 21:05:27.051827	200	4
2093	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:27.2259	2025-04-20 21:05:27.230959	200	5
2094	\N	/status	GET	2025-04-20 21:05:31.734517	2025-04-20 21:05:31.734518	200	0
2095	\N	/status	GET	2025-04-20 21:05:31.903596	2025-04-20 21:05:31.903597	200	0
2096	\N	/api/statistics/global	GET	2025-04-20 21:05:31.921678	2025-04-20 21:05:31.92335	200	1
2097	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:32.090094	2025-04-20 21:05:32.093715	200	3
2098	\N	/status	GET	2025-04-20 21:05:36.693814	2025-04-20 21:05:36.693816	200	0
2099	\N	/status	GET	2025-04-20 21:05:36.860328	2025-04-20 21:05:36.86033	200	0
2100	\N	/api/statistics/global	GET	2025-04-20 21:05:36.883634	2025-04-20 21:05:36.88504	200	1
2101	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:37.07404	2025-04-20 21:05:37.077198	200	3
2102	\N	/api/restaurant	GET	2025-04-20 21:05:41.447607	2025-04-20 21:05:41.447671	200	0
2103	\N	/api/restaurant	GET	2025-04-20 21:05:41.615973	2025-04-20 21:05:41.616028	200	0
2104	\N	/status	GET	2025-04-20 21:05:42.400624	2025-04-20 21:05:42.400626	200	0
2105	\N	/status	GET	2025-04-20 21:05:42.569146	2025-04-20 21:05:42.569148	200	0
2106	\N	/api/statistics/global	GET	2025-04-20 21:05:42.58798	2025-04-20 21:05:42.58972	200	1
2107	\N	/api/statistics/global	GET	2025-04-20 21:05:42.757689	2025-04-20 21:05:42.76008	200	2
2108	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:42.776427	2025-04-20 21:05:42.778708	200	2
2109	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:42.945062	2025-04-20 21:05:42.947283	200	2
2110	\N	/api/restaurant	GET	2025-04-20 21:05:43.711627	2025-04-20 21:05:43.711678	200	0
2111	\N	/api/restaurant	GET	2025-04-20 21:05:43.8812	2025-04-20 21:05:43.881253	200	0
2112	\N	/status	GET	2025-04-20 21:05:44.377227	2025-04-20 21:05:44.377228	200	0
2113	\N	/status	GET	2025-04-20 21:05:44.550981	2025-04-20 21:05:44.550983	200	0
2114	\N	/api/statistics/global	GET	2025-04-20 21:05:44.569813	2025-04-20 21:05:44.571442	200	1
2116	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:44.746583	2025-04-20 21:05:44.749473	200	2
2117	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:44.921972	2025-04-20 21:05:44.925164	200	3
2118	\N	/api/restaurant	GET	2025-04-20 21:05:45.555409	2025-04-20 21:05:45.555461	200	0
2119	\N	/api/restaurant	GET	2025-04-20 21:05:45.723654	2025-04-20 21:05:45.72371	200	0
2120	\N	/status	GET	2025-04-20 21:05:46.259849	2025-04-20 21:05:46.259851	200	0
2121	\N	/status	GET	2025-04-20 21:05:46.428518	2025-04-20 21:05:46.428519	200	0
2122	\N	/api/statistics/global	GET	2025-04-20 21:05:46.447405	2025-04-20 21:05:46.448728	200	1
2123	\N	/api/statistics/global	GET	2025-04-20 21:05:46.617214	2025-04-20 21:05:46.618435	200	1
2124	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:46.635959	2025-04-20 21:05:46.639036	200	3
2125	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:46.806524	2025-04-20 21:05:46.808624	200	2
2132	\N	/api/statistics/global	GET	2025-04-20 21:05:55.235742	2025-04-20 21:05:55.247125	200	11
2133	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:55.442146	2025-04-20 21:05:55.445691	200	3
2134	\N	/status	GET	2025-04-20 21:07:13.880313	2025-04-20 21:07:13.880314	200	0
2135	\N	/api/statistics/global	GET	2025-04-20 21:07:14.047775	2025-04-20 21:07:14.050242	200	2
2136	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:14.246923	2025-04-20 21:07:14.25433	200	7
2137	\N	/status	GET	2025-04-20 21:07:18.732923	2025-04-20 21:07:18.732924	200	0
2138	\N	/status	GET	2025-04-20 21:07:38.321737	2025-04-20 21:07:38.321739	200	0
2139	\N	/status	GET	2025-04-20 21:07:40.350853	2025-04-20 21:07:40.350854	200	0
2140	\N	/api/statistics/global	GET	2025-04-20 21:07:40.369663	2025-04-20 21:07:40.371325	200	1
2141	\N	/api/statistics/global	GET	2025-04-20 21:07:42.396415	2025-04-20 21:07:42.398357	200	1
2154	\N	/api/statistics/global	GET	2025-04-20 21:07:55.371828	2025-04-20 21:07:55.38917	200	17
2155	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:57.425452	2025-04-20 21:07:57.427842	200	2
2156	\N	/status	GET	2025-04-20 21:07:58.334704	2025-04-20 21:07:58.334705	200	0
2158	\N	/api/statistics/global	GET	2025-04-20 21:08:00.369015	2025-04-20 21:08:00.37083	200	1
2159	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:02.430564	2025-04-20 21:08:02.433614	200	3
2160	\N	/status	GET	2025-04-20 21:08:03.349441	2025-04-20 21:08:03.349442	200	0
2161	\N	/status	GET	2025-04-20 21:08:05.341104	2025-04-20 21:08:05.341106	200	0
2162	\N	/api/statistics/global	GET	2025-04-20 21:08:05.358055	2025-04-20 21:08:05.372354	200	14
2163	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:07.381931	2025-04-20 21:08:07.384567	200	2
2164	\N	/status	GET	2025-04-20 21:08:08.332273	2025-04-20 21:08:08.332275	200	0
2165	\N	/status	GET	2025-04-20 21:08:10.349639	2025-04-20 21:08:10.349641	200	0
2178	\N	/api/statistics/global	GET	2025-04-20 21:08:25.368325	2025-04-20 21:08:25.371366	200	3
2179	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:27.404106	2025-04-20 21:08:27.406936	200	2
2180	\N	/status	GET	2025-04-20 21:08:28.32772	2025-04-20 21:08:28.327721	200	0
2181	\N	/status	GET	2025-04-20 21:08:30.36531	2025-04-20 21:08:30.365312	200	0
2190	\N	/api/statistics/global	GET	2025-04-20 21:08:40.366889	2025-04-20 21:08:40.36927	200	2
2191	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:42.396735	2025-04-20 21:08:42.401724	200	4
2192	\N	/status	GET	2025-04-20 21:08:43.350186	2025-04-20 21:08:43.350188	200	0
2193	\N	/status	GET	2025-04-20 21:08:45.363522	2025-04-20 21:08:45.363524	200	0
2202	\N	/api/statistics/global	GET	2025-04-20 21:08:55.372136	2025-04-20 21:08:55.387207	200	15
2203	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:57.396432	2025-04-20 21:08:57.402507	200	6
2204	\N	/status	GET	2025-04-20 21:08:58.339838	2025-04-20 21:08:58.33984	200	0
2205	\N	/status	GET	2025-04-20 21:09:00.369188	2025-04-20 21:09:00.36919	200	0
2206	\N	/api/statistics/global	GET	2025-04-20 21:09:00.386294	2025-04-20 21:09:00.388774	200	2
2115	\N	/api/statistics/global	GET	2025-04-20 21:05:44.746745	2025-04-20 21:05:44.748509	200	1
2126	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:05:49.732945	2025-04-20 21:05:49.733027	304	0
2127	\N	/status	GET	2025-04-20 21:05:50.052266	2025-04-20 21:05:50.052267	200	0
2128	\N	/api/statistics/global	GET	2025-04-20 21:05:50.237117	2025-04-20 21:05:50.23878	200	1
2129	\N	/api/statistics/endpoints	GET	2025-04-20 21:05:50.423367	2025-04-20 21:05:50.426376	200	3
2130	\N	/status	GET	2025-04-20 21:05:55.041889	2025-04-20 21:05:55.04189	200	0
2131	\N	/status	GET	2025-04-20 21:05:55.217849	2025-04-20 21:05:55.21785	200	0
2142	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:42.396415	2025-04-20 21:07:42.399078	200	2
2143	\N	/status	GET	2025-04-20 21:07:43.334784	2025-04-20 21:07:43.334786	200	0
2144	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:44.464789	2025-04-20 21:07:44.467769	200	2
2145	\N	/status	GET	2025-04-20 21:07:45.363266	2025-04-20 21:07:45.363268	200	0
2146	\N	/api/statistics/global	GET	2025-04-20 21:07:45.381633	2025-04-20 21:07:45.383591	200	1
2147	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:47.436382	2025-04-20 21:07:47.439291	200	2
2148	\N	/status	GET	2025-04-20 21:07:48.334709	2025-04-20 21:07:48.33471	200	0
2149	\N	/status	GET	2025-04-20 21:07:50.371456	2025-04-20 21:07:50.371457	200	0
2150	\N	/api/statistics/global	GET	2025-04-20 21:07:50.38969	2025-04-20 21:07:50.391262	200	1
2151	\N	/api/statistics/endpoints	GET	2025-04-20 21:07:52.41309	2025-04-20 21:07:52.415755	200	2
2152	\N	/status	GET	2025-04-20 21:07:53.316869	2025-04-20 21:07:53.316871	200	0
2153	\N	/status	GET	2025-04-20 21:07:55.354274	2025-04-20 21:07:55.354275	200	0
2157	\N	/status	GET	2025-04-20 21:08:00.36902	2025-04-20 21:08:00.369022	200	0
2166	\N	/api/statistics/global	GET	2025-04-20 21:08:10.367972	2025-04-20 21:08:10.37399	200	6
2167	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:12.402765	2025-04-20 21:08:12.405952	200	3
2168	\N	/status	GET	2025-04-20 21:08:13.356717	2025-04-20 21:08:13.356719	200	0
2169	\N	/status	GET	2025-04-20 21:08:15.367571	2025-04-20 21:08:15.367572	200	0
2170	\N	/api/statistics/global	GET	2025-04-20 21:08:15.384559	2025-04-20 21:08:15.387373	200	2
2171	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:17.443797	2025-04-20 21:08:17.446293	200	2
2172	\N	/status	GET	2025-04-20 21:08:18.35489	2025-04-20 21:08:18.354891	200	0
2173	\N	/status	GET	2025-04-20 21:08:20.376634	2025-04-20 21:08:20.376636	200	0
2174	\N	/api/statistics/global	GET	2025-04-20 21:08:20.394534	2025-04-20 21:08:20.396179	200	1
2175	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:22.409662	2025-04-20 21:08:22.415186	200	5
2176	\N	/status	GET	2025-04-20 21:08:23.35491	2025-04-20 21:08:23.354912	200	0
2177	\N	/status	GET	2025-04-20 21:08:25.349704	2025-04-20 21:08:25.349706	200	0
2182	\N	/api/statistics/global	GET	2025-04-20 21:08:30.384062	2025-04-20 21:08:30.410569	200	26
2183	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:32.428461	2025-04-20 21:08:32.431811	200	3
2184	\N	/status	GET	2025-04-20 21:08:33.358569	2025-04-20 21:08:33.35857	200	0
2185	\N	/status	GET	2025-04-20 21:08:35.375418	2025-04-20 21:08:35.37542	200	0
2186	\N	/api/statistics/global	GET	2025-04-20 21:08:35.393874	2025-04-20 21:08:35.395369	200	1
2187	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:37.409289	2025-04-20 21:08:37.412639	200	3
2188	\N	/status	GET	2025-04-20 21:08:38.339626	2025-04-20 21:08:38.339627	200	0
2189	\N	/status	GET	2025-04-20 21:08:40.348742	2025-04-20 21:08:40.348744	200	0
2194	\N	/api/statistics/global	GET	2025-04-20 21:08:45.363523	2025-04-20 21:08:45.366943	200	3
2195	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:47.433748	2025-04-20 21:08:47.436793	200	3
2196	\N	/status	GET	2025-04-20 21:08:48.356715	2025-04-20 21:08:48.356716	200	0
2197	\N	/status	GET	2025-04-20 21:08:50.375841	2025-04-20 21:08:50.375842	200	0
2198	\N	/api/statistics/global	GET	2025-04-20 21:08:50.396196	2025-04-20 21:08:50.398268	200	2
2199	\N	/api/statistics/endpoints	GET	2025-04-20 21:08:52.450126	2025-04-20 21:08:52.453329	200	3
2200	\N	/status	GET	2025-04-20 21:08:53.329597	2025-04-20 21:08:53.3296	200	0
2201	\N	/status	GET	2025-04-20 21:08:55.350973	2025-04-20 21:08:55.350975	200	0
2218	\N	/api/statistics/global	GET	2025-04-20 21:09:15.384186	2025-04-20 21:09:15.391657	200	7
2219	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:17.432959	2025-04-20 21:09:17.436255	200	3
2220	\N	/status	GET	2025-04-20 21:09:18.324178	2025-04-20 21:09:18.32418	200	0
2221	\N	/status	GET	2025-04-20 21:09:20.337162	2025-04-20 21:09:20.337164	200	0
2222	\N	/api/statistics/global	GET	2025-04-20 21:09:20.35514	2025-04-20 21:09:20.357146	200	2
2223	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:22.389467	2025-04-20 21:09:22.392463	200	2
2860	\N	/status	GET	2025-04-20 23:57:27.978882	2025-04-20 23:57:27.978883	200	0
2861	\N	/status	GET	2025-04-20 23:57:28.147967	2025-04-20 23:57:28.147969	200	0
2862	\N	/api/statistics/global	GET	2025-04-20 23:57:28.166639	2025-04-20 23:57:28.172655	200	6
2863	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:28.376411	2025-04-20 23:57:28.379058	200	2
2864	\N	/status	GET	2025-04-20 23:57:33.035001	2025-04-20 23:57:33.035002	200	0
2865	\N	/status	GET	2025-04-20 23:57:33.238291	2025-04-20 23:57:33.238293	200	0
2922	\N	/api/statistics/global	GET	2025-04-20 23:58:43.137029	2025-04-20 23:58:43.143911	200	6
2923	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:43.344709	2025-04-20 23:58:43.355228	200	10
2924	\N	/status	GET	2025-04-20 23:58:48.00779	2025-04-20 23:58:48.007792	200	0
2925	\N	/status	GET	2025-04-20 23:58:48.180085	2025-04-20 23:58:48.180087	200	0
2969	\N	/api/restaurant	GET	2025-04-21 00:01:35.268735	2025-04-21 00:01:35.268801	200	0
2970	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 00:12:32.768235	2025-04-21 00:12:32.768328	304	0
2971	\N	/api/washingmachines	GET	2025-04-21 00:12:32.743892	2025-04-21 00:12:33.438642	200	694
2972	\N	/api/restaurant	GET	2025-04-21 00:13:01.234947	2025-04-21 00:13:01.235031	200	0
2974	\N	/api/restaurant	GET	2025-04-21 09:21:16.864237	2025-04-21 09:21:16.864305	200	0
2975	\N	/status	GET	2025-04-21 09:21:23.386328	2025-04-21 09:21:23.386333	200	0
2976	\N	/api/statistics/global	GET	2025-04-21 09:21:23.874365	2025-04-21 09:21:23.876385	200	2
2977	\N	/api/statistics/endpoints	GET	2025-04-21 09:21:24.240029	2025-04-21 09:21:24.243287	200	3
2978	\N	/status	GET	2025-04-21 09:21:28.451879	2025-04-21 09:21:28.451881	200	0
2979	\N	/api/statistics/global	GET	2025-04-21 09:21:28.887709	2025-04-21 09:21:28.889767	200	2
2980	\N	/status	GET	2025-04-21 09:21:29.045837	2025-04-21 09:21:29.045839	200	0
2981	\N	/api/statistics/endpoints	GET	2025-04-21 09:21:29.304594	2025-04-21 09:21:29.307974	200	3
2982	\N	/status	GET	2025-04-21 09:21:33.396715	2025-04-21 09:21:33.396716	200	0
2983	\N	/status	GET	2025-04-21 09:21:33.8039	2025-04-21 09:21:33.803901	200	0
2984	\N	/api/statistics/global	GET	2025-04-21 09:21:33.805126	2025-04-21 09:21:33.807093	200	1
2985	\N	/api/statistics/endpoints	GET	2025-04-21 09:21:34.219871	2025-04-21 09:21:34.223633	200	3
2986	\N	/status	GET	2025-04-21 09:21:38.407922	2025-04-21 09:21:38.407923	200	0
2987	\N	/status	GET	2025-04-21 09:21:38.779456	2025-04-21 09:21:38.779457	200	0
2988	\N	/api/statistics/global	GET	2025-04-21 09:21:38.782426	2025-04-21 09:21:38.78424	200	1
2989	\N	/api/statistics/endpoints	GET	2025-04-21 09:21:39.124065	2025-04-21 09:21:39.128318	200	4
2990	\N	/status	GET	2025-04-21 09:21:43.396022	2025-04-21 09:21:43.396024	200	0
2991	\N	/status	GET	2025-04-21 09:21:43.731165	2025-04-21 09:21:43.731166	200	0
2992	\N	/api/statistics/global	GET	2025-04-21 09:21:43.73426	2025-04-21 09:21:43.736397	200	2
2993	\N	/api/statistics/endpoints	GET	2025-04-21 09:21:44.14864	2025-04-21 09:21:44.151613	200	2
2994	\N	/api/restaurant	GET	2025-04-21 09:21:47.964061	2025-04-21 09:21:47.964126	200	0
2995	\N	/api/washingmachines	GET	2025-04-21 09:21:54.633927	2025-04-21 09:21:55.40554	200	771
2996	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 09:22:12.459147	2025-04-21 09:22:12.46077	200	1
2997	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 09:22:12.938687	2025-04-21 09:22:12.940595	200	1
2998	\N	/api/restaurant	GET	2025-04-21 09:22:13.111625	2025-04-21 09:22:13.111682	200	0
2999	\N	/api/washingmachines	GET	2025-04-21 09:26:56.138071	2025-04-21 09:26:56.860111	200	722
3000	\N	/api/washingmachines	GET	2025-04-21 09:38:22.73131	2025-04-21 09:38:23.415552	200	684
3001	\N	/api/washingmachines	GET	2025-04-21 11:39:59.493211	2025-04-21 11:40:00.266517	200	773
3002	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 11:47:05.584402	2025-04-21 11:47:05.584493	304	0
3003	\N	/api/washingmachines	GET	2025-04-21 11:47:05.946203	2025-04-21 11:47:06.719764	200	773
3004	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 11:47:39.115871	2025-04-21 11:47:39.119792	200	3
3005	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 11:47:39.622102	2025-04-21 11:47:39.623784	200	1
3006	\N	/api/restaurant	GET	2025-04-21 11:47:39.782585	2025-04-21 11:47:39.782643	200	0
2207	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:02.40525	2025-04-20 21:09:02.40982	200	4
2208	\N	/status	GET	2025-04-20 21:09:03.341613	2025-04-20 21:09:03.341615	200	0
2209	\N	/status	GET	2025-04-20 21:09:05.337662	2025-04-20 21:09:05.337663	200	0
2210	\N	/api/statistics/global	GET	2025-04-20 21:09:05.355613	2025-04-20 21:09:05.368054	200	12
2211	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:07.387077	2025-04-20 21:09:07.393653	200	6
2212	\N	/status	GET	2025-04-20 21:09:08.3267	2025-04-20 21:09:08.326702	200	0
2213	\N	/status	GET	2025-04-20 21:09:10.353349	2025-04-20 21:09:10.35335	200	0
2214	\N	/api/statistics/global	GET	2025-04-20 21:09:10.371522	2025-04-20 21:09:10.373349	200	1
2215	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:12.420141	2025-04-20 21:09:12.425642	200	5
2216	\N	/status	GET	2025-04-20 21:09:13.350451	2025-04-20 21:09:13.350452	200	0
2217	\N	/status	GET	2025-04-20 21:09:15.36623	2025-04-20 21:09:15.366232	200	0
2224	\N	/status	GET	2025-04-20 21:09:23.315551	2025-04-20 21:09:23.315553	200	0
2225	\N	/status	GET	2025-04-20 21:09:25.358736	2025-04-20 21:09:25.358737	200	0
2226	\N	/api/statistics/global	GET	2025-04-20 21:09:25.358739	2025-04-20 21:09:25.373828	200	15
2227	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:27.393055	2025-04-20 21:09:27.39519	200	2
2228	\N	/status	GET	2025-04-20 21:09:28.351182	2025-04-20 21:09:28.351184	200	0
2229	\N	/status	GET	2025-04-20 21:09:30.359226	2025-04-20 21:09:30.359227	200	0
2230	\N	/api/statistics/global	GET	2025-04-20 21:09:30.37842	2025-04-20 21:09:30.380171	200	1
2231	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:32.403478	2025-04-20 21:09:32.405957	200	2
2232	\N	/status	GET	2025-04-20 21:09:33.314236	2025-04-20 21:09:33.314238	200	0
2233	\N	/status	GET	2025-04-20 21:09:35.340271	2025-04-20 21:09:35.340273	200	0
2234	\N	/api/statistics/global	GET	2025-04-20 21:09:35.359046	2025-04-20 21:09:35.36104	200	1
2235	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:37.385958	2025-04-20 21:09:37.399265	200	13
2236	\N	/status	GET	2025-04-20 21:09:38.343183	2025-04-20 21:09:38.343184	200	0
2237	\N	/status	GET	2025-04-20 21:09:40.354965	2025-04-20 21:09:40.354966	200	0
2238	\N	/api/statistics/global	GET	2025-04-20 21:09:40.37322	2025-04-20 21:09:40.385354	200	12
2239	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:42.39669	2025-04-20 21:09:42.399887	200	3
2240	\N	/status	GET	2025-04-20 21:09:43.348141	2025-04-20 21:09:43.348142	200	0
2241	\N	/status	GET	2025-04-20 21:09:45.362119	2025-04-20 21:09:45.36212	200	0
2242	\N	/api/statistics/global	GET	2025-04-20 21:09:45.391679	2025-04-20 21:09:45.393874	200	2
2243	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:47.435603	2025-04-20 21:09:47.437868	200	2
2244	\N	/status	GET	2025-04-20 21:09:48.35139	2025-04-20 21:09:48.351392	200	0
2245	\N	/status	GET	2025-04-20 21:09:50.38631	2025-04-20 21:09:50.386312	200	0
2246	\N	/api/statistics/global	GET	2025-04-20 21:09:50.395774	2025-04-20 21:09:50.397454	200	1
2247	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:52.411705	2025-04-20 21:09:52.416293	200	4
2248	\N	/status	GET	2025-04-20 21:09:53.33751	2025-04-20 21:09:53.337511	200	0
2249	\N	/status	GET	2025-04-20 21:09:55.349421	2025-04-20 21:09:55.349423	200	0
2250	\N	/api/statistics/global	GET	2025-04-20 21:09:55.366899	2025-04-20 21:09:55.380481	200	13
2251	\N	/api/statistics/endpoints	GET	2025-04-20 21:09:57.393307	2025-04-20 21:09:57.398192	200	4
2252	\N	/status	GET	2025-04-20 21:09:58.330431	2025-04-20 21:09:58.330432	200	0
2253	\N	/status	GET	2025-04-20 21:10:00.333984	2025-04-20 21:10:00.333985	200	0
2254	\N	/api/statistics/global	GET	2025-04-20 21:10:00.351896	2025-04-20 21:10:00.353529	200	1
2255	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:02.371912	2025-04-20 21:10:02.379503	200	7
2256	\N	/status	GET	2025-04-20 21:10:03.32422	2025-04-20 21:10:03.324222	200	0
2257	\N	/status	GET	2025-04-20 21:10:05.340192	2025-04-20 21:10:05.340194	200	0
2258	\N	/api/statistics/global	GET	2025-04-20 21:10:05.35817	2025-04-20 21:10:05.362813	200	4
2259	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:07.384914	2025-04-20 21:10:07.3887	200	3
2260	\N	/status	GET	2025-04-20 21:10:08.34475	2025-04-20 21:10:08.344752	200	0
2261	\N	/status	GET	2025-04-20 21:10:10.35392	2025-04-20 21:10:10.353921	200	0
2262	\N	/api/statistics/global	GET	2025-04-20 21:10:10.353925	2025-04-20 21:10:10.355533	200	1
2263	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:12.395394	2025-04-20 21:10:12.397538	200	2
2264	\N	/status	GET	2025-04-20 21:10:13.332065	2025-04-20 21:10:13.332067	200	0
2265	\N	/status	GET	2025-04-20 21:10:15.365586	2025-04-20 21:10:15.365588	200	0
2266	\N	/api/statistics/global	GET	2025-04-20 21:10:15.384279	2025-04-20 21:10:15.387519	200	3
2267	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:17.405215	2025-04-20 21:10:17.408009	200	2
2268	\N	/status	GET	2025-04-20 21:10:18.349442	2025-04-20 21:10:18.349444	200	0
2269	\N	/status	GET	2025-04-20 21:10:20.343364	2025-04-20 21:10:20.343365	200	0
2270	\N	/api/statistics/global	GET	2025-04-20 21:10:20.354657	2025-04-20 21:10:20.357349	200	2
2271	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:22.380952	2025-04-20 21:10:22.383888	200	2
2272	\N	/status	GET	2025-04-20 21:10:23.337635	2025-04-20 21:10:23.337638	200	0
2273	\N	/status	GET	2025-04-20 21:10:25.349568	2025-04-20 21:10:25.34957	200	0
2274	\N	/api/statistics/global	GET	2025-04-20 21:10:25.367727	2025-04-20 21:10:25.426566	200	58
2275	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:27.392781	2025-04-20 21:10:27.396593	200	3
2276	\N	/status	GET	2025-04-20 21:10:28.329973	2025-04-20 21:10:28.329974	200	0
2277	\N	/status	GET	2025-04-20 21:10:30.363547	2025-04-20 21:10:30.363549	200	0
2278	\N	/api/statistics/global	GET	2025-04-20 21:10:30.381985	2025-04-20 21:10:30.384108	200	2
2279	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:32.399867	2025-04-20 21:10:32.418692	200	18
2280	\N	/status	GET	2025-04-20 21:10:33.348933	2025-04-20 21:10:33.348935	200	0
2281	\N	/status	GET	2025-04-20 21:10:35.33678	2025-04-20 21:10:35.336782	200	0
2282	\N	/api/statistics/global	GET	2025-04-20 21:10:35.354758	2025-04-20 21:10:35.359233	200	4
2283	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:37.38113	2025-04-20 21:10:37.383975	200	2
2284	\N	/status	GET	2025-04-20 21:10:38.339256	2025-04-20 21:10:38.339258	200	0
2285	\N	/status	GET	2025-04-20 21:10:40.351162	2025-04-20 21:10:40.351163	200	0
2286	\N	/api/statistics/global	GET	2025-04-20 21:10:40.351076	2025-04-20 21:10:40.352756	200	1
2287	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:42.394777	2025-04-20 21:10:42.403909	200	9
2288	\N	/status	GET	2025-04-20 21:10:43.345712	2025-04-20 21:10:43.345713	200	0
2289	\N	/status	GET	2025-04-20 21:10:45.366574	2025-04-20 21:10:45.366575	200	0
2290	\N	/api/statistics/global	GET	2025-04-20 21:10:45.390774	2025-04-20 21:10:45.393726	200	2
2291	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:47.410122	2025-04-20 21:10:47.414837	200	4
2292	\N	/status	GET	2025-04-20 21:10:48.334296	2025-04-20 21:10:48.334298	200	0
2293	\N	/status	GET	2025-04-20 21:10:50.339056	2025-04-20 21:10:50.339057	200	0
2294	\N	/api/statistics/global	GET	2025-04-20 21:10:50.35749	2025-04-20 21:10:50.370095	200	12
2295	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:52.386756	2025-04-20 21:10:52.393959	200	7
2296	\N	/status	GET	2025-04-20 21:10:53.338465	2025-04-20 21:10:53.338467	200	0
2297	\N	/status	GET	2025-04-20 21:10:55.349318	2025-04-20 21:10:55.34932	200	0
2298	\N	/api/statistics/global	GET	2025-04-20 21:10:55.36724	2025-04-20 21:10:55.382979	200	15
2299	\N	/api/statistics/endpoints	GET	2025-04-20 21:10:57.39625	2025-04-20 21:10:57.398639	200	2
2300	\N	/status	GET	2025-04-20 21:10:58.344002	2025-04-20 21:10:58.344003	200	0
2301	\N	/status	GET	2025-04-20 21:11:00.360288	2025-04-20 21:11:00.36029	200	0
2302	\N	/api/statistics/global	GET	2025-04-20 21:11:00.378478	2025-04-20 21:11:00.384326	200	5
2303	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:02.413673	2025-04-20 21:11:02.416909	200	3
2304	\N	/status	GET	2025-04-20 21:11:03.340593	2025-04-20 21:11:03.340594	200	0
2305	\N	/status	GET	2025-04-20 21:11:05.338799	2025-04-20 21:11:05.3388	200	0
2306	\N	/api/statistics/global	GET	2025-04-20 21:11:05.338735	2025-04-20 21:11:05.340193	200	1
2307	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:07.390019	2025-04-20 21:11:07.396882	200	6
2308	\N	/status	GET	2025-04-20 21:11:08.336225	2025-04-20 21:11:08.336227	200	0
2309	\N	/status	GET	2025-04-20 21:11:10.347486	2025-04-20 21:11:10.347487	200	0
2310	\N	/api/statistics/global	GET	2025-04-20 21:11:10.365857	2025-04-20 21:11:10.373934	200	8
2311	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:12.394758	2025-04-20 21:11:12.398402	200	3
2312	\N	/status	GET	2025-04-20 21:11:13.330613	2025-04-20 21:11:13.330614	200	0
2313	\N	/status	GET	2025-04-20 21:11:15.329158	2025-04-20 21:11:15.32916	200	0
2314	\N	/api/statistics/global	GET	2025-04-20 21:11:15.351044	2025-04-20 21:11:15.362422	200	11
2315	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:17.374184	2025-04-20 21:11:17.380558	200	6
2316	\N	/status	GET	2025-04-20 21:11:18.337537	2025-04-20 21:11:18.337539	200	0
2317	\N	/status	GET	2025-04-20 21:11:20.342488	2025-04-20 21:11:20.342489	200	0
2318	\N	/api/statistics/global	GET	2025-04-20 21:11:20.342487	2025-04-20 21:11:20.368827	200	26
2319	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:22.383026	2025-04-20 21:11:22.442998	200	59
2320	\N	/status	GET	2025-04-20 21:11:23.330053	2025-04-20 21:11:23.330055	200	0
2321	\N	/status	GET	2025-04-20 21:11:25.35633	2025-04-20 21:11:25.356331	200	0
2322	\N	/api/statistics/global	GET	2025-04-20 21:11:25.37034	2025-04-20 21:11:25.372219	200	1
2323	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:27.424571	2025-04-20 21:11:27.426968	200	2
2324	\N	/status	GET	2025-04-20 21:11:28.343446	2025-04-20 21:11:28.343447	200	0
2325	\N	/status	GET	2025-04-20 21:11:30.331305	2025-04-20 21:11:30.331306	200	0
2330	\N	/api/statistics/global	GET	2025-04-20 21:11:35.357372	2025-04-20 21:11:35.369791	200	12
2331	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:37.382844	2025-04-20 21:11:37.385245	200	2
2332	\N	/status	GET	2025-04-20 21:11:38.337199	2025-04-20 21:11:38.3372	200	0
2333	\N	/status	GET	2025-04-20 21:11:40.349819	2025-04-20 21:11:40.34982	200	0
2334	\N	/api/statistics/global	GET	2025-04-20 21:11:40.367756	2025-04-20 21:11:40.370453	200	2
2335	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:42.403224	2025-04-20 21:11:42.405633	200	2
2336	\N	/status	GET	2025-04-20 21:11:43.353961	2025-04-20 21:11:43.353963	200	0
2337	\N	/status	GET	2025-04-20 21:11:45.36086	2025-04-20 21:11:45.360861	200	0
2338	\N	/api/statistics/global	GET	2025-04-20 21:11:45.379845	2025-04-20 21:11:45.381648	200	1
2339	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:47.410048	2025-04-20 21:11:47.417234	200	7
2340	\N	/status	GET	2025-04-20 21:11:48.332805	2025-04-20 21:11:48.332806	200	0
2341	\N	/status	GET	2025-04-20 21:11:50.377285	2025-04-20 21:11:50.377286	200	0
2349	\N	/status	GET	2025-04-20 21:11:57.500148	2025-04-20 21:11:57.500149	200	0
2350	\N	/api/statistics/global	GET	2025-04-20 21:11:57.674869	2025-04-20 21:11:57.676385	200	1
2351	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:57.849096	2025-04-20 21:11:57.851733	200	2
2352	\N	/status	GET	2025-04-20 21:12:02.507005	2025-04-20 21:12:02.507006	200	0
2353	\N	/status	GET	2025-04-20 21:12:02.673997	2025-04-20 21:12:02.673999	200	0
2354	\N	/api/statistics/global	GET	2025-04-20 21:12:02.691654	2025-04-20 21:12:02.69339	200	1
2355	\N	/api/statistics/endpoints	GET	2025-04-20 21:12:02.864571	2025-04-20 21:12:02.86753	200	2
2356	\N	/status	GET	2025-04-20 21:12:07.510108	2025-04-20 21:12:07.510109	200	0
2357	\N	/status	GET	2025-04-20 21:12:07.68889	2025-04-20 21:12:07.688891	200	0
2358	\N	/api/statistics/global	GET	2025-04-20 21:12:07.708041	2025-04-20 21:12:07.70985	200	1
2359	\N	/api/statistics/endpoints	GET	2025-04-20 21:12:07.891512	2025-04-20 21:12:07.893709	200	2
2360	\N	/status	GET	2025-04-20 21:12:12.502068	2025-04-20 21:12:12.502069	200	0
2361	\N	/status	GET	2025-04-20 21:12:12.67297	2025-04-20 21:12:12.672972	200	0
2362	\N	/api/restaurant	GET	2025-04-20 21:12:12.691071	2025-04-20 21:12:12.691159	200	0
2363	\N	/api/statistics/global	GET	2025-04-20 21:12:12.810667	2025-04-20 21:12:12.812829	200	2
2364	\N	/api/restaurant	GET	2025-04-20 21:12:12.867528	2025-04-20 21:12:12.867584	200	0
2365	\N	/api/statistics/endpoints	GET	2025-04-20 21:12:12.98088	2025-04-20 21:12:12.984226	200	3
2366	\N	/status	GET	2025-04-20 21:13:05.416579	2025-04-20 21:13:05.416581	200	0
2367	\N	/status	GET	2025-04-20 21:13:05.584603	2025-04-20 21:13:05.584604	200	0
2368	\N	/api/statistics/global	GET	2025-04-20 21:13:05.602636	2025-04-20 21:13:05.604465	200	1
2369	\N	/api/statistics/global	GET	2025-04-20 21:13:05.775855	2025-04-20 21:13:05.777778	200	1
2370	\N	/api/statistics/endpoints	GET	2025-04-20 21:13:05.793907	2025-04-20 21:13:05.79645	200	2
2371	\N	/api/statistics/endpoints	GET	2025-04-20 21:13:05.962087	2025-04-20 21:13:05.964779	200	2
2372	\N	/api/restaurant	GET	2025-04-20 21:13:07.635004	2025-04-20 21:13:07.635063	200	0
2373	\N	/api/restaurant	GET	2025-04-20 21:13:07.802712	2025-04-20 21:13:07.802764	200	0
2374	\N	/api/restaurant	GET	2025-04-20 21:14:53.85841	2025-04-20 21:14:53.858472	200	0
2375	\N	/api/restaurant	GET	2025-04-20 21:17:16.277451	2025-04-20 21:17:16.277531	200	0
2376	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:19:54.334283	2025-04-20 21:19:54.334379	304	0
2401	\N	/api/washingmachines	GET	2025-04-20 21:32:17.560227	2025-04-20 21:32:18.252934	200	692
2402	\N	/api/washingmachines	GET	2025-04-20 21:32:41.196324	2025-04-20 21:32:41.989397	200	793
2403	\N	/api/washingmachines	GET	2025-04-20 21:35:16.538787	2025-04-20 21:35:17.318123	200	779
2404	\N	/api/washingmachines	GET	2025-04-20 21:36:57.102798	2025-04-20 21:36:57.729067	200	626
2405	\N	/api/washingmachines	GET	2025-04-20 21:36:58.679365	2025-04-20 21:36:58.852221	200	172
2406	\N	/api/washingmachines	GET	2025-04-20 21:38:39.915813	2025-04-20 21:38:40.59336	200	677
2407	\N	/api/restaurant	GET	2025-04-20 21:40:51.020547	2025-04-20 21:40:51.02062	200	0
2408	\N	/api/restaurant	GET	2025-04-20 21:40:51.191079	2025-04-20 21:40:51.19115	200	0
2409	\N	/api/washingmachines	GET	2025-04-20 21:40:52.011607	2025-04-20 21:40:52.772687	200	761
2410	\N	/api/washingmachines	GET	2025-04-20 21:40:52.940189	2025-04-20 21:40:53.112707	200	172
2411	\N	/	GET	2025-04-20 21:41:40.511637	2025-04-20 21:41:40.51164	500	0
2412	\N	/api/washingmachines	GET	2025-04-20 21:42:31.70277	2025-04-20 21:42:32.464708	200	761
2413	\N	/api/washingmachines	GET	2025-04-20 21:42:32.635384	2025-04-20 21:42:32.804786	200	169
2414	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:42:35.530373	2025-04-20 21:42:35.530469	304	0
2435	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 22:14:23.04073	2025-04-20 22:14:23.040841	304	0
2441	\N	/api/statistics/endpoints	GET	2025-04-20 22:14:31.116057	2025-04-20 22:14:31.122393	200	6
2442	\N	/api/statistics/endpoints	GET	2025-04-20 22:14:31.310118	2025-04-20 22:14:31.313508	200	3
2443	\N	/api/washingmachines	GET	2025-04-20 22:14:33.575443	2025-04-20 22:14:34.298353	200	722
2444	\N	/api/washingmachines	GET	2025-04-20 22:14:34.557522	2025-04-20 22:14:34.745897	200	188
2445	\N	/api/washingmachines	GET	2025-04-20 22:15:03.413876	2025-04-20 22:15:04.020558	200	606
2446	\N	/api/washingmachines	GET	2025-04-20 22:19:57.858651	2025-04-20 22:19:58.559166	200	700
2447	\N	/api/washingmachines	GET	2025-04-20 22:20:15.393172	2025-04-20 22:20:16.082692	200	689
2448	\N	/api/washingmachines	GET	2025-04-20 22:24:00.788637	2025-04-20 22:24:01.40063	200	611
2449	\N	/api/washingmachines	GET	2025-04-20 22:25:21.704379	2025-04-20 22:25:22.409555	200	705
2450	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 22:28:13.319729	2025-04-20 22:28:13.319825	304	0
2468	\N	/api/statistics/endpoints	GET	2025-04-20 22:48:29.427088	2025-04-20 22:48:29.463079	200	35
2470	\N	/api/statistics/endpoints	GET	2025-04-20 22:48:29.655619	2025-04-20 22:48:29.658896	200	3
2471	\N	/status	GET	2025-04-20 22:48:33.581049	2025-04-20 22:48:33.58105	200	0
2472	\N	/status	GET	2025-04-20 22:48:33.753369	2025-04-20 22:48:33.75337	200	0
2473	\N	/api/statistics/global	GET	2025-04-20 22:48:33.771581	2025-04-20 22:48:33.772861	200	1
2474	\N	/api/statistics/endpoints	GET	2025-04-20 22:48:33.947342	2025-04-20 22:48:33.950504	200	3
2475	\N	/status	GET	2025-04-20 22:48:38.659511	2025-04-20 22:48:38.659514	200	0
2476	\N	/status	GET	2025-04-20 22:48:38.826836	2025-04-20 22:48:38.826837	200	0
2477	\N	/api/statistics/global	GET	2025-04-20 22:48:38.963481	2025-04-20 22:48:38.995975	200	32
2478	\N	/api/statistics/endpoints	GET	2025-04-20 22:48:39.182824	2025-04-20 22:48:39.185611	200	2
2479	\N	/api/restaurant	GET	2025-04-20 22:48:43.250315	2025-04-20 22:48:43.25038	200	0
2480	\N	/api/restaurant	GET	2025-04-20 22:48:43.467806	2025-04-20 22:48:43.467871	200	0
2481	\N	/api/washingmachines	GET	2025-04-20 22:48:51.136404	2025-04-20 22:48:51.925444	200	789
2482	\N	/api/washingmachines	GET	2025-04-20 22:48:52.093238	2025-04-20 22:48:52.266193	200	172
2483	\N	/api/washingmachines	GET	2025-04-20 22:49:47.193926	2025-04-20 22:49:47.879829	200	685
2866	\N	/api/statistics/global	GET	2025-04-20 23:57:33.260393	2025-04-20 23:57:33.276402	200	16
2867	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:33.503555	2025-04-20 23:57:33.508338	200	4
2868	\N	/status	GET	2025-04-20 23:57:37.938729	2025-04-20 23:57:37.938731	200	0
2869	\N	/status	GET	2025-04-20 23:57:38.112368	2025-04-20 23:57:38.11237	200	0
2870	\N	/api/statistics/global	GET	2025-04-20 23:57:38.13253	2025-04-20 23:57:38.134568	200	2
2871	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:38.341404	2025-04-20 23:57:38.34471	200	3
2872	\N	/status	GET	2025-04-20 23:57:43.006516	2025-04-20 23:57:43.006518	200	0
2873	\N	/status	GET	2025-04-20 23:57:43.180675	2025-04-20 23:57:43.180677	200	0
2874	\N	/api/statistics/global	GET	2025-04-20 23:57:43.199207	2025-04-20 23:57:43.200816	200	1
2875	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:43.395378	2025-04-20 23:57:43.397778	200	2
2876	\N	/status	GET	2025-04-20 23:57:47.949287	2025-04-20 23:57:47.949289	200	0
2877	\N	/status	GET	2025-04-20 23:57:48.120569	2025-04-20 23:57:48.120571	200	0
3649	\N	/status	GET	2025-04-28 12:21:49.680315	2025-04-28 12:21:49.680316	200	0
2326	\N	/api/statistics/global	GET	2025-04-20 21:11:30.349227	2025-04-20 21:11:30.361604	200	12
2327	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:32.372853	2025-04-20 21:11:32.376398	200	3
2328	\N	/status	GET	2025-04-20 21:11:33.334861	2025-04-20 21:11:33.334863	200	0
2329	\N	/status	GET	2025-04-20 21:11:35.344025	2025-04-20 21:11:35.344027	200	0
2342	\N	/api/statistics/global	GET	2025-04-20 21:11:50.377286	2025-04-20 21:11:50.392984	200	15
2343	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:52.446955	2025-04-20 21:11:52.449048	200	2
2344	\N	/status	GET	2025-04-20 21:11:53.314576	2025-04-20 21:11:53.314577	200	0
2345	\N	/status	GET	2025-04-20 21:11:55.351236	2025-04-20 21:11:55.351239	200	0
2346	\N	/api/statistics/global	GET	2025-04-20 21:11:55.368884	2025-04-20 21:11:55.371728	200	2
2347	\N	/api/statistics/endpoints	GET	2025-04-20 21:11:56.878126	2025-04-20 21:11:56.88118	200	3
2348	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:11:57.265244	2025-04-20 21:11:57.265328	304	0
2377	\N	/api/restaurant	GET	2025-04-20 21:19:54.382535	2025-04-20 21:19:54.382598	200	0
2378	\N	/status	GET	2025-04-20 21:23:03.192573	2025-04-20 21:23:03.192574	200	0
2379	\N	/status	GET	2025-04-20 21:23:03.36592	2025-04-20 21:23:03.365922	200	0
2380	\N	/api/statistics/global	GET	2025-04-20 21:23:03.428159	2025-04-20 21:23:03.429995	200	1
2381	\N	/api/statistics/global	GET	2025-04-20 21:23:03.601922	2025-04-20 21:23:03.603576	200	1
2382	\N	/api/statistics/endpoints	GET	2025-04-20 21:23:03.640495	2025-04-20 21:23:03.644245	200	3
2383	\N	/api/statistics/endpoints	GET	2025-04-20 21:23:03.810877	2025-04-20 21:23:03.813347	200	2
2384	\N	/api/restaurant	GET	2025-04-20 21:23:04.764266	2025-04-20 21:23:04.764343	200	0
2385	\N	/api/restaurant	GET	2025-04-20 21:23:04.936392	2025-04-20 21:23:04.936444	200	0
2386	\N	/api/washingmachines	GET	2025-04-20 21:23:06.4875	2025-04-20 21:23:07.252928	200	765
2387	\N	/api/washingmachines	GET	2025-04-20 21:23:07.4212	2025-04-20 21:23:07.604409	200	183
2388	\N	/api/restaurant	GET	2025-04-20 21:23:19.888261	2025-04-20 21:23:19.888324	200	0
2389	\N	/api/restaurant	GET	2025-04-20 21:23:20.055679	2025-04-20 21:23:20.055736	200	0
2390	\N	/api/restaurant	GET	2025-04-20 21:23:40.009747	2025-04-20 21:23:40.009823	200	0
2391	\N	/api/restaurant	GET	2025-04-20 21:24:28.213868	2025-04-20 21:24:28.213921	200	0
2392	\N	/api/restaurant	GET	2025-04-20 21:25:23.585155	2025-04-20 21:25:23.585216	200	0
2393	\N	/api/restaurant	GET	2025-04-20 21:25:30.43498	2025-04-20 21:25:30.435057	200	0
2394	\N	/api/washingmachines	GET	2025-04-20 21:29:25.030932	2025-04-20 21:29:25.802887	200	771
2395	\N	/api/washingmachines	GET	2025-04-20 21:29:25.970202	2025-04-20 21:29:26.155978	200	185
2396	\N	/api/restaurant	GET	2025-04-20 21:31:48.203186	2025-04-20 21:31:48.203258	200	0
2397	\N	/api/restaurant	GET	2025-04-20 21:31:48.380301	2025-04-20 21:31:48.38035	200	0
2398	\N	/api/washingmachines	GET	2025-04-20 21:31:50.592686	2025-04-20 21:31:51.288725	200	696
2399	\N	/api/washingmachines	GET	2025-04-20 21:31:51.460234	2025-04-20 21:31:51.635973	200	175
2400	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:32:17.479412	2025-04-20 21:32:17.479507	304	0
2415	\N	/api/washingmachines	GET	2025-04-20 21:42:35.729291	2025-04-20 21:42:35.903377	200	174
2416	\N	/api/washingmachines	GET	2025-04-20 21:42:37.482709	2025-04-20 21:42:37.663704	200	180
2417	\N	/api/washingmachines	GET	2025-04-20 21:42:38.665954	2025-04-20 21:42:38.837249	200	171
2418	\N	/api/washingmachines	GET	2025-04-20 21:43:42.928339	2025-04-20 21:43:43.622506	200	694
2419	\N	/api/restaurant	GET	2025-04-20 21:44:23.544337	2025-04-20 21:44:23.544393	200	0
2420	\N	/api/restaurant	GET	2025-04-20 21:44:23.716336	2025-04-20 21:44:23.716393	200	0
2421	\N	/api/washingmachines	GET	2025-04-20 21:44:24.835568	2025-04-20 21:44:25.519659	200	684
2422	\N	/api/washingmachines	GET	2025-04-20 21:44:25.728194	2025-04-20 21:44:25.900946	200	172
2423	\N	/api/washingmachines	GET	2025-04-20 21:49:11.52851	2025-04-20 21:49:12.239573	200	711
2424	\N	/api/washingmachines	GET	2025-04-20 21:49:11.94463	2025-04-20 21:49:12.41512	200	470
2425	\N	/api/washingmachines	GET	2025-04-20 21:54:12.404538	2025-04-20 21:54:13.107679	200	703
2426	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 21:55:22.993514	2025-04-20 21:55:22.993611	304	0
2427	\N	/api/washingmachines	GET	2025-04-20 21:55:23.359318	2025-04-20 21:55:24.050945	200	691
2428	\N	/api/washingmachines	GET	2025-04-20 22:00:24.178093	2025-04-20 22:00:24.896453	200	718
2429	\N	/api/washingmachines	GET	2025-04-20 22:03:15.719059	2025-04-20 22:03:16.441511	200	722
2430	\N	/api/washingmachines	GET	2025-04-20 22:03:47.178925	2025-04-20 22:03:47.884376	200	705
2431	\N	/api/washingmachines	GET	2025-04-20 22:08:47.447827	2025-04-20 22:08:48.154921	200	707
2432	\N	/api/washingmachines	GET	2025-04-20 22:11:10.189079	2025-04-20 22:11:10.876391	200	687
2433	\N	/api/restaurant	GET	2025-04-20 22:14:19.100606	2025-04-20 22:14:19.1007	200	0
2434	\N	/api/restaurant	GET	2025-04-20 22:14:19.268379	2025-04-20 22:14:19.268435	200	0
2436	\N	/api/restaurant	GET	2025-04-20 22:14:23.230694	2025-04-20 22:14:23.230744	200	0
2437	\N	/status	GET	2025-04-20 22:14:30.618487	2025-04-20 22:14:30.618488	200	0
2438	\N	/status	GET	2025-04-20 22:14:30.791562	2025-04-20 22:14:30.791563	200	0
2439	\N	/api/statistics/global	GET	2025-04-20 22:14:30.929333	2025-04-20 22:14:30.932294	200	2
2440	\N	/api/statistics/global	GET	2025-04-20 22:14:31.098224	2025-04-20 22:14:31.107289	200	9
2451	\N	/api/washingmachines	GET	2025-04-20 22:28:13.683685	2025-04-20 22:28:14.366911	200	683
2452	\N	/api/washingmachines	GET	2025-04-20 22:30:03.375892	2025-04-20 22:30:04.149452	200	773
2453	\N	/api/washingmachines	GET	2025-04-20 22:30:18.848612	2025-04-20 22:30:19.478131	200	629
2454	\N	/api/washingmachines	GET	2025-04-20 22:30:30.926159	2025-04-20 22:30:31.629224	200	703
2455	\N	/api/washingmachines	GET	2025-04-20 22:32:39.766435	2025-04-20 22:32:40.423561	200	657
2456	\N	/api/washingmachines	GET	2025-04-20 22:36:08.577997	2025-04-20 22:36:09.262801	200	684
2457	\N	/api/washingmachines	GET	2025-04-20 22:40:56.336856	2025-04-20 22:40:57.122119	200	785
2458	\N	/api/washingmachines	GET	2025-04-20 22:45:56.144217	2025-04-20 22:45:56.843393	200	699
2459	\N	/api/washingmachines	GET	2025-04-20 22:47:46.837063	2025-04-20 22:47:47.62621	200	789
2460	\N	/api/washingmachines	GET	2025-04-20 22:47:47.804984	2025-04-20 22:47:47.977918	200	172
2461	\N	/api/restaurant	GET	2025-04-20 22:47:59.800406	2025-04-20 22:47:59.800481	200	0
2462	\N	/api/restaurant	GET	2025-04-20 22:47:59.969585	2025-04-20 22:47:59.969641	200	0
2463	\N	/api/washingmachines	GET	2025-04-20 22:48:11.423957	2025-04-20 22:48:12.192973	200	769
2464	\N	/api/washingmachines	GET	2025-04-20 22:48:12.373209	2025-04-20 22:48:12.556283	200	183
2465	\N	/status	GET	2025-04-20 22:48:28.80645	2025-04-20 22:48:28.806482	200	0
2466	\N	/status	GET	2025-04-20 22:48:29.004465	2025-04-20 22:48:29.004467	200	0
2467	\N	/api/statistics/global	GET	2025-04-20 22:48:29.144776	2025-04-20 22:48:29.17575	200	30
2469	\N	/api/statistics/global	GET	2025-04-20 22:48:29.406627	2025-04-20 22:48:29.462873	200	56
2484	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 22:50:01.248891	2025-04-20 22:50:01.248984	304	0
2485	\N	/api/washingmachines	GET	2025-04-20 22:50:01.622553	2025-04-20 22:50:02.243853	200	621
2486	\N	/api/washingmachines	GET	2025-04-20 22:50:37.539878	2025-04-20 22:50:38.243743	200	703
2487	\N	/api/washingmachines	GET	2025-04-20 22:51:46.519913	2025-04-20 22:51:47.222185	200	702
2488	\N	/api/washingmachines	GET	2025-04-20 22:53:01.405492	2025-04-20 22:53:02.03075	200	625
2489	\N	/api/washingmachines	GET	2025-04-20 22:54:34.016954	2025-04-20 22:54:34.639947	200	622
2490	\N	/api/washingmachines	GET	2025-04-20 22:55:42.652755	2025-04-20 22:55:43.342427	200	689
2491	\N	/api/washingmachines	GET	2025-04-20 22:58:11.766108	2025-04-20 22:58:12.467605	200	701
2492	\N	/api/washingmachines	GET	2025-04-20 22:58:57.41608	2025-04-20 22:58:58.031428	200	615
2493	\N	/api/washingmachines	GET	2025-04-20 22:59:17.989105	2025-04-20 22:59:18.596328	200	607
2494	\N	/api/washingmachines	GET	2025-04-20 22:59:58.116665	2025-04-20 22:59:58.74027	200	623
2495	\N	/api/washingmachines	GET	2025-04-20 23:00:23.768663	2025-04-20 23:00:24.473577	200	704
2496	\N	/api/washingmachines	GET	2025-04-20 23:05:08.180769	2025-04-20 23:05:08.87898	200	698
2497	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 23:07:08.413579	2025-04-20 23:07:08.413689	304	0
2498	\N	/api/washingmachines	GET	2025-04-20 23:07:08.515567	2025-04-20 23:07:09.237817	200	722
2499	\N	/api/washingmachines	GET	2025-04-20 23:12:08.767076	2025-04-20 23:12:09.462532	200	695
2500	\N	/api/washingmachines	GET	2025-04-20 23:12:52.587145	2025-04-20 23:12:53.264703	200	677
2501	\N	/api/washingmachines	GET	2025-04-20 23:13:22.862011	2025-04-20 23:13:23.618243	200	756
2502	\N	/api/washingmachines	GET	2025-04-20 23:18:23.029902	2025-04-20 23:18:23.725016	200	695
2503	\N	/api/washingmachines	GET	2025-04-20 23:19:36.685137	2025-04-20 23:19:37.398792	200	713
3657	\N	/status	GET	2025-04-28 12:23:32.659677	2025-04-28 12:23:32.659679	200	0
2504	\N	/api/washingmachines	GET	2025-04-20 23:19:41.483243	2025-04-20 23:19:41.660405	200	177
2505	\N	/api/washingmachines	GET	2025-04-20 23:19:46.681688	2025-04-20 23:19:47.372236	200	690
2506	\N	/api/washingmachines	GET	2025-04-20 23:19:51.727352	2025-04-20 23:19:51.901774	200	174
2507	\N	/api/washingmachines	GET	2025-04-20 23:19:56.354463	2025-04-20 23:19:56.537097	200	182
2508	\N	/api/washingmachines	GET	2025-04-20 23:20:01.352259	2025-04-20 23:20:01.517048	200	164
2509	\N	/api/washingmachines	GET	2025-04-20 23:20:06.468773	2025-04-20 23:20:06.505812	500	37
2510	\N	/api/washingmachines	GET	2025-04-20 23:20:11.330422	2025-04-20 23:20:11.950383	200	619
2511	\N	/api/washingmachines	GET	2025-04-20 23:20:16.36732	2025-04-20 23:20:16.540183	200	172
2512	\N	/api/washingmachines	GET	2025-04-20 23:20:21.326619	2025-04-20 23:20:21.500923	200	174
2513	\N	/api/washingmachines	GET	2025-04-20 23:20:26.362378	2025-04-20 23:20:26.544963	200	182
2514	\N	/api/washingmachines	GET	2025-04-20 23:20:31.355342	2025-04-20 23:20:31.540405	200	185
2515	\N	/api/washingmachines	GET	2025-04-20 23:20:36.483874	2025-04-20 23:20:36.565194	500	81
2516	\N	/api/washingmachines	GET	2025-04-20 23:20:41.363843	2025-04-20 23:20:42.052578	200	688
2517	\N	/api/washingmachines	GET	2025-04-20 23:20:46.355542	2025-04-20 23:20:46.541423	200	185
2518	\N	/api/washingmachines	GET	2025-04-20 23:20:51.727227	2025-04-20 23:20:52.408786	200	681
2519	\N	/api/washingmachines	GET	2025-04-20 23:20:56.349414	2025-04-20 23:20:56.532701	200	183
2520	\N	/api/washingmachines	GET	2025-04-20 23:21:01.365313	2025-04-20 23:21:01.537961	200	172
2521	\N	/api/washingmachines	GET	2025-04-20 23:21:06.343515	2025-04-20 23:21:06.529695	200	186
2522	\N	/api/washingmachines	GET	2025-04-20 23:21:11.339875	2025-04-20 23:21:11.513545	200	173
2523	\N	/api/washingmachines	GET	2025-04-20 23:21:16.353492	2025-04-20 23:21:16.527657	200	174
2524	\N	/api/washingmachines	GET	2025-04-20 23:21:21.345638	2025-04-20 23:21:21.530184	200	184
2525	\N	/api/washingmachines	GET	2025-04-20 23:21:26.346696	2025-04-20 23:21:26.482142	500	135
2526	\N	/api/washingmachines	GET	2025-04-20 23:21:31.342485	2025-04-20 23:21:31.993253	200	650
2527	\N	/api/washingmachines	GET	2025-04-20 23:21:36.351989	2025-04-20 23:21:36.529258	200	177
2528	\N	/api/washingmachines	GET	2025-04-20 23:21:41.355582	2025-04-20 23:21:41.470115	500	114
2529	\N	/api/washingmachines	GET	2025-04-20 23:21:46.341299	2025-04-20 23:21:47.046384	200	705
2530	\N	/api/washingmachines	GET	2025-04-20 23:21:51.369288	2025-04-20 23:21:51.54533	200	176
2531	\N	/api/washingmachines	GET	2025-04-20 23:21:56.338916	2025-04-20 23:21:56.482206	500	143
2532	\N	/api/washingmachines	GET	2025-04-20 23:22:01.356015	2025-04-20 23:22:01.969965	200	613
2533	\N	/api/washingmachines	GET	2025-04-20 23:22:02.149293	2025-04-20 23:22:02.317031	200	167
2534	\N	/api/washingmachines	GET	2025-04-20 23:22:06.521115	2025-04-20 23:22:06.688295	200	167
2535	\N	/api/washingmachines	GET	2025-04-20 23:22:11.56674	2025-04-20 23:22:11.62843	500	61
2536	\N	/api/washingmachines	GET	2025-04-20 23:22:16.519207	2025-04-20 23:22:17.140102	200	620
2537	\N	/api/washingmachines	GET	2025-04-20 23:22:21.544411	2025-04-20 23:22:21.717855	200	173
2538	\N	/api/washingmachines	GET	2025-04-20 23:22:26.602791	2025-04-20 23:22:26.70881	500	106
2539	\N	/api/washingmachines	GET	2025-04-20 23:22:31.52461	2025-04-20 23:22:32.149548	200	624
2540	\N	/api/washingmachines	GET	2025-04-20 23:22:36.531879	2025-04-20 23:22:36.719894	200	188
2541	\N	/api/washingmachines	GET	2025-04-20 23:22:41.54121	2025-04-20 23:22:41.730648	200	189
2542	\N	/api/washingmachines	GET	2025-04-20 23:22:46.561711	2025-04-20 23:22:46.676371	500	114
2543	\N	/api/washingmachines	GET	2025-04-20 23:22:51.541792	2025-04-20 23:22:52.153522	200	611
2544	\N	/api/washingmachines	GET	2025-04-20 23:22:56.536688	2025-04-20 23:22:56.722436	200	185
2545	\N	/api/washingmachines	GET	2025-04-20 23:23:01.549005	2025-04-20 23:23:01.677139	500	128
2546	\N	/api/washingmachines	GET	2025-04-20 23:23:02.193318	2025-04-20 23:23:02.811912	200	618
2547	\N	/api/washingmachines	GET	2025-04-20 23:24:29.206466	2025-04-20 23:24:29.901937	200	695
2548	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 23:26:17.004203	2025-04-20 23:26:17.004324	304	0
2549	\N	/api/washingmachines	GET	2025-04-20 23:26:17.305781	2025-04-20 23:26:17.929808	200	624
2550	\N	/api/washingmachines	GET	2025-04-20 23:29:12.577738	2025-04-20 23:29:13.280166	200	702
2551	\N	/api/washingmachines	GET	2025-04-20 23:30:33.585921	2025-04-20 23:30:34.334778	200	748
2552	\N	/api/restaurant	GET	2025-04-20 23:30:38.094514	2025-04-20 23:30:38.094588	200	0
2553	\N	/api/restaurant	GET	2025-04-20 23:30:38.269271	2025-04-20 23:30:38.269341	200	0
2554	\N	/status	GET	2025-04-20 23:30:40.610449	2025-04-20 23:30:40.610454	200	0
2555	\N	/status	GET	2025-04-20 23:30:40.781199	2025-04-20 23:30:40.7812	200	0
2556	\N	/api/statistics/global	GET	2025-04-20 23:30:40.917131	2025-04-20 23:30:40.919041	200	1
2557	\N	/api/statistics/global	GET	2025-04-20 23:30:41.085815	2025-04-20 23:30:41.112456	200	26
2566	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 23:38:10.179125	2025-04-20 23:38:10.179226	200	0
2568	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 23:38:25.495604	2025-04-20 23:38:25.495717	200	0
2569	\N	/api/washingmachines	GET	2025-04-20 23:38:40.108091	2025-04-20 23:38:40.783417	200	675
2570	\N	/api/restaurant	GET	2025-04-20 23:38:57.638434	2025-04-20 23:38:57.638489	200	0
2571	\N	/status	GET	2025-04-20 23:39:05.008259	2025-04-20 23:39:05.00826	200	0
2572	\N	/api/statistics/global	GET	2025-04-20 23:39:05.224367	2025-04-20 23:39:05.226273	200	1
2573	\N	/api/statistics/endpoints	GET	2025-04-20 23:39:05.408366	2025-04-20 23:39:05.411278	200	2
2574	\N	/status	GET	2025-04-20 23:39:10.010213	2025-04-20 23:39:10.010215	200	0
2575	\N	/status	GET	2025-04-20 23:39:10.196329	2025-04-20 23:39:10.19633	200	0
2576	\N	/api/statistics/global	GET	2025-04-20 23:39:10.330248	2025-04-20 23:39:10.338304	200	8
2577	\N	/api/statistics/endpoints	GET	2025-04-20 23:39:10.528491	2025-04-20 23:39:10.53186	200	3
2578	\N	/status	GET	2025-04-20 23:39:15.018046	2025-04-20 23:39:15.01805	200	0
2579	\N	/status	GET	2025-04-20 23:39:15.215869	2025-04-20 23:39:15.215871	200	0
2583	\N	/status	GET	2025-04-20 23:39:20.195868	2025-04-20 23:39:20.195869	200	0
2592	\N	/api/washingmachines	GET	2025-04-20 23:47:06.186253	2025-04-20 23:47:06.885897	200	699
2593	\N	/api/restaurant	GET	2025-04-20 23:47:11.902537	2025-04-20 23:47:11.902597	200	0
2594	\N	/api/restaurant	GET	2025-04-20 23:47:12.126267	2025-04-20 23:47:12.126328	200	0
2595	\N	/api/restaurant	GET	2025-04-20 23:47:22.351396	2025-04-20 23:47:22.351466	200	0
2596	\N	/status	GET	2025-04-20 23:47:25.090043	2025-04-20 23:47:25.090051	200	0
2597	\N	/status	GET	2025-04-20 23:47:25.265608	2025-04-20 23:47:25.265609	200	0
2598	\N	/api/statistics/global	GET	2025-04-20 23:47:25.414623	2025-04-20 23:47:25.422677	200	8
2599	\N	/api/statistics/global	GET	2025-04-20 23:47:25.588711	2025-04-20 23:47:25.590035	200	1
2600	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:25.607026	2025-04-20 23:47:25.609578	200	2
2601	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:25.776811	2025-04-20 23:47:25.779732	200	2
2602	\N	/status	GET	2025-04-20 23:47:30.161234	2025-04-20 23:47:30.161237	200	0
2603	\N	/status	GET	2025-04-20 23:47:30.335358	2025-04-20 23:47:30.335359	200	0
2604	\N	/api/statistics/global	GET	2025-04-20 23:47:30.474791	2025-04-20 23:47:30.476196	200	1
2605	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:30.646967	2025-04-20 23:47:30.649854	200	2
2606	\N	/status	GET	2025-04-20 23:47:35.336424	2025-04-20 23:47:35.336427	200	0
2607	\N	/status	GET	2025-04-20 23:47:35.507978	2025-04-20 23:47:35.507979	200	0
2608	\N	/api/statistics/global	GET	2025-04-20 23:47:35.644864	2025-04-20 23:47:35.647512	200	2
2609	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:35.842284	2025-04-20 23:47:35.84597	200	3
2610	\N	/status	GET	2025-04-20 23:47:40.164226	2025-04-20 23:47:40.164229	200	0
2611	\N	/status	GET	2025-04-20 23:47:40.34189	2025-04-20 23:47:40.341892	200	0
2612	\N	/api/statistics/global	GET	2025-04-20 23:47:40.477607	2025-04-20 23:47:40.478981	200	1
2613	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:40.645342	2025-04-20 23:47:40.64774	200	2
2614	\N	/status	GET	2025-04-20 23:47:45.191016	2025-04-20 23:47:45.191018	200	0
2615	\N	/status	GET	2025-04-20 23:47:45.359757	2025-04-20 23:47:45.359759	200	0
2616	\N	/api/statistics/global	GET	2025-04-20 23:47:45.389654	2025-04-20 23:47:45.43057	200	40
2617	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:45.611355	2025-04-20 23:47:45.614751	200	3
2618	\N	/status	GET	2025-04-20 23:47:50.18814	2025-04-20 23:47:50.188142	200	0
2619	\N	/status	GET	2025-04-20 23:47:50.357941	2025-04-20 23:47:50.357942	200	0
2620	\N	/api/statistics/global	GET	2025-04-20 23:47:50.375735	2025-04-20 23:47:50.37919	200	3
2621	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:50.564757	2025-04-20 23:47:50.567318	200	2
2622	\N	/status	GET	2025-04-20 23:47:55.183866	2025-04-20 23:47:55.183869	200	0
2623	\N	/status	GET	2025-04-20 23:47:55.352011	2025-04-20 23:47:55.352012	200	0
2558	\N	/api/statistics/endpoints	GET	2025-04-20 23:30:41.108378	2025-04-20 23:30:41.112648	200	4
2559	\N	/api/statistics/endpoints	GET	2025-04-20 23:30:41.283068	2025-04-20 23:30:41.286163	200	3
2560	\N	/api/washingmachines	GET	2025-04-20 23:30:42.463607	2025-04-20 23:30:43.168379	200	704
2561	\N	/api/washingmachines	GET	2025-04-20 23:30:43.353295	2025-04-20 23:30:43.538872	200	185
2562	\N	/api/washingmachines	GET	2025-04-20 23:33:38.147545	2025-04-20 23:33:38.86863	200	721
2563	\N	/api/washingmachines	GET	2025-04-20 23:34:34.562081	2025-04-20 23:34:35.256221	200	694
2564	\N	/api/washingmachines	GET	2025-04-20 23:34:59.551454	2025-04-20 23:35:00.197298	200	645
2565	\N	/api/washingmachines	GET	2025-04-20 23:36:32.605729	2025-04-20 23:36:33.236043	200	630
2567	\N	/api/washingmachines	GET	2025-04-20 23:38:14.31388	2025-04-20 23:38:15.005062	200	691
2580	\N	/api/statistics/global	GET	2025-04-20 23:39:15.215968	2025-04-20 23:39:15.218281	200	2
2581	\N	/api/statistics/endpoints	GET	2025-04-20 23:39:15.415827	2025-04-20 23:39:15.420031	200	4
2582	\N	/status	GET	2025-04-20 23:39:19.996584	2025-04-20 23:39:19.996585	200	0
2584	\N	/api/statistics/global	GET	2025-04-20 23:39:20.195856	2025-04-20 23:39:20.19737	200	1
2585	\N	/api/statistics/endpoints	GET	2025-04-20 23:39:20.387913	2025-04-20 23:39:20.390167	200	2
2586	\N	/api/washingmachines	GET	2025-04-20 23:41:33.380228	2025-04-20 23:41:34.096058	200	715
2587	\N	/api/washingmachines	GET	2025-04-20 23:42:28.65684	2025-04-20 23:42:29.577703	200	920
2588	\N	/api/washingmachines	GET	2025-04-20 23:45:20.8126	2025-04-20 23:45:21.418495	200	605
2589	\N	/api/restaurant	GET	2025-04-20 23:45:47.688064	2025-04-20 23:45:47.68813	200	0
2590	\N	/api/washingmachines	GET	2025-04-20 23:46:33.144138	2025-04-20 23:46:33.762849	200	618
2591	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-20 23:47:05.778797	2025-04-20 23:47:05.778882	304	0
2644	\N	/api/statistics/global	GET	2025-04-20 23:48:20.357227	2025-04-20 23:48:20.359166	200	1
2645	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:20.581591	2025-04-20 23:48:20.585674	200	4
2646	\N	/status	GET	2025-04-20 23:48:25.151517	2025-04-20 23:48:25.151518	200	0
2647	\N	/status	GET	2025-04-20 23:48:25.318814	2025-04-20 23:48:25.318815	200	0
2648	\N	/api/statistics/global	GET	2025-04-20 23:48:25.336883	2025-04-20 23:48:25.338472	200	1
2649	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:25.542351	2025-04-20 23:48:25.544645	200	2
2650	\N	/status	GET	2025-04-20 23:48:30.167391	2025-04-20 23:48:30.167392	200	0
2651	\N	/status	GET	2025-04-20 23:48:30.348195	2025-04-20 23:48:30.348197	200	0
2652	\N	/api/statistics/global	GET	2025-04-20 23:48:30.366737	2025-04-20 23:48:30.368167	200	1
2653	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:30.563366	2025-04-20 23:48:30.566804	200	3
2654	\N	/status	GET	2025-04-20 23:48:35.170717	2025-04-20 23:48:35.170718	200	0
2655	\N	/status	GET	2025-04-20 23:48:35.343259	2025-04-20 23:48:35.34326	200	0
2656	\N	/api/statistics/global	GET	2025-04-20 23:48:35.363271	2025-04-20 23:48:35.364624	200	1
2657	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:35.55666	2025-04-20 23:48:35.559556	200	2
2658	\N	/status	GET	2025-04-20 23:48:40.160901	2025-04-20 23:48:40.160904	200	0
2659	\N	/api/statistics/global	GET	2025-04-20 23:48:40.359887	2025-04-20 23:48:40.361855	200	1
2685	\N	/api/statistics/endpoints	GET	2025-04-20 23:51:44.09043	2025-04-20 23:51:44.093932	200	3
2686	\N	/api/statistics/endpoints	GET	2025-04-20 23:51:44.26383	2025-04-20 23:51:44.266992	200	3
2687	\N	/api/washingmachines	GET	2025-04-20 23:51:48.514915	2025-04-20 23:51:49.277242	200	762
2688	\N	/api/washingmachines	GET	2025-04-20 23:51:49.464559	2025-04-20 23:51:49.638246	200	173
2689	\N	/status	GET	2025-04-20 23:53:41.71792	2025-04-20 23:53:41.717938	200	0
2690	\N	/status	GET	2025-04-20 23:53:41.886137	2025-04-20 23:53:41.886138	200	0
2691	\N	/api/statistics/global	GET	2025-04-20 23:53:42.022695	2025-04-20 23:53:42.024525	200	1
2692	\N	/api/statistics/global	GET	2025-04-20 23:53:42.197211	2025-04-20 23:53:42.20423	200	7
2720	\N	/api/washingmachines	GET	2025-04-20 23:54:33.971823	2025-04-20 23:54:34.156751	200	184
2721	\N	/api/statistics/global	GET	2025-04-20 23:54:34.266672	2025-04-20 23:54:34.308118	200	41
2722	\N	/api/statistics/global	GET	2025-04-20 23:54:34.572019	2025-04-20 23:54:34.575777	200	3
2723	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:34.590019	2025-04-20 23:54:34.593165	200	3
2724	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:34.759601	2025-04-20 23:54:34.770843	200	11
2725	\N	/status	GET	2025-04-20 23:54:39.174248	2025-04-20 23:54:39.174251	200	0
2726	\N	/status	GET	2025-04-20 23:54:39.407853	2025-04-20 23:54:39.407855	200	0
2727	\N	/api/statistics/global	GET	2025-04-20 23:54:39.546309	2025-04-20 23:54:39.547844	200	1
2728	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:39.728885	2025-04-20 23:54:39.733146	200	4
2729	\N	/status	GET	2025-04-20 23:54:44.041935	2025-04-20 23:54:44.041937	200	0
2730	\N	/status	GET	2025-04-20 23:54:44.2181	2025-04-20 23:54:44.218102	200	0
2731	\N	/api/statistics/global	GET	2025-04-20 23:54:44.236878	2025-04-20 23:54:44.243539	200	6
2732	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:44.441757	2025-04-20 23:54:44.44405	200	2
2733	\N	/status	GET	2025-04-20 23:54:49.037265	2025-04-20 23:54:49.037269	200	0
2734	\N	/status	GET	2025-04-20 23:54:49.20894	2025-04-20 23:54:49.208942	200	0
2735	\N	/api/statistics/global	GET	2025-04-20 23:54:49.226677	2025-04-20 23:54:49.233561	200	6
2736	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:49.421711	2025-04-20 23:54:49.424712	200	3
2737	\N	/status	GET	2025-04-20 23:54:54.033559	2025-04-20 23:54:54.033561	200	0
2738	\N	/status	GET	2025-04-20 23:54:54.204623	2025-04-20 23:54:54.204625	200	0
2739	\N	/api/statistics/global	GET	2025-04-20 23:54:54.223188	2025-04-20 23:54:54.22564	200	2
2740	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:54.413588	2025-04-20 23:54:54.41684	200	3
2741	\N	/status	GET	2025-04-20 23:54:59.010341	2025-04-20 23:54:59.010343	200	0
2742	\N	/status	GET	2025-04-20 23:54:59.185461	2025-04-20 23:54:59.185463	200	0
2746	\N	/api/statistics/global	GET	2025-04-20 23:55:04.22002	2025-04-20 23:55:04.224695	200	4
2751	\N	/api/statistics/global	GET	2025-04-20 23:55:09.2166	2025-04-20 23:55:09.21903	200	2
2752	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:09.427604	2025-04-20 23:55:09.429919	200	2
2753	\N	/status	GET	2025-04-20 23:55:14.049099	2025-04-20 23:55:14.049101	200	0
2754	\N	/status	GET	2025-04-20 23:55:14.226618	2025-04-20 23:55:14.22662	200	0
2755	\N	/api/statistics/global	GET	2025-04-20 23:55:14.244556	2025-04-20 23:55:14.246794	200	2
2756	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:14.469204	2025-04-20 23:55:14.473016	200	3
2757	\N	/status	GET	2025-04-20 23:55:19.033223	2025-04-20 23:55:19.033224	200	0
2758	\N	/status	GET	2025-04-20 23:55:19.205029	2025-04-20 23:55:19.205031	200	0
2759	\N	/api/statistics/global	GET	2025-04-20 23:55:19.223216	2025-04-20 23:55:19.225317	200	2
2760	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:19.418517	2025-04-20 23:55:19.421233	200	2
2761	\N	/status	GET	2025-04-20 23:55:24.033762	2025-04-20 23:55:24.033764	200	0
2762	\N	/status	GET	2025-04-20 23:55:24.266077	2025-04-20 23:55:24.266079	200	0
2763	\N	/api/statistics/global	GET	2025-04-20 23:55:24.278353	2025-04-20 23:55:24.280119	200	1
2764	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:24.484597	2025-04-20 23:55:24.487174	200	2
2878	\N	/api/statistics/global	GET	2025-04-20 23:57:48.13685	2025-04-20 23:57:48.138217	200	1
2879	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:48.342571	2025-04-20 23:57:48.34547	200	2
2880	\N	/status	GET	2025-04-20 23:57:53.069298	2025-04-20 23:57:53.0693	200	0
2881	\N	/status	GET	2025-04-20 23:57:53.355507	2025-04-20 23:57:53.355509	200	0
2882	\N	/api/statistics/global	GET	2025-04-20 23:57:53.373163	2025-04-20 23:57:53.374848	200	1
2883	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:53.569907	2025-04-20 23:57:53.57338	200	3
2884	\N	/status	GET	2025-04-20 23:57:57.944543	2025-04-20 23:57:57.944545	200	0
2885	\N	/status	GET	2025-04-20 23:57:58.121196	2025-04-20 23:57:58.121197	200	0
2886	\N	/api/statistics/global	GET	2025-04-20 23:57:58.138899	2025-04-20 23:57:58.141293	200	2
2887	\N	/api/statistics/endpoints	GET	2025-04-20 23:57:58.333193	2025-04-20 23:57:58.346844	200	13
2888	\N	/status	GET	2025-04-20 23:58:03.024508	2025-04-20 23:58:03.024509	200	0
2889	\N	/status	GET	2025-04-20 23:58:03.192166	2025-04-20 23:58:03.192167	200	0
2890	\N	/api/statistics/global	GET	2025-04-20 23:58:03.209779	2025-04-20 23:58:03.211994	200	2
2891	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:03.407384	2025-04-20 23:58:03.409802	200	2
2892	\N	/status	GET	2025-04-20 23:58:07.942321	2025-04-20 23:58:07.942323	200	0
2893	\N	/status	GET	2025-04-20 23:58:08.109598	2025-04-20 23:58:08.1096	200	0
2894	\N	/api/statistics/global	GET	2025-04-20 23:58:08.12723	2025-04-20 23:58:08.130773	200	3
2895	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:08.331509	2025-04-20 23:58:08.3366	200	5
2896	\N	/status	GET	2025-04-20 23:58:13.019404	2025-04-20 23:58:13.019406	200	0
2624	\N	/api/statistics/global	GET	2025-04-20 23:47:55.368164	2025-04-20 23:47:55.369458	200	1
2625	\N	/api/statistics/endpoints	GET	2025-04-20 23:47:55.571056	2025-04-20 23:47:55.573244	200	2
2626	\N	/status	GET	2025-04-20 23:48:00.17528	2025-04-20 23:48:00.175282	200	0
2627	\N	/status	GET	2025-04-20 23:48:00.367952	2025-04-20 23:48:00.367954	200	0
2628	\N	/api/statistics/global	GET	2025-04-20 23:48:00.386217	2025-04-20 23:48:00.388218	200	2
2629	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:00.583518	2025-04-20 23:48:00.586055	200	2
2630	\N	/status	GET	2025-04-20 23:48:05.18866	2025-04-20 23:48:05.188661	200	0
2631	\N	/status	GET	2025-04-20 23:48:05.357474	2025-04-20 23:48:05.357475	200	0
2632	\N	/api/statistics/global	GET	2025-04-20 23:48:05.375956	2025-04-20 23:48:05.377296	200	1
2633	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:05.572112	2025-04-20 23:48:05.574353	200	2
2634	\N	/status	GET	2025-04-20 23:48:10.165302	2025-04-20 23:48:10.165304	200	0
2635	\N	/status	GET	2025-04-20 23:48:10.342582	2025-04-20 23:48:10.342584	200	0
2636	\N	/api/statistics/global	GET	2025-04-20 23:48:10.360466	2025-04-20 23:48:10.361816	200	1
2637	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:10.559395	2025-04-20 23:48:10.561622	200	2
2638	\N	/status	GET	2025-04-20 23:48:15.167888	2025-04-20 23:48:15.167892	200	0
2639	\N	/status	GET	2025-04-20 23:48:15.340366	2025-04-20 23:48:15.340368	200	0
2640	\N	/api/statistics/global	GET	2025-04-20 23:48:15.477249	2025-04-20 23:48:15.479622	200	2
2641	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:15.660625	2025-04-20 23:48:15.665183	200	4
2642	\N	/status	GET	2025-04-20 23:48:20.184663	2025-04-20 23:48:20.184667	200	0
2643	\N	/status	GET	2025-04-20 23:48:20.357226	2025-04-20 23:48:20.357227	200	0
2660	\N	/status	GET	2025-04-20 23:48:40.359887	2025-04-20 23:48:40.359889	200	0
2661	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:40.536641	2025-04-20 23:48:40.54012	200	3
2662	\N	/status	GET	2025-04-20 23:48:45.180526	2025-04-20 23:48:45.180528	200	0
2663	\N	/status	GET	2025-04-20 23:48:45.356188	2025-04-20 23:48:45.356189	200	0
2664	\N	/api/statistics/global	GET	2025-04-20 23:48:45.374042	2025-04-20 23:48:45.375464	200	1
2665	\N	/api/statistics/endpoints	GET	2025-04-20 23:48:45.573861	2025-04-20 23:48:45.576492	200	2
2666	\N	/api/washingmachines	GET	2025-04-20 23:48:47.782093	2025-04-20 23:48:48.558915	200	776
2667	\N	/api/washingmachines	GET	2025-04-20 23:48:48.728785	2025-04-20 23:48:48.914009	200	185
2668	\N	/api/restaurant	GET	2025-04-20 23:48:52.968145	2025-04-20 23:48:52.968206	200	0
2669	\N	/api/restaurant	GET	2025-04-20 23:48:53.135626	2025-04-20 23:48:53.135699	200	0
2670	\N	/api/washingmachines	GET	2025-04-20 23:49:09.030464	2025-04-20 23:49:09.717743	200	687
2671	\N	/api/washingmachines	GET	2025-04-20 23:49:09.885224	2025-04-20 23:49:10.06287	200	177
2672	\N	/api/washingmachines	GET	2025-04-20 23:49:50.913412	2025-04-20 23:49:51.534546	200	621
2673	\N	/api/washingmachines	GET	2025-04-20 23:49:53.235528	2025-04-20 23:49:53.436582	200	201
2674	\N	/api/washingmachines	GET	2025-04-20 23:49:54.445851	2025-04-20 23:49:54.617859	200	172
2675	\N	/api/restaurant	GET	2025-04-20 23:50:09.237893	2025-04-20 23:50:09.237945	200	0
2676	\N	/api/washingmachines	GET	2025-04-20 23:50:11.933657	2025-04-20 23:50:12.720617	200	786
2677	\N	/api/washingmachines	GET	2025-04-20 23:51:15.793545	2025-04-20 23:51:16.499711	200	706
2678	\N	/api/washingmachines	GET	2025-04-20 23:51:37.34768	2025-04-20 23:51:38.114547	200	766
2679	\N	/api/restaurant	GET	2025-04-20 23:51:42.009624	2025-04-20 23:51:42.009699	200	0
2680	\N	/api/restaurant	GET	2025-04-20 23:51:42.323683	2025-04-20 23:51:42.323737	200	0
2681	\N	/status	GET	2025-04-20 23:51:43.724551	2025-04-20 23:51:43.724556	200	0
2682	\N	/status	GET	2025-04-20 23:51:43.895197	2025-04-20 23:51:43.895199	200	0
2683	\N	/api/statistics/global	GET	2025-04-20 23:51:43.913431	2025-04-20 23:51:43.915559	200	2
2684	\N	/api/statistics/global	GET	2025-04-20 23:51:44.090358	2025-04-20 23:51:44.092578	200	2
2693	\N	/api/statistics/endpoints	GET	2025-04-20 23:53:42.21622	2025-04-20 23:53:42.225617	200	9
2694	\N	/api/statistics/endpoints	GET	2025-04-20 23:53:42.410228	2025-04-20 23:53:42.412815	200	2
2695	\N	/status	GET	2025-04-20 23:53:52.341924	2025-04-20 23:53:52.341926	200	0
2696	\N	/api/statistics/global	GET	2025-04-20 23:53:52.521614	2025-04-20 23:53:52.52319	200	1
2697	\N	/api/statistics/endpoints	GET	2025-04-20 23:53:52.703935	2025-04-20 23:53:52.706547	200	2
2698	\N	/status	GET	2025-04-20 23:54:10.949449	2025-04-20 23:54:10.949455	200	0
2699	\N	/api/statistics/global	GET	2025-04-20 23:54:11.142658	2025-04-20 23:54:11.171553	200	28
2700	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:11.393583	2025-04-20 23:54:11.407315	200	13
2701	\N	/status	GET	2025-04-20 23:54:15.988164	2025-04-20 23:54:15.988165	200	0
2702	\N	/status	GET	2025-04-20 23:54:16.156091	2025-04-20 23:54:16.156092	200	0
2703	\N	/api/statistics/global	GET	2025-04-20 23:54:16.174278	2025-04-20 23:54:16.176872	200	2
2704	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:16.35731	2025-04-20 23:54:16.359609	200	2
2705	\N	/status	GET	2025-04-20 23:54:20.961392	2025-04-20 23:54:20.96145	200	0
2706	\N	/status	GET	2025-04-20 23:54:21.179881	2025-04-20 23:54:21.179882	200	0
2707	\N	/api/statistics/global	GET	2025-04-20 23:54:21.1978	2025-04-20 23:54:21.199687	200	1
2708	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:21.392947	2025-04-20 23:54:21.395538	200	2
2709	\N	/status	GET	2025-04-20 23:54:25.955214	2025-04-20 23:54:25.955219	200	0
2710	\N	/status	GET	2025-04-20 23:54:26.130973	2025-04-20 23:54:26.130975	200	0
2711	\N	/api/statistics/global	GET	2025-04-20 23:54:26.270078	2025-04-20 23:54:26.321639	200	51
2712	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:26.608816	2025-04-20 23:54:26.611673	200	2
2713	\N	/status	GET	2025-04-20 23:54:30.972723	2025-04-20 23:54:30.972724	200	0
2714	\N	/status	GET	2025-04-20 23:54:31.139721	2025-04-20 23:54:31.139723	200	0
2715	\N	/api/statistics/global	GET	2025-04-20 23:54:31.157966	2025-04-20 23:54:31.16025	200	2
2716	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:31.348092	2025-04-20 23:54:31.351148	200	3
2717	\N	/api/washingmachines	GET	2025-04-20 23:54:33.112644	2025-04-20 23:54:33.804685	200	692
2718	\N	/status	GET	2025-04-20 23:54:33.947546	2025-04-20 23:54:33.947547	200	0
2719	\N	/status	GET	2025-04-20 23:54:34.11437	2025-04-20 23:54:34.114372	200	0
2743	\N	/api/statistics/global	GET	2025-04-20 23:54:59.203086	2025-04-20 23:54:59.206167	200	3
2744	\N	/api/statistics/endpoints	GET	2025-04-20 23:54:59.391821	2025-04-20 23:54:59.394359	200	2
2745	\N	/status	GET	2025-04-20 23:55:04.026294	2025-04-20 23:55:04.026298	200	0
2747	\N	/status	GET	2025-04-20 23:55:04.203855	2025-04-20 23:55:04.203857	200	0
2748	\N	/api/statistics/endpoints	GET	2025-04-20 23:55:04.456553	2025-04-20 23:55:04.459452	200	2
2749	\N	/status	GET	2025-04-20 23:55:09.025418	2025-04-20 23:55:09.02542	200	0
2750	\N	/status	GET	2025-04-20 23:55:09.198212	2025-04-20 23:55:09.198214	200	0
2897	\N	/status	GET	2025-04-20 23:58:13.191286	2025-04-20 23:58:13.191287	200	0
2898	\N	/api/statistics/global	GET	2025-04-20 23:58:13.209521	2025-04-20 23:58:13.211324	200	1
2899	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:13.402752	2025-04-20 23:58:13.406262	200	3
2900	\N	/status	GET	2025-04-20 23:58:17.939594	2025-04-20 23:58:17.939595	200	0
2901	\N	/status	GET	2025-04-20 23:58:18.113903	2025-04-20 23:58:18.113904	200	0
2902	\N	/api/statistics/global	GET	2025-04-20 23:58:18.127542	2025-04-20 23:58:18.129651	200	2
2903	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:18.327285	2025-04-20 23:58:18.358524	200	31
2904	\N	/status	GET	2025-04-20 23:58:23.286057	2025-04-20 23:58:23.286058	200	0
2905	\N	/status	GET	2025-04-20 23:58:23.456083	2025-04-20 23:58:23.456084	200	0
2906	\N	/api/statistics/global	GET	2025-04-20 23:58:23.473916	2025-04-20 23:58:23.475819	200	1
2907	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:23.669257	2025-04-20 23:58:23.671971	200	2
2908	\N	/status	GET	2025-04-20 23:58:27.954617	2025-04-20 23:58:27.954619	200	0
2909	\N	/status	GET	2025-04-20 23:58:28.125244	2025-04-20 23:58:28.125246	200	0
2910	\N	/api/statistics/global	GET	2025-04-20 23:58:28.145097	2025-04-20 23:58:28.150004	200	4
2911	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:28.360014	2025-04-20 23:58:28.363295	200	3
2912	\N	/status	GET	2025-04-20 23:58:33.00726	2025-04-20 23:58:33.007262	200	0
2913	\N	/status	GET	2025-04-20 23:58:33.176517	2025-04-20 23:58:33.176518	200	0
2914	\N	/api/statistics/global	GET	2025-04-20 23:58:33.194357	2025-04-20 23:58:33.202366	200	8
2915	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:33.396759	2025-04-20 23:58:33.399371	200	2
2916	\N	/status	GET	2025-04-20 23:58:38.015241	2025-04-20 23:58:38.015242	200	0
2917	\N	/status	GET	2025-04-20 23:58:38.183821	2025-04-20 23:58:38.183823	200	0
2918	\N	/api/statistics/global	GET	2025-04-20 23:58:38.20272	2025-04-20 23:58:38.205052	200	2
2919	\N	/api/statistics/endpoints	GET	2025-04-20 23:58:38.405086	2025-04-20 23:58:38.408098	200	3
2920	\N	/status	GET	2025-04-20 23:58:42.949593	2025-04-20 23:58:42.949595	200	0
3007	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 11:48:07.247888	2025-04-21 11:48:07.260389	200	12
3008	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 11:48:08.969111	2025-04-21 11:48:08.971034	200	1
3009	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 11:48:11.405469	2025-04-21 11:48:11.40741	200	1
3010	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-04-21 11:48:13.74633	2025-04-21 11:48:13.74762	200	1
3011	\N	/	GET	2025-04-21 11:52:55.144832	2025-04-21 11:52:55.144834	500	0
3012	\N	/api/washingmachines	GET	2025-04-21 12:04:28.620487	2025-04-21 12:04:29.406359	200	785
3013	\N	/api/auth/login	POST	2025-04-21 13:01:48.134598	2025-04-21 13:01:48.204458	200	69
3014	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:01:48.420906	2025-04-21 13:01:48.422389	200	1
3015	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:01:48.604174	2025-04-21 13:01:48.607801	200	3
3016	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:01:48.966683	2025-04-21 13:01:48.968183	200	1
3017	\N	/api/restaurant	GET	2025-04-21 13:01:49.101002	2025-04-21 13:01:50.444634	200	1343
3019	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-04-21 13:01:56.521012	2025-04-21 13:01:56.5211	200	0
3020	\N	/api/restaurant	GET	2025-04-21 13:02:00.097189	2025-04-21 13:02:00.09725	200	0
3021	\N	/api/restaurant	GET	2025-04-21 13:02:44.109797	2025-04-21 13:02:44.109859	200	0
3022	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:02:54.01358	2025-04-21 13:02:54.017582	200	4
3023	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-04-21 13:02:56.33869	2025-04-21 13:02:56.347445	200	8
3024	\N	/api/restaurant	GET	2025-04-21 13:03:00.940376	2025-04-21 13:03:00.940445	200	0
3025	\N	/api/restaurant	GET	2025-04-21 13:21:27.718499	2025-04-21 13:21:27.718584	200	0
3026	\N	/api/restaurant	GET	2025-04-21 13:21:30.079878	2025-04-21 13:21:30.07994	200	0
3027	\N	/api/restaurant	GET	2025-04-21 13:21:37.87402	2025-04-21 13:21:37.874078	200	0
3028	\N	/api/restaurant	GET	2025-04-21 13:21:57.85843	2025-04-21 13:21:57.858492	200	0
3029	\N	/api/restaurant	GET	2025-04-21 13:22:00.435021	2025-04-21 13:22:00.435074	200	0
3030	\N	/api/restaurant	GET	2025-04-21 13:25:16.072615	2025-04-21 13:25:16.072697	200	0
3031	\N	/api/restaurant	GET	2025-04-21 13:25:20.602192	2025-04-21 13:25:20.602255	200	0
3032	\N	/api/restaurant	GET	2025-04-21 13:25:22.233302	2025-04-21 13:25:22.233356	200	0
3033	\N	/api/auth/login	POST	2025-04-21 13:51:07.016623	2025-04-21 13:51:07.08538	200	68
3034	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:51:07.292238	2025-04-21 13:51:07.307333	200	15
3035	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:51:07.641997	2025-04-21 13:51:07.643572	200	1
3036	\N	/api/restaurant	GET	2025-04-21 13:51:08.948724	2025-04-21 13:51:08.948778	200	0
3037	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:51:09.099966	2025-04-21 13:51:09.102035	200	2
3038	\N	/api/traq/	GET	2025-04-21 13:51:16.341191	2025-04-21 13:51:16.343742	200	2
3039	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:51:36.22787	2025-04-21 13:51:36.229503	200	1
3040	\N	/api/restaurant	GET	2025-04-21 13:51:57.61166	2025-04-21 13:51:57.611712	200	0
3041	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:52:03.98206	2025-04-21 13:52:03.983942	200	1
3042	\N	/api/auth/login	POST	2025-04-21 13:56:06.685989	2025-04-21 13:56:06.754282	200	68
3043	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:56:07.251342	2025-04-21 13:56:07.253342	200	1
3044	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:56:07.850156	2025-04-21 13:56:07.851816	200	1
3045	\N	/api/restaurant	GET	2025-04-21 13:56:08.988998	2025-04-21 13:56:08.98908	200	0
3046	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:56:09.007903	2025-04-21 13:56:09.010547	200	2
3047	\N	/api/traq/	GET	2025-04-21 13:56:17.031726	2025-04-21 13:56:17.032654	200	0
3048	test@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:56:35.497238	2025-04-21 13:56:35.498579	200	1
3049	\N	/api/restaurant	GET	2025-04-21 13:57:21.614205	2025-04-21 13:57:21.614263	200	0
3050	\N	/api/restaurant	GET	2025-04-21 13:57:52.469109	2025-04-21 13:57:52.469185	200	0
3052	\N	/api/restaurant	GET	2025-04-21 14:39:15.499686	2025-04-21 14:39:15.499774	200	0
3064	\N	/api/restaurant	GET	2025-04-21 15:59:25.502736	2025-04-21 15:59:25.5028	200	0
3110	\N	/api/restaurant	GET	2025-04-21 16:45:00.946245	2025-04-21 16:45:00.946289	200	0
3120	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:56:23.351804	2025-04-21 16:56:23.353749	200	1
3121	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:13.624501	2025-04-21 16:59:13.626107	200	1
3122	\N	/api/restaurant	GET	2025-04-21 16:59:14.01017	2025-04-21 16:59:14.01026	200	0
3123	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:14.145955	2025-04-21 16:59:14.148234	200	2
3124	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:31.264212	2025-04-21 16:59:31.284811	200	20
3125	\N	/api/restaurant	GET	2025-04-21 16:59:31.692889	2025-04-21 16:59:31.692941	200	0
3126	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:31.827306	2025-04-21 16:59:31.828676	200	1
3127	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:42.318212	2025-04-21 16:59:42.345596	200	27
3129	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:59:42.801539	2025-04-21 16:59:42.802992	200	1
3130	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:00:19.46145	2025-04-21 17:00:19.462974	200	1
3132	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:00:19.84451	2025-04-21 17:00:19.84597	200	1
3133	\N	/api/washingmachines	GET	2025-04-21 17:00:46.059636	2025-04-21 17:00:46.83175	200	772
3134	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:03:07.325084	2025-04-21 17:03:07.327102	200	2
3136	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:03:07.717941	2025-04-21 17:03:07.719482	200	1
3137	\N	/api/washingmachines	GET	2025-04-21 17:05:46.580482	2025-04-21 17:05:47.351848	200	771
3138	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:06:10.461134	2025-04-21 17:06:10.462709	200	1
3139	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:06:10.878796	2025-04-21 17:06:10.880387	200	1
3155	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 23:50:12.194697	2025-04-21 23:50:12.194798	200	0
3159	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 06:08:23.849138	2025-04-22 06:08:23.849255	200	0
3161	\N	/status	GET	2025-04-22 06:08:37.375607	2025-04-22 06:08:37.37561	200	0
3162	\N	/api/statistics/global	GET	2025-04-22 06:08:37.611991	2025-04-22 06:08:37.614365	200	2
3163	\N	/api/statistics/endpoints	GET	2025-04-22 06:08:37.780783	2025-04-22 06:08:37.783469	200	2
3164	\N	/status	GET	2025-04-22 06:08:41.975079	2025-04-22 06:08:41.975083	200	0
3165	\N	/status	GET	2025-04-22 06:08:42.147229	2025-04-22 06:08:42.14723	200	0
3166	\N	/api/statistics/global	GET	2025-04-22 06:08:42.285196	2025-04-22 06:08:42.287212	200	2
3167	\N	/api/statistics/endpoints	GET	2025-04-22 06:08:42.457918	2025-04-22 06:08:42.461409	200	3
3168	\N	/api/restaurant	GET	2025-04-22 06:08:42.849969	2025-04-22 06:08:42.850042	200	0
3169	\N	/api/washingmachines	GET	2025-04-22 06:08:44.46485	2025-04-22 06:08:45.258086	200	793
3170	\N	/api/washingmachines	GET	2025-04-22 06:09:12.2797	2025-04-22 06:09:13.055156	200	775
3171	\N	/api/restaurant	GET	2025-04-22 06:09:19.596254	2025-04-22 06:09:19.59632	200	0
3172	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-04-22 09:34:30.317873	2025-04-22 09:34:30.321934	200	4
3173	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-04-22 09:34:30.569517	2025-04-22 09:34:30.571671	200	2
3175	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 09:53:23.783739	2025-04-22 09:53:23.783833	200	0
3176	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 10:00:47.646794	2025-04-22 10:00:47.64692	200	0
3177	\N	/api/restaurant	GET	2025-04-22 10:00:57.684589	2025-04-22 10:00:57.684675	200	0
3178	\N	/api/auth/login	POST	2025-04-22 10:01:22.83606	2025-04-22 10:01:22.904142	200	68
3179	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:23.075986	2025-04-22 10:01:23.088388	200	12
3180	\N	/api/auth/login	POST	2025-04-22 10:01:23.234424	2025-04-22 10:01:23.314247	200	79
3181	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:23.386418	2025-04-22 10:01:23.39102	200	4
3182	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:23.492394	2025-04-22 10:01:23.494323	200	1
3183	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:23.908989	2025-04-22 10:01:23.910449	200	1
3184	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:24.938762	2025-04-22 10:01:24.940398	200	1
3665	\N	/status	GET	2025-04-28 12:25:32.698429	2025-04-28 12:25:32.69843	200	0
3018	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 13:01:56.20997	2025-04-21 13:01:56.212207	200	2
3051	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 14:39:08.860521	2025-04-21 14:39:08.860652	200	0
3053	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-21 15:09:25.879401	2025-04-21 15:09:25.879514	304	0
3054	\N	/api/washingmachines	GET	2025-04-21 15:09:28.830808	2025-04-21 15:09:29.612463	200	781
3055	\N	/api/restaurant	GET	2025-04-21 15:09:37.888898	2025-04-21 15:09:37.88895	200	0
3056	\N	/api/washingmachines	GET	2025-04-21 15:21:15.380728	2025-04-21 15:21:16.077399	200	696
3057	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:39:01.111936	2025-04-21 15:39:01.115411	200	3
3058	\N	/api/restaurant	GET	2025-04-21 15:39:01.520904	2025-04-21 15:39:01.520974	200	0
3059	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:39:01.65465	2025-04-21 15:39:01.656198	200	1
3060	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:56:54.452551	2025-04-21 15:56:54.454216	200	1
3061	\N	/api/restaurant	GET	2025-04-21 15:56:54.835283	2025-04-21 15:56:54.835371	200	0
3062	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:56:54.969657	2025-04-21 15:56:54.971503	200	1
3063	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:59:25.103851	2025-04-21 15:59:25.105823	200	1
3065	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 15:59:25.502728	2025-04-21 15:59:25.504323	200	1
3066	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:01:51.662784	2025-04-21 16:01:51.664445	200	1
3067	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:01:52.413581	2025-04-21 16:01:52.415151	200	1
3068	\N	/api/restaurant	GET	2025-04-21 16:01:52.548318	2025-04-21 16:01:52.548377	200	0
3069	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:02:40.157572	2025-04-21 16:02:40.159668	200	2
3070	\N	/api/restaurant	GET	2025-04-21 16:02:40.519259	2025-04-21 16:02:40.519312	200	0
3071	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:02:40.52255	2025-04-21 16:02:40.52676	200	4
3072	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:04:33.480135	2025-04-21 16:04:33.481854	200	1
3073	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:04:33.884325	2025-04-21 16:04:33.886261	200	1
3074	\N	/api/restaurant	GET	2025-04-21 16:04:34.018867	2025-04-21 16:04:34.018936	200	0
3075	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:07:19.480634	2025-04-21 16:07:19.482483	200	1
3076	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:07:19.904404	2025-04-21 16:07:19.908383	200	3
3077	\N	/api/restaurant	GET	2025-04-21 16:07:20.039089	2025-04-21 16:07:20.03917	200	0
3078	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:16:09.442558	2025-04-21 16:16:09.444082	200	1
3079	\N	/api/restaurant	GET	2025-04-21 16:16:09.856546	2025-04-21 16:16:09.856615	200	0
3080	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:16:09.998848	2025-04-21 16:16:10.000997	200	2
3081	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:18:54.761592	2025-04-21 16:18:54.763178	200	1
3082	\N	/api/restaurant	GET	2025-04-21 16:18:55.148322	2025-04-21 16:18:55.148391	200	0
3083	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:18:55.283279	2025-04-21 16:18:55.304004	200	20
3084	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:19:33.21003	2025-04-21 16:19:33.219921	200	9
3085	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:19:33.60593	2025-04-21 16:19:33.607945	200	2
3086	\N	/api/restaurant	GET	2025-04-21 16:19:33.741952	2025-04-21 16:19:33.742033	200	0
3087	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:21:47.623643	2025-04-21 16:21:47.626357	200	2
3088	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:21:48.184725	2025-04-21 16:21:48.206456	200	21
3089	\N	/api/restaurant	GET	2025-04-21 16:21:48.319987	2025-04-21 16:21:48.320045	200	0
3090	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:22:53.713003	2025-04-21 16:22:53.714629	200	1
3091	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:22:54.205202	2025-04-21 16:22:54.219764	200	14
3092	\N	/api/restaurant	GET	2025-04-21 16:22:54.339825	2025-04-21 16:22:54.339888	200	0
3093	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:23:57.705783	2025-04-21 16:23:57.707235	200	1
3094	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:23:58.097279	2025-04-21 16:23:58.09969	200	2
3095	\N	/api/restaurant	GET	2025-04-21 16:23:58.23091	2025-04-21 16:23:58.230971	200	0
3096	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:30:23.114176	2025-04-21 16:30:23.117354	200	3
3097	\N	/api/restaurant	GET	2025-04-21 16:30:23.581734	2025-04-21 16:30:23.581818	200	0
3098	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:30:23.716367	2025-04-21 16:30:23.717858	200	1
3099	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:38:55.729226	2025-04-21 16:38:55.731327	200	2
3100	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:38:56.140971	2025-04-21 16:38:56.154431	200	13
3101	\N	/api/restaurant	GET	2025-04-21 16:38:56.275265	2025-04-21 16:38:56.275323	200	0
3102	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:40:38.11434	2025-04-21 16:40:38.129903	200	15
3103	\N	/api/restaurant	GET	2025-04-21 16:40:38.417025	2025-04-21 16:40:38.417087	200	0
3104	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:40:38.556554	2025-04-21 16:40:38.557984	200	1
3105	\N	/api/restaurant	GET	2025-04-21 16:40:45.942158	2025-04-21 16:40:45.942213	200	0
3106	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:42:07.831536	2025-04-21 16:42:07.890519	200	58
3107	\N	/api/restaurant	GET	2025-04-21 16:42:08.302598	2025-04-21 16:42:08.302665	200	0
3108	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:42:08.448512	2025-04-21 16:42:08.450226	200	1
3109	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:45:00.45907	2025-04-21 16:45:00.460921	200	1
3111	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:45:00.94612	2025-04-21 16:45:00.949168	200	3
3112	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:47:46.046776	2025-04-21 16:47:46.050372	200	3
3113	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:47:46.476364	2025-04-21 16:47:46.477799	200	1
3114	\N	/api/restaurant	GET	2025-04-21 16:47:46.610545	2025-04-21 16:47:46.610614	200	0
3115	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:52:15.392824	2025-04-21 16:52:15.39486	200	2
3116	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:52:15.92127	2025-04-21 16:52:15.922976	200	1
3117	\N	/api/restaurant	GET	2025-04-21 16:52:16.054767	2025-04-21 16:52:16.05484	200	0
3118	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 16:56:22.965633	2025-04-21 16:56:22.96741	200	1
3119	\N	/api/restaurant	GET	2025-04-21 16:56:23.351666	2025-04-21 16:56:23.351728	200	0
3128	\N	/api/restaurant	GET	2025-04-21 16:59:42.801991	2025-04-21 16:59:42.802037	200	0
3131	\N	/api/restaurant	GET	2025-04-21 17:00:19.845084	2025-04-21 17:00:19.845178	200	0
3135	\N	/api/restaurant	GET	2025-04-21 17:03:07.717977	2025-04-21 17:03:07.718084	200	0
3140	\N	/api/restaurant	GET	2025-04-21 17:06:10.878945	2025-04-21 17:06:10.878982	200	0
3141	\N	/api/washingmachines	GET	2025-04-21 17:10:46.571814	2025-04-21 17:10:47.375457	200	803
3142	\N	/api/washingmachines	GET	2025-04-21 17:15:46.425138	2025-04-21 17:15:47.196097	200	770
3143	\N	/api/washingmachines	GET	2025-04-21 17:20:46.764044	2025-04-21 17:20:47.465871	200	701
3144	\N	/api/washingmachines	GET	2025-04-21 17:26:41.586412	2025-04-21 17:26:42.266075	200	679
3145	\N	/api/auth/login	POST	2025-04-21 17:28:42.249896	2025-04-21 17:28:42.252619	401	2
3146	\N	/api/auth/login	POST	2025-04-21 17:28:45.525421	2025-04-21 17:28:45.593569	200	68
3147	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:28:45.826092	2025-04-21 17:28:45.827643	200	1
3148	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:28:46.040859	2025-04-21 17:28:46.049585	200	8
3149	\N	/api/restaurant	GET	2025-04-21 17:28:46.408353	2025-04-21 17:28:46.408429	200	0
3150	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-21 17:28:46.545134	2025-04-21 17:28:46.54663	200	1
3151	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-04-21 17:28:47.730926	2025-04-21 17:28:47.764279	200	33
3152	\N	/api/washingmachines	GET	2025-04-21 18:01:38.650161	2025-04-21 18:01:39.417096	200	766
3153	\N	/api/washingmachines	GET	2025-04-21 18:24:10.662215	2025-04-21 18:24:11.435827	200	773
3154	\N	/api/washingmachines	GET	2025-04-21 18:29:19.527007	2025-04-21 18:29:20.334468	200	807
3853	\N	/status	GET	2025-04-28 13:33:21.358225	2025-04-28 13:33:21.358227	200	0
3156	\N	/api/washingmachines	GET	2025-04-21 23:50:12.391504	2025-04-21 23:50:13.210319	200	818
3157	\N	/	GET	2025-04-22 02:26:13.071167	2025-04-22 02:26:13.071172	500	0
3158	\N	/favicon.ico	GET	2025-04-22 02:26:18.09922	2025-04-22 02:26:18.099223	500	0
3160	\N	/api/restaurant	GET	2025-04-22 06:08:24.698795	2025-04-22 06:08:25.629626	200	930
3174	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 09:50:43.959399	2025-04-22 09:50:43.959505	200	0
3188	\N	/api/restaurant	GET	2025-04-22 10:01:31.251101	2025-04-22 10:01:31.251186	200	0
3189	\N	/api/restaurant	GET	2025-04-22 10:01:39.385464	2025-04-22 10:01:39.385517	200	0
3190	\N	/api/restaurant	GET	2025-04-22 10:01:52.16421	2025-04-22 10:01:52.164259	200	0
3191	\N	/api/restaurant	GET	2025-04-22 10:02:07.479202	2025-04-22 10:02:07.479255	200	0
3192	\N	/api/restaurant	GET	2025-04-22 10:02:12.989872	2025-04-22 10:02:12.989924	200	0
3193	\N	/api/restaurant	GET	2025-04-22 10:02:25.022522	2025-04-22 10:02:25.022604	200	0
3194	\N	/api/traq/	GET	2025-04-22 10:02:33.410584	2025-04-22 10:02:33.445937	200	35
3195	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:02:42.247447	2025-04-22 10:02:42.249732	200	2
4669	\N	/.env.www	GET	2025-05-08 12:16:07.337483	2025-05-08 12:16:07.337485	500	0
4670	\N	/.docker/.env	GET	2025-05-08 12:16:07.494605	2025-05-08 12:16:07.494607	500	0
4671	\N	/.env.dev	GET	2025-05-08 12:16:07.650578	2025-05-08 12:16:07.65058	500	0
4672	\N	/.env.example	GET	2025-05-08 12:16:07.807175	2025-05-08 12:16:07.807177	500	0
4673	\N	/config.env	GET	2025-05-08 12:16:07.963819	2025-05-08 12:16:07.963822	500	0
4674	\N	/.environment	GET	2025-05-08 12:16:08.120066	2025-05-08 12:16:08.120068	500	0
4675	\N	/.env.production	GET	2025-05-08 12:16:08.281669	2025-05-08 12:16:08.281672	500	0
4676	\N	/.env.production.local	GET	2025-05-08 12:16:08.444007	2025-05-08 12:16:08.444011	500	0
4677	\N	/.env.prod	GET	2025-05-08 12:16:08.601241	2025-05-08 12:16:08.601244	500	0
4678	\N	/.env.test	GET	2025-05-08 12:16:08.757763	2025-05-08 12:16:08.757766	500	0
4679	\N	/.env.php	GET	2025-05-08 12:16:08.914481	2025-05-08 12:16:08.914485	500	0
4680	\N	/config/settings.ini	GET	2025-05-08 12:16:09.070487	2025-05-08 12:16:09.07049	500	0
4681	\N	/info	GET	2025-05-08 12:16:09.227342	2025-05-08 12:16:09.227345	500	0
4682	\N	/info.php	GET	2025-05-08 12:16:09.39028	2025-05-08 12:16:09.390282	500	0
4683	\N	/linusadmin-phpinfo.php	GET	2025-05-08 12:16:09.546306	2025-05-08 12:16:09.546309	500	0
4684	\N	/portal/.env	GET	2025-05-08 12:16:09.702875	2025-05-08 12:16:09.702877	500	0
4685	\N	/env/.env	GET	2025-05-08 12:16:09.85877	2025-05-08 12:16:09.858772	500	0
4686	\N	/dev/.env	GET	2025-05-08 12:16:10.020276	2025-05-08 12:16:10.020278	500	0
4687	\N	/new/.env	GET	2025-05-08 12:16:10.176326	2025-05-08 12:16:10.176328	500	0
4688	\N	/new/.env.local	GET	2025-05-08 12:16:10.33224	2025-05-08 12:16:10.332243	500	0
4689	\N	/new/.env.production	GET	2025-05-08 12:16:10.488034	2025-05-08 12:16:10.488035	500	0
4690	\N	/new/.env.staging	GET	2025-05-08 12:16:10.643733	2025-05-08 12:16:10.643735	500	0
4691	\N	/_phpinfo.php	GET	2025-05-08 12:16:10.807398	2025-05-08 12:16:10.8074	500	0
4692	\N	/_profiler/phpinfo/info.php	GET	2025-05-08 12:16:10.963079	2025-05-08 12:16:10.963082	500	0
4693	\N	/_profiler/phpinfo/phpinfo.php	GET	2025-05-08 12:16:11.119046	2025-05-08 12:16:11.119048	500	0
4694	\N	/wp-config	GET	2025-05-08 12:16:11.278175	2025-05-08 12:16:11.278178	500	0
4695	\N	/aws-secret.yaml	GET	2025-05-08 12:16:11.439383	2025-05-08 12:16:11.439385	500	0
4696	\N	/awstats/.env	GET	2025-05-08 12:16:11.655651	2025-05-08 12:16:11.655653	500	0
4697	\N	/conf/.env	GET	2025-05-08 12:16:11.811575	2025-05-08 12:16:11.811578	500	0
4698	\N	/cron/.env	GET	2025-05-08 12:16:11.967383	2025-05-08 12:16:11.967386	500	0
4699	\N	/www/.env	GET	2025-05-08 12:16:12.123258	2025-05-08 12:16:12.12326	500	0
4700	\N	/docker/.env	GET	2025-05-08 12:16:12.27978	2025-05-08 12:16:12.279782	500	0
4701	\N	/docker/app/.env	GET	2025-05-08 12:16:12.50023	2025-05-08 12:16:12.500234	500	0
4702	\N	/env.backup	GET	2025-05-08 12:16:12.656148	2025-05-08 12:16:12.656151	500	0
4703	\N	/xampp/phpinfo.php	GET	2025-05-08 12:16:12.81258	2025-05-08 12:16:12.812582	500	0
4704	\N	/lara/info.php	GET	2025-05-08 12:16:12.968351	2025-05-08 12:16:12.968353	500	0
4705	\N	/lara/phpinfo.php	GET	2025-05-08 12:16:13.124569	2025-05-08 12:16:13.124571	500	0
4706	\N	/laravel/info.php	GET	2025-05-08 12:16:13.330369	2025-05-08 12:16:13.330372	500	0
4707	\N	/.vscode/.env	GET	2025-05-08 12:16:13.486833	2025-05-08 12:16:13.486835	500	0
4708	\N	/js/.env	GET	2025-05-08 12:16:13.645919	2025-05-08 12:16:13.645921	500	0
4709	\N	/laravel/core/.env	GET	2025-05-08 12:16:13.806382	2025-05-08 12:16:13.806385	500	0
4710	\N	/mail/.env	GET	2025-05-08 12:16:13.962476	2025-05-08 12:16:13.962479	500	0
4711	\N	/mailer/.env	GET	2025-05-08 12:16:14.118848	2025-05-08 12:16:14.11885	500	0
4712	\N	/nginx/.env	GET	2025-05-08 12:16:14.279984	2025-05-08 12:16:14.279987	500	0
4713	\N	/public/.env	GET	2025-05-08 12:16:14.436859	2025-05-08 12:16:14.436862	500	0
4714	\N	/site/.env	GET	2025-05-08 12:16:14.592678	2025-05-08 12:16:14.59268	500	0
4715	\N	/xampp/.env	GET	2025-05-08 12:16:14.755951	2025-05-08 12:16:14.755954	500	0
4716	\N	/.docker/laravel/app/.env	GET	2025-05-08 12:16:14.912441	2025-05-08 12:16:14.912443	500	0
4717	\N	/laravel/.env.local	GET	2025-05-08 12:16:15.068586	2025-05-08 12:16:15.068588	500	0
4718	\N	/laravel/.env.production	GET	2025-05-08 12:16:15.236553	2025-05-08 12:16:15.236555	500	0
4719	\N	/laravel/.env.staging	GET	2025-05-08 12:16:15.392432	2025-05-08 12:16:15.392434	500	0
4720	\N	/laravel/core/.env.local	GET	2025-05-08 12:16:15.554448	2025-05-08 12:16:15.55445	500	0
4721	\N	/laravel/core/.env.production	GET	2025-05-08 12:16:15.718395	2025-05-08 12:16:15.718397	500	0
4722	\N	/laravel/core/.env.staging	GET	2025-05-08 12:16:15.876039	2025-05-08 12:16:15.876042	500	0
4723	\N	/main/.env	GET	2025-05-08 12:16:16.036079	2025-05-08 12:16:16.036081	500	0
4724	\N	/node_modules/.env	GET	2025-05-08 12:16:16.197101	2025-05-08 12:16:16.197103	500	0
4725	\N	/php.ini	GET	2025-05-08 12:16:16.35389	2025-05-08 12:16:16.353892	500	0
4726	\N	/.aws/credentials	GET	2025-05-08 12:16:16.525438	2025-05-08 12:16:16.525441	500	0
4727	\N	/config.json	GET	2025-05-08 12:16:16.684579	2025-05-08 12:16:16.684581	500	0
4728	\N	/config.php	GET	2025-05-08 12:16:16.840377	2025-05-08 12:16:16.840379	500	0
4729	\N	/server/config.json	GET	2025-05-08 12:16:17.000691	2025-05-08 12:16:17.000694	500	0
4730	\N	/.aws/config	GET	2025-05-08 12:16:17.158891	2025-05-08 12:16:17.158895	500	0
4731	\N	/config.ini	GET	2025-05-08 12:16:17.330708	2025-05-08 12:16:17.33071	500	0
4732	\N	/settings.php	GET	2025-05-08 12:16:17.532189	2025-05-08 12:16:17.532191	500	0
4733	\N	/docker_run.sh	GET	2025-05-08 12:16:17.688088	2025-05-08 12:16:17.688091	500	0
4734	\N	/secrets.json	GET	2025-05-08 12:16:17.85454	2025-05-08 12:16:17.854543	500	0
4735	\N	/awsconfig.js	GET	2025-05-08 12:16:18.010696	2025-05-08 12:16:18.010699	500	0
4736	\N	/startup.sh	GET	2025-05-08 12:16:18.173276	2025-05-08 12:16:18.173279	500	0
4737	\N	/config.inc.php	GET	2025-05-08 12:16:18.348367	2025-05-08 12:16:18.34837	500	0
4738	\N	/.travis.yml	GET	2025-05-08 12:16:18.509241	2025-05-08 12:16:18.509244	500	0
4739	\N	/parameters.yml	GET	2025-05-08 12:16:18.665139	2025-05-08 12:16:18.665142	500	0
4740	\N	/sendgrid.json	GET	2025-05-08 12:16:18.820947	2025-05-08 12:16:18.820949	500	0
4741	\N	/config/appsettings.json	GET	2025-05-08 12:16:18.976733	2025-05-08 12:16:18.976736	500	0
4742	\N	/bootstrap/cache/config.php	GET	2025-05-08 12:16:19.135724	2025-05-08 12:16:19.135727	500	0
4743	\N	/storage/app/private/.env	GET	2025-05-08 12:16:19.292378	2025-05-08 12:16:19.292381	500	0
4744	\N	/storage/logs/laravel.log	GET	2025-05-08 12:16:19.454872	2025-05-08 12:16:19.454875	500	0
4745	\N	/config/services.php	GET	2025-05-08 12:16:19.610731	2025-05-08 12:16:19.610733	500	0
4746	\N	/composer.lock	GET	2025-05-08 12:16:19.766789	2025-05-08 12:16:19.766791	500	0
4747	\N	/public/_debugbar/	GET	2025-05-08 12:16:19.9227	2025-05-08 12:16:19.922702	500	0
4748	\N	/info.php.bak	GET	2025-05-08 12:16:20.078568	2025-05-08 12:16:20.078571	500	0
4749	\N	/php-info.php	GET	2025-05-08 12:16:20.246171	2025-05-08 12:16:20.246174	500	0
4750	\N	/newinfo.php	GET	2025-05-08 12:16:20.413149	2025-05-08 12:16:20.413151	500	0
4751	\N	/siteinfo.php	GET	2025-05-08 12:16:20.573625	2025-05-08 12:16:20.573627	500	0
4752	\N	/info.php.1	GET	2025-05-08 12:16:20.735295	2025-05-08 12:16:20.735297	500	0
4753	\N	/admin/phpinfo.php	GET	2025-05-08 12:16:20.891063	2025-05-08 12:16:20.891065	500	0
4754	\N	/pageinfo.php	GET	2025-05-08 12:16:21.049341	2025-05-08 12:16:21.049343	500	0
4755	\N	/info.php.back	GET	2025-05-08 12:16:21.21222	2025-05-08 12:16:21.212223	500	0
4756	\N	/phpbb/phpinfo.php	GET	2025-05-08 12:16:21.387822	2025-05-08 12:16:21.387824	500	0
4924	\N	/api/auth/login	POST	2025-05-09 14:16:02.845112	2025-05-09 14:16:03.263026	200	417
3185	test@imt-atlantique.net	/api/newf/me	GET	2025-04-22 10:01:25.944938	2025-04-22 10:01:25.946525	200	1
3186	\N	/api/restaurant	GET	2025-04-22 10:01:24.908416	2025-04-22 10:01:26.480881	200	1572
3187	\N	/api/restaurant	GET	2025-04-22 10:01:25.732724	2025-04-22 10:01:26.480966	200	748
3196	\N	/api/restaurant	GET	2025-04-22 10:03:03.765438	2025-04-22 10:03:03.765504	200	0
3197	\N	/api/restaurant	GET	2025-04-22 10:03:14.468766	2025-04-22 10:03:14.46883	200	0
3198	\N	/api/restaurant	GET	2025-04-22 10:03:35.609428	2025-04-22 10:03:35.609495	200	0
3199	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 10:19:42.192241	2025-04-22 10:19:42.192334	200	0
3200	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-22 10:19:42.369521	2025-04-22 10:19:42.369662	200	0
3201	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:31:41.590194	2025-04-22 11:31:41.592586	200	2
3202	\N	/api/auth/login	POST	2025-04-22 11:48:19.678801	2025-04-22 11:48:19.684047	401	5
3203	\N	/api/auth/login	POST	2025-04-22 11:48:24.534475	2025-04-22 11:48:24.602118	200	67
3204	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:25.015031	2025-04-22 11:48:25.093334	200	78
3205	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:25.525571	2025-04-22 11:48:25.527459	200	1
3206	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:26.491187	2025-04-22 11:48:26.494971	200	3
3207	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-04-22 11:48:27.461686	2025-04-22 11:48:27.464835	200	3
3208	\N	/api/restaurant	GET	2025-04-22 11:48:26.172809	2025-04-22 11:48:27.598031	200	1425
3209	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:31.723669	2025-04-22 11:48:31.725367	200	1
3210	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:33.255928	2025-04-22 11:48:33.279359	200	23
3211	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-04-22 11:48:33.396233	2025-04-22 11:48:33.396321	200	0
3212	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-04-22 11:48:33.937112	2025-04-22 11:48:33.937172	200	0
3213	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-04-22 11:48:36.058432	2025-04-22 11:48:36.061371	200	2
3214	\N	/api/restaurant	GET	2025-04-22 11:48:39.753061	2025-04-22 11:48:39.753113	200	0
3215	\N	/api/restaurant	GET	2025-04-22 11:48:41.14705	2025-04-22 11:48:41.1471	200	0
3216	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:50.355712	2025-04-22 11:48:50.365262	200	9
3217	\N	/api/restaurant	GET	2025-04-22 11:48:50.824975	2025-04-22 11:48:50.825026	200	0
3218	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:48:50.987068	2025-04-22 11:48:50.989427	200	2
3219	\N	/api/restaurant	GET	2025-04-22 11:50:25.049231	2025-04-22 11:50:25.049302	200	0
3220	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:50:37.975563	2025-04-22 11:50:37.977514	200	1
3221	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-22 11:50:50.087331	2025-04-22 11:50:50.092372	200	5
3222	\N	/api/restaurant	GET	2025-04-22 18:33:46.23248	2025-04-22 18:33:46.232565	200	0
3223	\N	/favicon.ico	GET	2025-04-22 18:33:46.971012	2025-04-22 18:33:46.971015	500	0
3224	\N	/api/restaurant	GET	2025-04-22 18:34:03.756149	2025-04-22 18:34:03.756227	200	0
3225	\N	/favicon.ico	GET	2025-04-22 18:34:03.937804	2025-04-22 18:34:03.937809	500	0
3226	\N	/api/restaurant	GET	2025-04-22 18:34:04.333592	2025-04-22 18:34:04.333646	200	0
3227	\N	/favicon.ico	GET	2025-04-22 18:34:04.51541	2025-04-22 18:34:04.515416	500	0
3228	\N	/api/restaurant	GET	2025-04-22 18:34:04.862425	2025-04-22 18:34:04.862498	200	0
3229	\N	/api/restaurant	GET	2025-04-22 18:34:05.317467	2025-04-22 18:34:05.31752	200	0
3230	\N	/favicon.ico	GET	2025-04-22 18:34:05.634666	2025-04-22 18:34:05.63467	500	0
3231	\N	/.git/config	GET	2025-04-23 03:54:49.608648	2025-04-23 03:54:49.608651	500	0
3232	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-23 05:06:49.503917	2025-04-23 05:06:49.504034	200	0
3233	\N	/robots.txt	GET	2025-04-23 05:11:15.434049	2025-04-23 05:11:15.434052	500	0
3234	\N	/	GET	2025-04-23 05:11:15.518038	2025-04-23 05:11:15.518041	500	0
3235	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-23 08:15:24.129533	2025-04-23 08:15:24.129624	200	0
3236	\N	/api/restaurant	GET	2025-04-23 08:15:24.141083	2025-04-23 08:15:25.10523	200	964
3237	\N	/	GET	2025-04-23 08:15:28.880302	2025-04-23 08:15:28.880305	500	0
3238	\N	/favicon.ico	GET	2025-04-23 08:15:29.130126	2025-04-23 08:15:29.130131	500	0
3239	\N	/status	GET	2025-04-23 08:16:18.332444	2025-04-23 08:16:18.36151	200	29
3240	\N	/api/statistics/global	GET	2025-04-23 08:16:18.538596	2025-04-23 08:16:18.540594	200	1
3241	\N	/api/statistics/endpoints	GET	2025-04-23 08:16:18.79318	2025-04-23 08:16:18.798671	200	5
3242	\N	/status	GET	2025-04-23 08:16:22.984455	2025-04-23 08:16:22.984457	200	0
3243	\N	/status	GET	2025-04-23 08:16:23.181606	2025-04-23 08:16:23.181607	200	0
3244	\N	/api/statistics/global	GET	2025-04-23 08:16:23.316876	2025-04-23 08:16:23.319658	200	2
3245	\N	/api/statistics/endpoints	GET	2025-04-23 08:16:23.524629	2025-04-23 08:16:23.527879	200	3
3246	\N	/status	GET	2025-04-23 08:16:27.984446	2025-04-23 08:16:27.98445	200	0
3247	\N	/status	GET	2025-04-23 08:16:28.176095	2025-04-23 08:16:28.176096	200	0
3248	\N	/api/statistics/global	GET	2025-04-23 08:16:28.176082	2025-04-23 08:16:28.178459	200	2
3249	\N	/api/statistics/endpoints	GET	2025-04-23 08:16:28.361488	2025-04-23 08:16:28.365403	200	3
3250	\N	/status	GET	2025-04-23 08:16:33.000662	2025-04-23 08:16:33.000663	200	0
3251	\N	/status	GET	2025-04-23 08:16:33.185714	2025-04-23 08:16:33.185716	200	0
3252	\N	/api/statistics/global	GET	2025-04-23 08:16:33.187597	2025-04-23 08:16:33.190043	200	2
3253	\N	/api/statistics/endpoints	GET	2025-04-23 08:16:33.365634	2025-04-23 08:16:33.368879	200	3
3254	\N	/api/restaurant	GET	2025-04-23 08:16:36.614726	2025-04-23 08:16:36.614788	200	0
3255	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-23 09:44:21.539004	2025-04-23 09:44:21.539102	200	0
3256	\N	/api/auth/login	POST	2025-04-23 10:11:04.68003	2025-04-23 10:11:04.749519	200	69
3257	test@imt-atlantique.net	/api/newf/me	GET	2025-04-23 10:11:04.887773	2025-04-23 10:11:04.889669	200	1
3258	test@imt-atlantique.net	/api/newf/me	GET	2025-04-23 10:11:05.072879	2025-04-23 10:11:05.078539	200	5
3259	test@imt-atlantique.net	/api/newf/me	GET	2025-04-23 10:11:06.15159	2025-04-23 10:11:06.155999	200	4
3260	\N	/api/restaurant	GET	2025-04-23 10:11:06.479478	2025-04-23 10:11:07.129937	200	650
3261	\N	/api/traq/	GET	2025-04-23 10:11:12.93483	2025-04-23 10:11:12.936391	200	1
3262	test@imt-atlantique.net	/api/newf/me	GET	2025-04-23 10:11:28.343481	2025-04-23 10:11:28.363473	200	19
3263	\N	/api/restaurant	GET	2025-04-23 10:11:44.267166	2025-04-23 10:11:44.267213	200	0
3264	\N	/api/traq/	GET	2025-04-23 10:12:04.551356	2025-04-23 10:12:04.552542	200	1
3265	\N	/api/restaurant	GET	2025-04-23 10:12:27.672675	2025-04-23 10:12:27.672739	200	0
3266	test@imt-atlantique.net	/api/newf/me	GET	2025-04-23 10:12:52.984559	2025-04-23 10:12:52.986241	200	1
3267	\N	/api/restaurant	GET	2025-04-23 10:13:50.590323	2025-04-23 10:13:50.590383	200	0
3268	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-23 14:39:31.840236	2025-04-23 14:39:31.842003	200	1
3269	\N	/	GET	2025-04-24 01:09:57.713364	2025-04-24 01:09:57.83548	500	122
3270	\N	/favicon.ico	GET	2025-04-24 01:10:00.447085	2025-04-24 01:10:00.44709	500	0
3271	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-24 03:03:20.799095	2025-04-24 03:03:20.947203	200	148
3272	\N	/favicon.ico	GET	2025-04-24 05:56:30.286462	2025-04-24 05:56:30.286465	500	0
3273	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-24 07:06:26.420714	2025-04-24 07:06:26.420809	200	0
3274	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-24 07:06:26.472308	2025-04-24 07:06:26.472377	200	0
3275	\N	/	GET	2025-04-24 08:55:50.517474	2025-04-24 08:55:50.517479	500	0
3276	\N	/favicon.ico	GET	2025-04-24 08:56:07.150209	2025-04-24 08:56:07.150213	500	0
3277	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-24 09:02:33.758193	2025-04-24 09:02:33.758316	200	0
3278	\N	/wp-includes/js/jquery/jquery.js	GET	2025-04-24 14:10:38.920779	2025-04-24 14:10:38.920808	500	0
3279	\N	/	GET	2025-04-24 15:48:08.762914	2025-04-24 15:48:08.762918	500	0
3280	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-24 16:22:11.274744	2025-04-24 16:22:11.274862	200	0
3281	\N	/.env	GET	2025-04-24 18:31:41.788037	2025-04-24 18:31:41.78804	500	0
3282	\N	/config/.env	GET	2025-04-24 18:31:41.94473	2025-04-24 18:31:41.944733	500	0
3283	\N	/.env.production	GET	2025-04-24 18:31:42.099167	2025-04-24 18:31:42.099171	500	0
3284	\N	/api/.env	GET	2025-04-24 18:31:42.259912	2025-04-24 18:31:42.259923	500	0
3285	\N	/settings/.env	GET	2025-04-24 18:31:42.43775	2025-04-24 18:31:42.437755	500	0
3286	\N	/config/app.php	GET	2025-04-24 18:31:42.592181	2025-04-24 18:31:42.592186	500	0
3287	\N	/application.yml	GET	2025-04-24 18:31:42.746803	2025-04-24 18:31:42.746817	500	0
3288	\N	/config/database.yml	GET	2025-04-24 18:31:42.919353	2025-04-24 18:31:42.919356	500	0
3289	\N	/secrets.json	GET	2025-04-24 18:31:43.073678	2025-04-24 18:31:43.073682	500	0
3290	\N	/src/config.js	GET	2025-04-24 18:31:43.228157	2025-04-24 18:31:43.228165	500	0
3291	\N	/db.ini	GET	2025-04-24 18:31:43.387843	2025-04-24 18:31:43.387847	500	0
3292	\N	/api/credentials	GET	2025-04-24 18:31:43.549113	2025-04-24 18:31:43.549121	500	0
3293	\N	/.aws/credentials	GET	2025-04-24 18:31:43.703754	2025-04-24 18:31:43.703759	500	0
3294	\N	/secure-config.json	GET	2025-04-24 18:31:43.868332	2025-04-24 18:31:43.868336	500	0
3295	\N	/local_settings.py	GET	2025-04-24 18:31:44.02369	2025-04-24 18:31:44.023694	500	0
3296	\N	/config/default.json	GET	2025-04-24 18:31:44.18155	2025-04-24 18:31:44.181555	500	0
3297	\N	/config/production.json	GET	2025-04-24 18:31:44.362641	2025-04-24 18:31:44.362646	500	0
3298	\N	/bootstrap/cache/config.php	GET	2025-04-24 18:31:44.536213	2025-04-24 18:31:44.536218	500	0
3299	\N	/config/secrets.yml	GET	2025-04-24 18:31:44.69042	2025-04-24 18:31:44.690424	500	0
3300	\N	/settings.yaml	GET	2025-04-24 18:31:44.844886	2025-04-24 18:31:44.844891	500	0
3301	\N	/auth.json	GET	2025-04-24 18:31:45.005364	2025-04-24 18:31:45.00537	500	0
3302	\N	/helm/values.yaml	GET	2025-04-24 18:31:45.218888	2025-04-24 18:31:45.218892	500	0
3303	\N	/docker/.env	GET	2025-04-24 18:31:45.376857	2025-04-24 18:31:45.376862	500	0
3304	\N	/wp-config.php	GET	2025-04-24 18:31:45.531594	2025-04-24 18:31:45.531597	500	0
3305	\N	/config.json	GET	2025-04-24 18:31:45.686251	2025-04-24 18:31:45.686258	500	0
3306	\N	/database.json	GET	2025-04-24 18:31:45.840773	2025-04-24 18:31:45.840778	500	0
3307	\N	/config/secrets.json	GET	2025-04-24 18:31:45.99535	2025-04-24 18:31:45.995354	500	0
3308	\N	/env.backup	GET	2025-04-24 18:31:46.155978	2025-04-24 18:31:46.155983	500	0
3309	\N	/settings.bak	GET	2025-04-24 18:31:46.313404	2025-04-24 18:31:46.313409	500	0
3310	\N	/backup.env	GET	2025-04-24 18:31:46.478711	2025-04-24 18:31:46.478716	500	0
3311	\N	/old/.env	GET	2025-04-24 18:31:46.635378	2025-04-24 18:31:46.635383	500	0
3312	\N	/phpinfo.php	GET	2025-04-24 18:31:46.789691	2025-04-24 18:31:46.789697	500	0
3313	\N	/info.php	GET	2025-04-24 18:31:46.94413	2025-04-24 18:31:46.944134	500	0
3314	\N	/test.php	GET	2025-04-24 18:31:47.098635	2025-04-24 18:31:47.09864	500	0
3315	\N	/laravel/.env	GET	2025-04-24 18:31:47.286619	2025-04-24 18:31:47.286624	500	0
3316	\N	/app/config/.env	GET	2025-04-24 18:31:47.453577	2025-04-24 18:31:47.453589	500	0
3317	\N	/.git/config	GET	2025-04-24 18:31:47.607555	2025-04-24 18:31:47.607559	500	0
3318	\N	/.svn/entries	GET	2025-04-24 18:31:47.772187	2025-04-24 18:31:47.772192	500	0
3319	\N	/.git/HEAD	GET	2025-04-24 18:31:47.927561	2025-04-24 18:31:47.927567	500	0
3320	\N	/.git/index	GET	2025-04-24 18:31:48.081714	2025-04-24 18:31:48.081718	500	0
3321	\N	/.git/logs/HEAD	GET	2025-04-24 18:31:48.256315	2025-04-24 18:31:48.25632	500	0
3322	\N	/.gitignore	GET	2025-04-24 18:31:48.425175	2025-04-24 18:31:48.42518	500	0
3323	\N	/administrator/index.php	GET	2025-04-24 18:31:48.579305	2025-04-24 18:31:48.579309	500	0
3324	\N	/wp-admin/install.php	GET	2025-04-24 18:31:48.733837	2025-04-24 18:31:48.733843	500	0
3325	\N	/joomla/configuration.php-dist	GET	2025-04-24 18:31:48.889215	2025-04-24 18:31:48.88922	500	0
3326	\N	/sites/default/settings.php	GET	2025-04-24 18:31:49.054353	2025-04-24 18:31:49.054358	500	0
3327	\N	/bitrix/php_interface/dbconn.php	GET	2025-04-24 18:31:49.216554	2025-04-24 18:31:49.216558	500	0
3328	\N	/typo3conf/localconf.php	GET	2025-04-24 18:31:49.387855	2025-04-24 18:31:49.387859	500	0
3329	\N	/config.inc.php	GET	2025-04-24 18:31:49.542054	2025-04-24 18:31:49.542058	500	0
3330	\N	/config.old.php	GET	2025-04-24 18:31:49.703032	2025-04-24 18:31:49.703036	500	0
3331	\N	/php.ini	GET	2025-04-24 18:31:49.857296	2025-04-24 18:31:49.8573	500	0
3332	\N	/cgi-bin/phpinfo.php	GET	2025-04-24 18:31:50.011149	2025-04-24 18:31:50.011153	500	0
3333	\N	/debug.php	GET	2025-04-24 18:31:50.170753	2025-04-24 18:31:50.170757	500	0
3334	\N	/server-status	GET	2025-04-24 18:31:50.34459	2025-04-24 18:31:50.344594	500	0
3335	\N	/phpinfo1.php	GET	2025-04-24 18:31:50.519159	2025-04-24 18:31:50.519164	500	0
3336	\N	/phpinfo2.php	GET	2025-04-24 18:31:50.676947	2025-04-24 18:31:50.676951	500	0
3337	\N	/env.txt	GET	2025-04-24 18:31:50.83416	2025-04-24 18:31:50.834164	500	0
3338	\N	/prod.env	GET	2025-04-24 18:31:50.994529	2025-04-24 18:31:50.994534	500	0
3339	\N	/stage.env	GET	2025-04-24 18:31:51.157635	2025-04-24 18:31:51.15764	500	0
3340	\N	/development.env	GET	2025-04-24 18:31:51.314986	2025-04-24 18:31:51.314992	500	0
3341	\N	/credentials.env	GET	2025-04-24 18:31:51.483818	2025-04-24 18:31:51.483822	500	0
3342	\N	/public/.env	GET	2025-04-24 18:31:51.637738	2025-04-24 18:31:51.637742	500	0
3343	\N	/api/config.json	GET	2025-04-24 18:31:51.79271	2025-04-24 18:31:51.792719	500	0
3344	\N	/composer.json	GET	2025-04-24 18:31:51.947071	2025-04-24 18:31:51.947076	500	0
3345	\N	/api/v1/.env	GET	2025-04-24 18:31:52.101622	2025-04-24 18:31:52.101634	500	0
3346	\N	/staging.env	GET	2025-04-24 18:31:52.279251	2025-04-24 18:31:52.279256	500	0
3347	\N	/phpmyadmin/index.php	GET	2025-04-24 18:31:52.454718	2025-04-24 18:31:52.454724	500	0
3348	\N	/backup/config.php	GET	2025-04-24 18:31:52.613063	2025-04-24 18:31:52.613067	500	0
3349	\N	/.env.example	GET	2025-04-24 18:31:52.926773	2025-04-24 18:31:52.926777	500	0
3350	\N	/storage/logs/laravel.log	GET	2025-04-24 18:31:53.081156	2025-04-24 18:31:53.081164	500	0
3351	\N	/storage/framework/sessions/	GET	2025-04-24 18:31:53.240775	2025-04-24 18:31:53.240779	500	0
3352	\N	/storage/framework/cache/	GET	2025-04-24 18:31:53.429156	2025-04-24 18:31:53.429161	500	0
3353	\N	/storage/framework/views/	GET	2025-04-24 18:31:53.583291	2025-04-24 18:31:53.583297	500	0
3354	\N	/nova-api/styles	GET	2025-04-24 18:31:53.748729	2025-04-24 18:31:53.748737	500	0
3355	\N	/.env.local	GET	2025-04-24 18:31:53.905167	2025-04-24 18:31:53.90517	500	0
3356	\N	/.env.dev	GET	2025-04-24 18:31:54.059089	2025-04-24 18:31:54.059093	500	0
3357	\N	/.env.test	GET	2025-04-24 18:31:54.218275	2025-04-24 18:31:54.21828	500	0
3358	\N	/var/logs/dev.log	GET	2025-04-24 18:31:54.378191	2025-04-24 18:31:54.378197	500	0
3359	\N	/var/logs/prod.log	GET	2025-04-24 18:31:54.541908	2025-04-24 18:31:54.541913	500	0
3360	\N	/config/packages/	GET	2025-04-24 18:31:54.699204	2025-04-24 18:31:54.699208	500	0
3361	\N	/web/config.php	GET	2025-04-24 18:31:54.854087	2025-04-24 18:31:54.854092	500	0
3362	\N	/config/routes.yaml	GET	2025-04-24 18:31:55.008096	2025-04-24 18:31:55.0081	500	0
3363	\N	/web.config	GET	2025-04-24 18:31:55.166805	2025-04-24 18:31:55.166809	500	0
3364	\N	/.htaccess	GET	2025-04-24 18:31:55.322045	2025-04-24 18:31:55.322049	500	0
3365	\N	/sites/all/modules/	GET	2025-04-24 18:31:55.502063	2025-04-24 18:31:55.502067	500	0
3366	\N	/sites/all/themes/	GET	2025-04-24 18:31:55.658972	2025-04-24 18:31:55.658976	500	0
3367	\N	/core/install.php	GET	2025-04-24 18:31:55.815798	2025-04-24 18:31:55.815802	500	0
3368	\N	/CHANGELOG.txt	GET	2025-04-24 18:31:55.978488	2025-04-24 18:31:55.978494	500	0
3369	\N	/app/etc/local.xml	GET	2025-04-24 18:31:56.132991	2025-04-24 18:31:56.132998	500	0
3370	\N	/app/etc/env.php	GET	2025-04-24 18:31:56.295625	2025-04-24 18:31:56.295632	500	0
3371	\N	/var/log/system.log	GET	2025-04-24 18:31:56.498941	2025-04-24 18:31:56.498945	500	0
3372	\N	/var/log/exception.log	GET	2025-04-24 18:31:56.655648	2025-04-24 18:31:56.655652	500	0
3373	\N	/.wp-config.php.swp	GET	2025-04-24 18:31:56.816108	2025-04-24 18:31:56.816141	500	0
3374	\N	/wp-config-sample.php	GET	2025-04-24 18:31:56.978803	2025-04-24 18:31:56.978808	500	0
3375	\N	/wp-content/debug.log	GET	2025-04-24 18:31:57.147387	2025-04-24 18:31:57.147391	500	0
3376	\N	/xmlrpc.php	GET	2025-04-24 18:31:57.445838	2025-04-24 18:31:57.445842	500	0
3377	\N	/wp-json/wp/v2/users	GET	2025-04-24 18:31:57.649853	2025-04-24 18:31:57.649857	500	0
3378	\N	/configuration.php~	GET	2025-04-24 18:31:57.805829	2025-04-24 18:31:57.805833	500	0
3379	\N	/logs/error.php	GET	2025-04-24 18:31:58.215375	2025-04-24 18:31:58.21538	500	0
3380	\N	/package.json	GET	2025-04-24 18:31:58.390929	2025-04-24 18:31:58.390933	500	0
3381	\N	/yarn.lock	GET	2025-04-24 18:31:58.556738	2025-04-24 18:31:58.556742	500	0
3382	\N	/.npmrc	GET	2025-04-24 18:31:58.782436	2025-04-24 18:31:58.78244	500	0
3383	\N	/server.js	GET	2025-04-24 18:31:58.937128	2025-04-24 18:31:58.937131	500	0
3384	\N	/app.js	GET	2025-04-24 18:31:59.093543	2025-04-24 18:31:59.093549	500	0
3385	\N	/config.js	GET	2025-04-24 18:31:59.255609	2025-04-24 18:31:59.255614	500	0
3386	\N	/.dockerignore	GET	2025-04-24 18:31:59.425037	2025-04-24 18:31:59.42504	500	0
3387	\N	/Dockerfile	GET	2025-04-24 18:31:59.592146	2025-04-24 18:31:59.592151	500	0
3388	\N	/.git/logs/	GET	2025-04-24 18:31:59.74735	2025-04-24 18:31:59.747354	500	0
3389	\N	/.git/refs/	GET	2025-04-24 18:31:59.903521	2025-04-24 18:31:59.903525	500	0
3390	\N	/.git/objects/	GET	2025-04-24 18:32:00.060426	2025-04-24 18:32:00.060432	500	0
3391	\N	/.git/packed-refs	GET	2025-04-24 18:32:00.220802	2025-04-24 18:32:00.220806	500	0
3392	\N	/.git/branches/	GET	2025-04-24 18:32:00.394256	2025-04-24 18:32:00.394267	500	0
3393	\N	/api-docs	GET	2025-04-24 18:32:00.552835	2025-04-24 18:32:00.552841	500	0
3394	\N	/swagger.json	GET	2025-04-24 18:32:00.70826	2025-04-24 18:32:00.708264	500	0
3395	\N	/swagger-ui.html	GET	2025-04-24 18:32:00.866479	2025-04-24 18:32:00.866483	500	0
3396	\N	/openapi.json	GET	2025-04-24 18:32:01.021426	2025-04-24 18:32:01.021431	500	0
3397	\N	/backup.sql	GET	2025-04-24 18:32:01.191957	2025-04-24 18:32:01.191963	500	0
3398	\N	/db_backup.sql	GET	2025-04-24 18:32:01.374741	2025-04-24 18:32:01.374744	500	0
3399	\N	/.well-known/security.txt	GET	2025-04-24 18:32:01.540332	2025-04-24 18:32:01.540336	500	0
3400	\N	/phpinfo	GET	2025-04-25 00:13:17.084965	2025-04-25 00:13:17.08497	500	0
3401	\N	/phpinfo.php	GET	2025-04-25 00:13:20.292683	2025-04-25 00:13:20.292687	500	0
3402	\N	/test.php	GET	2025-04-25 00:13:23.334416	2025-04-25 00:13:23.334421	500	0
3403	\N	/_profiler/phpinfo	GET	2025-04-25 00:13:26.678534	2025-04-25 00:13:26.678539	500	0
3404	\N	/info.php	GET	2025-04-25 00:13:29.740052	2025-04-25 00:13:29.740057	500	0
3405	\N	/php.php	GET	2025-04-25 00:13:32.803709	2025-04-25 00:13:32.803716	500	0
3406	\N	/php_info.php	GET	2025-04-25 00:13:34.692687	2025-04-25 00:13:34.692691	500	0
3407	\N	/i.php	GET	2025-04-25 00:13:37.884284	2025-04-25 00:13:37.884288	500	0
3408	\N	/pi.php	GET	2025-04-25 00:13:40.90865	2025-04-25 00:13:40.908654	500	0
3409	\N	/config.phpinfo	GET	2025-04-25 00:13:44.25995	2025-04-25 00:13:44.259955	500	0
3410	\N	/admin/phpinfo.php	GET	2025-04-25 00:13:47.457709	2025-04-25 00:13:47.457713	500	0
3411	\N	/.aws/credentials	GET	2025-04-25 00:13:50.668178	2025-04-25 00:13:50.668182	500	0
3413	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-25 07:31:00.764414	2025-04-25 07:31:00.764518	200	0
3417	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-25 15:26:02.466279	2025-04-25 15:26:02.466388	200	0
3418	\N	/	GET	2025-04-25 16:15:13.229566	2025-04-25 16:15:13.229569	500	0
3419	\N	/	GET	2025-04-26 02:50:47.96467	2025-04-26 02:50:47.964673	500	0
3420	\N	/favicon.ico	GET	2025-04-26 02:50:49.255839	2025-04-26 02:50:49.255844	500	0
3541	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-26 06:51:40.31485	2025-04-26 06:51:40.314948	200	0
3542	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-26 06:51:40.37963	2025-04-26 06:51:40.379701	200	0
3543	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-26 10:40:50.243901	2025-04-26 10:40:50.244023	200	0
3544	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-26 23:39:54.871498	2025-04-26 23:39:54.871662	200	0
3545	\N	/	GET	2025-04-27 02:24:08.671928	2025-04-27 02:24:08.67193	500	0
3546	\N	/wp	GET	2025-04-27 02:24:12.258868	2025-04-27 02:24:12.258871	500	0
3547	\N	/wordpress	GET	2025-04-27 02:24:14.742049	2025-04-27 02:24:14.742052	500	0
3548	\N	/transat	GET	2025-04-27 02:24:19.353656	2025-04-27 02:24:19.353661	500	0
3549	\N	/wordpress	GET	2025-04-27 02:24:23.713928	2025-04-27 02:24:23.713932	500	0
3550	\N	/old	GET	2025-04-27 02:24:25.819249	2025-04-27 02:24:25.819251	500	0
3551	\N	/backup	GET	2025-04-27 02:24:29.347549	2025-04-27 02:24:29.347554	500	0
3552	\N	/new	GET	2025-04-27 02:24:33.234213	2025-04-27 02:24:33.234216	500	0
3553	\N	/test	GET	2025-04-27 02:24:35.539407	2025-04-27 02:24:35.539411	500	0
3554	\N	/wp	GET	2025-04-27 02:24:37.075495	2025-04-27 02:24:37.075501	500	0
3555	\N	/temp	GET	2025-04-27 02:24:39.885468	2025-04-27 02:24:39.885473	500	0
3556	\N	/blog	GET	2025-04-27 02:24:40.650767	2025-04-27 02:24:40.650772	500	0
3560	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-27 13:37:51.337404	2025-04-27 13:37:51.337507	200	0
3561	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:07:46.709484	2025-04-27 16:07:46.870195	200	160
3562	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:07:47.387549	2025-04-27 16:07:47.389052	200	1
3563	\N	/api/restaurant	GET	2025-04-27 16:07:47.554131	2025-04-27 16:07:48.595845	200	1041
3564	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:52:37.066855	2025-04-27 16:52:37.082594	200	15
3565	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:52:37.427749	2025-04-27 16:52:37.429906	200	2
3566	\N	/api/restaurant	GET	2025-04-27 16:52:37.565909	2025-04-27 16:52:37.56598	200	0
4757	\N	/administrator/components/com_admin/models/sysinfo.php	GET	2025-05-08 12:16:21.545297	2025-05-08 12:16:21.545301	500	0
4758	\N	/release_info.php	GET	2025-05-08 12:16:21.701384	2025-05-08 12:16:21.701387	500	0
4759	\N	/test/info.php	GET	2025-05-08 12:16:21.857481	2025-05-08 12:16:21.857483	500	0
4760	\N	/info	GET	2025-05-08 12:16:22.019315	2025-05-08 12:16:22.019318	500	0
4761	\N	/info.php	GET	2025-05-08 12:16:22.175617	2025-05-08 12:16:22.17562	500	0
4762	\N	/phpinfo.php3	GET	2025-05-08 12:16:22.336865	2025-05-08 12:16:22.336869	500	0
4763	\N	/phpinfo.php.txt	GET	2025-05-08 12:16:22.494355	2025-05-08 12:16:22.494357	500	0
4764	\N	/info.php.save.1	GET	2025-05-08 12:16:22.703599	2025-05-08 12:16:22.703601	500	0
4765	\N	/info.php_	GET	2025-05-08 12:16:22.86705	2025-05-08 12:16:22.867053	500	0
4766	\N	/_info.php	GET	2025-05-08 12:16:23.022525	2025-05-08 12:16:23.022528	500	0
4767	\N	/admin_phpinfo.php	GET	2025-05-08 12:16:23.178554	2025-05-08 12:16:23.178556	500	0
4768	\N	/dbinfo.php	GET	2025-05-08 12:16:23.338197	2025-05-08 12:16:23.3382	500	0
4769	\N	/admin_info.php	GET	2025-05-08 12:16:23.494235	2025-05-08 12:16:23.494237	500	0
4770	\N	/.aws/credentials	GET	2025-05-08 12:16:23.652295	2025-05-08 12:16:23.652297	500	0
4771	\N	/config.php	GET	2025-05-08 12:16:23.807805	2025-05-08 12:16:23.807806	500	0
4772	\N	/phpinfo.php	GET	2025-05-08 12:16:23.963322	2025-05-08 12:16:23.963324	500	0
4773	\N	/phpinfo.php.bak	GET	2025-05-08 12:16:24.119296	2025-05-08 12:16:24.119299	500	0
4774	\N	/admin/admin_phpinfo.php4	GET	2025-05-08 12:16:24.291258	2025-05-08 12:16:24.291261	500	0
4775	\N	/pinfo.php	GET	2025-05-08 12:16:24.447758	2025-05-08 12:16:24.44776	500	0
4776	\N	/core/dataaccess/tablemetadata.php	GET	2025-05-08 12:16:24.609455	2025-05-08 12:16:24.609458	500	0
4777	\N	/.env.bak	GET	2025-05-08 12:16:24.765333	2025-05-08 12:16:24.765335	500	0
4778	\N	/.env	GET	2025-05-08 12:16:24.92115	2025-05-08 12:16:24.921153	500	0
4779	\N	/.env.backup	GET	2025-05-08 12:16:25.07721	2025-05-08 12:16:25.077213	500	0
4780	\N	/.env_sample	GET	2025-05-08 12:16:25.234009	2025-05-08 12:16:25.234012	500	0
4781	\N	/.env.old	GET	2025-05-08 12:16:25.390903	2025-05-08 12:16:25.390906	500	0
4782	\N	/.env.www	GET	2025-05-08 12:16:25.547554	2025-05-08 12:16:25.547558	500	0
4783	\N	/.docker/.env	GET	2025-05-08 12:16:25.709882	2025-05-08 12:16:25.709885	500	0
4784	\N	/.env.dev	GET	2025-05-08 12:16:25.879907	2025-05-08 12:16:25.879911	500	0
4785	\N	/.env.example	GET	2025-05-08 12:16:26.03654	2025-05-08 12:16:26.036542	500	0
4786	\N	/.env_1	GET	2025-05-08 12:16:26.193218	2025-05-08 12:16:26.193221	500	0
4787	\N	/.env.stage	GET	2025-05-08 12:16:26.350137	2025-05-08 12:16:26.350139	500	0
4788	\N	/config.env	GET	2025-05-08 12:16:26.509263	2025-05-08 12:16:26.509265	500	0
4789	\N	/.pam_environment	GET	2025-05-08 12:16:26.665307	2025-05-08 12:16:26.665309	500	0
4790	\N	/.environment	GET	2025-05-08 12:16:26.821093	2025-05-08 12:16:26.821095	500	0
4791	\N	/.env.production	GET	2025-05-08 12:16:26.976974	2025-05-08 12:16:26.976977	500	0
4792	\N	/.env.production.local	GET	2025-05-08 12:16:27.132767	2025-05-08 12:16:27.13277	500	0
4793	\N	/env.js	GET	2025-05-08 12:16:27.291325	2025-05-08 12:16:27.291328	500	0
4794	\N	/.powenv	GET	2025-05-08 12:16:27.453873	2025-05-08 12:16:27.453887	500	0
4795	\N	/.rbenv-gemsets	GET	2025-05-08 12:16:27.61229	2025-05-08 12:16:27.612293	500	0
4796	\N	/.envs	GET	2025-05-08 12:16:27.768043	2025-05-08 12:16:27.768045	500	0
4797	\N	/.env.docker	GET	2025-05-08 12:16:27.933408	2025-05-08 12:16:27.933412	500	0
4798	\N	/.env.sample	GET	2025-05-08 12:16:28.089145	2025-05-08 12:16:28.089149	500	0
4799	\N	/.env_bak	GET	2025-05-08 12:16:28.247256	2025-05-08 12:16:28.24726	500	0
4800	\N	/.env.php	GET	2025-05-08 12:16:28.419226	2025-05-08 12:16:28.419229	500	0
4801	\N	/.env.development.local	GET	2025-05-08 12:16:28.575853	2025-05-08 12:16:28.575855	500	0
4802	\N	/.env-example	GET	2025-05-08 12:16:28.731261	2025-05-08 12:16:28.731263	500	0
4803	\N	/.env.prod	GET	2025-05-08 12:16:28.90267	2025-05-08 12:16:28.902672	500	0
4804	\N	/.env.docker.dev	GET	2025-05-08 12:16:29.058253	2025-05-08 12:16:29.058257	500	0
4805	\N	/.env.test	GET	2025-05-08 12:16:29.215943	2025-05-08 12:16:29.215945	500	0
3412	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-25 02:28:23.203153	2025-04-25 02:28:23.203286	200	0
3414	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-25 07:31:00.865491	2025-04-25 07:31:00.865564	200	0
3415	\N	/robots.txt	GET	2025-04-25 07:55:33.501205	2025-04-25 07:55:33.501209	500	0
3416	\N	/	GET	2025-04-25 07:55:33.80554	2025-04-25 07:55:33.805544	500	0
3421	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-26 02:58:11.801941	2025-04-26 02:58:11.802081	200	0
3422	\N	/.env	GET	2025-04-26 06:34:14.985654	2025-04-26 06:34:14.985658	500	0
3423	\N	/config/.env	GET	2025-04-26 06:34:15.140911	2025-04-26 06:34:15.140915	500	0
3424	\N	/.env.production	GET	2025-04-26 06:34:15.292653	2025-04-26 06:34:15.292656	500	0
3425	\N	/api/.env	GET	2025-04-26 06:34:15.450547	2025-04-26 06:34:15.450564	500	0
3426	\N	/settings/.env	GET	2025-04-26 06:34:15.910697	2025-04-26 06:34:15.910703	500	0
3427	\N	/config/app.php	GET	2025-04-26 06:34:16.062127	2025-04-26 06:34:16.062133	500	0
3428	\N	/application.yml	GET	2025-04-26 06:34:16.220218	2025-04-26 06:34:16.220226	500	0
3429	\N	/config/database.yml	GET	2025-04-26 06:34:16.373561	2025-04-26 06:34:16.373565	500	0
3430	\N	/secrets.json	GET	2025-04-26 06:34:16.525411	2025-04-26 06:34:16.525416	500	0
3431	\N	/src/config.js	GET	2025-04-26 06:34:16.67712	2025-04-26 06:34:16.677124	500	0
3432	\N	/db.ini	GET	2025-04-26 06:34:16.828692	2025-04-26 06:34:16.828697	500	0
3433	\N	/api/credentials	GET	2025-04-26 06:34:16.98141	2025-04-26 06:34:16.981418	500	0
3434	\N	/.aws/credentials	GET	2025-04-26 06:34:17.144896	2025-04-26 06:34:17.144902	500	0
3435	\N	/secure-config.json	GET	2025-04-26 06:34:17.299459	2025-04-26 06:34:17.299463	500	0
3436	\N	/local_settings.py	GET	2025-04-26 06:34:17.468887	2025-04-26 06:34:17.468893	500	0
3437	\N	/config/default.json	GET	2025-04-26 06:34:17.620312	2025-04-26 06:34:17.620318	500	0
3438	\N	/config/production.json	GET	2025-04-26 06:34:17.885511	2025-04-26 06:34:17.885517	500	0
3439	\N	/bootstrap/cache/config.php	GET	2025-04-26 06:34:18.113701	2025-04-26 06:34:18.113706	500	0
3440	\N	/config/secrets.yml	GET	2025-04-26 06:34:18.265972	2025-04-26 06:34:18.265977	500	0
3441	\N	/settings.yaml	GET	2025-04-26 06:34:18.432965	2025-04-26 06:34:18.43297	500	0
3442	\N	/auth.json	GET	2025-04-26 06:34:18.845519	2025-04-26 06:34:18.845524	500	0
3443	\N	/helm/values.yaml	GET	2025-04-26 06:34:18.997264	2025-04-26 06:34:18.997268	500	0
3444	\N	/docker/.env	GET	2025-04-26 06:34:19.148885	2025-04-26 06:34:19.14889	500	0
3445	\N	/wp-config.php	GET	2025-04-26 06:34:19.313121	2025-04-26 06:34:19.313125	500	0
3446	\N	/config.json	GET	2025-04-26 06:34:19.484772	2025-04-26 06:34:19.484777	500	0
3447	\N	/database.json	GET	2025-04-26 06:34:19.636368	2025-04-26 06:34:19.636373	500	0
3448	\N	/config/secrets.json	GET	2025-04-26 06:34:19.788042	2025-04-26 06:34:19.788048	500	0
3449	\N	/env.backup	GET	2025-04-26 06:34:20.335243	2025-04-26 06:34:20.335248	500	0
3450	\N	/settings.bak	GET	2025-04-26 06:34:20.48973	2025-04-26 06:34:20.489735	500	0
3451	\N	/backup.env	GET	2025-04-26 06:34:20.640981	2025-04-26 06:34:20.640986	500	0
3452	\N	/old/.env	GET	2025-04-26 06:34:20.793304	2025-04-26 06:34:20.793311	500	0
3453	\N	/phpinfo.php	GET	2025-04-26 06:34:20.949041	2025-04-26 06:34:20.949046	500	0
3454	\N	/info.php	GET	2025-04-26 06:34:21.100848	2025-04-26 06:34:21.100853	500	0
3455	\N	/test.php	GET	2025-04-26 06:34:21.264484	2025-04-26 06:34:21.264491	500	0
3456	\N	/laravel/.env	GET	2025-04-26 06:34:21.440179	2025-04-26 06:34:21.440184	500	0
3457	\N	/app/config/.env	GET	2025-04-26 06:34:21.671626	2025-04-26 06:34:21.671634	500	0
3458	\N	/.git/config	GET	2025-04-26 06:34:21.830928	2025-04-26 06:34:21.830933	500	0
3459	\N	/.svn/entries	GET	2025-04-26 06:34:21.999959	2025-04-26 06:34:21.999963	500	0
3460	\N	/.git/HEAD	GET	2025-04-26 06:34:22.152628	2025-04-26 06:34:22.152631	500	0
3461	\N	/.git/index	GET	2025-04-26 06:34:22.304436	2025-04-26 06:34:22.304439	500	0
3462	\N	/.git/logs/HEAD	GET	2025-04-26 06:34:22.456391	2025-04-26 06:34:22.456395	500	0
3463	\N	/.gitignore	GET	2025-04-26 06:34:22.609265	2025-04-26 06:34:22.60927	500	0
3464	\N	/administrator/index.php	GET	2025-04-26 06:34:22.76125	2025-04-26 06:34:22.761255	500	0
3465	\N	/wp-admin/install.php	GET	2025-04-26 06:34:22.913296	2025-04-26 06:34:22.913301	500	0
3466	\N	/joomla/configuration.php-dist	GET	2025-04-26 06:34:23.071806	2025-04-26 06:34:23.07181	500	0
3467	\N	/sites/default/settings.php	GET	2025-04-26 06:34:23.34029	2025-04-26 06:34:23.340295	500	0
3468	\N	/bitrix/php_interface/dbconn.php	GET	2025-04-26 06:34:23.493092	2025-04-26 06:34:23.493097	500	0
3469	\N	/typo3conf/localconf.php	GET	2025-04-26 06:34:23.674384	2025-04-26 06:34:23.674389	500	0
3470	\N	/config.inc.php	GET	2025-04-26 06:34:23.826351	2025-04-26 06:34:23.826356	500	0
3471	\N	/config.old.php	GET	2025-04-26 06:34:23.97792	2025-04-26 06:34:23.977925	500	0
3472	\N	/php.ini	GET	2025-04-26 06:34:24.129379	2025-04-26 06:34:24.129384	500	0
3473	\N	/cgi-bin/phpinfo.php	GET	2025-04-26 06:34:24.302103	2025-04-26 06:34:24.302108	500	0
3474	\N	/debug.php	GET	2025-04-26 06:34:24.459224	2025-04-26 06:34:24.459228	500	0
3475	\N	/server-status	GET	2025-04-26 06:34:24.612475	2025-04-26 06:34:24.612482	500	0
3476	\N	/phpinfo1.php	GET	2025-04-26 06:34:24.764137	2025-04-26 06:34:24.76414	500	0
3477	\N	/phpinfo2.php	GET	2025-04-26 06:34:24.916073	2025-04-26 06:34:24.916077	500	0
3478	\N	/env.txt	GET	2025-04-26 06:34:25.067844	2025-04-26 06:34:25.067849	500	0
3479	\N	/prod.env	GET	2025-04-26 06:34:25.270878	2025-04-26 06:34:25.270881	500	0
3480	\N	/stage.env	GET	2025-04-26 06:34:25.431855	2025-04-26 06:34:25.43186	500	0
3481	\N	/development.env	GET	2025-04-26 06:34:25.585857	2025-04-26 06:34:25.585862	500	0
3482	\N	/credentials.env	GET	2025-04-26 06:34:25.743704	2025-04-26 06:34:25.743709	500	0
3483	\N	/public/.env	GET	2025-04-26 06:34:25.896156	2025-04-26 06:34:25.896159	500	0
3484	\N	/api/config.json	GET	2025-04-26 06:34:26.049411	2025-04-26 06:34:26.049418	500	0
3485	\N	/composer.json	GET	2025-04-26 06:34:26.206238	2025-04-26 06:34:26.206242	500	0
3486	\N	/api/v1/.env	GET	2025-04-26 06:34:26.371304	2025-04-26 06:34:26.371315	500	0
3487	\N	/staging.env	GET	2025-04-26 06:34:26.523203	2025-04-26 06:34:26.523209	500	0
3488	\N	/phpmyadmin/index.php	GET	2025-04-26 06:34:26.674951	2025-04-26 06:34:26.674956	500	0
3489	\N	/backup/config.php	GET	2025-04-26 06:34:26.828432	2025-04-26 06:34:26.828436	500	0
3490	\N	/.env.example	GET	2025-04-26 06:34:26.980378	2025-04-26 06:34:26.980384	500	0
3491	\N	/storage/logs/laravel.log	GET	2025-04-26 06:34:27.132992	2025-04-26 06:34:27.132997	500	0
3492	\N	/storage/framework/sessions/	GET	2025-04-26 06:34:27.295367	2025-04-26 06:34:27.295373	500	0
3493	\N	/storage/framework/cache/	GET	2025-04-26 06:34:27.453738	2025-04-26 06:34:27.453744	500	0
3494	\N	/storage/framework/views/	GET	2025-04-26 06:34:27.605649	2025-04-26 06:34:27.605655	500	0
3495	\N	/nova-api/styles	GET	2025-04-26 06:34:27.758063	2025-04-26 06:34:27.758067	500	0
3496	\N	/.env.local	GET	2025-04-26 06:34:27.958947	2025-04-26 06:34:27.958953	500	0
3497	\N	/.env.dev	GET	2025-04-26 06:34:28.110316	2025-04-26 06:34:28.11032	500	0
3498	\N	/.env.test	GET	2025-04-26 06:34:28.34602	2025-04-26 06:34:28.346024	500	0
3499	\N	/var/logs/dev.log	GET	2025-04-26 06:34:28.507002	2025-04-26 06:34:28.507006	500	0
3500	\N	/var/logs/prod.log	GET	2025-04-26 06:34:28.683668	2025-04-26 06:34:28.683674	500	0
3501	\N	/config/packages/	GET	2025-04-26 06:34:28.834939	2025-04-26 06:34:28.834943	500	0
3502	\N	/web/config.php	GET	2025-04-26 06:34:28.988339	2025-04-26 06:34:28.988348	500	0
3503	\N	/config/routes.yaml	GET	2025-04-26 06:34:29.140207	2025-04-26 06:34:29.140213	500	0
3504	\N	/web.config	GET	2025-04-26 06:34:29.334642	2025-04-26 06:34:29.334647	500	0
3505	\N	/.htaccess	GET	2025-04-26 06:34:29.505182	2025-04-26 06:34:29.505187	500	0
3506	\N	/sites/all/modules/	GET	2025-04-26 06:34:29.656708	2025-04-26 06:34:29.656714	500	0
3507	\N	/sites/all/themes/	GET	2025-04-26 06:34:29.808166	2025-04-26 06:34:29.808171	500	0
3508	\N	/core/install.php	GET	2025-04-26 06:34:29.969498	2025-04-26 06:34:29.969502	500	0
3509	\N	/CHANGELOG.txt	GET	2025-04-26 06:34:30.142157	2025-04-26 06:34:30.142164	500	0
3510	\N	/app/etc/local.xml	GET	2025-04-26 06:34:30.298653	2025-04-26 06:34:30.29866	500	0
3511	\N	/app/etc/env.php	GET	2025-04-26 06:34:30.457365	2025-04-26 06:34:30.457372	500	0
3512	\N	/var/log/system.log	GET	2025-04-26 06:34:30.611466	2025-04-26 06:34:30.611471	500	0
3513	\N	/var/log/exception.log	GET	2025-04-26 06:34:30.762965	2025-04-26 06:34:30.76297	500	0
3514	\N	/.wp-config.php.swp	GET	2025-04-26 06:34:30.914539	2025-04-26 06:34:30.914543	500	0
3515	\N	/wp-config-sample.php	GET	2025-04-26 06:34:31.089371	2025-04-26 06:34:31.089377	500	0
3516	\N	/wp-content/debug.log	GET	2025-04-26 06:34:31.251362	2025-04-26 06:34:31.251367	500	0
3517	\N	/xmlrpc.php	GET	2025-04-26 06:34:31.422629	2025-04-26 06:34:31.422633	500	0
3518	\N	/wp-json/wp/v2/users	GET	2025-04-26 06:34:31.575283	2025-04-26 06:34:31.575287	500	0
3519	\N	/configuration.php~	GET	2025-04-26 06:34:31.726994	2025-04-26 06:34:31.727	500	0
3520	\N	/logs/error.php	GET	2025-04-26 06:34:31.878557	2025-04-26 06:34:31.878562	500	0
3521	\N	/package.json	GET	2025-04-26 06:34:32.033375	2025-04-26 06:34:32.033379	500	0
3522	\N	/yarn.lock	GET	2025-04-26 06:34:32.253903	2025-04-26 06:34:32.253908	500	0
3523	\N	/.npmrc	GET	2025-04-26 06:34:32.411996	2025-04-26 06:34:32.412	500	0
3524	\N	/server.js	GET	2025-04-26 06:34:32.566433	2025-04-26 06:34:32.566438	500	0
3525	\N	/app.js	GET	2025-04-26 06:34:33.303789	2025-04-26 06:34:33.303796	500	0
3526	\N	/config.js	GET	2025-04-26 06:34:33.496325	2025-04-26 06:34:33.49633	500	0
3527	\N	/.dockerignore	GET	2025-04-26 06:34:33.647725	2025-04-26 06:34:33.64773	500	0
3528	\N	/Dockerfile	GET	2025-04-26 06:34:33.799686	2025-04-26 06:34:33.799691	500	0
3529	\N	/.git/logs/	GET	2025-04-26 06:34:33.951404	2025-04-26 06:34:33.951408	500	0
3530	\N	/.git/refs/	GET	2025-04-26 06:34:34.103144	2025-04-26 06:34:34.103149	500	0
3531	\N	/.git/objects/	GET	2025-04-26 06:34:34.283868	2025-04-26 06:34:34.283872	500	0
3532	\N	/.git/packed-refs	GET	2025-04-26 06:34:34.435844	2025-04-26 06:34:34.435848	500	0
3533	\N	/.git/branches/	GET	2025-04-26 06:34:34.593124	2025-04-26 06:34:34.593129	500	0
3534	\N	/api-docs	GET	2025-04-26 06:34:34.747083	2025-04-26 06:34:34.747092	500	0
3535	\N	/swagger.json	GET	2025-04-26 06:34:34.968999	2025-04-26 06:34:34.969003	500	0
3536	\N	/swagger-ui.html	GET	2025-04-26 06:34:35.176517	2025-04-26 06:34:35.176521	500	0
3537	\N	/openapi.json	GET	2025-04-26 06:34:35.333191	2025-04-26 06:34:35.333197	500	0
3538	\N	/backup.sql	GET	2025-04-26 06:34:35.50812	2025-04-26 06:34:35.508124	500	0
3539	\N	/db_backup.sql	GET	2025-04-26 06:34:35.675209	2025-04-26 06:34:35.675213	500	0
3540	\N	/.well-known/security.txt	GET	2025-04-26 06:34:35.831299	2025-04-26 06:34:35.831303	500	0
3557	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-27 07:21:23.258387	2025-04-27 07:21:23.25851	200	0
3558	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-27 07:21:23.321687	2025-04-27 07:21:23.321753	200	0
3559	\N	/.git/config	GET	2025-04-27 09:59:31.753832	2025-04-27 09:59:31.753836	500	0
3567	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:54:52.524391	2025-04-27 16:54:52.530938	200	6
3568	\N	/api/restaurant	GET	2025-04-27 16:54:52.954429	2025-04-27 16:54:52.954489	200	0
3569	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-27 16:54:53.090656	2025-04-27 16:54:53.103968	200	13
3570	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-28 05:29:30.968244	2025-04-28 05:29:30.968334	200	0
3571	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-28 06:57:43.321077	2025-04-28 06:57:43.321189	200	0
3572	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-28 06:57:43.478907	2025-04-28 06:57:43.478966	200	0
3573	\N	/	GET	2025-04-28 07:39:49.183157	2025-04-28 07:39:49.183161	500	0
3574	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-28 12:19:57.086562	2025-04-28 12:19:57.086677	304	0
3575	\N	/api/washingmachines	GET	2025-04-28 12:20:01.604642	2025-04-28 12:20:02.65508	200	1050
3576	\N	/api/restaurant	GET	2025-04-28 12:20:11.938754	2025-04-28 12:20:11.938842	200	0
3577	\N	/status	GET	2025-04-28 12:20:20.33974	2025-04-28 12:20:20.364168	200	24
3578	\N	/api/statistics/global	GET	2025-04-28 12:20:20.591222	2025-04-28 12:20:20.639344	200	48
3579	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:20.862769	2025-04-28 12:20:20.868482	200	5
3580	\N	/status	GET	2025-04-28 12:20:24.151953	2025-04-28 12:20:24.15198	200	0
3581	\N	/status	GET	2025-04-28 12:20:24.475746	2025-04-28 12:20:24.475747	200	0
3582	\N	/api/statistics/global	GET	2025-04-28 12:20:24.611021	2025-04-28 12:20:24.643802	200	32
3583	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:24.842748	2025-04-28 12:20:24.848486	200	5
3584	\N	/status	GET	2025-04-28 12:20:28.945913	2025-04-28 12:20:28.945917	200	0
3585	\N	/status	GET	2025-04-28 12:20:29.242673	2025-04-28 12:20:29.242675	200	0
3586	\N	/api/statistics/global	GET	2025-04-28 12:20:29.395175	2025-04-28 12:20:29.428171	200	32
3587	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:29.737658	2025-04-28 12:20:29.786356	200	48
3588	\N	/status	GET	2025-04-28 12:20:33.954722	2025-04-28 12:20:33.954726	200	0
3589	\N	/status	GET	2025-04-28 12:20:34.345946	2025-04-28 12:20:34.345948	200	0
3590	\N	/api/statistics/global	GET	2025-04-28 12:20:34.346061	2025-04-28 12:20:34.3536	200	7
3591	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:34.651927	2025-04-28 12:20:34.694736	200	42
3592	\N	/status	GET	2025-04-28 12:20:39.353634	2025-04-28 12:20:39.353636	200	0
3593	\N	/api/statistics/global	GET	2025-04-28 12:20:39.580851	2025-04-28 12:20:39.583233	200	2
3594	\N	/status	GET	2025-04-28 12:20:39.71539	2025-04-28 12:20:39.715408	200	0
3595	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:39.803157	2025-04-28 12:20:39.808453	200	5
3596	\N	/status	GET	2025-04-28 12:20:43.927824	2025-04-28 12:20:43.927827	200	0
3597	\N	/status	GET	2025-04-28 12:20:44.139991	2025-04-28 12:20:44.139993	200	0
3598	\N	/api/statistics/global	GET	2025-04-28 12:20:44.273895	2025-04-28 12:20:44.279357	200	5
3599	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:44.517634	2025-04-28 12:20:44.522669	200	5
3600	\N	/status	GET	2025-04-28 12:20:49.21382	2025-04-28 12:20:49.21385	200	0
3601	\N	/status	GET	2025-04-28 12:20:49.426908	2025-04-28 12:20:49.42691	200	0
3602	\N	/api/statistics/global	GET	2025-04-28 12:20:49.57439	2025-04-28 12:20:49.576921	200	2
3603	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:49.772104	2025-04-28 12:20:49.777422	200	5
3604	\N	/status	GET	2025-04-28 12:20:54.434122	2025-04-28 12:20:54.434126	200	0
3605	\N	/status	GET	2025-04-28 12:20:54.638235	2025-04-28 12:20:54.638237	200	0
3606	\N	/api/statistics/global	GET	2025-04-28 12:20:54.650096	2025-04-28 12:20:54.699834	200	49
3607	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:54.998322	2025-04-28 12:20:55.00831	200	9
3608	\N	/status	GET	2025-04-28 12:20:59.477859	2025-04-28 12:20:59.477861	200	0
3609	\N	/status	GET	2025-04-28 12:20:59.672387	2025-04-28 12:20:59.672389	200	0
3610	\N	/api/statistics/global	GET	2025-04-28 12:20:59.687344	2025-04-28 12:20:59.689305	200	1
3611	\N	/api/statistics/endpoints	GET	2025-04-28 12:20:59.902591	2025-04-28 12:20:59.91134	200	8
3612	\N	/status	GET	2025-04-28 12:21:04.428118	2025-04-28 12:21:04.428125	200	0
3613	\N	/status	GET	2025-04-28 12:21:04.644359	2025-04-28 12:21:04.64436	200	0
3614	\N	/api/statistics/global	GET	2025-04-28 12:21:04.644274	2025-04-28 12:21:04.687555	200	43
3615	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:04.955912	2025-04-28 12:21:04.961573	200	5
3616	\N	/status	GET	2025-04-28 12:21:09.666107	2025-04-28 12:21:09.666111	200	0
3617	\N	/status	GET	2025-04-28 12:21:09.878099	2025-04-28 12:21:09.878101	200	0
3618	\N	/api/statistics/global	GET	2025-04-28 12:21:09.878098	2025-04-28 12:21:09.960001	200	81
3619	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:10.166507	2025-04-28 12:21:10.172008	200	5
3620	\N	/status	GET	2025-04-28 12:21:14.431842	2025-04-28 12:21:14.431844	200	0
3621	\N	/status	GET	2025-04-28 12:21:14.649435	2025-04-28 12:21:14.649437	200	0
3622	\N	/api/statistics/global	GET	2025-04-28 12:21:14.649439	2025-04-28 12:21:14.652107	200	2
3623	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:14.857985	2025-04-28 12:21:14.862742	200	4
3624	\N	/status	GET	2025-04-28 12:21:20.796144	2025-04-28 12:21:20.796146	200	0
3625	\N	/api/statistics/global	GET	2025-04-28 12:21:21.048619	2025-04-28 12:21:21.052435	200	3
3626	\N	/status	GET	2025-04-28 12:21:21.183802	2025-04-28 12:21:21.183806	200	0
3627	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:21.282975	2025-04-28 12:21:21.310412	200	27
3628	\N	/status	GET	2025-04-28 12:21:24.464377	2025-04-28 12:21:24.464379	200	0
3629	\N	/status	GET	2025-04-28 12:21:24.683405	2025-04-28 12:21:24.683406	200	0
3630	\N	/api/statistics/global	GET	2025-04-28 12:21:24.683394	2025-04-28 12:21:24.686893	200	3
3631	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:24.895659	2025-04-28 12:21:24.902015	200	6
3632	\N	/status	GET	2025-04-28 12:21:29.44772	2025-04-28 12:21:29.447722	200	0
3633	\N	/status	GET	2025-04-28 12:21:29.739076	2025-04-28 12:21:29.739077	200	0
3634	\N	/api/statistics/global	GET	2025-04-28 12:21:29.739071	2025-04-28 12:21:29.741213	200	2
3635	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:29.963983	2025-04-28 12:21:29.969694	200	5
3636	\N	/status	GET	2025-04-28 12:21:34.442773	2025-04-28 12:21:34.442781	200	0
3637	\N	/status	GET	2025-04-28 12:21:34.646206	2025-04-28 12:21:34.646207	200	0
3638	\N	/api/statistics/global	GET	2025-04-28 12:21:34.783091	2025-04-28 12:21:34.785103	200	2
3639	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:34.992319	2025-04-28 12:21:34.997283	200	4
3640	\N	/status	GET	2025-04-28 12:21:39.804928	2025-04-28 12:21:39.80493	200	0
3641	\N	/status	GET	2025-04-28 12:21:40.027958	2025-04-28 12:21:40.02796	200	0
3642	\N	/api/statistics/global	GET	2025-04-28 12:21:40.169687	2025-04-28 12:21:40.174827	200	5
3643	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:40.415431	2025-04-28 12:21:40.420092	200	4
3644	\N	/status	GET	2025-04-28 12:21:44.453805	2025-04-28 12:21:44.453807	200	0
3645	\N	/status	GET	2025-04-28 12:21:44.672168	2025-04-28 12:21:44.67217	200	0
3646	\N	/api/statistics/global	GET	2025-04-28 12:21:44.682192	2025-04-28 12:21:44.684465	200	2
3647	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:44.886818	2025-04-28 12:21:44.893828	200	7
3648	\N	/status	GET	2025-04-28 12:21:49.464949	2025-04-28 12:21:49.464953	200	0
3650	\N	/api/statistics/global	GET	2025-04-28 12:21:49.680315	2025-04-28 12:21:49.683065	200	2
3651	\N	/api/statistics/endpoints	GET	2025-04-28 12:21:49.899484	2025-04-28 12:21:49.904448	200	4
3652	\N	/status	GET	2025-04-28 12:22:32.471414	2025-04-28 12:22:32.471415	200	0
3653	\N	/status	GET	2025-04-28 12:22:32.672619	2025-04-28 12:22:32.67262	200	0
3654	\N	/api/statistics/global	GET	2025-04-28 12:22:32.677692	2025-04-28 12:22:32.679917	200	2
3655	\N	/api/statistics/endpoints	GET	2025-04-28 12:22:32.886522	2025-04-28 12:22:32.893215	200	6
3656	\N	/status	GET	2025-04-28 12:23:32.443869	2025-04-28 12:23:32.443871	200	0
3658	\N	/api/statistics/global	GET	2025-04-28 12:23:32.659692	2025-04-28 12:23:32.663503	200	3
3659	\N	/api/statistics/endpoints	GET	2025-04-28 12:23:32.862603	2025-04-28 12:23:32.866549	200	3
3660	\N	/status	GET	2025-04-28 12:24:32.43489	2025-04-28 12:24:32.434892	200	0
3661	\N	/status	GET	2025-04-28 12:24:32.629741	2025-04-28 12:24:32.629742	200	0
3662	\N	/api/statistics/global	GET	2025-04-28 12:24:32.638057	2025-04-28 12:24:32.640726	200	2
3663	\N	/api/statistics/endpoints	GET	2025-04-28 12:24:32.859563	2025-04-28 12:24:32.863497	200	3
3664	\N	/status	GET	2025-04-28 12:25:32.470075	2025-04-28 12:25:32.470076	200	0
3666	\N	/api/statistics/global	GET	2025-04-28 12:25:32.698362	2025-04-28 12:25:32.70007	200	1
3667	\N	/api/statistics/endpoints	GET	2025-04-28 12:25:32.934182	2025-04-28 12:25:32.939688	200	5
3668	\N	/status	GET	2025-04-28 12:26:32.487799	2025-04-28 12:26:32.4878	200	0
3669	\N	/status	GET	2025-04-28 12:26:32.708051	2025-04-28 12:26:32.708053	200	0
3670	\N	/api/statistics/global	GET	2025-04-28 12:26:33.747935	2025-04-28 12:26:33.750148	200	2
3671	\N	/api/statistics/endpoints	GET	2025-04-28 12:26:33.986807	2025-04-28 12:26:33.991709	200	4
3672	\N	/status	GET	2025-04-28 12:27:32.440765	2025-04-28 12:27:32.440767	200	0
3673	\N	/status	GET	2025-04-28 12:27:32.710055	2025-04-28 12:27:32.710057	200	0
3682	\N	/api/statistics/global	GET	2025-04-28 12:29:32.65924	2025-04-28 12:29:32.662306	200	3
3683	\N	/api/statistics/endpoints	GET	2025-04-28 12:29:32.914811	2025-04-28 12:29:32.920191	200	5
3684	\N	/status	GET	2025-04-28 12:30:32.461355	2025-04-28 12:30:32.461357	200	0
3685	\N	/status	GET	2025-04-28 12:30:32.672785	2025-04-28 12:30:32.672786	200	0
3686	\N	/api/statistics/global	GET	2025-04-28 12:30:32.677864	2025-04-28 12:30:32.68087	200	3
3687	\N	/api/statistics/endpoints	GET	2025-04-28 12:30:32.894852	2025-04-28 12:30:32.900325	200	5
3688	\N	/status	GET	2025-04-28 12:31:32.690103	2025-04-28 12:31:32.690122	200	0
3689	\N	/api/statistics/global	GET	2025-04-28 12:31:33.006608	2025-04-28 12:31:33.009499	200	2
3690	\N	/status	GET	2025-04-28 12:31:33.141146	2025-04-28 12:31:33.141149	200	0
3691	\N	/api/statistics/endpoints	GET	2025-04-28 12:31:33.209504	2025-04-28 12:31:33.214211	200	4
3692	\N	/status	GET	2025-04-28 12:32:32.423057	2025-04-28 12:32:32.423059	200	0
3693	\N	/status	GET	2025-04-28 12:32:32.627406	2025-04-28 12:32:32.627407	200	0
3694	\N	/api/statistics/global	GET	2025-04-28 12:32:32.632406	2025-04-28 12:32:32.634585	200	2
3695	\N	/api/statistics/endpoints	GET	2025-04-28 12:32:32.832709	2025-04-28 12:32:32.838095	200	5
3696	\N	/status	GET	2025-04-28 12:33:33.781253	2025-04-28 12:33:33.781276	200	0
3697	\N	/status	GET	2025-04-28 12:33:33.985524	2025-04-28 12:33:33.985526	200	0
3698	\N	/api/statistics/global	GET	2025-04-28 12:33:34.122916	2025-04-28 12:33:34.125397	200	2
3699	\N	/api/statistics/endpoints	GET	2025-04-28 12:33:34.35228	2025-04-28 12:33:34.356801	200	4
3700	\N	/status	GET	2025-04-28 12:34:32.732213	2025-04-28 12:34:32.732215	200	0
3701	\N	/status	GET	2025-04-28 12:34:33.010994	2025-04-28 12:34:33.010996	200	0
3714	\N	/api/statistics/global	GET	2025-04-28 12:37:32.71621	2025-04-28 12:37:32.718105	200	1
3715	\N	/api/statistics/endpoints	GET	2025-04-28 12:37:32.93711	2025-04-28 12:37:32.941887	200	4
3716	\N	/status	GET	2025-04-28 12:38:32.451412	2025-04-28 12:38:32.451418	200	0
3717	\N	/status	GET	2025-04-28 12:38:32.660551	2025-04-28 12:38:32.660553	200	0
3722	\N	/api/statistics/global	GET	2025-04-28 12:38:55.822286	2025-04-28 12:38:55.8249	200	2
3723	\N	/api/statistics/endpoints	GET	2025-04-28 12:38:56.038347	2025-04-28 12:38:56.043881	200	5
3724	\N	/status	GET	2025-04-28 12:38:59.29088	2025-04-28 12:38:59.290885	200	0
3726	\N	/api/statistics/global	GET	2025-04-28 12:38:59.517521	2025-04-28 12:38:59.520359	200	2
3727	\N	/api/statistics/endpoints	GET	2025-04-28 12:38:59.726492	2025-04-28 12:38:59.732198	200	5
3728	\N	/status	GET	2025-04-28 12:39:04.296108	2025-04-28 12:39:04.29611	200	0
3729	\N	/status	GET	2025-04-28 12:39:04.511442	2025-04-28 12:39:04.511443	200	0
3730	\N	/api/statistics/global	GET	2025-04-28 12:39:04.530287	2025-04-28 12:39:04.532733	200	2
3731	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:04.736607	2025-04-28 12:39:04.742369	200	5
3732	\N	/status	GET	2025-04-28 12:39:09.297783	2025-04-28 12:39:09.297784	200	0
3734	\N	/api/statistics/global	GET	2025-04-28 12:39:09.508906	2025-04-28 12:39:09.510608	200	1
3735	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:09.737115	2025-04-28 12:39:09.742815	200	5
3736	\N	/status	GET	2025-04-28 12:39:14.292702	2025-04-28 12:39:14.292704	200	0
3738	\N	/status	GET	2025-04-28 12:39:14.517248	2025-04-28 12:39:14.51725	200	0
3739	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:14.731697	2025-04-28 12:39:14.736237	200	4
3740	\N	/status	GET	2025-04-28 12:39:19.296593	2025-04-28 12:39:19.296595	200	0
3741	\N	/status	GET	2025-04-28 12:39:19.525003	2025-04-28 12:39:19.525005	200	0
3742	\N	/api/statistics/global	GET	2025-04-28 12:39:19.535071	2025-04-28 12:39:19.537356	200	2
3743	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:19.739833	2025-04-28 12:39:19.746178	200	6
3744	\N	/status	GET	2025-04-28 12:39:24.299561	2025-04-28 12:39:24.299563	200	0
3745	\N	/status	GET	2025-04-28 12:39:24.533264	2025-04-28 12:39:24.533266	200	0
3746	\N	/api/statistics/global	GET	2025-04-28 12:39:24.540414	2025-04-28 12:39:24.542284	200	1
3747	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:24.751036	2025-04-28 12:39:24.756385	200	5
3748	\N	/status	GET	2025-04-28 12:39:29.300738	2025-04-28 12:39:29.30074	200	0
3750	\N	/api/statistics/global	GET	2025-04-28 12:39:29.521272	2025-04-28 12:39:29.524325	200	3
3751	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:29.747211	2025-04-28 12:39:29.75208	200	4
3752	\N	/status	GET	2025-04-28 12:39:34.302496	2025-04-28 12:39:34.302498	200	0
3753	\N	/status	GET	2025-04-28 12:39:34.550911	2025-04-28 12:39:34.550913	200	0
3757	\N	/status	GET	2025-04-28 12:39:39.53604	2025-04-28 12:39:39.536041	200	0
3761	\N	/status	GET	2025-04-28 12:39:44.521627	2025-04-28 12:39:44.521629	200	0
3766	\N	/api/statistics/global	GET	2025-04-28 12:39:49.504796	2025-04-28 12:39:49.507258	200	2
3767	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:49.720617	2025-04-28 12:39:49.725911	200	5
3768	\N	/status	GET	2025-04-28 12:39:54.302054	2025-04-28 12:39:54.302056	200	0
3769	\N	/status	GET	2025-04-28 12:39:54.510413	2025-04-28 12:39:54.510415	200	0
3778	\N	/api/statistics/global	GET	2025-04-28 12:41:32.506053	2025-04-28 12:41:32.508341	200	2
3779	\N	/api/statistics/endpoints	GET	2025-04-28 12:41:32.711857	2025-04-28 12:41:32.715937	200	4
3780	\N	/status	GET	2025-04-28 12:42:32.294698	2025-04-28 12:42:32.2947	200	0
3782	\N	/api/statistics/global	GET	2025-04-28 12:42:32.509399	2025-04-28 12:42:32.512451	200	3
3783	\N	/api/statistics/endpoints	GET	2025-04-28 12:42:32.718837	2025-04-28 12:42:32.726225	200	7
3784	\N	/status	GET	2025-04-28 12:43:32.320682	2025-04-28 12:43:32.320684	200	0
3785	\N	/status	GET	2025-04-28 12:43:32.538383	2025-04-28 12:43:32.538385	200	0
3786	\N	/api/statistics/global	GET	2025-04-28 12:43:32.548579	2025-04-28 12:43:32.55116	200	2
3787	\N	/api/statistics/endpoints	GET	2025-04-28 12:43:32.757321	2025-04-28 12:43:32.763227	200	5
3788	\N	/status	GET	2025-04-28 12:44:32.307046	2025-04-28 12:44:32.307047	200	0
3789	\N	/status	GET	2025-04-28 12:44:33.130129	2025-04-28 12:44:33.130131	200	0
3801	\N	/status	GET	2025-04-28 12:47:32.531008	2025-04-28 12:47:32.531009	200	0
3814	\N	/api/statistics/global	GET	2025-04-28 12:50:32.725892	2025-04-28 12:50:32.729205	200	3
3815	\N	/api/statistics/endpoints	GET	2025-04-28 12:50:32.937614	2025-04-28 12:50:32.943009	200	5
3816	\N	/status	GET	2025-04-28 12:51:32.299645	2025-04-28 12:51:32.299647	200	0
3817	\N	/status	GET	2025-04-28 12:51:32.529518	2025-04-28 12:51:32.529519	200	0
3833	\N	/status	GET	2025-04-28 13:07:35.237306	2025-04-28 13:07:35.237308	200	0
3674	\N	/api/statistics/global	GET	2025-04-28 12:27:32.710151	2025-04-28 12:27:32.713036	200	2
3675	\N	/api/statistics/endpoints	GET	2025-04-28 12:27:33.070401	2025-04-28 12:27:33.087524	200	17
3676	\N	/status	GET	2025-04-28 12:28:32.587298	2025-04-28 12:28:32.587319	200	0
3677	\N	/status	GET	2025-04-28 12:28:33.071763	2025-04-28 12:28:33.071765	200	0
3678	\N	/api/statistics/global	GET	2025-04-28 12:28:33.215642	2025-04-28 12:28:33.218599	200	2
3679	\N	/api/statistics/endpoints	GET	2025-04-28 12:28:33.558836	2025-04-28 12:28:33.563984	200	5
3680	\N	/status	GET	2025-04-28 12:29:32.442768	2025-04-28 12:29:32.44277	200	0
3681	\N	/status	GET	2025-04-28 12:29:32.656578	2025-04-28 12:29:32.656579	200	0
3702	\N	/api/statistics/global	GET	2025-04-28 12:34:33.157973	2025-04-28 12:34:33.186114	200	28
3703	\N	/api/statistics/endpoints	GET	2025-04-28 12:34:33.387595	2025-04-28 12:34:33.420313	200	32
3704	\N	/status	GET	2025-04-28 12:35:32.590674	2025-04-28 12:35:32.590676	200	0
3705	\N	/status	GET	2025-04-28 12:35:32.801716	2025-04-28 12:35:32.801718	200	0
3706	\N	/api/statistics/global	GET	2025-04-28 12:35:32.93704	2025-04-28 12:35:32.939459	200	2
3707	\N	/api/statistics/endpoints	GET	2025-04-28 12:35:33.143833	2025-04-28 12:35:33.166721	200	22
3708	\N	/status	GET	2025-04-28 12:36:32.593213	2025-04-28 12:36:32.593241	200	0
3709	\N	/status	GET	2025-04-28 12:36:32.820351	2025-04-28 12:36:32.820352	200	0
3710	\N	/api/statistics/global	GET	2025-04-28 12:36:32.955754	2025-04-28 12:36:32.957569	200	1
3711	\N	/api/statistics/endpoints	GET	2025-04-28 12:36:33.172224	2025-04-28 12:36:33.188501	200	16
3712	\N	/status	GET	2025-04-28 12:37:32.488829	2025-04-28 12:37:32.488832	200	0
3713	\N	/status	GET	2025-04-28 12:37:32.714783	2025-04-28 12:37:32.714784	200	0
3718	\N	/api/statistics/global	GET	2025-04-28 12:38:32.660661	2025-04-28 12:38:32.663786	200	3
3719	\N	/api/statistics/endpoints	GET	2025-04-28 12:38:32.866212	2025-04-28 12:38:32.871187	200	4
3720	\N	/status	GET	2025-04-28 12:38:55.610084	2025-04-28 12:38:55.610087	200	0
3721	\N	/status	GET	2025-04-28 12:38:55.822186	2025-04-28 12:38:55.822187	200	0
3725	\N	/status	GET	2025-04-28 12:38:59.517523	2025-04-28 12:38:59.517525	200	0
3733	\N	/status	GET	2025-04-28 12:39:09.50891	2025-04-28 12:39:09.508911	200	0
3737	\N	/api/statistics/global	GET	2025-04-28 12:39:14.517249	2025-04-28 12:39:14.520832	200	3
3749	\N	/status	GET	2025-04-28 12:39:29.521276	2025-04-28 12:39:29.521277	200	0
3754	\N	/api/statistics/global	GET	2025-04-28 12:39:34.552708	2025-04-28 12:39:34.555296	200	2
3755	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:34.773711	2025-04-28 12:39:34.778192	200	4
3756	\N	/status	GET	2025-04-28 12:39:39.298845	2025-04-28 12:39:39.298847	200	0
3758	\N	/api/statistics/global	GET	2025-04-28 12:39:39.535969	2025-04-28 12:39:39.538591	200	2
3759	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:39.769227	2025-04-28 12:39:39.774137	200	4
3760	\N	/status	GET	2025-04-28 12:39:44.296412	2025-04-28 12:39:44.296415	200	0
3762	\N	/api/statistics/global	GET	2025-04-28 12:39:44.521616	2025-04-28 12:39:44.523988	200	2
3763	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:44.754502	2025-04-28 12:39:44.76005	200	5
3764	\N	/status	GET	2025-04-28 12:39:49.28997	2025-04-28 12:39:49.289972	200	0
3765	\N	/status	GET	2025-04-28 12:39:49.504677	2025-04-28 12:39:49.504678	200	0
3770	\N	/api/statistics/global	GET	2025-04-28 12:39:54.512178	2025-04-28 12:39:54.51456	200	2
3771	\N	/api/statistics/endpoints	GET	2025-04-28 12:39:54.718709	2025-04-28 12:39:54.722791	200	4
3772	\N	/status	GET	2025-04-28 12:40:32.298029	2025-04-28 12:40:32.298031	200	0
3773	\N	/status	GET	2025-04-28 12:40:32.515443	2025-04-28 12:40:32.515445	200	0
3774	\N	/api/statistics/global	GET	2025-04-28 12:40:32.530036	2025-04-28 12:40:32.532159	200	2
3775	\N	/api/statistics/endpoints	GET	2025-04-28 12:40:32.740842	2025-04-28 12:40:32.824628	200	83
3776	\N	/status	GET	2025-04-28 12:41:32.295805	2025-04-28 12:41:32.295806	200	0
3777	\N	/status	GET	2025-04-28 12:41:32.505608	2025-04-28 12:41:32.505609	200	0
3781	\N	/status	GET	2025-04-28 12:42:32.509423	2025-04-28 12:42:32.509425	200	0
3790	\N	/api/statistics/global	GET	2025-04-28 12:44:33.140919	2025-04-28 12:44:33.14408	200	3
3791	\N	/api/statistics/endpoints	GET	2025-04-28 12:44:34.017871	2025-04-28 12:44:34.022804	200	4
3792	\N	/status	GET	2025-04-28 12:45:32.295978	2025-04-28 12:45:32.29598	200	0
3793	\N	/status	GET	2025-04-28 12:45:32.496745	2025-04-28 12:45:32.496746	200	0
3794	\N	/api/statistics/global	GET	2025-04-28 12:45:32.506509	2025-04-28 12:45:32.509038	200	2
3795	\N	/api/statistics/endpoints	GET	2025-04-28 12:45:32.70602	2025-04-28 12:45:32.715764	200	9
3796	\N	/status	GET	2025-04-28 12:46:32.316486	2025-04-28 12:46:32.316487	200	0
3797	\N	/status	GET	2025-04-28 12:46:32.533299	2025-04-28 12:46:32.533301	200	0
3798	\N	/api/statistics/global	GET	2025-04-28 12:46:32.542897	2025-04-28 12:46:32.546561	200	3
3799	\N	/api/statistics/endpoints	GET	2025-04-28 12:46:32.767527	2025-04-28 12:46:32.773317	200	5
3800	\N	/status	GET	2025-04-28 12:47:32.303632	2025-04-28 12:47:32.303634	200	0
3802	\N	/api/statistics/global	GET	2025-04-28 12:47:32.531007	2025-04-28 12:47:32.533244	200	2
3803	\N	/api/statistics/endpoints	GET	2025-04-28 12:47:32.743971	2025-04-28 12:47:32.749346	200	5
3804	\N	/status	GET	2025-04-28 12:48:32.302953	2025-04-28 12:48:32.302955	200	0
3805	\N	/status	GET	2025-04-28 12:48:32.50129	2025-04-28 12:48:32.501291	200	0
3806	\N	/api/statistics/global	GET	2025-04-28 12:48:32.505614	2025-04-28 12:48:32.508427	200	2
3807	\N	/api/statistics/endpoints	GET	2025-04-28 12:48:32.729389	2025-04-28 12:48:32.734378	200	4
3808	\N	/status	GET	2025-04-28 12:49:32.307296	2025-04-28 12:49:32.307298	200	0
3809	\N	/status	GET	2025-04-28 12:49:32.52597	2025-04-28 12:49:32.525971	200	0
3810	\N	/api/statistics/global	GET	2025-04-28 12:49:32.539677	2025-04-28 12:49:32.549818	200	10
3811	\N	/api/statistics/endpoints	GET	2025-04-28 12:49:32.767538	2025-04-28 12:49:32.774136	200	6
3812	\N	/status	GET	2025-04-28 12:50:32.349719	2025-04-28 12:50:32.349721	200	0
3813	\N	/status	GET	2025-04-28 12:50:32.725885	2025-04-28 12:50:32.725887	200	0
3818	\N	/api/statistics/global	GET	2025-04-28 12:51:32.530204	2025-04-28 12:51:32.532549	200	2
3819	\N	/api/statistics/endpoints	GET	2025-04-28 12:51:32.739429	2025-04-28 12:51:32.745476	200	6
3820	\N	/status	GET	2025-04-28 12:57:08.98335	2025-04-28 12:57:08.983369	200	0
3821	\N	/status	GET	2025-04-28 12:57:09.594565	2025-04-28 12:57:09.594567	200	0
3822	\N	/api/statistics/global	GET	2025-04-28 12:57:09.740975	2025-04-28 12:57:09.743598	200	2
3823	\N	/api/statistics/endpoints	GET	2025-04-28 12:57:10.521954	2025-04-28 12:57:10.526307	200	4
3824	\N	/status	GET	2025-04-28 13:02:18.420711	2025-04-28 13:02:18.420741	200	0
3825	\N	/status	GET	2025-04-28 13:02:19.131435	2025-04-28 13:02:19.131436	200	0
3826	\N	/api/statistics/global	GET	2025-04-28 13:02:19.270721	2025-04-28 13:02:19.275393	200	4
3827	\N	/api/statistics/endpoints	GET	2025-04-28 13:02:20.026491	2025-04-28 13:02:20.037887	200	11
3828	\N	/status	GET	2025-04-28 13:07:29.151487	2025-04-28 13:07:29.151506	200	0
3829	\N	/status	GET	2025-04-28 13:07:29.752696	2025-04-28 13:07:29.752698	200	0
3830	\N	/api/statistics/global	GET	2025-04-28 13:07:29.888831	2025-04-28 13:07:29.894193	200	5
3831	\N	/api/statistics/endpoints	GET	2025-04-28 13:07:30.499767	2025-04-28 13:07:30.506174	200	6
3832	\N	/status	GET	2025-04-28 13:07:33.539612	2025-04-28 13:07:33.539614	200	0
3834	\N	/api/statistics/global	GET	2025-04-28 13:07:35.237215	2025-04-28 13:07:35.245909	200	8
3835	\N	/api/statistics/endpoints	GET	2025-04-28 13:07:36.55446	2025-04-28 13:07:36.560082	200	5
3836	\N	/status	GET	2025-04-28 13:12:40.260816	2025-04-28 13:12:40.26084	200	0
3837	\N	/status	GET	2025-04-28 13:12:40.968572	2025-04-28 13:12:40.968573	200	0
3838	\N	/api/statistics/global	GET	2025-04-28 13:12:41.10251	2025-04-28 13:12:41.105996	200	3
3839	\N	/api/statistics/endpoints	GET	2025-04-28 13:12:42.241182	2025-04-28 13:12:42.247752	200	6
3840	\N	/status	GET	2025-04-28 13:17:49.256045	2025-04-28 13:17:49.256064	200	0
3841	\N	/status	GET	2025-04-28 13:17:50.529492	2025-04-28 13:17:50.529493	200	0
3842	\N	/api/statistics/global	GET	2025-04-28 13:17:50.66592	2025-04-28 13:17:50.667835	200	1
3843	\N	/api/statistics/endpoints	GET	2025-04-28 13:17:52.082377	2025-04-28 13:17:52.094208	200	11
3844	\N	/status	GET	2025-04-28 13:22:59.946265	2025-04-28 13:22:59.946267	200	0
3845	\N	/status	GET	2025-04-28 13:23:01.590231	2025-04-28 13:23:01.590233	200	0
3846	\N	/api/statistics/global	GET	2025-04-28 13:23:01.723782	2025-04-28 13:23:01.725703	200	1
3847	\N	/api/statistics/endpoints	GET	2025-04-28 13:23:03.462589	2025-04-28 13:23:03.468738	200	6
3848	\N	/status	GET	2025-04-28 13:28:09.245743	2025-04-28 13:28:09.245764	200	0
3849	\N	/api/statistics/global	GET	2025-04-28 13:28:10.01202	2025-04-28 13:28:10.055276	200	43
3850	\N	/status	GET	2025-04-28 13:28:10.145525	2025-04-28 13:28:10.145528	200	0
3851	\N	/api/statistics/endpoints	GET	2025-04-28 13:28:11.079035	2025-04-28 13:28:11.090732	200	11
3852	\N	/status	GET	2025-04-28 13:33:19.670033	2025-04-28 13:33:19.670034	200	0
3854	\N	/api/statistics/global	GET	2025-04-28 13:33:21.495121	2025-04-28 13:33:21.519019	200	23
3855	\N	/api/statistics/endpoints	GET	2025-04-28 13:33:22.799042	2025-04-28 13:33:22.804207	200	5
3856	\N	/status	GET	2025-04-28 13:38:30.434069	2025-04-28 13:38:30.434071	200	0
3857	\N	/api/statistics/global	GET	2025-04-28 13:38:31.777578	2025-04-28 13:38:31.780495	200	2
3858	\N	/status	GET	2025-04-28 13:38:31.911243	2025-04-28 13:38:31.911245	200	0
3859	\N	/api/statistics/endpoints	GET	2025-04-28 13:38:33.238793	2025-04-28 13:38:33.247019	200	8
3860	\N	/status	GET	2025-04-28 13:38:33.406479	2025-04-28 13:38:33.406481	200	0
3861	\N	/status	GET	2025-04-28 13:38:34.465036	2025-04-28 13:38:34.465037	200	0
3872	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-28 14:40:06.596796	2025-04-28 14:40:06.596887	200	0
3873	\N	/	GET	2025-04-28 18:01:36.353527	2025-04-28 18:01:36.35353	500	0
4806	\N	/env	GET	2025-05-08 12:16:29.385354	2025-05-08 12:16:29.385356	500	0
4807	\N	/env.list	GET	2025-05-08 12:16:29.541205	2025-05-08 12:16:29.541208	500	0
4808	\N	/.ghc.environment	GET	2025-05-08 12:16:29.697136	2025-05-08 12:16:29.697141	500	0
4809	\N	/env.txt	GET	2025-05-08 12:16:29.852979	2025-05-08 12:16:29.852981	500	0
4810	\N	/printenv	GET	2025-05-08 12:16:30.02032	2025-05-08 12:16:30.020322	500	0
4811	\N	/.hsenv	GET	2025-05-08 12:16:30.177256	2025-05-08 12:16:30.17726	500	0
4812	\N	/environment	GET	2025-05-08 12:16:30.337235	2025-05-08 12:16:30.337238	500	0
4813	\N	/env.json	GET	2025-05-08 12:16:30.505025	2025-05-08 12:16:30.505027	500	0
4814	\N	/envrc	GET	2025-05-08 12:16:30.668046	2025-05-08 12:16:30.668049	500	0
4815	\N	/.env.dev.local	GET	2025-05-08 12:16:30.837274	2025-05-08 12:16:30.837277	500	0
4816	\N	/environment.ts	GET	2025-05-08 12:16:30.99428	2025-05-08 12:16:30.994283	500	0
4817	\N	/.env_old	GET	2025-05-08 12:16:31.157711	2025-05-08 12:16:31.157713	500	0
4818	\N	/.zshenv	GET	2025-05-08 12:16:31.375378	2025-05-08 12:16:31.37538	500	0
4819	\N	/.env.development.sample	GET	2025-05-08 12:16:31.531953	2025-05-08 12:16:31.531955	500	0
4820	\N	/.rbenv-version	GET	2025-05-08 12:16:31.689035	2025-05-08 12:16:31.689037	500	0
4821	\N	/.env-sample	GET	2025-05-08 12:16:31.845928	2025-05-08 12:16:31.845931	500	0
4822	\N	/.env.prod.local	GET	2025-05-08 12:16:32.011345	2025-05-08 12:16:32.011348	500	0
4823	\N	/.env.travis	GET	2025-05-08 12:16:32.168254	2025-05-08 12:16:32.168257	500	0
4824	\N	/.env.test.sample	GET	2025-05-08 12:16:32.326079	2025-05-08 12:16:32.326082	500	0
4825	\N	/.env.2	GET	2025-05-08 12:16:32.485107	2025-05-08 12:16:32.485109	500	0
4826	\N	/.vscode/.env	GET	2025-05-08 12:16:32.642118	2025-05-08 12:16:32.64212	500	0
4827	\N	/.env.local	GET	2025-05-08 12:16:32.804702	2025-05-08 12:16:32.804705	500	0
4828	\N	/.env.txt	GET	2025-05-08 12:16:32.961687	2025-05-08 12:16:32.961689	500	0
4829	\N	/.env.save	GET	2025-05-08 12:16:33.118268	2025-05-08 12:16:33.11827	500	0
4830	\N	/printenv.tmp	GET	2025-05-08 12:16:33.285327	2025-05-08 12:16:33.28533	500	0
4831	\N	/.flaskenv	GET	2025-05-08 12:16:33.46692	2025-05-08 12:16:33.466922	500	0
4832	\N	/.env.test.local	GET	2025-05-08 12:16:33.623732	2025-05-08 12:16:33.623734	500	0
4833	\N	/.jenv-version	GET	2025-05-08 12:16:33.780892	2025-05-08 12:16:33.780894	500	0
4834	\N	/.envrc	GET	2025-05-08 12:16:33.937262	2025-05-08 12:16:33.937265	500	0
4835	\N	/.env.dist	GET	2025-05-08 12:16:34.138247	2025-05-08 12:16:34.138249	500	0
4836	\N	/.env.sample.php	GET	2025-05-08 12:16:34.308348	2025-05-08 12:16:34.308351	500	0
4837	\N	/.env1	GET	2025-05-08 12:16:34.465186	2025-05-08 12:16:34.465188	500	0
4838	\N	/.venv	GET	2025-05-08 12:16:34.622069	2025-05-08 12:16:34.622071	500	0
4839	\N	/env.prod.js	GET	2025-05-08 12:16:34.778826	2025-05-08 12:16:34.778828	500	0
4840	\N	/env.test.js	GET	2025-05-08 12:16:34.935754	2025-05-08 12:16:34.935758	500	0
4841	\N	/env.bak	GET	2025-05-08 12:16:35.092918	2025-05-08 12:16:35.09292	500	0
4842	\N	/env.php	GET	2025-05-08 12:16:35.266766	2025-05-08 12:16:35.266768	500	0
4850	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 07:51:35.491583	2025-05-09 07:51:35.491665	200	0
4855	\N	/status	GET	2025-05-09 12:07:10.618076	2025-05-09 12:07:10.678132	200	60
4856	\N	/api/statistics/global	GET	2025-05-09 12:07:12.986333	2025-05-09 12:07:12.989143	200	2
4857	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:14.893181	2025-05-09 12:07:14.898521	200	5
4858	\N	/status	GET	2025-05-09 12:07:15.291024	2025-05-09 12:07:15.291042	200	0
4860	\N	/api/statistics/global	GET	2025-05-09 12:07:17.615437	2025-05-09 12:07:17.617154	200	1
4861	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:19.538927	2025-05-09 12:07:19.546084	200	7
4862	\N	/status	GET	2025-05-09 12:07:20.156121	2025-05-09 12:07:20.156128	200	0
4863	\N	/status	GET	2025-05-09 12:07:22.651621	2025-05-09 12:07:22.651653	200	0
4871	\N	/status	GET	2025-05-09 12:07:32.648131	2025-05-09 12:07:32.648152	200	0
4875	\N	/status	GET	2025-05-09 12:07:37.256564	2025-05-09 12:07:37.256566	200	0
4877	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:39.308897	2025-05-09 12:07:39.34963	200	40
4878	\N	/status	GET	2025-05-09 12:07:40.325117	2025-05-09 12:07:40.325118	200	0
4879	\N	/status	GET	2025-05-09 12:07:42.693245	2025-05-09 12:07:42.693264	200	0
4883	\N	/status	GET	2025-05-09 12:07:47.251761	2025-05-09 12:07:47.251765	200	0
4888	\N	/api/statistics/global	GET	2025-05-09 12:07:52.260097	2025-05-09 12:07:52.266587	200	6
4889	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:54.3109	2025-05-09 12:07:54.339654	200	28
4890	\N	/status	GET	2025-05-09 12:07:55.164984	2025-05-09 12:07:55.164986	200	0
4891	\N	/status	GET	2025-05-09 12:07:57.267163	2025-05-09 12:07:57.267165	200	0
4896	\N	/status	GET	2025-05-09 12:08:02.242745	2025-05-09 12:08:02.242746	200	0
4902	\N	/api/statistics/global	GET	2025-05-09 12:08:07.258318	2025-05-09 12:08:07.269107	200	10
4903	\N	/api/statistics/endpoints	GET	2025-05-09 12:08:09.312691	2025-05-09 12:08:09.319888	200	7
4904	\N	/status	GET	2025-05-09 12:08:10.157344	2025-05-09 12:08:10.157345	200	0
4906	\N	/api/statistics/global	GET	2025-05-09 12:08:12.282618	2025-05-09 12:08:12.284846	200	2
4907	\N	/api/statistics/endpoints	GET	2025-05-09 12:08:14.348934	2025-05-09 12:08:14.408094	200	59
4908	\N	/status	GET	2025-05-09 12:08:15.157174	2025-05-09 12:08:15.157182	200	0
4910	\N	/api/statistics/global	GET	2025-05-09 12:08:17.250627	2025-05-09 12:08:17.25702	200	6
4913	\N	/status/	GET	2025-05-09 12:09:50.154549	2025-05-09 12:09:50.154551	200	0
4916	\N	/status	GET	2025-05-09 12:10:45.215699	2025-05-09 12:10:45.215703	200	0
4917	\N	/status	GET	2025-05-09 12:10:45.380484	2025-05-09 12:10:45.380485	200	0
4918	\N	/api/statistics/global	GET	2025-05-09 12:10:45.515174	2025-05-09 12:10:45.516779	200	1
4919	\N	/api/statistics/endpoints	GET	2025-05-09 12:10:45.678838	2025-05-09 12:10:45.683537	200	4
4920	\N	/api/restaurant	GET	2025-05-09 12:10:49.325702	2025-05-09 12:10:49.325757	200	0
4921	\N	/api/washingmachines	GET	2025-05-09 12:10:58.571259	2025-05-09 12:10:59.374138	200	802
4922	\N	/api/washingmachines	GET	2025-05-09 12:11:23.252949	2025-05-09 12:11:23.995289	200	742
4923	\N	/api/washingmachines	GET	2025-05-09 12:11:46.513039	2025-05-09 12:11:47.152583	200	639
4925	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 14:16:16.447262	2025-05-09 14:16:16.773103	200	325
4926	\N	/api/washingmachines	GET	2025-05-09 12:16:46.653511	2025-05-09 12:16:47.333484	200	679
4927	\N	/status	GET	2025-05-09 12:17:49.506144	2025-05-09 12:17:49.506146	200	0
4928	\N	/api/statistics/global	GET	2025-05-09 12:17:49.682379	2025-05-09 12:17:49.68396	200	1
4929	\N	/api/statistics/endpoints	GET	2025-05-09 12:17:49.849966	2025-05-09 12:17:49.854854	200	4
4930	\N	/status	GET	2025-05-09 12:17:54.333302	2025-05-09 12:17:54.33332	200	0
4931	\N	/status	GET	2025-05-09 12:17:54.533872	2025-05-09 12:17:54.533874	200	0
4932	\N	/api/statistics/global	GET	2025-05-09 12:17:54.670146	2025-05-09 12:17:54.671726	200	1
4933	\N	/api/statistics/endpoints	GET	2025-05-09 12:17:54.837857	2025-05-09 12:17:54.842474	200	4
4934	\N	/status	GET	2025-05-09 12:17:59.493822	2025-05-09 12:17:59.49384	200	0
4935	\N	/status	GET	2025-05-09 12:17:59.670464	2025-05-09 12:17:59.670466	200	0
4936	\N	/api/statistics/global	GET	2025-05-09 12:17:59.806382	2025-05-09 12:17:59.807888	200	1
4937	\N	/api/statistics/endpoints	GET	2025-05-09 12:17:59.972774	2025-05-09 12:17:59.979094	200	6
4938	\N	/api/washingmachines	GET	2025-05-09 12:18:01.845737	2025-05-09 12:18:02.606625	200	760
4939	\N	/api/restaurant	GET	2025-05-09 12:18:05.359993	2025-05-09 12:18:05.36004	200	0
4940	\N	/api/washingmachines	GET	2025-05-09 12:18:07.491269	2025-05-09 12:18:07.513455	500	22
4941	\N	/api/washingmachines	GET	2025-05-09 12:23:07.801632	2025-05-09 12:23:08.47476	200	673
4944	\N	/api/statistics/endpoints	GET	2025-05-09 12:30:49.618688	2025-05-09 12:30:49.628523	200	9
4945	\N	/api/auth/login	POST	2025-05-09 12:30:53.634818	2025-05-09 12:30:53.703555	200	68
4946	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 12:30:58.429337	2025-05-09 12:30:58.430894	200	1
3862	\N	/api/statistics/global	GET	2025-04-28 13:38:34.465034	2025-04-28 13:38:34.468173	200	3
3863	\N	/api/statistics/endpoints	GET	2025-04-28 13:38:35.37199	2025-04-28 13:38:35.376825	200	4
3864	\N	/status	GET	2025-04-28 13:43:39.458437	2025-04-28 13:43:39.458462	200	0
3865	\N	/status	GET	2025-04-28 13:43:41.075683	2025-04-28 13:43:41.075685	200	0
3866	\N	/api/statistics/global	GET	2025-04-28 13:43:41.211392	2025-04-28 13:43:41.230821	200	19
3867	\N	/api/statistics/endpoints	GET	2025-04-28 13:43:43.257866	2025-04-28 13:43:43.284112	200	26
3868	\N	/status	GET	2025-04-28 13:48:49.879039	2025-04-28 13:48:49.879059	200	0
3869	\N	/status	GET	2025-04-28 13:48:52.591038	2025-04-28 13:48:52.591039	200	0
3870	\N	/api/statistics/global	GET	2025-04-28 13:48:52.727305	2025-04-28 13:48:52.72998	200	2
3871	\N	/api/statistics/endpoints	GET	2025-04-28 13:48:54.407085	2025-04-28 13:48:54.411558	200	4
3874	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-29 02:38:09.503134	2025-04-29 02:38:09.503239	200	0
3875	\N	/.env	GET	2025-04-29 04:18:48.458626	2025-04-29 04:18:48.458629	500	0
3876	\N	/config/.env	GET	2025-04-29 04:18:48.596792	2025-04-29 04:18:48.596795	500	0
3877	\N	/.env.production	GET	2025-04-29 04:18:48.733104	2025-04-29 04:18:48.733109	500	0
3878	\N	/api/.env	GET	2025-04-29 04:18:48.868639	2025-04-29 04:18:48.868652	500	0
3879	\N	/settings/.env	GET	2025-04-29 04:18:49.004886	2025-04-29 04:18:49.004889	500	0
3880	\N	/config/app.php	GET	2025-04-29 04:18:49.295716	2025-04-29 04:18:49.295721	500	0
3881	\N	/application.yml	GET	2025-04-29 04:18:49.446221	2025-04-29 04:18:49.446228	500	0
3882	\N	/config/database.yml	GET	2025-04-29 04:18:49.582467	2025-04-29 04:18:49.582472	500	0
3883	\N	/secrets.json	GET	2025-04-29 04:18:49.724766	2025-04-29 04:18:49.724771	500	0
3884	\N	/src/config.js	GET	2025-04-29 04:18:49.862335	2025-04-29 04:18:49.862339	500	0
3885	\N	/db.ini	GET	2025-04-29 04:18:49.997464	2025-04-29 04:18:49.997469	500	0
3886	\N	/api/credentials	GET	2025-04-29 04:18:50.133569	2025-04-29 04:18:50.133581	500	0
3887	\N	/.aws/credentials	GET	2025-04-29 04:18:50.603237	2025-04-29 04:18:50.60324	500	0
3888	\N	/secure-config.json	GET	2025-04-29 04:18:50.73847	2025-04-29 04:18:50.738476	500	0
3889	\N	/local_settings.py	GET	2025-04-29 04:18:50.874098	2025-04-29 04:18:50.874102	500	0
3890	\N	/config/default.json	GET	2025-04-29 04:18:51.009509	2025-04-29 04:18:51.009512	500	0
3891	\N	/config/production.json	GET	2025-04-29 04:18:51.149249	2025-04-29 04:18:51.149253	500	0
3892	\N	/bootstrap/cache/config.php	GET	2025-04-29 04:18:51.315554	2025-04-29 04:18:51.315558	500	0
3893	\N	/config/secrets.yml	GET	2025-04-29 04:18:51.471493	2025-04-29 04:18:51.471498	500	0
3894	\N	/settings.yaml	GET	2025-04-29 04:18:51.606937	2025-04-29 04:18:51.60694	500	0
3895	\N	/auth.json	GET	2025-04-29 04:18:51.742011	2025-04-29 04:18:51.742015	500	0
3896	\N	/helm/values.yaml	GET	2025-04-29 04:18:51.878831	2025-04-29 04:18:51.878836	500	0
3897	\N	/docker/.env	GET	2025-04-29 04:18:52.018367	2025-04-29 04:18:52.018371	500	0
3898	\N	/wp-config.php	GET	2025-04-29 04:18:52.162915	2025-04-29 04:18:52.16292	500	0
3899	\N	/config.json	GET	2025-04-29 04:18:52.301259	2025-04-29 04:18:52.301263	500	0
3900	\N	/database.json	GET	2025-04-29 04:18:52.444838	2025-04-29 04:18:52.444842	500	0
3901	\N	/config/secrets.json	GET	2025-04-29 04:18:52.581893	2025-04-29 04:18:52.581898	500	0
3902	\N	/env.backup	GET	2025-04-29 04:18:53.012243	2025-04-29 04:18:53.012248	500	0
3903	\N	/settings.bak	GET	2025-04-29 04:18:53.190883	2025-04-29 04:18:53.190886	500	0
3904	\N	/backup.env	GET	2025-04-29 04:18:53.344234	2025-04-29 04:18:53.344239	500	0
3905	\N	/old/.env	GET	2025-04-29 04:18:53.486071	2025-04-29 04:18:53.486076	500	0
3906	\N	/phpinfo.php	GET	2025-04-29 04:18:53.622915	2025-04-29 04:18:53.62292	500	0
3907	\N	/info.php	GET	2025-04-29 04:18:53.760614	2025-04-29 04:18:53.760618	500	0
3908	\N	/test.php	GET	2025-04-29 04:18:53.900124	2025-04-29 04:18:53.900128	500	0
3909	\N	/laravel/.env	GET	2025-04-29 04:18:54.041391	2025-04-29 04:18:54.041396	500	0
3910	\N	/app/config/.env	GET	2025-04-29 04:18:54.18291	2025-04-29 04:18:54.182918	500	0
3911	\N	/.git/config	GET	2025-04-29 04:18:54.332712	2025-04-29 04:18:54.332717	500	0
3912	\N	/.svn/entries	GET	2025-04-29 04:18:54.570685	2025-04-29 04:18:54.570689	500	0
3913	\N	/.git/HEAD	GET	2025-04-29 04:18:54.70659	2025-04-29 04:18:54.706594	500	0
3914	\N	/.git/index	GET	2025-04-29 04:18:54.865089	2025-04-29 04:18:54.865094	500	0
3915	\N	/.git/logs/HEAD	GET	2025-04-29 04:18:55.000726	2025-04-29 04:18:55.00073	500	0
3916	\N	/.gitignore	GET	2025-04-29 04:18:55.170028	2025-04-29 04:18:55.170032	500	0
3917	\N	/administrator/index.php	GET	2025-04-29 04:18:55.307154	2025-04-29 04:18:55.307158	500	0
3918	\N	/wp-admin/install.php	GET	2025-04-29 04:18:55.550335	2025-04-29 04:18:55.550339	500	0
3919	\N	/joomla/configuration.php-dist	GET	2025-04-29 04:18:55.685482	2025-04-29 04:18:55.685487	500	0
3920	\N	/sites/default/settings.php	GET	2025-04-29 04:18:55.822821	2025-04-29 04:18:55.822827	500	0
3921	\N	/bitrix/php_interface/dbconn.php	GET	2025-04-29 04:18:55.962099	2025-04-29 04:18:55.962103	500	0
3922	\N	/typo3conf/localconf.php	GET	2025-04-29 04:18:56.103985	2025-04-29 04:18:56.10399	500	0
3923	\N	/config.inc.php	GET	2025-04-29 04:18:56.2394	2025-04-29 04:18:56.239404	500	0
3924	\N	/config.old.php	GET	2025-04-29 04:18:56.388633	2025-04-29 04:18:56.388641	500	0
3925	\N	/php.ini	GET	2025-04-29 04:18:56.572087	2025-04-29 04:18:56.572091	500	0
3926	\N	/cgi-bin/phpinfo.php	GET	2025-04-29 04:18:56.7116	2025-04-29 04:18:56.711605	500	0
3927	\N	/debug.php	GET	2025-04-29 04:18:56.847364	2025-04-29 04:18:56.847368	500	0
3928	\N	/server-status	GET	2025-04-29 04:18:57.000464	2025-04-29 04:18:57.000469	500	0
3929	\N	/phpinfo1.php	GET	2025-04-29 04:18:57.13679	2025-04-29 04:18:57.136795	500	0
3930	\N	/phpinfo2.php	GET	2025-04-29 04:18:57.274788	2025-04-29 04:18:57.274793	500	0
3931	\N	/env.txt	GET	2025-04-29 04:18:57.441497	2025-04-29 04:18:57.441501	500	0
3932	\N	/prod.env	GET	2025-04-29 04:18:57.674354	2025-04-29 04:18:57.674359	500	0
3933	\N	/stage.env	GET	2025-04-29 04:18:58.336456	2025-04-29 04:18:58.336461	500	0
3934	\N	/development.env	GET	2025-04-29 04:18:58.491882	2025-04-29 04:18:58.491886	500	0
3935	\N	/credentials.env	GET	2025-04-29 04:18:58.631406	2025-04-29 04:18:58.631412	500	0
3936	\N	/public/.env	GET	2025-04-29 04:18:58.776878	2025-04-29 04:18:58.776883	500	0
3937	\N	/api/config.json	GET	2025-04-29 04:18:58.914078	2025-04-29 04:18:58.914088	500	0
3938	\N	/composer.json	GET	2025-04-29 04:18:59.058242	2025-04-29 04:18:59.058247	500	0
3939	\N	/api/v1/.env	GET	2025-04-29 04:18:59.221231	2025-04-29 04:18:59.221238	500	0
3940	\N	/staging.env	GET	2025-04-29 04:18:59.378601	2025-04-29 04:18:59.378607	500	0
3941	\N	/phpmyadmin/index.php	GET	2025-04-29 04:18:59.514313	2025-04-29 04:18:59.514316	500	0
3942	\N	/backup/config.php	GET	2025-04-29 04:18:59.651424	2025-04-29 04:18:59.651429	500	0
3943	\N	/.env.example	GET	2025-04-29 04:18:59.787673	2025-04-29 04:18:59.787678	500	0
3944	\N	/storage/logs/laravel.log	GET	2025-04-29 04:18:59.922944	2025-04-29 04:18:59.922949	500	0
3945	\N	/storage/framework/sessions/	GET	2025-04-29 04:19:00.059299	2025-04-29 04:19:00.059304	500	0
3946	\N	/storage/framework/cache/	GET	2025-04-29 04:19:00.198493	2025-04-29 04:19:00.198498	500	0
3947	\N	/storage/framework/views/	GET	2025-04-29 04:19:00.337385	2025-04-29 04:19:00.337391	500	0
3948	\N	/nova-api/styles	GET	2025-04-29 04:19:00.477489	2025-04-29 04:19:00.477493	500	0
3949	\N	/.env.local	GET	2025-04-29 04:19:00.614625	2025-04-29 04:19:00.614631	500	0
3950	\N	/.env.dev	GET	2025-04-29 04:19:00.753712	2025-04-29 04:19:00.753717	500	0
3951	\N	/.env.test	GET	2025-04-29 04:19:00.889382	2025-04-29 04:19:00.889387	500	0
3952	\N	/var/logs/dev.log	GET	2025-04-29 04:19:01.025949	2025-04-29 04:19:01.025953	500	0
3953	\N	/var/logs/prod.log	GET	2025-04-29 04:19:01.171969	2025-04-29 04:19:01.171974	500	0
3954	\N	/config/packages/	GET	2025-04-29 04:19:01.318671	2025-04-29 04:19:01.318676	500	0
3955	\N	/web/config.php	GET	2025-04-29 04:19:01.484602	2025-04-29 04:19:01.484607	500	0
3956	\N	/config/routes.yaml	GET	2025-04-29 04:19:01.62263	2025-04-29 04:19:01.622635	500	0
3957	\N	/web.config	GET	2025-04-29 04:19:01.759472	2025-04-29 04:19:01.759476	500	0
3958	\N	/.htaccess	GET	2025-04-29 04:19:01.899388	2025-04-29 04:19:01.899394	500	0
3959	\N	/sites/all/modules/	GET	2025-04-29 04:19:02.046238	2025-04-29 04:19:02.046242	500	0
3960	\N	/sites/all/themes/	GET	2025-04-29 04:19:02.190784	2025-04-29 04:19:02.190788	500	0
3961	\N	/core/install.php	GET	2025-04-29 04:19:02.337744	2025-04-29 04:19:02.337749	500	0
3962	\N	/CHANGELOG.txt	GET	2025-04-29 04:19:02.482243	2025-04-29 04:19:02.482248	500	0
3963	\N	/app/etc/local.xml	GET	2025-04-29 04:19:02.621772	2025-04-29 04:19:02.62178	500	0
3964	\N	/app/etc/env.php	GET	2025-04-29 04:19:02.939392	2025-04-29 04:19:02.939399	500	0
3965	\N	/var/log/system.log	GET	2025-04-29 04:19:03.147585	2025-04-29 04:19:03.14759	500	0
3966	\N	/var/log/exception.log	GET	2025-04-29 04:19:03.322569	2025-04-29 04:19:03.322573	500	0
3967	\N	/.wp-config.php.swp	GET	2025-04-29 04:19:03.52154	2025-04-29 04:19:03.521544	500	0
3968	\N	/wp-config-sample.php	GET	2025-04-29 04:19:03.660508	2025-04-29 04:19:03.660512	500	0
3969	\N	/wp-content/debug.log	GET	2025-04-29 04:19:03.797477	2025-04-29 04:19:03.797481	500	0
3970	\N	/xmlrpc.php	GET	2025-04-29 04:19:03.937966	2025-04-29 04:19:03.937971	500	0
3971	\N	/wp-json/wp/v2/users	GET	2025-04-29 04:19:04.191846	2025-04-29 04:19:04.191852	500	0
3972	\N	/configuration.php~	GET	2025-04-29 04:19:04.346471	2025-04-29 04:19:04.346475	500	0
3973	\N	/logs/error.php	GET	2025-04-29 04:19:04.489071	2025-04-29 04:19:04.489075	500	0
3974	\N	/package.json	GET	2025-04-29 04:19:04.635435	2025-04-29 04:19:04.63544	500	0
3975	\N	/yarn.lock	GET	2025-04-29 04:19:04.7714	2025-04-29 04:19:04.771405	500	0
3976	\N	/.npmrc	GET	2025-04-29 04:19:04.910979	2025-04-29 04:19:04.910985	500	0
3977	\N	/server.js	GET	2025-04-29 04:19:05.052043	2025-04-29 04:19:05.052049	500	0
3978	\N	/app.js	GET	2025-04-29 04:19:05.48115	2025-04-29 04:19:05.48116	500	0
3979	\N	/config.js	GET	2025-04-29 04:19:05.618949	2025-04-29 04:19:05.618953	500	0
3980	\N	/.dockerignore	GET	2025-04-29 04:19:05.75488	2025-04-29 04:19:05.754884	500	0
3981	\N	/Dockerfile	GET	2025-04-29 04:19:05.89032	2025-04-29 04:19:05.890326	500	0
3982	\N	/.git/logs/	GET	2025-04-29 04:19:06.030935	2025-04-29 04:19:06.030939	500	0
3983	\N	/.git/refs/	GET	2025-04-29 04:19:06.171845	2025-04-29 04:19:06.17185	500	0
3984	\N	/.git/objects/	GET	2025-04-29 04:19:06.32601	2025-04-29 04:19:06.326016	500	0
3985	\N	/.git/packed-refs	GET	2025-04-29 04:19:06.663729	2025-04-29 04:19:06.663735	500	0
3986	\N	/.git/branches/	GET	2025-04-29 04:19:06.810534	2025-04-29 04:19:06.81054	500	0
3987	\N	/api-docs	GET	2025-04-29 04:19:06.954866	2025-04-29 04:19:06.954877	500	0
3988	\N	/swagger.json	GET	2025-04-29 04:19:07.091271	2025-04-29 04:19:07.091276	500	0
3989	\N	/swagger-ui.html	GET	2025-04-29 04:19:07.229556	2025-04-29 04:19:07.229562	500	0
3990	\N	/openapi.json	GET	2025-04-29 04:19:07.366182	2025-04-29 04:19:07.366187	500	0
3991	\N	/backup.sql	GET	2025-04-29 04:19:07.501554	2025-04-29 04:19:07.50156	500	0
3992	\N	/db_backup.sql	GET	2025-04-29 04:19:07.645717	2025-04-29 04:19:07.645722	500	0
3993	\N	/.well-known/security.txt	GET	2025-04-29 04:19:07.829114	2025-04-29 04:19:07.829119	500	0
3994	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-29 08:19:21.845692	2025-04-29 08:19:21.845791	200	0
3995	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-29 08:19:22.382948	2025-04-29 08:19:22.383027	200	0
3996	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-29 10:00:45.333743	2025-04-29 10:00:45.354087	200	20
3997	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-29 10:00:45.669032	2025-04-29 10:00:45.670847	200	1
3998	\N	/api/restaurant	GET	2025-04-29 10:00:45.806913	2025-04-29 10:00:45.807003	200	0
3999	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-29 11:01:38.011252	2025-04-29 11:01:38.011392	200	0
4000	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-29 16:59:41.035088	2025-04-29 16:59:41.037036	200	1
4001	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-29 16:59:41.368313	2025-04-29 16:59:41.369705	200	1
4002	\N	/api/restaurant	GET	2025-04-29 16:59:41.517198	2025-04-29 16:59:41.517261	200	0
4003	\N	/api/restaurant	GET	2025-04-29 16:59:44.225928	2025-04-29 16:59:44.225988	200	0
4004	\N	/	GET	2025-04-29 19:32:19.187258	2025-04-29 19:32:19.187262	500	0
4005	\N	/	GET	2025-04-29 23:08:10.021477	2025-04-29 23:08:10.02148	500	0
4006	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-30 02:42:35.042422	2025-04-30 02:42:35.042538	200	0
4007	\N	/	GET	2025-04-30 05:13:22.832428	2025-04-30 05:13:22.832432	500	0
4008	\N	/favicon.ico	GET	2025-04-30 05:13:25.579332	2025-04-30 05:13:25.579334	500	0
4009	\N	/ads.txt	GET	2025-04-30 05:13:25.783587	2025-04-30 05:13:25.783591	500	0
4010	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-30 07:02:10.76015	2025-04-30 07:02:10.760273	200	0
4011	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-30 07:02:11.01136	2025-04-30 07:02:11.011451	200	0
4012	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-30 09:59:18.878592	2025-04-30 09:59:18.87998	200	1
4013	\N	/api/restaurant	GET	2025-04-30 09:59:19.424151	2025-04-30 09:59:19.42421	200	0
4014	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-30 09:59:19.622822	2025-04-30 09:59:19.624042	200	1
4015	\N	/api/restaurant	GET	2025-04-30 09:59:20.553176	2025-04-30 09:59:20.553233	200	0
4016	\N	/api/restaurant	GET	2025-04-30 09:59:21.88042	2025-04-30 09:59:21.880474	200	0
4017	\N	/api/restaurant	GET	2025-04-30 09:59:22.490749	2025-04-30 09:59:22.490842	200	0
4018	\N	/api/restaurant	GET	2025-04-30 09:59:23.925541	2025-04-30 09:59:23.925605	200	0
4019	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-30 09:59:35.464159	2025-04-30 09:59:35.464267	200	0
4020	\N	/api/restaurant	GET	2025-04-30 09:59:38.248566	2025-04-30 09:59:38.24868	200	0
4021	\N	/	GET	2025-04-30 10:57:36.44469	2025-04-30 10:57:36.444693	500	0
4022	\N	/favicon.ico	GET	2025-04-30 10:57:51.221671	2025-04-30 10:57:51.221673	500	0
4023	\N	/.env	GET	2025-04-30 11:12:10.052623	2025-04-30 11:12:10.052627	500	0
4024	\N	/config/.env	GET	2025-04-30 11:12:10.217095	2025-04-30 11:12:10.217098	500	0
4025	\N	/.env.production	GET	2025-04-30 11:12:10.443764	2025-04-30 11:12:10.443767	500	0
4026	\N	/api/.env	GET	2025-04-30 11:12:10.602804	2025-04-30 11:12:10.602817	500	0
4027	\N	/settings/.env	GET	2025-04-30 11:12:10.763544	2025-04-30 11:12:10.763548	500	0
4028	\N	/config/app.php	GET	2025-04-30 11:12:10.929285	2025-04-30 11:12:10.92929	500	0
4029	\N	/application.yml	GET	2025-04-30 11:12:11.09223	2025-04-30 11:12:11.092238	500	0
4030	\N	/config/database.yml	GET	2025-04-30 11:12:11.254376	2025-04-30 11:12:11.254381	500	0
4031	\N	/secrets.json	GET	2025-04-30 11:12:11.411818	2025-04-30 11:12:11.411823	500	0
4032	\N	/src/config.js	GET	2025-04-30 11:12:11.570162	2025-04-30 11:12:11.570167	500	0
4033	\N	/db.ini	GET	2025-04-30 11:12:11.744168	2025-04-30 11:12:11.744172	500	0
4034	\N	/api/credentials	GET	2025-04-30 11:12:11.900515	2025-04-30 11:12:11.900522	500	0
4035	\N	/.aws/credentials	GET	2025-04-30 11:12:12.056793	2025-04-30 11:12:12.056797	500	0
4036	\N	/secure-config.json	GET	2025-04-30 11:12:12.213769	2025-04-30 11:12:12.213775	500	0
4037	\N	/local_settings.py	GET	2025-04-30 11:12:12.376558	2025-04-30 11:12:12.376562	500	0
4038	\N	/config/default.json	GET	2025-04-30 11:12:12.557002	2025-04-30 11:12:12.557006	500	0
4039	\N	/config/production.json	GET	2025-04-30 11:12:12.713471	2025-04-30 11:12:12.713475	500	0
4040	\N	/bootstrap/cache/config.php	GET	2025-04-30 11:12:12.934643	2025-04-30 11:12:12.934648	500	0
4041	\N	/config/secrets.yml	GET	2025-04-30 11:12:13.093037	2025-04-30 11:12:13.093042	500	0
4042	\N	/settings.yaml	GET	2025-04-30 11:12:13.250985	2025-04-30 11:12:13.25099	500	0
4043	\N	/auth.json	GET	2025-04-30 11:12:13.415431	2025-04-30 11:12:13.415436	500	0
4044	\N	/helm/values.yaml	GET	2025-04-30 11:12:13.588182	2025-04-30 11:12:13.588185	500	0
4045	\N	/docker/.env	GET	2025-04-30 11:12:13.759191	2025-04-30 11:12:13.759196	500	0
4046	\N	/wp-config.php	GET	2025-04-30 11:12:13.9241	2025-04-30 11:12:13.924105	500	0
4047	\N	/config.json	GET	2025-04-30 11:12:14.081317	2025-04-30 11:12:14.081321	500	0
4048	\N	/database.json	GET	2025-04-30 11:12:14.240286	2025-04-30 11:12:14.240291	500	0
4049	\N	/config/secrets.json	GET	2025-04-30 11:12:14.407283	2025-04-30 11:12:14.407287	500	0
4050	\N	/env.backup	GET	2025-04-30 11:12:14.564539	2025-04-30 11:12:14.564543	500	0
4051	\N	/settings.bak	GET	2025-04-30 11:12:14.721032	2025-04-30 11:12:14.721043	500	0
4052	\N	/backup.env	GET	2025-04-30 11:12:14.877086	2025-04-30 11:12:14.87709	500	0
4053	\N	/old/.env	GET	2025-04-30 11:12:15.437195	2025-04-30 11:12:15.4372	500	0
4054	\N	/phpinfo.php	GET	2025-04-30 11:12:15.593645	2025-04-30 11:12:15.593649	500	0
4055	\N	/info.php	GET	2025-04-30 11:12:15.75616	2025-04-30 11:12:15.756164	500	0
4056	\N	/test.php	GET	2025-04-30 11:12:15.913889	2025-04-30 11:12:15.913893	500	0
4057	\N	/laravel/.env	GET	2025-04-30 11:12:16.07022	2025-04-30 11:12:16.070224	500	0
4058	\N	/app/config/.env	GET	2025-04-30 11:12:16.254393	2025-04-30 11:12:16.254401	500	0
4059	\N	/.git/config	GET	2025-04-30 11:12:16.411564	2025-04-30 11:12:16.411567	500	0
4060	\N	/.svn/entries	GET	2025-04-30 11:12:16.568676	2025-04-30 11:12:16.568681	500	0
4061	\N	/.git/HEAD	GET	2025-04-30 11:12:16.725043	2025-04-30 11:12:16.725046	500	0
4062	\N	/.git/index	GET	2025-04-30 11:12:16.880935	2025-04-30 11:12:16.880938	500	0
4063	\N	/.git/logs/HEAD	GET	2025-04-30 11:12:17.038427	2025-04-30 11:12:17.038432	500	0
4064	\N	/.gitignore	GET	2025-04-30 11:12:17.196285	2025-04-30 11:12:17.19629	500	0
4065	\N	/administrator/index.php	GET	2025-04-30 11:12:17.605548	2025-04-30 11:12:17.605554	500	0
4066	\N	/wp-admin/install.php	GET	2025-04-30 11:12:17.773467	2025-04-30 11:12:17.773472	500	0
4067	\N	/joomla/configuration.php-dist	GET	2025-04-30 11:12:17.929749	2025-04-30 11:12:17.929753	500	0
4068	\N	/sites/default/settings.php	GET	2025-04-30 11:12:18.088446	2025-04-30 11:12:18.088449	500	0
4069	\N	/bitrix/php_interface/dbconn.php	GET	2025-04-30 11:12:18.251525	2025-04-30 11:12:18.251529	500	0
4070	\N	/typo3conf/localconf.php	GET	2025-04-30 11:12:18.419534	2025-04-30 11:12:18.419538	500	0
4071	\N	/config.inc.php	GET	2025-04-30 11:12:18.575896	2025-04-30 11:12:18.575903	500	0
4072	\N	/config.old.php	GET	2025-04-30 11:12:18.732274	2025-04-30 11:12:18.73228	500	0
4073	\N	/php.ini	GET	2025-04-30 11:12:18.888596	2025-04-30 11:12:18.8886	500	0
4074	\N	/cgi-bin/phpinfo.php	GET	2025-04-30 11:12:19.055205	2025-04-30 11:12:19.055211	500	0
4075	\N	/debug.php	GET	2025-04-30 11:12:19.217538	2025-04-30 11:12:19.217542	500	0
4076	\N	/server-status	GET	2025-04-30 11:12:19.394504	2025-04-30 11:12:19.39451	500	0
4077	\N	/phpinfo1.php	GET	2025-04-30 11:12:19.557658	2025-04-30 11:12:19.557662	500	0
4078	\N	/phpinfo2.php	GET	2025-04-30 11:12:19.715375	2025-04-30 11:12:19.715379	500	0
4079	\N	/env.txt	GET	2025-04-30 11:12:19.872061	2025-04-30 11:12:19.872064	500	0
4080	\N	/prod.env	GET	2025-04-30 11:12:20.028077	2025-04-30 11:12:20.028081	500	0
4081	\N	/stage.env	GET	2025-04-30 11:12:20.190997	2025-04-30 11:12:20.191001	500	0
4082	\N	/development.env	GET	2025-04-30 11:12:20.347935	2025-04-30 11:12:20.347938	500	0
4083	\N	/credentials.env	GET	2025-04-30 11:12:20.53616	2025-04-30 11:12:20.536164	500	0
4084	\N	/public/.env	GET	2025-04-30 11:12:20.706543	2025-04-30 11:12:20.706547	500	0
4085	\N	/api/config.json	GET	2025-04-30 11:12:20.868094	2025-04-30 11:12:20.868101	500	0
4086	\N	/composer.json	GET	2025-04-30 11:12:21.036283	2025-04-30 11:12:21.036288	500	0
4087	\N	/api/v1/.env	GET	2025-04-30 11:12:21.196288	2025-04-30 11:12:21.196296	500	0
4088	\N	/staging.env	GET	2025-04-30 11:12:21.368617	2025-04-30 11:12:21.368621	500	0
4089	\N	/phpmyadmin/index.php	GET	2025-04-30 11:12:21.547603	2025-04-30 11:12:21.547609	500	0
4090	\N	/backup/config.php	GET	2025-04-30 11:12:21.703934	2025-04-30 11:12:21.703939	500	0
4091	\N	/.env.example	GET	2025-04-30 11:12:21.866597	2025-04-30 11:12:21.866601	500	0
4092	\N	/storage/logs/laravel.log	GET	2025-04-30 11:12:22.022811	2025-04-30 11:12:22.022817	500	0
4093	\N	/storage/framework/sessions/	GET	2025-04-30 11:12:22.182374	2025-04-30 11:12:22.18238	500	0
4094	\N	/storage/framework/cache/	GET	2025-04-30 11:12:22.340086	2025-04-30 11:12:22.340091	500	0
4095	\N	/storage/framework/views/	GET	2025-04-30 11:12:22.498626	2025-04-30 11:12:22.498633	500	0
4096	\N	/nova-api/styles	GET	2025-04-30 11:12:22.659812	2025-04-30 11:12:22.659816	500	0
4097	\N	/.env.local	GET	2025-04-30 11:12:22.816935	2025-04-30 11:12:22.816939	500	0
4098	\N	/.env.dev	GET	2025-04-30 11:12:22.973101	2025-04-30 11:12:22.973105	500	0
4099	\N	/.env.test	GET	2025-04-30 11:12:23.129408	2025-04-30 11:12:23.129412	500	0
4100	\N	/var/logs/dev.log	GET	2025-04-30 11:12:23.286416	2025-04-30 11:12:23.28642	500	0
4101	\N	/var/logs/prod.log	GET	2025-04-30 11:12:23.442665	2025-04-30 11:12:23.442668	500	0
4102	\N	/config/packages/	GET	2025-04-30 11:12:23.598819	2025-04-30 11:12:23.598824	500	0
4103	\N	/web/config.php	GET	2025-04-30 11:12:23.761913	2025-04-30 11:12:23.761916	500	0
4104	\N	/config/routes.yaml	GET	2025-04-30 11:12:23.91804	2025-04-30 11:12:23.918044	500	0
4105	\N	/web.config	GET	2025-04-30 11:12:24.074349	2025-04-30 11:12:24.074353	500	0
4106	\N	/.htaccess	GET	2025-04-30 11:12:24.230827	2025-04-30 11:12:24.230832	500	0
4107	\N	/sites/all/modules/	GET	2025-04-30 11:12:24.392814	2025-04-30 11:12:24.392819	500	0
4108	\N	/sites/all/themes/	GET	2025-04-30 11:12:24.54998	2025-04-30 11:12:24.549984	500	0
4109	\N	/core/install.php	GET	2025-04-30 11:12:24.715176	2025-04-30 11:12:24.71518	500	0
4110	\N	/CHANGELOG.txt	GET	2025-04-30 11:12:24.871539	2025-04-30 11:12:24.871544	500	0
4111	\N	/app/etc/local.xml	GET	2025-04-30 11:12:25.032562	2025-04-30 11:12:25.032569	500	0
4112	\N	/app/etc/env.php	GET	2025-04-30 11:12:25.196666	2025-04-30 11:12:25.196673	500	0
4113	\N	/var/log/system.log	GET	2025-04-30 11:12:25.353855	2025-04-30 11:12:25.353858	500	0
4114	\N	/var/log/exception.log	GET	2025-04-30 11:12:25.527618	2025-04-30 11:12:25.527622	500	0
4115	\N	/.wp-config.php.swp	GET	2025-04-30 11:12:25.684938	2025-04-30 11:12:25.684942	500	0
4116	\N	/wp-config-sample.php	GET	2025-04-30 11:12:25.840945	2025-04-30 11:12:25.840948	500	0
4117	\N	/wp-content/debug.log	GET	2025-04-30 11:12:25.999727	2025-04-30 11:12:25.999732	500	0
4118	\N	/xmlrpc.php	GET	2025-04-30 11:12:26.157537	2025-04-30 11:12:26.157541	500	0
4119	\N	/wp-json/wp/v2/users	GET	2025-04-30 11:12:26.33787	2025-04-30 11:12:26.337876	500	0
4120	\N	/configuration.php~	GET	2025-04-30 11:12:26.509849	2025-04-30 11:12:26.509853	500	0
4121	\N	/logs/error.php	GET	2025-04-30 11:12:26.669659	2025-04-30 11:12:26.669663	500	0
4122	\N	/package.json	GET	2025-04-30 11:12:26.826505	2025-04-30 11:12:26.826509	500	0
4123	\N	/yarn.lock	GET	2025-04-30 11:12:26.982828	2025-04-30 11:12:26.982831	500	0
4124	\N	/.npmrc	GET	2025-04-30 11:12:27.156365	2025-04-30 11:12:27.156369	500	0
4125	\N	/server.js	GET	2025-04-30 11:12:27.316204	2025-04-30 11:12:27.316208	500	0
4126	\N	/app.js	GET	2025-04-30 11:12:27.489892	2025-04-30 11:12:27.489901	500	0
4127	\N	/config.js	GET	2025-04-30 11:12:27.645952	2025-04-30 11:12:27.645956	500	0
4128	\N	/.dockerignore	GET	2025-04-30 11:12:27.802703	2025-04-30 11:12:27.802707	500	0
4129	\N	/Dockerfile	GET	2025-04-30 11:12:27.958758	2025-04-30 11:12:27.958762	500	0
4130	\N	/.git/logs/	GET	2025-04-30 11:12:28.115216	2025-04-30 11:12:28.115221	500	0
4131	\N	/.git/refs/	GET	2025-04-30 11:12:28.287322	2025-04-30 11:12:28.287326	500	0
4132	\N	/.git/objects/	GET	2025-04-30 11:12:28.456002	2025-04-30 11:12:28.456006	500	0
4133	\N	/.git/packed-refs	GET	2025-04-30 11:12:28.617229	2025-04-30 11:12:28.617233	500	0
4134	\N	/.git/branches/	GET	2025-04-30 11:12:28.773425	2025-04-30 11:12:28.773432	500	0
4135	\N	/api-docs	GET	2025-04-30 11:12:28.937597	2025-04-30 11:12:28.937603	500	0
4136	\N	/swagger.json	GET	2025-04-30 11:12:29.093905	2025-04-30 11:12:29.093909	500	0
4137	\N	/swagger-ui.html	GET	2025-04-30 11:12:29.255064	2025-04-30 11:12:29.255068	500	0
4138	\N	/openapi.json	GET	2025-04-30 11:12:29.430516	2025-04-30 11:12:29.430521	500	0
4139	\N	/backup.sql	GET	2025-04-30 11:12:29.59485	2025-04-30 11:12:29.594854	500	0
4140	\N	/db_backup.sql	GET	2025-04-30 11:12:29.752645	2025-04-30 11:12:29.75265	500	0
4141	\N	/.well-known/security.txt	GET	2025-04-30 11:12:29.909422	2025-04-30 11:12:29.909427	500	0
4142	\N	/favicon.ico	GET	2025-04-30 12:00:19.855971	2025-04-30 12:00:19.855975	500	0
4150	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-01 02:17:46.565898	2025-05-01 02:17:46.565992	200	0
4159	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-01 22:03:54.182232	2025-05-01 22:03:54.197549	200	15
4160	\N	/api/restaurant	GET	2025-05-01 22:04:00.442876	2025-05-01 22:04:00.442965	200	0
4161	\N	/api/restaurant	GET	2025-05-01 22:04:03.951008	2025-05-01 22:04:03.951073	200	0
4162	\N	/api/traq/	GET	2025-05-01 22:04:08.16823	2025-05-01 22:04:08.170164	200	1
4163	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:04:14.47484	2025-05-01 22:04:14.476734	200	1
4164	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-01 22:04:15.994464	2025-05-01 22:04:15.994566	200	0
4165	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:04:33.248624	2025-05-01 22:04:33.270179	200	21
4166	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-01 22:04:37.176822	2025-05-01 22:04:37.199925	200	23
4167	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:04:39.384961	2025-05-01 22:04:39.404696	200	19
4170	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 06:51:58.490861	2025-05-02 06:51:58.490956	200	0
4171	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 06:51:58.638336	2025-05-02 06:51:58.638401	200	0
4172	\N	/.git/config	GET	2025-05-02 07:51:06.09454	2025-05-02 07:51:06.094542	500	0
4173	\N	/.git/config	GET	2025-05-02 07:51:13.498232	2025-05-02 07:51:13.498236	500	0
4174	\N	/	GET	2025-05-02 09:35:43.987173	2025-05-02 09:35:43.987177	500	0
4175	\N	/	GET	2025-05-02 09:48:43.575586	2025-05-02 09:48:43.575588	500	0
4176	\N	/favicon.ico	GET	2025-05-02 09:49:23.747859	2025-05-02 09:49:23.747862	500	0
4182	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 15:15:50.512635	2025-05-02 15:15:50.512756	200	0
4183	\N	/	GET	2025-05-02 16:17:27.780625	2025-05-02 16:17:27.780627	500	0
4184	\N	/api/auth/login	POST	2025-05-02 21:24:48.479637	2025-05-02 21:24:48.547883	200	68
4143	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-04-30 13:39:19.834237	2025-04-30 13:39:19.834347	200	0
4144	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-30 16:46:55.107734	2025-04-30 16:46:55.109821	200	2
4145	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-30 16:46:55.648497	2025-04-30 16:46:55.649761	200	1
4146	\N	/api/restaurant	GET	2025-04-30 16:46:55.818038	2025-04-30 16:46:55.818104	200	0
4147	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-04-30 16:47:18.007095	2025-04-30 16:47:18.010097	200	3
4148	\N	/api/traq/	GET	2025-04-30 16:47:27.702198	2025-04-30 16:47:27.722133	200	19
4149	\N	/api/restaurant	GET	2025-04-30 16:47:30.00791	2025-04-30 16:47:30.008014	200	0
4151	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-01 07:27:13.120158	2025-05-01 07:27:13.120296	200	0
4152	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-01 07:27:13.219905	2025-05-01 07:27:13.219992	200	0
4153	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-01 15:05:10.957046	2025-05-01 15:05:10.957142	200	0
4154	\N	/api/auth/login	POST	2025-05-01 22:02:31.554984	2025-05-01 22:02:31.720583	200	165
4155	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:02:32.175195	2025-05-01 22:02:32.177084	200	1
4156	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:02:32.53779	2025-05-01 22:02:32.539042	200	1
4157	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-01 22:02:33.541471	2025-05-01 22:02:33.583305	200	41
4158	\N	/api/restaurant	GET	2025-05-01 22:02:33.218395	2025-05-01 22:02:34.863185	200	1644
4168	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 05:30:57.821016	2025-05-02 05:30:57.82112	200	0
4169	\N	/	GET	2025-05-02 05:48:35.054195	2025-05-02 05:48:35.054199	500	0
4177	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 10:18:35.586407	2025-05-02 10:18:35.586498	200	0
4178	\N	/favicon.ico	GET	2025-05-02 10:30:53.539346	2025-05-02 10:30:53.539349	500	0
4179	\N	/	GET	2025-05-02 13:28:21.375154	2025-05-02 13:28:21.375157	500	0
4180	\N	/.env	GET	2025-05-02 14:22:33.144331	2025-05-02 14:22:33.144333	500	0
4181	\N	/.git/config	GET	2025-05-02 15:02:21.956071	2025-05-02 15:02:21.956073	500	0
4189	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:24:52.202873	2025-05-02 21:24:52.205503	200	2
4190	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:25:00.901963	2025-05-02 21:25:00.903532	200	1
4258	\N	/api/traq/	GET	2025-05-02 21:38:28.970943	2025-05-02 21:38:28.971906	200	0
4259	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:42:22.390056	2025-05-02 21:42:22.40799	200	17
4260	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:42:24.43264	2025-05-02 21:42:24.436007	200	3
4261	\N	/api/traq/	GET	2025-05-02 21:42:45.956804	2025-05-02 21:42:45.957588	200	0
4262	\N	/api/traq/	GET	2025-05-02 21:43:54.547283	2025-05-02 21:43:54.548819	200	1
4263	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:44:24.534663	2025-05-02 21:44:24.538796	200	4
4264	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:52:33.675624	2025-05-02 21:52:33.676831	200	1
4265	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-02 21:52:33.714736	2025-05-02 21:52:33.715861	200	1
4266	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:52:38.117082	2025-05-02 21:52:38.19327	200	76
4267	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:56:56.014544	2025-05-02 21:56:56.01842	200	3
4268	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:56:56.731877	2025-05-02 21:56:56.773968	200	42
4269	\N	/api/restaurant	GET	2025-05-02 21:56:56.903233	2025-05-02 21:56:56.903287	200	0
4270	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:06:23.061958	2025-05-02 22:06:23.066063	200	4
4271	\N	/api/restaurant	GET	2025-05-02 22:06:23.85783	2025-05-02 22:06:23.857896	200	0
4272	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:06:24.019627	2025-05-02 22:06:24.021745	200	2
4273	\N	/api/traq/	GET	2025-05-02 22:07:05.550622	2025-05-02 22:07:05.551756	200	1
4274	\N	/api/restaurant	GET	2025-05-02 22:07:08.99886	2025-05-02 22:07:08.998935	200	0
4275	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:07:10.233433	2025-05-02 22:07:10.236178	200	2
4276	\N	/api/restaurant	GET	2025-05-02 22:07:12.004798	2025-05-02 22:07:12.00487	200	0
4277	\N	/api/restaurant	GET	2025-05-02 22:07:14.384671	2025-05-02 22:07:14.384725	200	0
4278	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:08:13.402969	2025-05-02 22:08:13.422854	200	19
4279	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:08:26.570088	2025-05-02 22:08:26.572124	200	2
4280	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:13:48.299281	2025-05-02 22:13:48.324939	200	25
4281	\N	/api/restaurant	GET	2025-05-02 22:13:49.080264	2025-05-02 22:13:49.080333	200	0
4282	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:13:49.250029	2025-05-02 22:13:49.262295	200	12
4283	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:13:51.660345	2025-05-02 22:13:51.661458	200	1
4284	\N	/api/restaurant	GET	2025-05-02 22:19:01.696759	2025-05-02 22:19:01.69683	200	0
4285	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:19:01.855292	2025-05-02 22:19:01.876833	200	21
4286	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:19:24.886067	2025-05-02 22:19:24.888348	200	2
4287	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:23:46.420875	2025-05-02 22:23:46.448245	200	27
4288	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:23:47.202503	2025-05-02 22:23:47.216724	200	14
4289	\N	/api/restaurant	GET	2025-05-02 22:23:47.401016	2025-05-02 22:23:47.401075	200	0
4290	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:23:51.376159	2025-05-02 22:23:51.38295	200	6
4291	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:24:04.35937	2025-05-02 22:24:04.361306	200	1
4292	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 22:24:07.137755	2025-05-02 22:24:07.14027	200	2
4293	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:24:22.928776	2025-05-02 22:24:22.92995	200	1
4294	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 22:24:25.073493	2025-05-02 22:24:25.078873	200	5
4295	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:25:37.822957	2025-05-02 22:25:37.824648	200	1
4296	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:01.583664	2025-05-02 22:26:01.585207	200	1
4297	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:17.149663	2025-05-02 22:26:17.151623	200	1
4298	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:25.043525	2025-05-02 22:26:25.046087	200	2
4299	\N	/api/restaurant	GET	2025-05-02 22:26:25.804508	2025-05-02 22:26:25.804561	200	0
4300	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:25.980701	2025-05-02 22:26:25.981975	200	1
4301	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:26.77047	2025-05-02 22:26:26.771766	200	1
4302	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:28.35345	2025-05-02 22:26:28.357801	200	4
4303	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:26:45.094978	2025-05-02 22:26:45.096475	200	1
4304	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:10.877408	2025-05-02 22:27:10.87915	200	1
4305	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:13.359378	2025-05-02 22:27:13.382225	200	22
4306	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:28.703796	2025-05-02 22:27:28.705216	200	1
4307	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:42.211244	2025-05-02 22:27:42.214843	200	3
4308	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:44.580161	2025-05-02 22:27:44.5837	200	3
4309	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:51.74354	2025-05-02 22:27:51.744969	200	1
4310	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:58.436502	2025-05-02 22:27:58.438817	200	2
4311	\N	/api/restaurant	GET	2025-05-02 22:27:59.216189	2025-05-02 22:27:59.21625	200	0
4312	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:27:59.406883	2025-05-02 22:27:59.43689	200	30
4313	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:30:07.504801	2025-05-02 22:30:07.506044	200	1
4314	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:31:01.271633	2025-05-02 22:31:01.293288	200	21
4315	\N	/api/restaurant	GET	2025-05-02 22:31:08.08053	2025-05-02 22:31:08.080588	200	0
4185	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:24:48.94624	2025-05-02 21:24:48.947476	200	1
4186	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:24:49.383841	2025-05-02 21:24:49.38738	200	3
4187	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:24:50.450707	2025-05-02 21:24:50.452103	200	1
4188	\N	/api/restaurant	GET	2025-05-02 21:24:50.123499	2025-05-02 21:24:51.28203	200	1158
4191	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-02 21:25:01.97428	2025-05-02 21:25:01.97439	200	0
4192	\N	/api/restaurant	GET	2025-05-02 21:25:04.109759	2025-05-02 21:25:04.109829	200	0
4193	\N	/api/traq/	GET	2025-05-02 21:25:06.038747	2025-05-02 21:25:06.039563	200	0
4194	\N	/api/restaurant	GET	2025-05-02 21:25:16.39543	2025-05-02 21:25:16.395483	200	0
4195	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:25:33.234511	2025-05-02 21:25:33.254436	200	19
4196	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:25:42.24012	2025-05-02 21:25:42.267294	200	27
4197	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:25:46.739651	2025-05-02 21:25:46.740982	200	1
4198	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:26:00.189498	2025-05-02 21:26:00.191569	200	2
4199	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-02 21:26:03.966829	2025-05-02 21:26:03.968018	200	1
4200	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:26:06.454463	2025-05-02 21:26:06.455822	200	1
4201	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-02 21:26:44.226136	2025-05-02 21:26:44.2414	200	15
4202	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-02 21:26:48.351587	2025-05-02 21:26:48.361654	200	10
4203	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:26:50.8289	2025-05-02 21:26:50.830255	200	1
4204	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:28:28.713623	2025-05-02 21:28:28.715357	200	1
4205	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:28:29.732181	2025-05-02 21:28:29.734066	200	1
4206	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:28:57.980663	2025-05-02 21:28:57.982762	200	2
4207	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:28:58.732375	2025-05-02 21:28:58.733626	200	1
4208	\N	/api/restaurant	GET	2025-05-02 21:28:58.901674	2025-05-02 21:28:58.901741	200	0
4209	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:28:59.549251	2025-05-02 21:28:59.550818	200	1
4210	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:29:01.148643	2025-05-02 21:29:01.151493	200	2
4211	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:29:16.041785	2025-05-02 21:29:16.044078	200	2
4212	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:30:17.626222	2025-05-02 21:30:17.630922	200	4
4213	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:30:19.945666	2025-05-02 21:30:19.948481	200	2
4214	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:30:27.922537	2025-05-02 21:30:27.924186	200	1
4215	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:30:35.609149	2025-05-02 21:30:35.610399	200	1
4216	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:00.455153	2025-05-02 21:31:00.457751	200	2
4217	\N	/api/restaurant	GET	2025-05-02 21:31:01.157042	2025-05-02 21:31:01.157104	200	0
4218	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:01.369879	2025-05-02 21:31:01.37603	200	6
4219	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:02.14196	2025-05-02 21:31:02.145697	200	3
4220	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:02.926755	2025-05-02 21:31:02.928114	200	1
4221	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:20.012665	2025-05-02 21:31:20.015139	200	2
4222	\N	/api/restaurant	GET	2025-05-02 21:31:20.742723	2025-05-02 21:31:20.742786	200	0
4223	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:20.941008	2025-05-02 21:31:20.942234	200	1
4224	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:21.56436	2025-05-02 21:31:21.565736	200	1
4225	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:22.949544	2025-05-02 21:31:22.950836	200	1
4226	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:29.838245	2025-05-02 21:31:29.840008	200	1
4227	\N	/api/restaurant	GET	2025-05-02 21:31:30.580827	2025-05-02 21:31:30.580881	200	0
4228	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:30.758736	2025-05-02 21:31:30.760071	200	1
4229	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:32.133727	2025-05-02 21:31:32.144991	200	11
4230	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:31:58.773199	2025-05-02 21:31:58.774911	200	1
4231	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:32:38.455704	2025-05-02 21:32:38.45836	200	2
4232	\N	/api/auth/login	POST	2025-05-02 21:34:22.258618	2025-05-02 21:34:22.328423	200	69
4233	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:22.764046	2025-05-02 21:34:22.765339	200	1
4234	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:23.178687	2025-05-02 21:34:23.198136	200	19
4235	\N	/api/restaurant	GET	2025-05-02 21:34:23.876545	2025-05-02 21:34:23.876615	200	0
4236	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:24.046095	2025-05-02 21:34:24.054854	200	8
4237	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:30.38674	2025-05-02 21:34:30.393321	200	6
4238	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:34.255398	2025-05-02 21:34:34.258541	200	3
4239	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:34:47.977682	2025-05-02 21:34:47.980016	200	2
4240	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:50.931436	2025-05-02 21:34:50.932596	200	1
4241	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:34:52.744621	2025-05-02 21:34:52.745957	200	1
4242	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:35:50.673417	2025-05-02 21:35:50.674917	200	1
4243	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:35:58.094412	2025-05-02 21:35:58.101616	200	7
4244	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:36:02.155848	2025-05-02 21:36:02.157257	200	1
4245	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:36:05.799092	2025-05-02 21:36:05.800367	200	1
4246	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:36:11.023861	2025-05-02 21:36:11.025096	200	1
4247	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:36:44.452237	2025-05-02 21:36:44.462496	200	10
4248	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:37:00.159498	2025-05-02 21:37:00.163598	200	4
4249	\N	/api/restaurant	GET	2025-05-02 21:37:06.686647	2025-05-02 21:37:06.686722	200	0
4250	\N	/api/traq/	GET	2025-05-02 21:37:08.268009	2025-05-02 21:37:08.269025	200	1
4251	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:37:25.205138	2025-05-02 21:37:25.208951	200	3
4252	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:37:26.116845	2025-05-02 21:37:26.118582	200	1
4253	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:37:31.61536	2025-05-02 21:37:31.617994	200	2
4254	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 21:37:34.05816	2025-05-02 21:37:34.069928	200	11
4255	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:37:38.631017	2025-05-02 21:37:38.633501	200	2
4256	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 21:38:01.479787	2025-05-02 21:38:01.483156	200	3
4257	\N	/api/restaurant	GET	2025-05-02 21:38:26.623956	2025-05-02 21:38:28.004774	200	1380
4317	\N	/api/washingmachines	GET	2025-05-02 22:31:49.066925	2025-05-02 22:31:49.83928	200	772
4318	\N	/api/restaurant	GET	2025-05-02 22:31:55.667333	2025-05-02 22:31:55.667399	200	0
4319	\N	/status	GET	2025-05-02 22:31:59.578308	2025-05-02 22:31:59.578328	200	0
4320	\N	/api/statistics/global	GET	2025-05-02 22:31:59.741697	2025-05-02 22:31:59.743731	200	2
4321	\N	/api/statistics/endpoints	GET	2025-05-02 22:31:59.903182	2025-05-02 22:31:59.906558	200	3
4322	\N	/status	GET	2025-05-02 22:32:04.291913	2025-05-02 22:32:04.291917	200	0
4323	\N	/status	GET	2025-05-02 22:32:04.45208	2025-05-02 22:32:04.452082	200	0
4324	\N	/api/statistics/global	GET	2025-05-02 22:32:04.589595	2025-05-02 22:32:04.591951	200	2
4325	\N	/api/statistics/endpoints	GET	2025-05-02 22:32:04.763272	2025-05-02 22:32:04.767979	200	4
4316	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 22:31:46.353303	2025-05-02 22:31:46.353401	304	0
4336	\N	/api/statistics/global	GET	2025-05-02 22:32:19.471979	2025-05-02 22:32:19.473496	200	1
4337	\N	/api/statistics/endpoints	GET	2025-05-02 22:32:19.636395	2025-05-02 22:32:19.639671	200	3
4338	\N	/api/washingmachines	GET	2025-05-02 22:32:19.607798	2025-05-02 22:32:20.345662	200	737
4339	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:33:21.775553	2025-05-02 22:33:21.776995	200	1
4340	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-02 22:33:24.516896	2025-05-02 22:33:24.558006	200	41
4341	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:35:20.735622	2025-05-02 22:35:20.737145	200	1
4342	\N	/api/restaurant	GET	2025-05-02 22:35:21.4349	2025-05-02 22:35:21.434974	200	0
4343	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:35:21.612032	2025-05-02 22:35:21.614146	200	2
4345	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:37:05.347156	2025-05-02 22:37:05.376546	200	29
4346	\N	/api/washingmachines	GET	2025-05-02 22:37:20.774695	2025-05-02 22:37:21.529108	200	754
4347	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:39:02.794772	2025-05-02 22:39:02.796207	200	1
4348	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:39:03.538711	2025-05-02 22:39:03.539898	200	1
4349	\N	/api/restaurant	GET	2025-05-02 22:39:03.749754	2025-05-02 22:39:03.749842	200	0
4350	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:39:50.362756	2025-05-02 22:39:50.364138	200	1
4351	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:39:51.568256	2025-05-02 22:39:51.569342	200	1
4352	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:40:21.494839	2025-05-02 22:40:21.496815	200	1
4353	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:40:22.319665	2025-05-02 22:40:22.32548	200	5
4354	\N	/api/restaurant	GET	2025-05-02 22:40:22.489169	2025-05-02 22:40:22.489232	200	0
4355	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:40:56.119364	2025-05-02 22:40:56.121695	200	2
4356	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:40:58.481105	2025-05-02 22:40:58.48245	200	1
4357	\N	/api/washingmachines	GET	2025-05-02 22:42:20.88819	2025-05-02 22:42:21.573735	200	685
4358	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:43:45.222605	2025-05-02 22:43:45.244009	200	21
4359	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:43:45.916922	2025-05-02 22:43:45.918075	200	1
4360	\N	/api/restaurant	GET	2025-05-02 22:43:46.086346	2025-05-02 22:43:46.086424	200	0
4361	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:43:51.551717	2025-05-02 22:43:51.553409	200	1
4362	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:43:52.839763	2025-05-02 22:43:52.841095	200	1
4363	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:44:36.765943	2025-05-02 22:44:36.767446	200	1
4364	\N	/api/restaurant	GET	2025-05-02 22:44:37.55222	2025-05-02 22:44:37.552277	200	0
4365	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:44:37.721713	2025-05-02 22:44:37.72296	200	1
4366	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:44:39.619243	2025-05-02 22:44:39.62045	200	1
4367	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:44:40.434269	2025-05-02 22:44:40.43592	200	1
4368	\N	/api/washingmachines	GET	2025-05-02 22:47:20.771843	2025-05-02 22:47:21.470185	200	698
4369	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:51:06.560976	2025-05-02 22:51:06.562749	200	1
4370	\N	/api/restaurant	GET	2025-05-02 22:51:07.320835	2025-05-02 22:51:07.320899	200	0
4371	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:51:07.506581	2025-05-02 22:51:07.507753	200	1
4372	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:51:12.90461	2025-05-02 22:51:12.905985	200	1
4373	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 22:51:15.773546	2025-05-02 22:51:15.774927	200	1
4374	\N	/api/washingmachines	GET	2025-05-02 22:52:20.751231	2025-05-02 22:52:21.439636	200	688
4375	\N	/api/washingmachines	GET	2025-05-02 22:57:46.785569	2025-05-02 22:57:47.465058	200	679
4376	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:00:15.333711	2025-05-02 23:00:15.343803	200	10
4377	\N	/api/restaurant	GET	2025-05-02 23:00:16.070048	2025-05-02 23:00:16.07011	200	0
4378	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:00:16.236759	2025-05-02 23:00:16.255425	200	18
4379	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:00:17.857687	2025-05-02 23:00:17.859157	200	1
4380	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:00:19.086812	2025-05-02 23:00:19.089111	200	2
4381	\N	/api/washingmachines	GET	2025-05-02 23:02:46.777804	2025-05-02 23:02:47.454192	200	676
4382	\N	/api/washingmachines	GET	2025-05-02 23:07:46.634612	2025-05-02 23:07:47.313271	200	678
4383	\N	/api/washingmachines	GET	2025-05-02 23:12:46.546157	2025-05-02 23:12:47.217496	200	671
4384	\N	/api/washingmachines	GET	2025-05-02 23:17:46.532892	2025-05-02 23:17:47.226242	200	693
4385	\N	/api/washingmachines	GET	2025-05-02 23:22:46.755758	2025-05-02 23:22:47.433791	200	678
4386	\N	/api/auth/login	POST	2025-05-02 23:24:24.010129	2025-05-02 23:24:24.080121	200	69
4387	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:24:24.281602	2025-05-02 23:24:24.297958	200	16
4388	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:24:25.122917	2025-05-02 23:24:25.12429	200	1
4389	\N	/api/restaurant	GET	2025-05-02 23:24:26.771309	2025-05-02 23:24:26.771383	200	0
4390	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:24:27.036336	2025-05-02 23:24:27.037723	200	1
4391	\N	/api/traq/	GET	2025-05-02 23:24:38.74257	2025-05-02 23:24:38.743726	200	1
4392	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:25:30.83944	2025-05-02 23:25:30.840783	200	1
4393	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:25:56.533372	2025-05-02 23:25:56.534857	200	1
4394	\N	/api/washingmachines	GET	2025-05-02 23:27:46.806443	2025-05-02 23:27:47.479324	200	672
4420	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 23:49:10.850066	2025-05-02 23:49:10.850162	200	0
4421	\N	/status	GET	2025-05-02 23:49:12.407292	2025-05-02 23:49:12.407295	200	0
4422	\N	/status	GET	2025-05-02 23:49:12.58985	2025-05-02 23:49:12.589852	200	0
4423	\N	/api/statistics/global	GET	2025-05-02 23:49:12.725186	2025-05-02 23:49:12.726529	200	1
4424	\N	/api/statistics/global	GET	2025-05-02 23:49:12.885816	2025-05-02 23:49:12.887165	200	1
4429	\N	/api/statistics/global	GET	2025-05-02 23:49:17.22811	2025-05-02 23:49:17.234756	200	6
4430	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:17.432332	2025-05-02 23:49:17.458054	200	25
4431	\N	/status	GET	2025-05-02 23:49:22.044509	2025-05-02 23:49:22.04451	200	0
4432	\N	/status	GET	2025-05-02 23:49:22.212512	2025-05-02 23:49:22.212514	200	0
4437	\N	/api/statistics/global	GET	2025-05-02 23:49:27.206564	2025-05-02 23:49:27.208181	200	1
4438	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:27.40832	2025-05-02 23:49:27.413867	200	5
4439	\N	/status	GET	2025-05-02 23:49:32.043031	2025-05-02 23:49:32.043033	200	0
4441	\N	/api/statistics/global	GET	2025-05-02 23:49:32.232609	2025-05-02 23:49:32.236623	200	4
4442	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:32.434271	2025-05-02 23:49:32.453093	200	18
4443	\N	/status	GET	2025-05-02 23:49:37.042774	2025-05-02 23:49:37.042775	200	0
4444	\N	/status	GET	2025-05-02 23:49:37.208884	2025-05-02 23:49:37.208885	200	0
4449	\N	/api/statistics/global	GET	2025-05-02 23:49:42.198818	2025-05-02 23:49:42.200456	200	1
4450	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:42.403978	2025-05-02 23:49:42.409458	200	5
4451	\N	/status	GET	2025-05-02 23:49:47.03907	2025-05-02 23:49:47.039071	200	0
4452	\N	/status	GET	2025-05-02 23:49:47.205408	2025-05-02 23:49:47.205409	200	0
4453	\N	/api/statistics/global	GET	2025-05-02 23:49:47.210417	2025-05-02 23:49:47.213833	200	3
4454	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:47.402646	2025-05-02 23:49:47.40685	200	4
4455	\N	/status	GET	2025-05-02 23:49:52.041538	2025-05-02 23:49:52.04154	200	0
4456	\N	/status	GET	2025-05-02 23:49:52.227668	2025-05-02 23:49:52.227669	200	0
4461	\N	/api/statistics/global	GET	2025-05-02 23:49:57.211688	2025-05-02 23:49:57.227321	200	15
4462	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:57.455959	2025-05-02 23:49:57.460909	200	4
4463	\N	/status	GET	2025-05-02 23:50:02.0391	2025-05-02 23:50:02.039102	200	0
4464	\N	/status	GET	2025-05-02 23:50:02.199676	2025-05-02 23:50:02.199677	200	0
4469	\N	/api/statistics/global	GET	2025-05-02 23:50:07.199339	2025-05-02 23:50:07.208718	200	9
4470	\N	/api/statistics/endpoints	GET	2025-05-02 23:50:07.40392	2025-05-02 23:50:07.407461	200	3
4471	\N	/status	GET	2025-05-02 23:50:12.037495	2025-05-02 23:50:12.037496	200	0
4326	\N	/status	GET	2025-05-02 22:32:09.281456	2025-05-02 22:32:09.281457	200	0
4327	\N	/status	GET	2025-05-02 22:32:09.470632	2025-05-02 22:32:09.470633	200	0
4328	\N	/api/statistics/global	GET	2025-05-02 22:32:09.606449	2025-05-02 22:32:09.607892	200	1
4329	\N	/api/statistics/endpoints	GET	2025-05-02 22:32:09.782435	2025-05-02 22:32:09.785545	200	3
4330	\N	/status	GET	2025-05-02 22:32:14.551053	2025-05-02 22:32:14.551075	200	0
4331	\N	/status	GET	2025-05-02 22:32:14.71445	2025-05-02 22:32:14.714451	200	0
4332	\N	/api/statistics/global	GET	2025-05-02 22:32:14.855329	2025-05-02 22:32:14.857611	200	2
4333	\N	/api/statistics/endpoints	GET	2025-05-02 22:32:15.022331	2025-05-02 22:32:15.033593	200	11
4334	\N	/status	GET	2025-05-02 22:32:19.281248	2025-05-02 22:32:19.28125	200	0
4335	\N	/status	GET	2025-05-02 22:32:19.47165	2025-05-02 22:32:19.471651	200	0
4344	\N	/api/restaurant	GET	2025-05-02 22:37:05.347187	2025-05-02 22:37:05.347261	200	0
4395	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 23:28:50.059731	2025-05-02 23:28:50.059824	200	0
4396	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:31:46.75117	2025-05-02 23:31:46.780149	200	28
4397	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:31:47.3786	2025-05-02 23:31:47.380764	200	2
4398	\N	/api/restaurant	GET	2025-05-02 23:31:47.54956	2025-05-02 23:31:47.549623	200	0
4399	\N	/api/washingmachines	GET	2025-05-02 23:32:56.947494	2025-05-02 23:32:57.698706	200	751
4400	\N	/api/auth/login	POST	2025-05-02 23:36:56.412736	2025-05-02 23:36:56.481578	200	68
4401	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:36:56.777124	2025-05-02 23:36:56.778478	200	1
4402	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:36:56.947751	2025-05-02 23:36:56.948935	200	1
4403	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:36:59.543372	2025-05-02 23:36:59.545355	200	1
4404	\N	/api/restaurant	GET	2025-05-02 23:37:00.284181	2025-05-02 23:37:00.284243	200	0
4405	\N	/api/traq/	GET	2025-05-02 23:37:39.380074	2025-05-02 23:37:39.381495	200	1
4406	\N	/api/washingmachines	GET	2025-05-02 23:37:46.548266	2025-05-02 23:37:47.145735	200	597
4407	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:38:45.237828	2025-05-02 23:38:45.274686	200	36
4408	\N	/api/auth/login	POST	2025-05-02 23:42:53.591755	2025-05-02 23:42:53.660298	200	68
4409	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:42:54.606824	2025-05-02 23:42:54.608092	200	1
4410	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:42:54.936329	2025-05-02 23:42:54.937785	200	1
4411	\N	/api/restaurant	GET	2025-05-02 23:42:57.602409	2025-05-02 23:42:57.602466	200	0
4412	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:42:57.756656	2025-05-02 23:42:57.757924	200	1
4413	\N	/api/auth/login	POST	2025-05-02 23:43:27.411703	2025-05-02 23:43:27.483709	200	72
4414	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:43:27.642208	2025-05-02 23:43:27.643573	200	1
4415	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:43:27.833632	2025-05-02 23:43:27.835	200	1
4416	\N	/api/restaurant	GET	2025-05-02 23:43:29.549382	2025-05-02 23:43:29.54944	200	0
4417	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:43:30.371733	2025-05-02 23:43:30.41588	200	44
4418	\N	/api/traq/	GET	2025-05-02 23:43:36.395045	2025-05-02 23:43:36.396207	200	1
4419	test@imt-atlantique.net	/api/newf/me	GET	2025-05-02 23:45:30.919287	2025-05-02 23:45:30.920709	200	1
4425	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:12.886765	2025-05-02 23:49:12.889885	200	3
4426	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:13.057751	2025-05-02 23:49:13.067777	200	10
4427	\N	/status	GET	2025-05-02 23:49:17.061819	2025-05-02 23:49:17.061825	200	0
4428	\N	/status	GET	2025-05-02 23:49:17.227087	2025-05-02 23:49:17.227089	200	0
4433	\N	/api/statistics/global	GET	2025-05-02 23:49:22.213402	2025-05-02 23:49:22.215174	200	1
4434	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:22.415077	2025-05-02 23:49:22.447983	200	32
4435	\N	/status	GET	2025-05-02 23:49:27.042864	2025-05-02 23:49:27.042865	200	0
4436	\N	/status	GET	2025-05-02 23:49:27.205576	2025-05-02 23:49:27.205577	200	0
4440	\N	/status	GET	2025-05-02 23:49:32.232601	2025-05-02 23:49:32.232603	200	0
4445	\N	/api/statistics/global	GET	2025-05-02 23:49:37.209774	2025-05-02 23:49:37.214858	200	5
4446	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:37.420866	2025-05-02 23:49:37.447437	200	26
4447	\N	/status	GET	2025-05-02 23:49:42.03712	2025-05-02 23:49:42.037121	200	0
4448	\N	/status	GET	2025-05-02 23:49:42.19774	2025-05-02 23:49:42.197742	200	0
4457	\N	/api/statistics/global	GET	2025-05-02 23:49:52.22878	2025-05-02 23:49:52.230546	200	1
4458	\N	/api/statistics/endpoints	GET	2025-05-02 23:49:52.43101	2025-05-02 23:49:52.450081	200	19
4459	\N	/status	GET	2025-05-02 23:49:57.042098	2025-05-02 23:49:57.042099	200	0
4460	\N	/status	GET	2025-05-02 23:49:57.210476	2025-05-02 23:49:57.210478	200	0
4465	\N	/api/statistics/global	GET	2025-05-02 23:50:02.200506	2025-05-02 23:50:02.203008	200	2
4466	\N	/api/statistics/endpoints	GET	2025-05-02 23:50:02.394118	2025-05-02 23:50:02.415791	200	21
4467	\N	/status	GET	2025-05-02 23:50:07.037654	2025-05-02 23:50:07.037655	200	0
4468	\N	/status	GET	2025-05-02 23:50:07.198564	2025-05-02 23:50:07.198565	200	0
4473	\N	/api/statistics/global	GET	2025-05-02 23:50:12.209583	2025-05-02 23:50:12.219296	200	9
4474	\N	/api/statistics/endpoints	GET	2025-05-02 23:50:12.421209	2025-05-02 23:50:12.424737	200	3
4843	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 00:53:53.719054	2025-05-09 00:53:53.719172	200	0
4844	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 07:40:00.753031	2025-05-09 07:40:00.812128	200	59
4845	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 07:40:01.29896	2025-05-09 07:40:01.300825	200	1
4846	\N	/api/restaurant	GET	2025-05-09 07:40:01.485282	2025-05-09 07:40:01.485358	200	0
4847	\N	/api/restaurant	GET	2025-05-09 07:40:05.667073	2025-05-09 07:40:05.667145	200	0
4848	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 07:40:09.039229	2025-05-09 07:40:09.049486	200	10
4849	\N	/api/restaurant	GET	2025-05-09 07:40:11.054692	2025-05-09 07:40:11.054742	200	0
4851	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 07:51:35.535559	2025-05-09 07:51:35.535611	200	0
4854	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 12:07:04.416938	2025-05-09 12:07:04.417039	304	0
4859	\N	/status	GET	2025-05-09 12:07:17.615437	2025-05-09 12:07:17.615455	200	0
4864	\N	/api/statistics/global	GET	2025-05-09 12:07:22.651825	2025-05-09 12:07:22.653705	200	1
4865	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:24.547624	2025-05-09 12:07:24.552509	200	4
4866	\N	/status	GET	2025-05-09 12:07:25.158635	2025-05-09 12:07:25.158638	200	0
4867	\N	/status	GET	2025-05-09 12:07:27.273245	2025-05-09 12:07:27.273246	200	0
4868	\N	/api/statistics/global	GET	2025-05-09 12:07:27.418814	2025-05-09 12:07:27.420767	200	1
4869	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:29.320415	2025-05-09 12:07:29.335033	200	14
4870	\N	/status	GET	2025-05-09 12:07:30.158775	2025-05-09 12:07:30.158781	200	0
4872	\N	/api/statistics/global	GET	2025-05-09 12:07:32.648129	2025-05-09 12:07:32.649636	200	1
4873	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:34.569035	2025-05-09 12:07:34.573956	200	4
4874	\N	/status	GET	2025-05-09 12:07:35.164058	2025-05-09 12:07:35.164062	200	0
4876	\N	/api/statistics/global	GET	2025-05-09 12:07:37.256561	2025-05-09 12:07:37.260948	200	4
4880	\N	/api/statistics/global	GET	2025-05-09 12:07:42.693313	2025-05-09 12:07:42.695621	200	2
4881	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:44.619507	2025-05-09 12:07:44.623787	200	4
4882	\N	/status	GET	2025-05-09 12:07:45.156934	2025-05-09 12:07:45.156936	200	0
4884	\N	/api/statistics/global	GET	2025-05-09 12:07:47.25176	2025-05-09 12:07:47.288224	200	36
4885	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:49.310084	2025-05-09 12:07:49.378543	200	68
4886	\N	/status	GET	2025-05-09 12:07:50.158767	2025-05-09 12:07:50.158768	200	0
4887	\N	/status	GET	2025-05-09 12:07:52.259909	2025-05-09 12:07:52.25991	200	0
4892	\N	/api/statistics/global	GET	2025-05-09 12:07:57.267229	2025-05-09 12:07:57.269361	200	2
4893	\N	/api/statistics/endpoints	GET	2025-05-09 12:07:59.322146	2025-05-09 12:07:59.376521	200	54
4894	\N	/status	GET	2025-05-09 12:08:00.158708	2025-05-09 12:08:00.158709	200	0
4895	\N	/api/statistics/global	GET	2025-05-09 12:08:00.250757	2025-05-09 12:08:00.25802	200	7
4897	\N	/api/statistics/global	GET	2025-05-09 12:08:02.242741	2025-05-09 12:08:02.244972	200	2
4899	\N	/api/statistics/endpoints	GET	2025-05-09 12:08:04.301607	2025-05-09 12:08:04.306508	200	4
4900	\N	/status	GET	2025-05-09 12:08:05.158952	2025-05-09 12:08:05.158953	200	0
4901	\N	/status	GET	2025-05-09 12:08:07.258201	2025-05-09 12:08:07.258202	200	0
4905	\N	/status	GET	2025-05-09 12:08:12.282625	2025-05-09 12:08:12.282627	200	0
4909	\N	/status	GET	2025-05-09 12:08:17.250714	2025-05-09 12:08:17.250715	200	0
4472	\N	/status	GET	2025-05-02 23:50:12.206862	2025-05-02 23:50:12.206863	200	0
4475	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 23:57:43.696107	2025-05-02 23:57:43.696213	200	0
4476	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-02 23:57:43.771541	2025-05-02 23:57:43.771619	200	0
4477	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 00:09:33.194514	2025-05-03 00:09:33.208526	200	14
4478	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 00:09:33.779457	2025-05-03 00:09:33.780877	200	1
4479	\N	/api/restaurant	GET	2025-05-03 00:09:33.996743	2025-05-03 00:09:34.918228	200	921
4480	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 00:09:36.70977	2025-05-03 00:09:36.711264	200	1
4481	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 00:09:39.219091	2025-05-03 00:09:39.221389	200	2
4482	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-03 01:09:51.818251	2025-05-03 01:09:51.818351	200	0
4483	\N	/.git/config	GET	2025-05-03 03:42:04.981605	2025-05-03 03:42:04.981608	500	0
4484	\N	/.git-credentials	GET	2025-05-03 03:42:05.138794	2025-05-03 03:42:05.138798	500	0
4485	\N	/	GET	2025-05-03 06:36:52.94884	2025-05-03 06:36:52.948842	500	0
4486	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-03 06:37:29.900285	2025-05-03 06:37:29.900373	200	0
4487	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-03 06:37:30.112682	2025-05-03 06:37:30.112737	200	0
4488	\N	/	GET	2025-05-03 06:46:56.947229	2025-05-03 06:46:56.947231	500	0
4489	\N	/	GET	2025-05-03 09:01:02.335849	2025-05-03 09:01:02.335853	500	0
4490	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-03 10:02:27.225394	2025-05-03 10:02:27.225486	200	0
4491	\N	/.git/config	GET	2025-05-03 11:31:17.166495	2025-05-03 11:31:17.166498	500	0
4492	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:15:01.047718	2025-05-03 13:15:01.053507	200	5
4493	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:15:01.441755	2025-05-03 13:15:01.456654	200	14
4494	\N	/api/restaurant	GET	2025-05-03 13:15:01.583616	2025-05-03 13:15:01.58371	200	0
4495	\N	/api/auth/login	POST	2025-05-03 13:27:04.757533	2025-05-03 13:27:04.75897	401	1
4496	\N	/api/auth/login	POST	2025-05-03 13:27:22.826371	2025-05-03 13:27:22.827625	401	1
4497	\N	/api/auth/login	POST	2025-05-03 13:31:11.189724	2025-05-03 13:31:11.196749	401	7
4498	\N	/api/auth/reset-password	POST	2025-05-03 13:31:27.989513	2025-05-03 13:31:27.989527	500	0
4499	\N	/api/auth/change-password	POST	2025-05-03 13:43:11.688414	2025-05-03 13:43:11.688455	500	0
4500	\N	/api/auth/verification-code	POST	2025-05-03 13:43:29.116438	2025-05-03 13:43:29.145672	200	29
4501	\N	/api/auth/change-password	POST	2025-05-03 13:43:55.169288	2025-05-03 13:43:55.169302	500	0
4502	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-03 13:49:38.57895	2025-05-03 13:49:38.579065	200	0
4503	\N	/api/auth/change-password	POST	2025-05-03 15:56:40.565724	2025-05-03 15:56:40.566206	500	0
4504	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:58:17.924569	2025-05-03 13:58:17.925906	200	1
4505	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:58:18.184529	2025-05-03 13:58:18.186004	200	1
4506	\N	/api/restaurant	GET	2025-05-03 13:58:18.324423	2025-05-03 13:58:18.324479	200	0
4507	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:58:39.161222	2025-05-03 13:58:39.163374	200	2
4508	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 13:58:43.379426	2025-05-03 13:58:43.407474	200	28
4509	\N	/api/auth/change-password	PATCH	2025-05-03 16:01:01.14529	2025-05-03 16:01:01.574461	400	429
4510	\N	/api/auth/verification-code	POST	2025-05-03 16:01:54.117149	2025-05-03 16:01:55.483388	200	1366
4511	\N	/api/auth/change-password	PATCH	2025-05-03 16:02:15.111899	2025-05-03 16:02:15.88114	200	769
4512	\N	/api/auth/login	POST	2025-05-03 14:02:51.710569	2025-05-03 14:02:51.778314	401	67
4513	\N	/api/auth/login	POST	2025-05-03 14:03:12.815125	2025-05-03 14:03:12.882116	401	66
4514	\N	/api/auth/login	POST	2025-05-03 14:03:42.560403	2025-05-03 14:03:42.628094	200	67
4515	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:03:43.033304	2025-05-03 14:03:43.042088	200	8
4516	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:03:43.494275	2025-05-03 14:03:43.495507	200	1
4517	\N	/api/restaurant	GET	2025-05-03 14:03:44.194614	2025-05-03 14:03:44.194687	200	0
4518	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:03:44.375776	2025-05-03 14:03:44.404628	200	28
4519	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-03 14:03:45.70114	2025-05-03 14:03:45.703303	200	2
4520	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:03:53.644502	2025-05-03 14:03:53.645802	200	1
4521	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:03:54.816453	2025-05-03 14:03:54.817823	200	1
4522	\N	/api/auth/login	POST	2025-05-03 14:04:19.010946	2025-05-03 14:04:19.08352	401	72
4523	\N	/api/auth/login	POST	2025-05-03 14:08:49.408639	2025-05-03 14:08:49.477465	401	68
4524	\N	/api/auth/login	POST	2025-05-03 14:09:52.794478	2025-05-03 14:09:52.861941	401	67
4525	\N	/	GET	2025-05-03 14:10:45.926923	2025-05-03 14:10:45.926927	500	0
4526	\N	/api/auth/login	POST	2025-05-03 14:13:25.508812	2025-05-03 14:13:25.577071	401	68
4527	\N	/api/auth/login	POST	2025-05-03 14:13:34.330463	2025-05-03 14:13:34.426112	401	95
4528	\N	/api/auth/login	POST	2025-05-03 14:13:35.528824	2025-05-03 14:13:35.596669	401	67
4529	\N	/api/auth/login	POST	2025-05-03 14:13:41.011198	2025-05-03 14:13:41.01122	429	0
4530	\N	/api/auth/login	POST	2025-05-03 14:13:42.033392	2025-05-03 14:13:42.033416	429	0
4531	\N	/api/auth/login	POST	2025-05-03 14:13:47.802168	2025-05-03 14:13:47.802188	429	0
4532	\N	/api/auth/login	POST	2025-05-03 14:13:48.925195	2025-05-03 14:13:48.925218	429	0
4533	\N	/api/auth/login	POST	2025-05-03 14:13:49.753588	2025-05-03 14:13:49.753608	429	0
4534	\N	/api/auth/login	POST	2025-05-03 14:13:52.16821	2025-05-03 14:13:52.168229	429	0
4535	\N	/api/auth/login	POST	2025-05-03 14:13:53.43116	2025-05-03 14:13:53.431181	429	0
4536	\N	/api/auth/login	POST	2025-05-03 14:13:58.099034	2025-05-03 14:13:58.09906	429	0
4537	\N	/api/auth/login	POST	2025-05-03 14:13:59.159214	2025-05-03 14:13:59.159246	429	0
4538	\N	/api/auth/login	POST	2025-05-03 14:14:00.39543	2025-05-03 14:14:00.395464	429	0
4539	\N	/api/auth/login	POST	2025-05-03 14:14:10.029099	2025-05-03 14:14:10.029123	429	0
4540	\N	/api/auth/login	POST	2025-05-03 14:14:26.723954	2025-05-03 14:14:26.724012	429	0
4541	\N	/api/auth/login	POST	2025-05-03 14:14:27.704604	2025-05-03 14:14:27.704647	429	0
4542	\N	/api/auth/login	POST	2025-05-03 14:14:28.798389	2025-05-03 14:14:28.798436	429	0
4543	\N	/api/auth/login	POST	2025-05-03 14:16:15.266268	2025-05-03 14:16:15.266327	429	0
4544	\N	/api/auth/login	POST	2025-05-03 14:16:29.978823	2025-05-03 14:16:29.978851	429	0
4545	\N	/api/auth/login	POST	2025-05-03 14:16:31.132534	2025-05-03 14:16:31.132556	429	0
4546	\N	/api/auth/login	POST	2025-05-03 14:16:32.069185	2025-05-03 14:16:32.069208	429	0
4547	\N	/api/auth/login	POST	2025-05-03 14:16:32.985369	2025-05-03 14:16:32.985408	429	0
4548	\N	/api/auth/login	POST	2025-05-03 14:16:35.167263	2025-05-03 14:16:35.167282	429	0
4549	\N	/api/auth/login	POST	2025-05-03 14:16:39.107016	2025-05-03 14:16:39.107035	429	0
4550	\N	/api/auth/login	POST	2025-05-03 14:16:40.446888	2025-05-03 14:16:40.446912	429	0
4551	\N	/api/auth/login	POST	2025-05-03 14:16:43.94472	2025-05-03 14:16:43.944763	429	0
4552	\N	/api/auth/register	POST	2025-05-03 14:18:11.211925	2025-05-03 14:18:11.381992	409	170
4553	\N	/api/auth/login	POST	2025-05-03 14:18:35.729879	2025-05-03 14:18:35.798114	401	68
4554	\N	/api/auth/reset-password	POST	2025-05-03 14:19:45.713564	2025-05-03 14:19:45.71358	500	0
4555	\N	/api/auth/reset-password	POST	2025-05-03 14:28:52.116691	2025-05-03 14:28:52.116809	500	0
4556	\N	/api/auth/login	POST	2025-05-03 14:33:13.646193	2025-05-03 14:33:13.715014	401	68
4557	\N	/api/auth/login	POST	2025-05-03 14:37:41.271304	2025-05-03 14:37:41.352873	401	81
4558	\N	/api/auth/login	POST	2025-05-03 14:37:53.011495	2025-05-03 14:37:53.083707	401	72
4559	\N	/api/auth/login	POST	2025-05-03 14:37:55.469694	2025-05-03 14:37:55.53955	401	69
4560	\N	/api/auth/login	POST	2025-05-03 14:37:56.028994	2025-05-03 14:37:56.029027	429	0
4561	\N	/api/auth/login	POST	2025-05-03 14:44:23.037148	2025-05-03 14:44:23.112381	401	75
4562	\N	/api/auth/verification-code	POST	2025-05-03 14:44:26.181467	2025-05-03 14:44:26.192977	200	11
4563	\N	/api/auth/change-password	PATCH	2025-05-03 14:44:50.491786	2025-05-03 14:44:50.598173	200	106
4564	\N	/api/auth/verification-code	POST	2025-05-03 14:45:00.303301	2025-05-03 14:45:00.303335	429	0
4565	\N	/api/auth/change-password	PATCH	2025-05-03 14:45:52.052373	2025-05-03 14:45:52.052416	429	0
4566	\N	/api/auth/login	POST	2025-05-03 14:45:57.085262	2025-05-03 14:45:57.085292	429	0
4567	\N	/api/auth/login	POST	2025-05-03 14:46:03.826023	2025-05-03 14:46:03.826046	429	0
4568	\N	/api/auth/verification-code	POST	2025-05-03 14:46:07.148317	2025-05-03 14:46:07.148345	429	0
4569	\N	/api/auth/change-password	PATCH	2025-05-03 14:46:33.305049	2025-05-03 14:46:33.30509	429	0
4570	\N	/api/auth/login	POST	2025-05-03 14:46:49.070749	2025-05-03 14:46:49.070778	429	0
4571	\N	/api/auth/login	POST	2025-05-03 14:50:34.48549	2025-05-03 14:50:34.552917	401	67
4572	\N	/api/auth/verification-code	POST	2025-05-03 14:50:41.236841	2025-05-03 14:50:41.305542	200	68
4573	\N	/api/auth/change-password	PATCH	2025-05-03 14:51:09.990537	2025-05-03 14:51:10.070381	200	79
4574	\N	/api/auth/login	POST	2025-05-03 14:51:18.902765	2025-05-03 14:51:18.902788	429	0
4575	\N	/api/auth/verification-code	POST	2025-05-03 14:53:01.039789	2025-05-03 14:53:01.046828	200	7
4576	\N	/api/auth/verification-code	POST	2025-05-03 14:54:23.559809	2025-05-03 14:54:23.566435	200	6
4577	\N	/api/auth/change-password	PATCH	2025-05-03 14:54:47.576157	2025-05-03 14:54:47.645769	200	69
4578	\N	/api/auth/login	POST	2025-05-03 14:54:50.598598	2025-05-03 14:54:50.667281	200	68
4579	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:54:51.035865	2025-05-03 14:54:51.037729	200	1
4580	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:54:51.473508	2025-05-03 14:54:51.475608	200	2
4581	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:54:52.160233	2025-05-03 14:54:52.162669	200	2
4582	\N	/api/restaurant	GET	2025-05-03 14:54:52.341353	2025-05-03 14:54:52.83977	200	498
4583	\N	/api/restaurant	GET	2025-05-03 14:55:33.144143	2025-05-03 14:55:33.144234	200	0
4584	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 14:55:33.144138	2025-05-03 14:55:33.156547	200	12
4585	test@imt-atlantique.net	/api/newf/me	GET	2025-05-03 15:02:07.205467	2025-05-03 15:02:07.21779	200	12
4586	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-03 15:24:42.126053	2025-05-03 15:24:42.138893	200	12
4587	\N	/.git/config	GET	2025-05-03 17:14:27.687733	2025-05-03 17:14:27.687765	500	0
4588	\N	/	GET	2025-05-03 19:40:10.765286	2025-05-03 19:40:10.76529	500	0
4589	\N	/.git/config	GET	2025-05-04 00:21:56.132291	2025-05-04 00:21:56.132294	500	0
4590	\N	/.git-credentials	GET	2025-05-04 00:21:56.282666	2025-05-04 00:21:56.28267	500	0
4591	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 01:10:57.563176	2025-05-04 01:10:57.567922	200	4
4592	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 06:43:43.911929	2025-05-04 06:43:43.912039	200	0
4593	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 06:43:44.183564	2025-05-04 06:43:44.183624	200	0
4594	\N	/.git/config	GET	2025-05-04 08:25:31.202631	2025-05-04 08:25:31.202634	500	0
4595	\N	/.env	GET	2025-05-04 08:25:31.37527	2025-05-04 08:25:31.375274	500	0
4596	\N	/.env.prod	GET	2025-05-04 08:25:31.546096	2025-05-04 08:25:31.546101	500	0
4597	\N	/phpinfo	GET	2025-05-04 08:25:31.702839	2025-05-04 08:25:31.702843	500	0
4598	\N	/info.php	GET	2025-05-04 08:25:31.859307	2025-05-04 08:25:31.859311	500	0
4599	\N	/_profler/phpinfo	GET	2025-05-04 08:25:32.016116	2025-05-04 08:25:32.01612	500	0
4600	\N	/man/.	GET	2025-05-04 08:25:32.17694	2025-05-04 08:25:32.176946	500	0
4601	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 09:05:03.307591	2025-05-04 09:05:03.307693	200	0
4602	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 13:17:47.474994	2025-05-04 13:17:47.475099	200	0
4603	\N	/	GET	2025-05-04 13:51:26.095441	2025-05-04 13:51:26.095444	500	0
4604	\N	//xmlrpc.php	GET	2025-05-04 13:51:26.329793	2025-05-04 13:51:26.329798	500	0
4605	\N	//blog/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:26.53295	2025-05-04 13:51:26.532953	500	0
4606	\N	//wordpress/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:26.736668	2025-05-04 13:51:26.736671	500	0
4607	\N	//wp/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:26.828596	2025-05-04 13:51:26.828601	500	0
4608	\N	//2018/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:27.047675	2025-05-04 13:51:27.047679	500	0
4609	\N	//shop/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:27.257374	2025-05-04 13:51:27.257377	500	0
4610	\N	//test/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:27.361052	2025-05-04 13:51:27.361059	500	0
4611	\N	//wp2/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:27.45133	2025-05-04 13:51:27.451334	500	0
4612	\N	//cms/wp-includes/wlwmanifest.xml	GET	2025-05-04 13:51:27.549699	2025-05-04 13:51:27.549704	500	0
4613	\N	/	GET	2025-05-04 16:31:50.792966	2025-05-04 16:31:50.79297	500	0
4614	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-04 21:19:09.399144	2025-05-04 21:19:09.399242	200	0
4615	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-05 07:21:59.655053	2025-05-05 07:21:59.655171	200	0
4616	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-05 07:21:59.994825	2025-05-05 07:21:59.994909	200	0
4617	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-05 10:20:42.221629	2025-05-05 10:20:42.221721	200	0
4618	\N	/sftp-config.json	GET	2025-05-05 11:49:45.518322	2025-05-05 11:49:45.518326	500	0
4619	\N	/.vscode/sftp.json	GET	2025-05-05 11:49:51.475088	2025-05-05 11:49:51.475091	500	0
4620	\N	/settings.json	GET	2025-05-05 12:50:05.547653	2025-05-05 12:50:05.547656	500	0
4621	\N	/sendgrid.env	GET	2025-05-05 12:50:05.58476	2025-05-05 12:50:05.584763	500	0
4622	\N	/.env	GET	2025-05-05 12:50:05.606573	2025-05-05 12:50:05.606576	500	0
4623	\N	/info.php	GET	2025-05-05 12:50:05.741773	2025-05-05 12:50:05.741776	500	0
4624	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-05 14:43:10.860657	2025-05-05 14:43:10.860773	200	0
4625	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-06 01:13:57.880176	2025-05-06 01:13:57.880276	200	0
4626	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-06 07:22:26.005578	2025-05-06 07:22:26.005707	200	0
4627	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-06 07:22:26.120269	2025-05-06 07:22:26.12033	200	0
4628	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-06 10:23:23.839725	2025-05-06 10:23:23.839854	200	0
4629	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-06 10:31:27.139219	2025-05-06 10:31:27.139325	200	0
4630	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-06 10:43:22.285813	2025-05-06 10:43:22.315863	200	30
4631	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-06 10:43:22.577536	2025-05-06 10:43:22.579365	200	1
4632	\N	/api/restaurant	GET	2025-05-06 10:43:22.714616	2025-05-06 10:43:22.762023	200	47
4633	\N	/api/restaurant	GET	2025-05-06 10:43:25.254456	2025-05-06 10:43:25.254523	200	0
4634	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-06 10:43:35.6433	2025-05-06 10:43:35.644695	200	1
4635	\N	/robots.txt	GET	2025-05-06 23:54:07.569419	2025-05-06 23:54:07.652517	500	83
4636	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-07 04:24:43.132043	2025-05-07 04:24:43.132149	200	0
4637	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-07 07:20:04.744843	2025-05-07 07:20:04.744978	200	0
4638	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-07 07:20:05.017091	2025-05-07 07:20:05.017162	200	0
4639	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-07 09:28:11.698247	2025-05-07 09:28:11.698367	200	0
4640	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-07 18:00:57.648312	2025-05-07 18:00:57.674933	200	26
4641	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-08 00:55:03.599959	2025-05-08 00:55:03.600058	200	0
4642	\N	/.git/config	GET	2025-05-08 01:39:59.423576	2025-05-08 01:39:59.423582	500	0
4643	\N	/	GET	2025-05-08 04:02:16.259333	2025-05-08 04:02:16.259354	500	0
4644	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-08 06:51:54.26509	2025-05-08 06:51:54.26517	200	0
4645	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-08 06:51:54.276615	2025-05-08 06:51:54.276739	200	0
4646	\N	/.git/config	GET	2025-05-08 08:47:20.892279	2025-05-08 08:47:20.892282	500	0
4647	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-08 09:15:06.223031	2025-05-08 09:15:06.223223	200	0
4648	\N	/	GET	2025-05-08 12:16:03.97938	2025-05-08 12:16:03.979401	500	0
4649	\N	/.env	GET	2025-05-08 12:16:04.138242	2025-05-08 12:16:04.138244	500	0
4650	\N	/.env_example	GET	2025-05-08 12:16:04.302066	2025-05-08 12:16:04.302068	500	0
4651	\N	/core/.env	GET	2025-05-08 12:16:04.459116	2025-05-08 12:16:04.459119	500	0
4652	\N	/app/.env	GET	2025-05-08 12:16:04.614777	2025-05-08 12:16:04.614781	500	0
4653	\N	/laravel/.env	GET	2025-05-08 12:16:04.771001	2025-05-08 12:16:04.771003	500	0
4654	\N	/web/.env	GET	2025-05-08 12:16:04.926778	2025-05-08 12:16:04.92678	500	0
4655	\N	/crm/.env	GET	2025-05-08 12:16:05.08341	2025-05-08 12:16:05.083413	500	0
4656	\N	/backend/.env	GET	2025-05-08 12:16:05.248731	2025-05-08 12:16:05.248734	500	0
4657	\N	/local/.env	GET	2025-05-08 12:16:05.414462	2025-05-08 12:16:05.414465	500	0
4658	\N	/api/.env	GET	2025-05-08 12:16:05.570114	2025-05-08 12:16:05.57012	500	0
4659	\N	/admin/.env	GET	2025-05-08 12:16:05.725956	2025-05-08 12:16:05.725959	500	0
4660	\N	/application/.env	GET	2025-05-08 12:16:05.881707	2025-05-08 12:16:05.881712	500	0
4661	\N	/env.bak	GET	2025-05-08 12:16:06.037852	2025-05-08 12:16:06.037855	500	0
4662	\N	/phpinfo	GET	2025-05-08 12:16:06.210892	2025-05-08 12:16:06.210894	500	0
4663	\N	/_profiler/phpinfo	GET	2025-05-08 12:16:06.382704	2025-05-08 12:16:06.382707	500	0
4664	\N	/phpinfo.php	GET	2025-05-08 12:16:06.540779	2025-05-08 12:16:06.540781	500	0
4665	\N	/.env.bak	GET	2025-05-08 12:16:06.69901	2025-05-08 12:16:06.699012	500	0
4666	\N	/.env.backup	GET	2025-05-08 12:16:06.854567	2025-05-08 12:16:06.854569	500	0
4667	\N	/.env_sample	GET	2025-05-08 12:16:07.010481	2025-05-08 12:16:07.010483	500	0
4668	\N	/.env.old	GET	2025-05-08 12:16:07.167723	2025-05-08 12:16:07.167725	500	0
4852	\N	/	GET	2025-05-09 12:07:56.402185	2025-05-09 12:07:56.402191	500	0
4853	\N	/favicon.ico	GET	2025-05-09 12:07:57.11682	2025-05-09 12:07:57.116823	500	0
4898	\N	/api/statistics/global	GET	2025-05-09 14:08:02.193746	2025-05-09 14:08:02.356217	200	162
4911	\N	/api/statistics/endpoints	GET	2025-05-09 14:08:42.982103	2025-05-09 14:08:43.431221	200	449
4912	\N	/status/	GET	2025-05-09 14:09:35.896241	2025-05-09 14:09:35.89625	500	0
4914	\N	/status/	GET	2025-05-09 14:09:52.994436	2025-05-09 14:09:52.99445	500	0
4915	\N	/status/	GET	2025-05-09 14:10:07.872324	2025-05-09 14:10:07.872327	500	0
4942	\N	/api/auth/login	POST	2025-05-09 14:23:55.504354	2025-05-09 14:23:55.915917	200	411
4943	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 14:24:06.968564	2025-05-09 14:24:07.297672	200	329
4947	\N	/api/washingmachines	GET	2025-05-09 12:33:07.80264	2025-05-09 12:33:07.802649	500	0
4948	\N	/api/newf/me	GET	2025-05-09 12:33:22.192263	2025-05-09 12:33:22.192374	401	0
4949	\N	/api/auth/login	POST	2025-05-09 12:34:03.158955	2025-05-09 12:34:03.173575	401	14
4950	\N	/api/auth/login	POST	2025-05-09 12:34:07.029708	2025-05-09 12:34:07.108223	200	78
4951	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 12:34:07.734289	2025-05-09 12:34:07.735637	200	1
4952	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 12:34:08.166103	2025-05-09 12:34:08.179018	200	12
4953	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 12:34:08.804021	2025-05-09 12:34:08.805294	200	1
4954	\N	/api/restaurant	GET	2025-05-09 12:34:08.968018	2025-05-09 12:34:09.93614	200	968
4955	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-09 12:34:16.544247	2025-05-09 12:34:16.546369	200	2
4956	\N	/api/washingmachines	GET	2025-05-09 12:38:07.827289	2025-05-09 12:38:07.827293	500	0
4957	\N	/api/washingmachines	GET	2025-05-09 12:43:07.779553	2025-05-09 12:43:07.779558	500	0
4958	\N	/api/washingmachines	GET	2025-05-09 12:48:17.818416	2025-05-09 12:48:17.818421	500	0
4959	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 14:04:00.696519	2025-05-09 14:04:00.698224	200	1
4960	\N	/api/auth/login	POST	2025-05-09 15:05:59.427187	2025-05-09 15:05:59.494708	200	67
4961	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 15:05:59.735223	2025-05-09 15:05:59.736417	200	1
4962	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 15:05:59.952522	2025-05-09 15:05:59.953796	200	1
4963	\N	/api/weather	GET	2025-05-09 15:06:00.373492	2025-05-09 15:06:00.373496	500	0
4964	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 15:06:00.432044	2025-05-09 15:06:00.433297	200	1
4965	\N	/api/weather	GET	2025-05-09 15:06:01.63906	2025-05-09 15:06:01.639065	500	0
4966	\N	/api/restaurant	GET	2025-05-09 15:06:00.430393	2025-05-09 15:06:01.66378	200	1233
4967	\N	/api/weather	GET	2025-05-09 15:06:03.805726	2025-05-09 15:06:03.80573	500	0
4968	\N	/api/weather	GET	2025-05-09 15:06:07.973144	2025-05-09 15:06:07.973151	500	0
4969	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 15:10:38.432396	2025-05-09 15:10:38.433882	200	1
4970	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-09 15:10:39.495084	2025-05-09 15:10:39.495217	200	0
4971	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 15:10:41.742354	2025-05-09 15:10:41.743993	200	1
4972	\N	/api/restaurant	GET	2025-05-09 15:10:46.488712	2025-05-09 15:10:46.488773	200	0
4973	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 16:10:23.916163	2025-05-09 16:10:23.918979	200	2
4974	\N	/api/weather	GET	2025-05-09 16:10:24.298976	2025-05-09 16:10:24.298991	500	0
4975	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-09 16:10:24.355049	2025-05-09 16:10:24.356655	200	1
4976	\N	/api/restaurant	GET	2025-05-09 16:10:24.356453	2025-05-09 16:10:24.923801	200	567
4977	\N	/api/weather	GET	2025-05-09 16:10:25.476719	2025-05-09 16:10:25.476724	500	0
4978	\N	/api/weather	GET	2025-05-09 16:10:27.649109	2025-05-09 16:10:27.649114	500	0
4979	\N	/api/weather	GET	2025-05-09 16:10:31.813799	2025-05-09 16:10:31.813804	500	0
4980	\N	/	GET	2025-05-09 16:10:36.248464	2025-05-09 16:10:36.248466	500	0
4981	\N	/api/	GET	2025-05-09 16:11:13.447077	2025-05-09 16:11:13.447082	500	0
4982	\N	/api/weather	GET	2025-05-09 16:11:15.915491	2025-05-09 16:11:15.915496	500	0
4983	\N	/api/weather	GET	2025-05-09 16:15:04.115444	2025-05-09 16:15:04.115452	500	0
4984	\N	/	GET	2025-05-09 21:01:42.480006	2025-05-09 21:01:42.48001	500	0
4985	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-09 22:58:39.878922	2025-05-09 22:58:39.899267	200	20
4986	\N	/	GET	2025-05-10 04:30:03.608631	2025-05-10 04:30:03.608634	500	0
4987	\N	/favicon.ico	GET	2025-05-10 04:30:24.523942	2025-05-10 04:30:24.523944	500	0
4988	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 06:48:11.560581	2025-05-10 06:48:11.56067	200	0
4989	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 06:48:11.568618	2025-05-10 06:48:11.568664	200	0
4990	\N	/favicon.ico	GET	2025-05-10 07:31:03.739947	2025-05-10 07:31:03.739949	500	0
4991	\N	/api/auth/register	POST	2025-05-10 10:28:41.050136	2025-05-10 10:28:41.134694	201	84
4992	\N	/api/auth/verify-account	POST	2025-05-10 10:29:10.34331	2025-05-10 10:29:10.351825	200	8
8512	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:12:40.047662	2025-05-21 19:12:40.049873	200	2
8513	\N	/api/washingmachines	GET	2025-05-21 19:12:40.353435	2025-05-21 19:12:40.4467	200	93
4995	\N	/api/weather	GET	2025-05-10 10:29:11.437541	2025-05-10 10:29:11.437545	500	0
8515	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:12:40.525091	2025-05-21 19:12:40.527072	200	1
4997	\N	/api/weather	GET	2025-05-10 10:29:12.63919	2025-05-10 10:29:12.639196	500	0
4998	\N	/api/restaurant	GET	2025-05-10 10:29:11.580902	2025-05-10 10:29:13.013736	200	1432
4999	\N	/api/weather	GET	2025-05-10 10:29:14.85356	2025-05-10 10:29:14.853565	500	0
5000	\N	/api/weather	GET	2025-05-10 10:29:19.059896	2025-05-10 10:29:19.059902	500	0
5001	\N	/api/auth/register	GET	2025-05-10 10:33:01.577709	2025-05-10 10:33:01.577717	500	0
5002	\N	/favicon.ico	GET	2025-05-10 10:33:02.130945	2025-05-10 10:33:02.130948	500	0
5003	\N	/api/weather	GET	2025-05-10 10:33:07.050424	2025-05-10 10:33:07.05043	500	0
5004	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 10:34:49.931691	2025-05-10 10:34:49.931786	200	0
5005	\N	/api/weather	GET	2025-05-10 10:35:31.269364	2025-05-10 10:35:31.269369	500	0
5006	\N	/api/weather	GET	2025-05-10 12:35:33.793094	2025-05-10 12:35:33.79314	500	0
5007	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 10:36:12.443073	2025-05-10 10:36:12.443181	200	0
8516	\N	/api/weather	GET	2025-05-21 19:12:40.481665	2025-05-21 19:12:40.53358	200	51
5009	\N	/api/weather	GET	2025-05-10 10:36:19.819099	2025-05-10 10:36:19.819119	500	0
5010	\N	/api/restaurant	GET	2025-05-10 10:36:19.969871	2025-05-10 10:36:19.96994	200	0
8517	\N	/api/washingmachines	GET	2025-05-21 19:12:42.353442	2025-05-21 19:12:42.382863	200	29
5012	\N	/api/weather	GET	2025-05-10 10:36:21.02723	2025-05-10 10:36:21.027235	500	0
5013	\N	/api/weather	GET	2025-05-10 10:36:23.37374	2025-05-10 10:36:23.373745	500	0
5014	\N	/api/weather	GET	2025-05-10 10:36:27.572642	2025-05-10 10:36:27.572647	500	0
5015	\N	/api/weather/	GET	2025-05-10 12:39:10.69786	2025-05-10 12:39:10.816539	200	118
5016	\N	/api/weather/	GET	2025-05-10 12:39:15.035939	2025-05-10 12:39:15.035958	200	0
5017	\N	/api/weather	GET	2025-05-10 12:39:16.840952	2025-05-10 12:39:16.840977	500	0
5018	\N	/api/restaurant	GET	2025-05-10 10:40:00.742022	2025-05-10 10:40:00.742089	200	0
5019	\N	/api/weather	GET	2025-05-10 10:40:00.741999	2025-05-10 10:40:00.742003	500	0
5020	\N	/api/weather	GET	2025-05-10 10:40:01.957924	2025-05-10 10:40:01.957929	500	0
5021	\N	/api/weather	GET	2025-05-10 10:40:04.154426	2025-05-10 10:40:04.154431	500	0
5022	\N	/api/weather	GET	2025-05-10 10:40:08.40542	2025-05-10 10:40:08.405424	500	0
8514	\N	/api/restaurant	GET	2025-05-21 19:12:40.525076	2025-05-21 19:12:40.525202	200	0
8523	\N	/api/restaurant	GET	2025-05-21 19:14:29.638482	2025-05-21 19:14:29.638599	200	0
8553	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:15:05.781253	2025-05-21 20:15:05.781392	200	0
8554	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:15:05.984132	2025-05-21 20:15:05.986109	200	1
8574	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:26:46.780538	2025-05-21 20:26:46.782577	200	2
8575	\N	/api/weather	GET	2025-05-21 20:26:46.750221	2025-05-21 20:26:46.819629	200	69
8576	\N	/api/washingmachines	GET	2025-05-21 20:26:49.97898	2025-05-21 20:26:50.017393	200	38
8577	\N	/status	GET	2025-05-21 20:26:50.565266	2025-05-21 20:26:50.565268	200	0
8578	\N	/api/restaurant	GET	2025-05-21 20:26:52.14758	2025-05-21 20:26:52.147684	200	0
8579	\N	/api/traq/	GET	2025-05-21 20:26:54.496491	2025-05-21 20:26:54.49821	200	1
8580	\N	/api/traq/	GET	2025-05-21 20:26:56.757151	2025-05-21 20:26:56.758518	200	1
8581	\N	/status	GET	2025-05-21 20:27:00.562999	2025-05-21 20:27:00.563	200	0
8582	\N	/status	GET	2025-05-21 20:27:09.568138	2025-05-21 20:27:09.56814	200	0
8583	\N	/status	GET	2025-05-21 20:27:10.57423	2025-05-21 20:27:10.574232	200	0
8584	\N	/status	GET	2025-05-21 20:27:12.461688	2025-05-21 20:27:12.46169	200	0
8585	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:27:15.388402	2025-05-21 20:27:15.390867	200	2
8586	\N	/status	GET	2025-05-21 20:27:20.572831	2025-05-21 20:27:20.572832	200	0
8587	\N	/status	GET	2025-05-21 20:27:29.671649	2025-05-21 20:27:29.67165	200	0
8588	\N	/api/statistics/global	GET	2025-05-21 20:27:30.084646	2025-05-21 20:27:30.087079	200	2
8589	\N	/api/statistics/top-users	GET	2025-05-21 20:27:30.487041	2025-05-21 20:27:30.489055	200	2
8590	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:27:35.745591	2025-05-21 20:27:35.748141	200	2
8591	\N	/status	GET	2025-05-21 20:27:39.671435	2025-05-21 20:27:39.671437	200	0
8592	\N	/status	GET	2025-05-21 20:27:40.361081	2025-05-21 20:27:40.361083	200	0
8593	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:27:52.85177	2025-05-21 20:27:52.853524	200	1
8594	\N	/api/washingmachines	GET	2025-05-21 20:27:53.453533	2025-05-21 20:27:53.480898	200	27
8595	\N	/api/weather	GET	2025-05-21 20:27:53.654753	2025-05-21 20:27:53.719564	200	64
8596	\N	/api/restaurant	GET	2025-05-21 20:27:53.786573	2025-05-21 20:27:53.786721	200	0
8598	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:27:55.149662	2025-05-21 20:27:55.149795	200	0
8606	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:28:14.32201	2025-05-21 20:28:14.323915	200	1
8646	\N	/api/restaurant	GET	2025-05-21 20:32:38.427944	2025-05-21 20:32:38.428123	200	0
8647	\N	/api/washingmachines	GET	2025-05-21 20:32:38.410126	2025-05-21 20:32:38.444494	200	34
8648	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:32:41.345794	2025-05-21 20:32:41.347462	200	1
8649	\N	/api/washingmachines	GET	2025-05-21 20:32:45.615513	2025-05-21 20:32:45.646721	200	31
8650	\N	/api/washingmachines	GET	2025-05-21 20:32:49.940298	2025-05-21 20:32:49.971066	200	30
8651	\N	/api/washingmachines	GET	2025-05-21 20:32:54.148446	2025-05-21 20:32:54.184883	200	36
8652	\N	/api/restaurant	GET	2025-05-21 20:32:57.533116	2025-05-21 20:32:57.533288	200	0
8713	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 21:01:24.64582	2025-05-21 21:01:24.647445	200	1
8714	\N	/api/weather	GET	2025-05-21 21:01:24.644875	2025-05-21 21:01:24.7183	200	73
8716	\N	/api/restaurant	GET	2025-05-21 21:02:03.806662	2025-05-21 21:02:03.806774	200	0
8790	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 21:57:15.71546	2025-05-21 21:57:15.715616	200	0
8791	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 22:03:00.004811	2025-05-21 22:03:00.007589	200	2
8792	\N	/api/restaurant	GET	2025-05-21 22:03:00.116568	2025-05-21 22:03:00.116691	200	0
8819	\N	/api/restaurant	GET	2025-05-21 22:11:46.545535	2025-05-21 22:11:46.545649	200	0
8822	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 22:26:10.802056	2025-05-21 22:26:10.802177	200	0
5023	\N	/api/weather	GET	2025-05-10 12:40:09.205687	2025-05-10 12:40:09.205718	500	0
5024	\N	/api/weather	GET	2025-05-10 12:40:13.723435	2025-05-10 12:40:13.723442	500	0
5025	\N	/api/weather	GET	2025-05-10 12:40:21.373589	2025-05-10 12:40:21.373596	500	0
5026	\N	/api/weather	GET	2025-05-10 12:40:24.630069	2025-05-10 12:40:24.630079	500	0
5027	\N	/api/weather	GET	2025-05-10 12:40:29.812872	2025-05-10 12:40:29.812895	500	0
5028	\N	/api/weather/	GET	2025-05-10 12:40:33.249534	2025-05-10 12:40:33.34936	200	99
5029	\N	/api/weather/	GET	2025-05-10 12:40:37.277377	2025-05-10 12:40:37.303005	200	25
5030	\N	/api/weather/	GET	2025-05-10 12:40:45.042609	2025-05-10 12:40:45.069413	200	26
5031	\N	/api/weather	GET	2025-05-10 12:41:42.725865	2025-05-10 12:41:42.725887	500	0
5032	\N	/api/weather	GET	2025-05-10 12:41:44.544208	2025-05-10 12:41:44.544216	500	0
5033	\N	/api/weather/	GET	2025-05-10 12:41:48.304946	2025-05-10 12:41:48.30498	200	0
5034	\N	/api/weather	GET	2025-05-10 12:41:51.163825	2025-05-10 12:41:51.163833	500	0
5035	\N	/status	GET	2025-05-10 12:42:23.705273	2025-05-10 12:42:23.705275	200	0
5036	\N	/api/weather	GET	2025-05-10 12:42:30.296761	2025-05-10 12:42:30.296804	500	0
5037	\N	/status	GET	2025-05-10 12:43:21.634905	2025-05-10 12:43:21.63492	200	0
5038	\N	/api/weather	GET	2025-05-10 12:43:26.106899	2025-05-10 12:43:26.199895	200	92
5039	\N	/api/weather	GET	2025-05-10 12:44:03.246876	2025-05-10 12:44:03.332587	200	85
5040	\N	/api/auth/login	POST	2025-05-10 10:45:56.726533	2025-05-10 10:45:56.796507	401	69
5041	\N	/api/auth/login	POST	2025-05-10 10:46:05.38573	2025-05-10 10:46:05.453218	401	67
5042	\N	/api/auth/login	POST	2025-05-10 10:46:21.842259	2025-05-10 10:46:21.910424	401	68
5043	\N	/api/auth/login	POST	2025-05-10 10:46:30.639952	2025-05-10 10:46:30.707435	401	67
5044	\N	/api/auth/login	POST	2025-05-10 10:46:51.542527	2025-05-10 10:46:51.609967	401	67
5045	\N	/api/auth/login	POST	2025-05-10 10:46:56.836952	2025-05-10 10:46:56.904632	401	67
5046	\N	/api/auth/login	POST	2025-05-10 10:47:05.550276	2025-05-10 10:47:05.617853	401	67
5047	\N	/api/weather	GET	2025-05-10 10:47:09.647682	2025-05-10 10:47:09.69887	200	51
5048	\N	/api/restaurant	GET	2025-05-10 10:47:09.663642	2025-05-10 10:47:10.90082	200	1237
5049	\N	/api/auth/login	POST	2025-05-10 10:47:17.700876	2025-05-10 10:47:17.770176	401	69
8518	\N	/api/washingmachines	GET	2025-05-21 19:12:50.080193	2025-05-21 19:12:50.111595	200	31
5051	\N	/api/auth/login	POST	2025-05-10 10:47:26.835863	2025-05-10 10:47:26.903054	401	67
5052	\N	/api/auth/login	POST	2025-05-10 10:47:46.772418	2025-05-10 10:47:46.840164	401	67
5053	\N	/api/auth/login	POST	2025-05-10 10:47:49.748579	2025-05-10 10:47:49.816419	401	67
5054	\N	/api/auth/login	POST	2025-05-10 10:48:12.885438	2025-05-10 10:48:12.952761	401	67
5055	\N	/api/auth/login	POST	2025-05-10 10:48:14.821387	2025-05-10 10:48:14.888688	401	67
5056	\N	/api/auth/login	POST	2025-05-10 10:48:35.812575	2025-05-10 10:48:35.879687	401	67
5057	\N	/api/auth/login	POST	2025-05-10 10:48:39.930877	2025-05-10 10:48:39.999749	401	68
5058	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:13:22.617074	2025-05-10 11:13:22.618389	200	1
5059	\N	/api/weather	GET	2025-05-10 11:13:23.058199	2025-05-10 11:13:23.102269	200	44
5060	\N	/api/restaurant	GET	2025-05-10 11:13:23.126495	2025-05-10 11:13:23.126561	200	0
5061	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:13:23.126506	2025-05-10 11:13:23.135474	200	8
5062	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:13:23.615284	2025-05-10 11:13:23.616451	200	1
5063	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:13:30.441388	2025-05-10 11:13:30.443274	200	1
5064	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 11:13:46.845437	2025-05-10 11:13:46.845779	200	0
5065	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 11:13:46.945978	2025-05-10 11:13:46.946051	200	0
5066	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 11:14:56.011555	2025-05-10 11:14:56.011641	200	0
5067	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:16:53.907231	2025-05-10 11:16:53.908636	200	1
5068	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:16:54.428594	2025-05-10 11:16:54.429786	200	1
5069	\N	/api/restaurant	GET	2025-05-10 11:16:54.608371	2025-05-10 11:16:55.55115	200	942
5070	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:16:58.602226	2025-05-10 11:16:58.619592	200	17
5071	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 11:17:03.894446	2025-05-10 11:17:03.896014	200	1
4993	\N	/api/newf/me	GET	2025-05-10 10:29:10.717761	2025-05-10 10:29:10.719691	200	1
4994	\N	/api/newf/me	GET	2025-05-10 10:29:10.924736	2025-05-10 10:29:10.925949	200	1
4996	\N	/api/newf/me	GET	2025-05-10 10:29:11.581631	2025-05-10 10:29:11.58287	200	1
5008	\N	/api/newf/me	GET	2025-05-10 10:36:19.364835	2025-05-10 10:36:19.366209	200	1
5011	\N	/api/newf/me	GET	2025-05-10 10:36:19.970063	2025-05-10 10:36:19.971305	200	1
5050	\N	/api/newf/me	GET	2025-05-10 10:47:18.775356	2025-05-10 10:47:18.777396	200	2
8519	\N	/api/washingmachines	GET	2025-05-21 19:13:22.824625	2025-05-21 19:13:22.854923	200	30
8520	\N	/api/washingmachines	GET	2025-05-21 19:13:35.298294	2025-05-21 19:13:35.332664	200	34
8521	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:14:29.312314	2025-05-21 19:14:29.314274	200	1
8522	\N	/api/washingmachines	GET	2025-05-21 19:14:29.529902	2025-05-21 19:14:29.559804	200	29
8524	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:14:29.638393	2025-05-21 19:14:29.640184	200	1
8525	\N	/api/weather	GET	2025-05-21 19:14:29.594587	2025-05-21 19:14:29.648147	200	53
8526	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:14:40.120093	2025-05-21 19:14:40.12223	200	2
5079	\N	/api/weather	GET	2025-05-10 11:43:04.613124	2025-05-10 11:43:04.676001	200	62
5080	\N	/api/restaurant	GET	2025-05-10 11:43:04.755225	2025-05-10 11:43:05.279097	200	523
5081	\N	/api/restaurant	GET	2025-05-10 11:43:07.423151	2025-05-10 11:43:07.423198	200	0
5082	\N	/api/weather	GET	2025-05-10 11:43:07.42318	2025-05-10 11:43:07.423195	200	0
8527	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 19:14:45.310411	2025-05-21 19:14:45.312179	200	1
5084	\N	/api/restaurant	GET	2025-05-10 11:43:10.607508	2025-05-10 11:43:10.607556	200	0
8528	\N	/status	GET	2025-05-21 19:14:51.710272	2025-05-21 19:14:51.710278	200	0
8529	\N	/api/statistics/global	GET	2025-05-21 19:14:51.800548	2025-05-21 19:14:51.802731	200	2
8530	\N	/api/statistics/top-users	GET	2025-05-21 19:14:51.900049	2025-05-21 19:14:51.901861	200	1
5088	\N	/api/auth/login	POST	2025-05-10 12:03:33.646256	2025-05-10 12:03:33.715851	200	69
5089	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:33.927701	2025-05-10 12:03:33.928924	200	1
5090	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:34.106736	2025-05-10 12:03:34.108306	200	1
5091	\N	/api/auth/login	POST	2025-05-10 12:03:34.532934	2025-05-10 12:03:34.600492	200	67
5092	\N	/api/restaurant	GET	2025-05-10 12:03:34.704844	2025-05-10 12:03:34.705022	200	0
5093	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:34.704457	2025-05-10 12:03:34.705892	200	1
5094	\N	/api/weather	GET	2025-05-10 12:03:34.704805	2025-05-10 12:03:34.755478	200	50
5095	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:34.776719	2025-05-10 12:03:34.777964	200	1
5096	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:35.048475	2025-05-10 12:03:35.058259	200	9
5097	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-10 12:03:44.361216	2025-05-10 12:03:44.361559	200	0
5098	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:03:59.55504	2025-05-10 12:03:59.556414	200	1
5099	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:04:00.001126	2025-05-10 12:04:00.005587	200	4
5100	\N	/api/weather	GET	2025-05-10 12:05:00.737089	2025-05-10 12:05:00.737108	200	0
5101	\N	/api/restaurant	GET	2025-05-10 12:05:00.738456	2025-05-10 12:05:00.73853	200	0
5102	\N	/api/restaurant	GET	2025-05-10 12:05:03.506377	2025-05-10 12:05:03.506425	200	0
5103	\N	/api/weather	GET	2025-05-10 12:05:03.506395	2025-05-10 12:05:03.55307	200	46
5104	\N	/api/restaurant	GET	2025-05-10 12:05:04.887825	2025-05-10 12:05:04.887879	200	0
5105	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:06:06.491813	2025-05-10 12:06:06.493089	200	1
5106	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-10 12:06:09.569421	2025-05-10 12:06:09.571852	200	2
5107	\N	/api/weather	GET	2025-05-10 12:06:16.55552	2025-05-10 12:06:16.588075	200	32
5108	\N	/api/restaurant	GET	2025-05-10 12:06:16.555522	2025-05-10 12:06:17.849167	200	1293
5109	\N	/api/restaurant	GET	2025-05-10 12:06:20.855012	2025-05-10 12:06:20.855061	200	0
8531	\N	/status	GET	2025-05-21 19:15:01.739694	2025-05-21 19:15:01.739696	200	0
5110	\N	/api/traq/	GET	2025-05-10 12:06:24.657548	2025-05-10 12:06:24.659065	200	1
5111	\N	/api/traq/	GET	2025-05-10 12:06:32.461191	2025-05-10 12:06:32.462711	200	1
5112	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:07:03.408325	2025-05-10 12:07:03.40963	200	1
5113	\N	/api/weather	GET	2025-05-10 12:07:03.836149	2025-05-10 12:07:03.836168	200	0
5114	\N	/api/restaurant	GET	2025-05-10 12:07:03.978387	2025-05-10 12:07:03.978476	200	0
5119	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:07:11.50807	2025-05-10 12:07:11.509564	200	1
5120	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-10 12:09:34.842032	2025-05-10 12:09:34.844727	200	2
5121	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:23:20.89727	2025-05-10 12:23:20.898704	200	1
5122	\N	/api/restaurant	GET	2025-05-10 12:23:21.515326	2025-05-10 12:23:21.515379	200	0
5123	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:23:21.65754	2025-05-10 12:23:21.658981	200	1
5124	\N	/api/weather	GET	2025-05-10 12:23:21.657542	2025-05-10 12:23:21.708523	200	50
5136	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 01:05:42.113004	2025-05-11 01:05:42.113086	200	0
5261	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 07:57:40.303151	2025-05-11 07:57:40.308847	200	5
5263	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 08:55:03.422034	2025-05-11 08:55:03.422084	200	0
5264	\N	/api/auth/login	POST	2025-05-11 09:04:37.270859	2025-05-11 09:04:37.339961	401	69
5265	\N	/api/auth/login	POST	2025-05-11 09:05:04.633519	2025-05-11 09:05:04.701501	401	67
5266	\N	/api/auth/login	POST	2025-05-11 09:05:32.543708	2025-05-11 09:05:32.622847	401	79
5267	\N	/api/auth/login	POST	2025-05-11 09:05:58.919752	2025-05-11 09:05:58.987537	401	67
5268	\N	/api/auth/login	POST	2025-05-11 09:06:27.90252	2025-05-11 09:06:27.972359	401	69
5269	\N	/api/auth/login	POST	2025-05-11 09:06:52.744778	2025-05-11 09:06:52.813585	401	68
5270	\N	/api/auth/login	POST	2025-05-11 09:07:02.124125	2025-05-11 09:07:02.20375	401	79
5271	\N	/api/auth/login	POST	2025-05-11 09:07:20.417013	2025-05-11 09:07:20.484293	401	67
5272	\N	/api/auth/login	POST	2025-05-11 09:07:29.641115	2025-05-11 09:07:29.709622	401	68
5273	\N	/api/auth/login	POST	2025-05-11 09:07:46.187163	2025-05-11 09:07:46.255831	401	68
5274	\N	/api/auth/login	POST	2025-05-11 09:08:12.209103	2025-05-11 09:08:12.306814	401	97
5275	\N	/api/auth/login	POST	2025-05-11 09:08:38.25425	2025-05-11 09:08:38.337224	401	82
5276	\N	/api/auth/login	POST	2025-05-11 09:09:45.320332	2025-05-11 09:09:45.389582	401	69
5277	\N	/api/auth/login	POST	2025-05-11 09:10:05.263535	2025-05-11 09:10:05.372528	401	108
5278	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:16:32.401608	2025-05-11 09:16:32.402869	200	1
5279	\N	/api/restaurant	GET	2025-05-11 09:16:32.880942	2025-05-11 09:16:32.880995	200	0
5280	\N	/api/weather	GET	2025-05-11 09:16:33.020859	2025-05-11 09:16:33.080468	200	59
5281	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:17:21.841798	2025-05-11 09:17:21.843299	200	1
5282	\N	/api/weather	GET	2025-05-11 09:17:22.322337	2025-05-11 09:17:22.322357	200	0
5283	\N	/api/restaurant	GET	2025-05-11 09:17:22.462759	2025-05-11 09:17:22.462854	200	0
5284	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:17:28.202693	2025-05-11 09:17:28.204704	200	2
5285	\N	/api/traq/	GET	2025-05-11 09:17:38.606671	2025-05-11 09:17:38.607594	200	0
5286	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:17:56.28405	2025-05-11 09:17:56.285554	200	1
5287	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:19:57.816045	2025-05-11 09:19:57.817518	200	1
5288	\N	/api/weather	GET	2025-05-11 09:19:58.345892	2025-05-11 09:19:58.345912	200	0
5289	\N	/api/restaurant	GET	2025-05-11 09:19:58.487057	2025-05-11 09:19:58.487246	200	0
5294	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:22:35.741098	2025-05-11 09:22:35.742383	200	1
5296	\N	/api/weather	GET	2025-05-11 09:22:40.241317	2025-05-11 09:22:40.241341	200	0
5300	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 09:25:25.13532	2025-05-11 09:25:25.13542	200	0
5302	\N	/api/traq/	GET	2025-05-11 09:27:03.437378	2025-05-11 09:27:03.438265	200	0
5303	\N	/api/traq/	GET	2025-05-11 09:28:31.717951	2025-05-11 09:28:31.721803	200	3
5304	\N	/api/traq/	GET	2025-05-11 09:28:32.91937	2025-05-11 09:28:32.921419	200	2
5305	\N	/api/traq/	GET	2025-05-11 09:28:43.355854	2025-05-11 09:28:43.369101	200	13
5306	\N	/api/traq/	GET	2025-05-11 09:28:46.594455	2025-05-11 09:28:46.595469	200	1
5323	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 13:57:18.889567	2025-05-11 13:57:18.889651	200	0
5327	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 08:26:47.128992	2025-05-12 08:26:47.129079	200	0
5328	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 08:26:47.151551	2025-05-12 08:26:47.151616	200	0
5354	\N	/api/restaurant	GET	2025-05-12 12:43:17.303453	2025-05-12 12:43:17.303515	200	0
8532	\N	/api/auth/login	POST	2025-05-21 20:08:58.860732	2025-05-21 20:08:58.934163	200	73
8533	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:08:59.359623	2025-05-21 20:08:59.361951	200	2
8534	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:08:59.774566	2025-05-21 20:08:59.777051	200	2
8535	\N	/api/washingmachines	GET	2025-05-21 20:09:00.324317	2025-05-21 20:09:00.414054	200	89
8536	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:09:00.587831	2025-05-21 20:09:00.589654	200	1
8537	\N	/api/restaurant	GET	2025-05-21 20:09:00.599315	2025-05-21 20:09:00.599412	200	0
8538	\N	/api/weather	GET	2025-05-21 20:09:00.599265	2025-05-21 20:09:00.653932	200	54
8539	\N	/api/restaurant	GET	2025-05-21 20:10:15.585334	2025-05-21 20:10:15.585418	200	0
8540	\N	/api/restaurant	GET	2025-05-21 20:10:47.859413	2025-05-21 20:10:47.859492	200	0
8541	\N	/api/restaurant	GET	2025-05-21 20:10:51.890054	2025-05-21 20:10:51.890133	200	0
8542	\N	/api/restaurant	GET	2025-05-21 20:11:49.865266	2025-05-21 20:11:49.865359	200	0
8543	\N	/api/restaurant	GET	2025-05-21 20:12:04.502715	2025-05-21 20:12:04.502797	200	0
8544	\N	/api/restaurant	GET	2025-05-21 20:12:15.184108	2025-05-21 20:12:15.184188	200	0
8545	\N	/api/restaurant	GET	2025-05-21 20:12:29.299036	2025-05-21 20:12:29.299112	200	0
8546	\N	/api/restaurant	GET	2025-05-21 20:12:38.774845	2025-05-21 20:12:38.774923	200	0
8547	\N	/api/restaurant	GET	2025-05-21 20:12:46.508175	2025-05-21 20:12:46.508244	200	0
8548	\N	/api/restaurant	GET	2025-05-21 20:13:32.88304	2025-05-21 20:13:32.883197	200	0
8549	\N	/api/restaurant	GET	2025-05-21 20:13:45.100968	2025-05-21 20:13:45.101064	200	0
8550	\N	/api/restaurant	GET	2025-05-21 20:13:51.056949	2025-05-21 20:13:51.057025	200	0
8551	\N	/api/restaurant	GET	2025-05-21 20:13:51.843286	2025-05-21 20:13:51.84337	200	0
8552	\N	/api/restaurant	GET	2025-05-21 20:14:29.050254	2025-05-21 20:14:29.050333	200	0
8555	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:15:24.863565	2025-05-21 20:15:24.863676	200	0
8556	lucie.delestre@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 20:15:27.529811	2025-05-21 20:15:27.53352	200	3
8557	lucie.delestre@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 20:15:30.653902	2025-05-21 20:15:30.658631	200	4
8558	lucie.delestre@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 20:15:34.18876	2025-05-21 20:15:34.194257	200	5
8559	lucie.delestre@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 20:15:39.789881	2025-05-21 20:15:39.794379	200	4
8560	\N	/api/restaurant	GET	2025-05-21 20:16:04.046838	2025-05-21 20:16:04.046935	200	0
8561	\N	/api/restaurant	GET	2025-05-21 20:16:12.60104	2025-05-21 20:16:12.601139	200	0
8562	\N	/api/restaurant	GET	2025-05-21 20:17:06.376423	2025-05-21 20:17:06.376514	200	0
8563	\N	/api/restaurant	GET	2025-05-21 20:17:17.046034	2025-05-21 20:17:17.046125	200	0
8564	\N	/api/restaurant	GET	2025-05-21 20:17:21.800866	2025-05-21 20:17:21.800938	200	0
8565	\N	/api/restaurant	GET	2025-05-21 20:17:24.039167	2025-05-21 20:17:24.039253	200	0
8566	\N	/api/restaurant	GET	2025-05-21 20:17:50.090147	2025-05-21 20:17:50.090224	200	0
8567	\N	/api/restaurant	GET	2025-05-21 20:18:24.120705	2025-05-21 20:18:24.120791	200	0
8568	\N	/status	GET	2025-05-21 20:26:21.247627	2025-05-21 20:26:21.24763	200	0
8569	\N	/status	GET	2025-05-21 20:26:30.561299	2025-05-21 20:26:30.5613	200	0
8570	\N	/status	GET	2025-05-21 20:26:40.559279	2025-05-21 20:26:40.55928	200	0
8571	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:26:46.358258	2025-05-21 20:26:46.360539	200	2
8572	\N	/api/washingmachines	GET	2025-05-21 20:26:46.650363	2025-05-21 20:26:46.73871	200	88
8573	\N	/api/restaurant	GET	2025-05-21 20:26:46.775528	2025-05-21 20:26:46.775601	200	0
8597	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:27:53.788495	2025-05-21 20:27:53.790679	200	2
5115	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:07:03.978386	2025-05-10 12:07:03.979766	200	1
5116	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 12:07:10.939731	2025-05-10 12:07:10.94102	200	1
5117	\N	/api/weather	GET	2025-05-10 12:07:11.364994	2025-05-10 12:07:11.365012	200	0
5118	\N	/api/restaurant	GET	2025-05-10 12:07:11.505766	2025-05-10 12:07:11.505832	200	0
5125	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-10 14:41:55.203867	2025-05-10 14:41:55.203964	200	0
5126	\N	/api/auth/login	POST	2025-05-10 16:56:02.775885	2025-05-10 16:56:02.843583	401	67
5127	\N	/api/auth/login	POST	2025-05-10 16:56:03.423459	2025-05-10 16:56:03.503565	401	80
5128	\N	/api/newf/me	GET	2025-05-10 21:43:58.238838	2025-05-10 21:43:58.238952	401	0
5129	\N	/api/auth/login	POST	2025-05-10 21:44:22.833041	2025-05-10 21:44:22.900777	200	67
5130	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 21:44:23.314838	2025-05-10 21:44:23.328715	200	13
5131	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 21:44:23.731168	2025-05-10 21:44:23.732381	200	1
5132	\N	/api/weather	GET	2025-05-10 21:44:24.296695	2025-05-10 21:44:24.348167	200	51
5133	\N	/api/restaurant	GET	2025-05-10 21:44:24.506449	2025-05-10 21:44:24.506501	200	0
5134	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 21:44:24.618286	2025-05-10 21:44:24.619631	200	1
5135	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-10 21:44:28.159231	2025-05-10 21:44:28.190925	200	31
5137	\N	/.env	GET	2025-05-11 04:56:55.732526	2025-05-11 04:56:55.732538	500	0
5138	\N	/config/.env	GET	2025-05-11 04:56:55.967184	2025-05-11 04:56:55.967186	500	0
5139	\N	/.env.production	GET	2025-05-11 04:56:56.174842	2025-05-11 04:56:56.174845	500	0
5140	\N	/api/.env	GET	2025-05-11 04:56:56.339218	2025-05-11 04:56:56.339223	500	0
5141	\N	/settings/.env	GET	2025-05-11 04:56:56.497107	2025-05-11 04:56:56.497109	500	0
5142	\N	/config/app.php	GET	2025-05-11 04:56:56.694343	2025-05-11 04:56:56.694345	500	0
5143	\N	/application.yml	GET	2025-05-11 04:56:56.850219	2025-05-11 04:56:56.850224	500	0
5144	\N	/config/database.yml	GET	2025-05-11 04:56:57.024344	2025-05-11 04:56:57.024347	500	0
5145	\N	/secrets.json	GET	2025-05-11 04:56:57.196644	2025-05-11 04:56:57.196647	500	0
5146	\N	/src/config.js	GET	2025-05-11 04:56:57.362634	2025-05-11 04:56:57.362636	500	0
5147	\N	/db.ini	GET	2025-05-11 04:56:57.548473	2025-05-11 04:56:57.548475	500	0
5148	\N	/api/credentials	GET	2025-05-11 04:56:57.71353	2025-05-11 04:56:57.713535	500	0
5149	\N	/.aws/credentials	GET	2025-05-11 04:56:57.901608	2025-05-11 04:56:57.901611	500	0
5150	\N	/secure-config.json	GET	2025-05-11 04:56:58.10623	2025-05-11 04:56:58.106233	500	0
5151	\N	/local_settings.py	GET	2025-05-11 04:56:58.285735	2025-05-11 04:56:58.285738	500	0
5152	\N	/config/default.json	GET	2025-05-11 04:56:58.458564	2025-05-11 04:56:58.458568	500	0
5153	\N	/config/production.json	GET	2025-05-11 04:56:58.616118	2025-05-11 04:56:58.616121	500	0
5154	\N	/bootstrap/cache/config.php	GET	2025-05-11 04:57:06.052692	2025-05-11 04:57:06.05284	500	0
5155	\N	/config/secrets.yml	GET	2025-05-11 04:57:06.229834	2025-05-11 04:57:06.229836	500	0
5156	\N	/settings.yaml	GET	2025-05-11 04:57:06.386391	2025-05-11 04:57:06.386394	500	0
5157	\N	/auth.json	GET	2025-05-11 04:57:06.544566	2025-05-11 04:57:06.544568	500	0
5158	\N	/helm/values.yaml	GET	2025-05-11 04:57:13.282042	2025-05-11 04:57:13.282044	500	0
5159	\N	/docker/.env	GET	2025-05-11 04:57:13.491214	2025-05-11 04:57:13.491216	500	0
5160	\N	/wp-config.php	GET	2025-05-11 04:57:13.647228	2025-05-11 04:57:13.647231	500	0
5161	\N	/config.json	GET	2025-05-11 04:57:13.828474	2025-05-11 04:57:13.828477	500	0
5162	\N	/database.json	GET	2025-05-11 04:57:13.996348	2025-05-11 04:57:13.996352	500	0
5163	\N	/config/secrets.json	GET	2025-05-11 04:57:14.162345	2025-05-11 04:57:14.162348	500	0
5164	\N	/env.backup	GET	2025-05-11 04:57:14.336076	2025-05-11 04:57:14.336079	500	0
5165	\N	/settings.bak	GET	2025-05-11 04:57:14.539333	2025-05-11 04:57:14.539335	500	0
5166	\N	/backup.env	GET	2025-05-11 04:57:14.706956	2025-05-11 04:57:14.706958	500	0
5167	\N	/old/.env	GET	2025-05-11 04:57:14.862971	2025-05-11 04:57:14.862973	500	0
5168	\N	/phpinfo.php	GET	2025-05-11 04:57:15.021222	2025-05-11 04:57:15.021225	500	0
5169	\N	/info.php	GET	2025-05-11 04:57:15.181479	2025-05-11 04:57:15.181481	500	0
5170	\N	/test.php	GET	2025-05-11 04:57:15.350167	2025-05-11 04:57:15.350169	500	0
5171	\N	/laravel/.env	GET	2025-05-11 04:57:15.578569	2025-05-11 04:57:15.578571	500	0
5172	\N	/app/config/.env	GET	2025-05-11 04:57:15.772657	2025-05-11 04:57:15.772662	500	0
5173	\N	/.git/config	GET	2025-05-11 04:57:15.938228	2025-05-11 04:57:15.938231	500	0
5174	\N	/.svn/entries	GET	2025-05-11 04:57:16.268964	2025-05-11 04:57:16.268966	500	0
5175	\N	/.git/HEAD	GET	2025-05-11 04:57:16.457994	2025-05-11 04:57:16.457996	500	0
5176	\N	/.git/index	GET	2025-05-11 04:57:16.642519	2025-05-11 04:57:16.642522	500	0
5177	\N	/.git/logs/HEAD	GET	2025-05-11 04:57:16.847161	2025-05-11 04:57:16.847164	500	0
5178	\N	/.gitignore	GET	2025-05-11 04:57:17.015354	2025-05-11 04:57:17.015356	500	0
5179	\N	/administrator/index.php	GET	2025-05-11 04:57:17.179068	2025-05-11 04:57:17.179071	500	0
5180	\N	/wp-admin/install.php	GET	2025-05-11 04:57:17.398371	2025-05-11 04:57:17.398373	500	0
5181	\N	/joomla/configuration.php-dist	GET	2025-05-11 04:57:17.59478	2025-05-11 04:57:17.594783	500	0
5182	\N	/sites/default/settings.php	GET	2025-05-11 04:57:17.755407	2025-05-11 04:57:17.755409	500	0
5183	\N	/bitrix/php_interface/dbconn.php	GET	2025-05-11 04:57:17.917703	2025-05-11 04:57:17.917705	500	0
5184	\N	/typo3conf/localconf.php	GET	2025-05-11 04:57:18.076666	2025-05-11 04:57:18.076668	500	0
5185	\N	/config.inc.php	GET	2025-05-11 04:57:18.254446	2025-05-11 04:57:18.254449	500	0
5186	\N	/config.old.php	GET	2025-05-11 04:57:18.434748	2025-05-11 04:57:18.434749	500	0
5187	\N	/php.ini	GET	2025-05-11 04:57:18.63548	2025-05-11 04:57:18.635483	500	0
5188	\N	/cgi-bin/phpinfo.php	GET	2025-05-11 04:57:18.814997	2025-05-11 04:57:18.814999	500	0
5189	\N	/debug.php	GET	2025-05-11 04:57:18.975962	2025-05-11 04:57:18.975964	500	0
5190	\N	/server-status	GET	2025-05-11 04:57:19.138298	2025-05-11 04:57:19.138302	500	0
5191	\N	/phpinfo1.php	GET	2025-05-11 04:57:19.495675	2025-05-11 04:57:19.495677	500	0
5192	\N	/phpinfo2.php	GET	2025-05-11 04:57:19.654983	2025-05-11 04:57:19.654985	500	0
5193	\N	/env.txt	GET	2025-05-11 04:57:19.814845	2025-05-11 04:57:19.814848	500	0
5194	\N	/prod.env	GET	2025-05-11 04:57:20.084608	2025-05-11 04:57:20.084612	500	0
5195	\N	/stage.env	GET	2025-05-11 04:57:20.253386	2025-05-11 04:57:20.253389	500	0
5196	\N	/development.env	GET	2025-05-11 04:57:20.435539	2025-05-11 04:57:20.435541	500	0
5197	\N	/credentials.env	GET	2025-05-11 04:57:20.603383	2025-05-11 04:57:20.603386	500	0
5198	\N	/public/.env	GET	2025-05-11 04:57:20.771798	2025-05-11 04:57:20.7718	500	0
5199	\N	/api/config.json	GET	2025-05-11 04:57:20.978942	2025-05-11 04:57:20.978948	500	0
5200	\N	/composer.json	GET	2025-05-11 04:57:21.138179	2025-05-11 04:57:21.138181	500	0
5201	\N	/api/v1/.env	GET	2025-05-11 04:57:21.333015	2025-05-11 04:57:21.33302	500	0
5202	\N	/staging.env	GET	2025-05-11 04:57:21.495741	2025-05-11 04:57:21.495744	500	0
5203	\N	/phpmyadmin/index.php	GET	2025-05-11 04:57:21.676293	2025-05-11 04:57:21.676295	500	0
5204	\N	/backup/config.php	GET	2025-05-11 04:57:21.835585	2025-05-11 04:57:21.835587	500	0
5205	\N	/.env.example	GET	2025-05-11 04:57:22.0032	2025-05-11 04:57:22.003203	500	0
5206	\N	/storage/logs/laravel.log	GET	2025-05-11 04:57:29.713543	2025-05-11 04:57:29.713547	500	0
5207	\N	/storage/framework/sessions/	GET	2025-05-11 04:57:29.872786	2025-05-11 04:57:29.872789	500	0
5208	\N	/storage/framework/cache/	GET	2025-05-11 04:57:30.032138	2025-05-11 04:57:30.032141	500	0
5209	\N	/storage/framework/views/	GET	2025-05-11 04:57:30.210557	2025-05-11 04:57:30.21056	500	0
5210	\N	/nova-api/styles	GET	2025-05-11 04:57:30.377938	2025-05-11 04:57:30.377941	500	0
5211	\N	/.env.local	GET	2025-05-11 04:57:30.563762	2025-05-11 04:57:30.563772	500	0
5212	\N	/.env.dev	GET	2025-05-11 04:57:30.723525	2025-05-11 04:57:30.723527	500	0
5213	\N	/.env.test	GET	2025-05-11 04:57:30.882973	2025-05-11 04:57:30.882975	500	0
5214	\N	/var/logs/dev.log	GET	2025-05-11 04:57:31.04335	2025-05-11 04:57:31.043354	500	0
5215	\N	/var/logs/prod.log	GET	2025-05-11 04:57:31.204671	2025-05-11 04:57:31.204673	500	0
5216	\N	/config/packages/	GET	2025-05-11 04:57:31.384603	2025-05-11 04:57:31.384605	500	0
5217	\N	/web/config.php	GET	2025-05-11 04:57:31.613724	2025-05-11 04:57:31.613727	500	0
5218	\N	/config/routes.yaml	GET	2025-05-11 04:57:31.900636	2025-05-11 04:57:31.900639	500	0
5219	\N	/web.config	GET	2025-05-11 04:57:32.170439	2025-05-11 04:57:32.170441	500	0
5220	\N	/.htaccess	GET	2025-05-11 04:57:32.339553	2025-05-11 04:57:32.339556	500	0
5221	\N	/sites/all/modules/	GET	2025-05-11 04:57:32.521198	2025-05-11 04:57:32.5212	500	0
5222	\N	/sites/all/themes/	GET	2025-05-11 04:57:32.695492	2025-05-11 04:57:32.695495	500	0
5223	\N	/core/install.php	GET	2025-05-11 04:57:40.42474	2025-05-11 04:57:40.424742	500	0
5224	\N	/CHANGELOG.txt	GET	2025-05-11 04:57:40.587641	2025-05-11 04:57:40.587643	500	0
5225	\N	/app/etc/local.xml	GET	2025-05-11 04:57:40.758412	2025-05-11 04:57:40.758417	500	0
5226	\N	/app/etc/env.php	GET	2025-05-11 04:57:40.918751	2025-05-11 04:57:40.918758	500	0
5227	\N	/var/log/system.log	GET	2025-05-11 04:57:41.135704	2025-05-11 04:57:41.135705	500	0
5228	\N	/var/log/exception.log	GET	2025-05-11 04:57:41.391869	2025-05-11 04:57:41.391872	500	0
5229	\N	/.wp-config.php.swp	GET	2025-05-11 04:57:41.555414	2025-05-11 04:57:41.555416	500	0
5230	\N	/wp-config-sample.php	GET	2025-05-11 04:57:41.721892	2025-05-11 04:57:41.721895	500	0
5231	\N	/wp-content/debug.log	GET	2025-05-11 04:57:41.926536	2025-05-11 04:57:41.926539	500	0
5232	\N	/xmlrpc.php	GET	2025-05-11 04:57:42.116825	2025-05-11 04:57:42.116829	500	0
5233	\N	/wp-json/wp/v2/users	GET	2025-05-11 04:57:42.296161	2025-05-11 04:57:42.296165	500	0
5234	\N	/configuration.php~	GET	2025-05-11 04:57:42.46375	2025-05-11 04:57:42.46376	500	0
5235	\N	/logs/error.php	GET	2025-05-11 04:57:42.624044	2025-05-11 04:57:42.624048	500	0
5236	\N	/package.json	GET	2025-05-11 04:57:42.805538	2025-05-11 04:57:42.80554	500	0
5237	\N	/yarn.lock	GET	2025-05-11 04:57:43.008021	2025-05-11 04:57:43.008023	500	0
5238	\N	/.npmrc	GET	2025-05-11 04:57:43.171329	2025-05-11 04:57:43.171333	500	0
5239	\N	/server.js	GET	2025-05-11 04:57:43.3827	2025-05-11 04:57:43.382703	500	0
5240	\N	/app.js	GET	2025-05-11 04:57:43.542403	2025-05-11 04:57:43.542408	500	0
5241	\N	/config.js	GET	2025-05-11 04:57:43.737315	2025-05-11 04:57:43.737317	500	0
5242	\N	/.dockerignore	GET	2025-05-11 04:57:43.924939	2025-05-11 04:57:43.924941	500	0
5243	\N	/Dockerfile	GET	2025-05-11 04:57:44.084932	2025-05-11 04:57:44.084935	500	0
5244	\N	/.git/logs/	GET	2025-05-11 04:57:44.246206	2025-05-11 04:57:44.246209	500	0
5245	\N	/.git/refs/	GET	2025-05-11 04:57:44.456225	2025-05-11 04:57:44.456227	500	0
5246	\N	/.git/objects/	GET	2025-05-11 04:57:44.637321	2025-05-11 04:57:44.637324	500	0
5247	\N	/.git/packed-refs	GET	2025-05-11 04:57:44.896221	2025-05-11 04:57:44.896223	500	0
5248	\N	/.git/branches/	GET	2025-05-11 04:57:45.065062	2025-05-11 04:57:45.065064	500	0
5249	\N	/api-docs	GET	2025-05-11 04:57:45.263287	2025-05-11 04:57:45.263292	500	0
5250	\N	/swagger.json	GET	2025-05-11 04:57:45.461563	2025-05-11 04:57:45.461565	500	0
5251	\N	/swagger-ui.html	GET	2025-05-11 04:57:45.622071	2025-05-11 04:57:45.622073	500	0
5252	\N	/openapi.json	GET	2025-05-11 04:57:45.788936	2025-05-11 04:57:45.788939	500	0
5253	\N	/backup.sql	GET	2025-05-11 04:57:45.963174	2025-05-11 04:57:45.963176	500	0
5254	\N	/db_backup.sql	GET	2025-05-11 04:57:46.153052	2025-05-11 04:57:46.153055	500	0
5255	\N	/.well-known/security.txt	GET	2025-05-11 04:57:46.328938	2025-05-11 04:57:46.32894	500	0
5256	\N	/	GET	2025-05-11 06:50:41.558639	2025-05-11 06:50:41.558641	500	0
5257	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 07:57:35.086521	2025-05-11 07:57:35.100986	200	14
5258	\N	/api/weather	GET	2025-05-11 07:57:35.679853	2025-05-11 07:57:35.727575	200	47
5259	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 07:57:36.00757	2025-05-11 07:57:36.039764	200	32
5260	\N	/api/restaurant	GET	2025-05-11 07:57:35.887866	2025-05-11 07:57:36.834463	200	946
5262	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 08:55:03.406961	2025-05-11 08:55:03.407046	200	0
5290	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:19:58.487059	2025-05-11 09:19:58.488336	200	1
5291	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 09:22:35.158787	2025-05-11 09:22:35.164543	200	5
5292	\N	/api/weather	GET	2025-05-11 09:22:35.594261	2025-05-11 09:22:35.639453	200	45
5293	\N	/api/restaurant	GET	2025-05-11 09:22:35.741096	2025-05-11 09:22:35.741175	200	0
5295	\N	/api/restaurant	GET	2025-05-11 09:22:40.241566	2025-05-11 09:22:40.241823	200	0
5297	\N	/api/traq/	GET	2025-05-11 09:24:28.234054	2025-05-11 09:24:28.238812	200	4
5298	\N	/api/traq/	GET	2025-05-11 09:24:32.752512	2025-05-11 09:24:32.753368	200	0
5299	\N	/api/traq/	GET	2025-05-11 09:25:04.890439	2025-05-11 09:25:04.891337	200	0
5301	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 09:25:25.316222	2025-05-11 09:25:25.316281	200	0
5307	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-11 09:34:04.497103	2025-05-11 09:34:04.49721	200	0
5308	\N	/	GET	2025-05-11 12:18:05.435643	2025-05-11 12:18:05.435645	500	0
5309	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:46:44.264939	2025-05-11 12:46:44.282311	200	17
5310	\N	/api/weather	GET	2025-05-11 12:46:44.80507	2025-05-11 12:46:44.859659	200	54
5311	\N	/api/restaurant	GET	2025-05-11 12:46:45.019536	2025-05-11 12:46:45.019603	200	0
5312	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:46:45.169794	2025-05-11 12:46:45.184747	200	14
5313	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:46:48.583822	2025-05-11 12:46:48.585166	200	1
5314	\N	/api/traq/	GET	2025-05-11 12:46:56.467024	2025-05-11 12:46:56.467944	200	0
5315	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:47:33.382527	2025-05-11 12:47:33.383827	200	1
5316	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:47:39.058916	2025-05-11 12:47:39.061829	200	2
5317	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:47:45.670266	2025-05-11 12:47:45.67172	200	1
5318	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:48:31.627281	2025-05-11 12:48:31.628698	200	1
5319	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-11 12:48:32.399777	2025-05-11 12:48:32.401754	200	1
5320	\N	/api/traq/	GET	2025-05-11 12:48:54.278811	2025-05-11 12:48:54.279894	200	1
5321	\N	/api/traq/	GET	2025-05-11 12:48:59.617209	2025-05-11 12:48:59.640911	200	23
5322	\N	/api/restaurant	GET	2025-05-11 12:49:20.528403	2025-05-11 12:49:20.528473	200	0
5324	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 01:14:03.829061	2025-05-12 01:14:03.829145	200	0
5325	\N	/	GET	2025-05-12 01:24:25.091614	2025-05-12 01:24:25.091616	500	0
5326	\N	/favicon.ico	GET	2025-05-12 01:24:26.392907	2025-05-12 01:24:26.392909	500	0
5329	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 09:19:43.517302	2025-05-12 09:19:43.517395	200	0
5330	\N	/api/newf/me	GET	2025-05-12 09:31:13.520462	2025-05-12 09:31:13.520574	401	0
5331	\N	/api/auth/login	POST	2025-05-12 09:31:34.919823	2025-05-12 09:31:34.988857	200	69
5332	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 09:31:35.232828	2025-05-12 09:31:35.249242	200	16
5333	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 09:31:35.586001	2025-05-12 09:31:35.587386	200	1
5334	\N	/api/weather	GET	2025-05-12 09:31:35.96503	2025-05-12 09:31:36.020577	200	55
5335	\N	/api/restaurant	GET	2025-05-12 09:31:36.167642	2025-05-12 09:31:36.167703	200	0
5336	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 09:31:36.200022	2025-05-12 09:31:36.221736	200	21
5337	\N	/api/restaurant	GET	2025-05-12 09:31:37.708398	2025-05-12 09:31:37.708458	200	0
5338	\N	/api/restaurant	GET	2025-05-12 09:31:41.525042	2025-05-12 09:31:41.525088	200	0
5339	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 09:32:17.12215	2025-05-12 09:32:17.138956	200	16
5340	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 09:32:37.160941	2025-05-12 09:32:37.161026	304	0
5341	\N	/api/restaurant	GET	2025-05-12 09:32:39.845291	2025-05-12 09:32:39.845382	200	0
5342	\N	/api/restaurant	GET	2025-05-12 10:47:51.152363	2025-05-12 10:47:51.15244	200	0
5343	\N	/api/auth/login	POST	2025-05-12 12:39:16.427011	2025-05-12 12:39:16.494025	401	67
5344	\N	/api/auth/login	POST	2025-05-12 12:39:39.119974	2025-05-12 12:39:39.198668	401	78
5345	\N	/api/auth/login	POST	2025-05-12 12:40:00.214009	2025-05-12 12:40:00.289418	401	75
5346	\N	/api/auth/login	POST	2025-05-12 12:40:02.060576	2025-05-12 12:40:02.129549	401	68
5347	\N	/api/auth/login	POST	2025-05-12 12:40:15.197448	2025-05-12 12:40:15.276175	401	78
5348	\N	/api/auth/login	POST	2025-05-12 12:40:33.539897	2025-05-12 12:40:33.610021	401	70
5349	\N	/api/auth/login	POST	2025-05-12 12:40:41.029108	2025-05-12 12:40:41.097959	401	68
5350	\N	/api/auth/login	POST	2025-05-12 12:43:16.197454	2025-05-12 12:43:16.271304	200	73
5351	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 12:43:16.560609	2025-05-12 12:43:16.562488	200	1
5352	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 12:43:16.747234	2025-05-12 12:43:16.7485	200	1
5353	\N	/api/weather	GET	2025-05-12 12:43:17.150294	2025-05-12 12:43:17.214093	200	63
5700	\N	/status	GET	2025-05-16 14:55:59.1424	2025-05-16 14:55:59.142402	200	0
5355	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 12:43:17.303238	2025-05-12 12:43:17.326339	200	23
5356	\N	/api/restaurant	GET	2025-05-12 12:45:49.429939	2025-05-12 12:45:49.43005	200	0
5357	\N	/api/restaurant	GET	2025-05-12 12:45:54.859298	2025-05-12 12:45:54.859372	200	0
5358	\N	/api/restaurant	GET	2025-05-12 12:45:55.119786	2025-05-12 12:45:55.119843	200	0
5359	\N	/api/restaurant	GET	2025-05-12 12:45:56.104127	2025-05-12 12:45:56.104176	200	0
5360	\N	/api/restaurant	GET	2025-05-12 12:49:28.824927	2025-05-12 12:49:29.115835	200	290
5361	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:21:59.283553	2025-05-12 13:21:59.28588	200	2
5362	\N	/api/restaurant	GET	2025-05-12 13:35:27.949447	2025-05-12 13:35:27.949552	200	0
5363	\N	/api/weather	GET	2025-05-12 13:35:27.943141	2025-05-12 13:35:27.989459	200	46
5364	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:35:28.093386	2025-05-12 13:35:28.102449	200	9
5365	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:36:08.297044	2025-05-12 13:36:08.310957	200	13
5366	\N	/api/weather	GET	2025-05-12 13:36:08.770801	2025-05-12 13:36:08.770825	200	0
5367	\N	/api/restaurant	GET	2025-05-12 13:36:08.771399	2025-05-12 13:36:08.771489	200	0
5368	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:36:08.915024	2025-05-12 13:36:08.91741	200	2
5369	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:36:23.799708	2025-05-12 13:36:23.801151	200	1
5370	\N	/api/restaurant	GET	2025-05-12 13:36:24.178004	2025-05-12 13:36:24.17805	200	0
5371	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:36:24.17805	2025-05-12 13:36:24.192478	200	14
5372	\N	/api/weather	GET	2025-05-12 13:36:24.177819	2025-05-12 13:36:24.18612	200	8
5373	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 13:37:12.069094	2025-05-12 13:37:12.072548	200	3
5374	\N	/api/weather	GET	2025-05-12 13:38:35.775661	2025-05-12 13:38:35.775681	200	0
5375	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-12 13:44:57.03855	2025-05-12 13:44:57.039847	200	1
5376	\N	/api/weather	GET	2025-05-12 14:06:20.53019	2025-05-12 14:06:20.57453	200	44
5377	\N	/api/weather	GET	2025-05-12 14:08:08.591722	2025-05-12 14:08:08.591747	200	0
5378	\N	/api/support	GET	2025-05-12 14:08:33.147382	2025-05-12 14:08:33.147391	500	0
5379	\N	/api/support	GET	2025-05-12 14:08:34.358599	2025-05-12 14:08:34.358605	500	0
5380	\N	/api/support	GET	2025-05-12 14:08:36.556288	2025-05-12 14:08:36.556294	500	0
5381	\N	/api/support	GET	2025-05-12 14:08:40.790918	2025-05-12 14:08:40.790923	500	0
5382	\N	/api/weather	GET	2025-05-12 17:03:55.167096	2025-05-12 17:03:55.302998	200	135
5383	\N	/api/auth/login	POST	2025-05-12 17:04:59.788493	2025-05-12 17:05:00.204798	200	416
5384	yohann.chavanel@imt-atlantique.net	/api/support	POST	2025-05-12 17:05:46.766991	2025-05-12 17:05:47.431166	201	664
5385	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 17:13:48.558446	2025-05-12 17:13:49.395215	200	836
5386	\N	/api/support	GET	2025-05-12 15:14:54.042975	2025-05-12 15:14:54.042982	500	0
5387	\N	/api/support	GET	2025-05-12 15:14:55.257586	2025-05-12 15:14:55.257592	500	0
5388	\N	/api/support	GET	2025-05-12 15:14:57.47766	2025-05-12 15:14:57.477666	500	0
5389	\N	/api/support	GET	2025-05-12 15:15:01.676084	2025-05-12 15:15:01.67609	500	0
5390	yohann.chavanel@imt-atlantique.net	/api/support	POST	2025-05-12 17:34:57.110447	2025-05-12 17:34:57.796085	201	685
5391	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 17:35:00.174217	2025-05-12 17:35:00.99616	200	821
5392	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 17:41:03.791586	2025-05-12 17:41:05.303734	200	1512
5393	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:44:24.132046	2025-05-12 15:44:24.190814	200	58
5394	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 15:44:48.925607	2025-05-12 15:44:48.92763	200	2
5395	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 15:44:49.592175	2025-05-12 15:44:49.59377	200	1
5396	\N	/api/weather	GET	2025-05-12 15:44:49.732223	2025-05-12 15:44:49.787875	200	55
5397	\N	/api/restaurant	GET	2025-05-12 15:44:49.732267	2025-05-12 15:44:49.988347	200	256
5398	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:44:56.441864	2025-05-12 15:44:56.482456	200	40
5399	yohann.chavanel@imt-atlantique.net	/api/support	POST	2025-05-12 17:45:17.278289	2025-05-12 17:45:17.940852	201	662
5400	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:45:20.98014	2025-05-12 15:45:21.014397	200	34
5401	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:45:31.109014	2025-05-12 15:45:31.11253	200	3
5402	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:45:33.307132	2025-05-12 15:45:33.323999	200	16
5403	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:45:44.114305	2025-05-12 15:45:44.118476	200	4
5404	yohann.chavanel@imt-atlantique.net	/api/support	POST	2025-05-12 15:46:36.818781	2025-05-12 15:46:36.824147	201	5
5405	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 15:46:39.73119	2025-05-12 15:46:39.733851	200	2
5406	yohann.chavanel@imt-atlantique.net	/api/support	GET	2025-05-12 17:46:52.777923	2025-05-12 17:46:53.936012	200	1158
5407	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 17:02:16.835759	2025-05-12 17:02:16.838577	200	2
5408	\N	/api/weather	GET	2025-05-12 17:02:17.308241	2025-05-12 17:02:17.361662	200	53
5409	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-12 17:02:17.629659	2025-05-12 17:02:17.631508	200	1
5410	\N	/api/restaurant	GET	2025-05-12 17:02:17.540369	2025-05-12 17:02:18.558386	200	1018
5411	\N	/api/restaurant	GET	2025-05-12 17:02:32.110813	2025-05-12 17:02:32.110903	200	0
5412	\N	/api/auth/login	POST	2025-05-12 18:10:07.363406	2025-05-12 18:10:07.44421	200	80
5413	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:10:07.873021	2025-05-12 18:10:07.874786	200	1
5414	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:10:08.279102	2025-05-12 18:10:08.305218	200	26
5415	\N	/api/weather	GET	2025-05-12 18:10:09.107902	2025-05-12 18:10:09.177641	200	69
5416	\N	/api/restaurant	GET	2025-05-12 18:10:09.28563	2025-05-12 18:10:09.285704	200	0
5417	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:10:09.286778	2025-05-12 18:10:09.293818	200	7
5418	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-12 18:10:12.293975	2025-05-12 18:10:12.301048	200	7
5419	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:10:31.276142	2025-05-12 18:10:31.299874	200	23
5420	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-12 18:10:31.879832	2025-05-12 18:10:31.879918	304	0
5421	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-12 18:10:33.509296	2025-05-12 18:10:33.509364	304	0
5422	\N	/api/traq/	GET	2025-05-12 18:10:58.122642	2025-05-12 18:10:58.124783	200	2
5423	\N	/api/restaurant	GET	2025-05-12 18:11:10.93539	2025-05-12 18:11:10.935438	200	0
5424	\N	/api/restaurant	GET	2025-05-12 18:15:25.296064	2025-05-12 18:15:25.296156	200	0
5425	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:15:36.709001	2025-05-12 18:15:36.710417	200	1
5426	\N	/api/restaurant	GET	2025-05-12 18:15:37.673875	2025-05-12 18:15:37.673971	200	0
5427	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-12 18:15:37.673862	2025-05-12 18:15:37.675319	200	1
5428	\N	/api/weather	GET	2025-05-12 18:15:37.612106	2025-05-12 18:15:37.719611	200	107
5429	\N	/api/restaurant	GET	2025-05-12 18:16:10.78983	2025-05-12 18:16:10.789882	200	0
5430	\N	/api/restaurant	GET	2025-05-12 18:16:24.714443	2025-05-12 18:16:24.714515	200	0
5431	\N	/api/restaurant	GET	2025-05-12 18:16:44.416201	2025-05-12 18:16:44.41627	200	0
5432	\N	/api/restaurant	GET	2025-05-12 18:17:02.74092	2025-05-12 18:17:02.740968	200	0
5433	\N	/api/restaurant	GET	2025-05-12 18:17:45.636864	2025-05-12 18:17:45.636951	200	0
5434	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-13 00:35:23.777624	2025-05-13 00:35:23.779466	200	1
5435	\N	/robots.txt	GET	2025-05-13 04:32:38.624934	2025-05-13 04:32:38.624941	500	0
5436	\N	/	GET	2025-05-13 04:32:38.688734	2025-05-13 04:32:38.688737	500	0
5437	\N	/robots.txt	GET	2025-05-13 06:08:02.471452	2025-05-13 06:08:02.471455	500	0
5438	\N	/api/statistics/endpoints	GET	2025-05-13 07:52:33.846448	2025-05-13 07:52:33.882831	200	36
5439	\N	/favicon.ico	GET	2025-05-13 07:52:34.504963	2025-05-13 07:52:34.504965	500	0
5440	\N	/api/statistics/endpoints	GET	2025-05-13 07:53:32.226433	2025-05-13 07:53:32.265129	200	38
5441	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-13 08:23:51.463722	2025-05-13 08:23:51.463796	200	0
5442	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-13 08:23:51.788147	2025-05-13 08:23:51.788212	200	0
5443	\N	/api/statistics/endpoints	GET	2025-05-13 08:58:22.445237	2025-05-13 08:58:22.451439	200	6
5444	\N	/api/newf/me	GET	2025-05-13 09:33:33.621377	2025-05-13 09:33:33.621486	401	0
5445	\N	/api/auth/login	POST	2025-05-13 09:33:58.530667	2025-05-13 09:33:58.599685	200	69
5446	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-13 09:33:58.862999	2025-05-13 09:33:58.864335	200	1
5447	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-13 09:33:59.089167	2025-05-13 09:33:59.097489	200	8
5448	\N	/api/weather	GET	2025-05-13 09:33:59.469306	2025-05-13 09:33:59.516466	200	47
5449	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-13 09:33:59.628166	2025-05-13 09:33:59.629525	200	1
5450	\N	/api/restaurant	GET	2025-05-13 09:33:59.628168	2025-05-13 09:34:00.571169	200	943
5453	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-13 10:48:04.005725	2025-05-13 10:48:04.005828	200	0
5454	\N	/api/statistics/endpoints	GET	2025-05-13 12:55:28.033656	2025-05-13 12:55:28.041192	200	7
5455	\N	/favicon.ico	GET	2025-05-13 12:55:28.690964	2025-05-13 12:55:28.690967	500	0
5457	\N	/api/statistics/endpoints	GET	2025-05-13 14:39:24.559503	2025-05-13 14:39:24.566978	200	7
5458	\N	/	GET	2025-05-14 00:31:53.032524	2025-05-14 00:31:53.03253	500	0
5459	\N	/favicon.ico	GET	2025-05-14 00:31:55.720581	2025-05-14 00:31:55.720584	500	0
5462	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-14 07:14:41.209888	2025-05-14 07:14:41.209968	200	0
5463	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-14 07:14:41.484687	2025-05-14 07:14:41.484752	200	0
5471	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:17.169502	2025-05-14 13:13:17.18595	200	16
5472	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:19.694304	2025-05-14 13:13:19.695699	200	1
5473	\N	/api/restaurant	GET	2025-05-14 13:13:23.56498	2025-05-14 13:13:23.565037	200	0
5474	\N	/api/restaurant	GET	2025-05-14 13:13:33.926614	2025-05-14 13:13:33.926662	200	0
5475	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:36.834713	2025-05-14 13:13:36.836646	200	1
5476	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:39.027187	2025-05-14 13:13:39.029307	200	2
5477	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:44.1687	2025-05-14 13:13:44.170101	200	1
5478	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-14 13:13:49.557092	2025-05-14 13:13:49.559618	200	2
5479	\N	/api/restaurant	GET	2025-05-14 13:13:53.473371	2025-05-14 13:13:55.448866	200	1975
5484	\N	/api/restaurant	GET	2025-05-14 13:14:26.427595	2025-05-14 13:14:26.427647	200	0
5485	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:14:31.906256	2025-05-14 13:14:31.910397	200	4
5486	\N	/api/weather	GET	2025-05-14 13:14:32.252214	2025-05-14 13:14:32.310675	200	58
5488	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:14:32.428169	2025-05-14 13:14:32.429999	200	1
5489	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:14:34.506228	2025-05-14 13:14:34.507615	200	1
5490	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:14:35.535561	2025-05-14 13:14:35.537002	200	1
5491	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-14 13:14:38.52665	2025-05-14 13:14:38.529265	200	2
5492	\N	/api/restaurant	GET	2025-05-14 13:14:40.203541	2025-05-14 13:14:40.203608	200	0
8599	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:27:56.125402	2025-05-21 20:27:56.125498	200	0
8600	\N	/status	GET	2025-05-21 20:27:58.282278	2025-05-21 20:27:58.282281	200	0
8601	\N	/api/status	GET	2025-05-21 20:28:05.585576	2025-05-21 20:28:05.585589	500	0
8602	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:28:13.389018	2025-05-21 20:28:13.391367	200	2
8603	\N	/api/washingmachines	GET	2025-05-21 20:28:13.981945	2025-05-21 20:28:14.011	200	29
8604	\N	/api/weather	GET	2025-05-21 20:28:14.180529	2025-05-21 20:28:14.18056	200	0
8605	\N	/api/restaurant	GET	2025-05-21 20:28:14.321914	2025-05-21 20:28:14.322029	200	0
8607	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:28:17.500119	2025-05-21 20:28:17.500244	200	0
8608	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-21 20:28:18.556489	2025-05-21 20:28:18.556572	200	0
8609	\N	/api/status	GET	2025-05-21 20:28:20.666386	2025-05-21 20:28:20.666397	500	0
8610	\N	/api/status	GET	2025-05-21 20:28:30.654773	2025-05-21 20:28:30.654814	500	0
8611	\N	/api/status	GET	2025-05-21 20:28:34.161925	2025-05-21 20:28:34.161935	500	0
8612	\N	/api/status	GET	2025-05-21 20:28:44.144796	2025-05-21 20:28:44.144805	500	0
8613	\N	/api/status	GET	2025-05-21 20:28:54.173123	2025-05-21 20:28:54.173131	500	0
8614	\N	/api/status	GET	2025-05-21 20:29:04.153224	2025-05-21 20:29:04.153232	500	0
8615	\N	/api/status	GET	2025-05-21 20:29:14.15124	2025-05-21 20:29:14.151256	500	0
8616	\N	/api/status	GET	2025-05-21 20:29:20.505509	2025-05-21 20:29:20.505518	500	0
8617	\N	/api/status	GET	2025-05-21 20:29:30.498866	2025-05-21 20:29:30.498873	500	0
8618	\N	/status	GET	2025-05-21 20:29:34.587546	2025-05-21 20:29:34.587547	200	0
8619	\N	/status	GET	2025-05-21 20:29:40.260536	2025-05-21 20:29:40.260539	200	0
8620	\N	/status	GET	2025-05-21 20:29:41.47998	2025-05-21 20:29:41.479982	200	0
8621	\N	/status	GET	2025-05-21 20:29:44.571571	2025-05-21 20:29:44.571577	200	0
8622	\N	/status	GET	2025-05-21 20:29:45.888577	2025-05-21 20:29:45.888579	200	0
8623	\N	/api/statistics/global	GET	2025-05-21 20:29:46.338559	2025-05-21 20:29:46.340843	200	2
8624	\N	/api/statistics/top-users	GET	2025-05-21 20:29:46.728252	2025-05-21 20:29:46.730513	200	2
8625	\N	/status	GET	2025-05-21 20:29:55.893922	2025-05-21 20:29:55.893924	200	0
8626	\N	/status	GET	2025-05-21 20:30:05.889519	2025-05-21 20:30:05.889521	200	0
8627	\N	/status	GET	2025-05-21 20:30:11.438575	2025-05-21 20:30:11.438576	200	0
8628	\N	/api/statistics/global	GET	2025-05-21 20:30:11.804432	2025-05-21 20:30:11.807537	200	3
8629	\N	/api/statistics/top-users	GET	2025-05-21 20:30:12.172146	2025-05-21 20:30:12.173956	200	1
8630	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:30:12.552188	2025-05-21 20:30:12.554027	200	1
8631	\N	/api/restaurant	GET	2025-05-21 20:30:19.122098	2025-05-21 20:30:19.122199	200	0
8632	\N	/api/restaurant	GET	2025-05-21 20:30:19.230769	2025-05-21 20:30:19.230888	200	0
8633	\N	/api/weather	GET	2025-05-21 20:30:19.230755	2025-05-21 20:30:19.297135	200	66
8634	\N	/api/washingmachines	GET	2025-05-21 20:30:19.218952	2025-05-21 20:30:19.313663	200	94
8635	\N	/api/restaurant	GET	2025-05-21 20:30:19.445051	2025-05-21 20:30:19.445131	200	0
8636	\N	/api/restaurant	GET	2025-05-21 20:30:19.741847	2025-05-21 20:30:19.74192	200	0
8637	\N	/api/restaurant	GET	2025-05-21 20:30:20.019492	2025-05-21 20:30:20.019567	200	0
8638	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:31:27.29756	2025-05-21 20:31:27.299546	200	1
8639	\N	/api/traq/	GET	2025-05-21 20:31:30.083074	2025-05-21 20:31:30.084535	200	1
8640	\N	/api/traq/	GET	2025-05-21 20:31:55.376506	2025-05-21 20:31:55.377616	200	1
8641	\N	/api/traq/	GET	2025-05-21 20:31:56.93898	2025-05-21 20:31:56.940112	200	1
8642	\N	/api/traq/	GET	2025-05-21 20:32:14.074569	2025-05-21 20:32:14.075873	200	1
8643	\N	/api/washingmachines	GET	2025-05-21 20:32:18.229572	2025-05-21 20:32:18.324334	200	94
8644	\N	/api/washingmachines	GET	2025-05-21 20:32:19.648945	2025-05-21 20:32:19.697362	200	48
8645	\N	/api/weather	GET	2025-05-21 20:32:38.427943	2025-05-21 20:32:38.42797	200	0
8653	\N	/api/weather	GET	2025-05-21 20:32:57.53333	2025-05-21 20:32:57.53335	200	0
8654	\N	/api/washingmachines	GET	2025-05-21 20:32:57.530961	2025-05-21 20:32:57.561751	200	30
8655	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:33:00.931813	2025-05-21 20:33:00.93382	200	2
8656	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 20:33:04.388397	2025-05-21 20:33:04.39046	200	2
8657	\N	/api/restaurant	GET	2025-05-21 20:34:28.384174	2025-05-21 20:34:28.384262	200	0
8658	\N	/api/restaurant	GET	2025-05-21 20:35:58.861609	2025-05-21 20:35:58.861712	200	0
8659	\N	/api/weather	GET	2025-05-21 20:35:58.861595	2025-05-21 20:35:58.90597	200	44
8660	\N	/api/washingmachines	GET	2025-05-21 20:35:58.855943	2025-05-21 20:35:58.925576	200	69
8661	\N	/api/restaurant	GET	2025-05-21 20:39:06.225405	2025-05-21 20:39:06.225534	200	0
8662	\N	/api/restaurant	GET	2025-05-21 20:39:10.211609	2025-05-21 20:39:10.211697	200	0
8663	\N	/api/restaurant	GET	2025-05-21 20:39:26.496061	2025-05-21 20:39:26.496148	200	0
8664	\N	/api/restaurant	GET	2025-05-21 20:40:52.000208	2025-05-21 20:40:52.00033	200	0
8665	\N	/api/restaurant	GET	2025-05-21 20:40:59.238894	2025-05-21 20:40:59.238979	200	0
8666	\N	/api/restaurant	GET	2025-05-21 20:41:06.230093	2025-05-21 20:41:06.230175	200	0
8667	\N	/api/restaurant	GET	2025-05-21 20:41:14.337273	2025-05-21 20:41:14.337357	200	0
8668	\N	/api/restaurant	GET	2025-05-21 20:41:26.386205	2025-05-21 20:41:26.386284	200	0
8669	\N	/api/restaurant	GET	2025-05-21 20:42:06.643823	2025-05-21 20:42:06.644001	200	0
8670	\N	/api/restaurant	GET	2025-05-21 20:42:27.312631	2025-05-21 20:42:27.312719	200	0
5451	\N	/api/restaurant	GET	2025-05-13 09:34:02.237521	2025-05-13 09:34:02.237575	200	0
5452	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-13 09:34:28.451095	2025-05-13 09:34:28.452676	200	1
5456	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-13 13:31:01.087793	2025-05-13 13:31:01.087965	200	0
5460	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-14 00:35:35.844435	2025-05-14 00:35:35.844518	200	0
5461	\N	/robots.txt	GET	2025-05-14 00:57:13.595836	2025-05-14 00:57:13.595839	500	0
5464	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-14 10:24:45.714971	2025-05-14 10:24:45.715056	200	0
5465	\N	/	GET	2025-05-14 10:41:16.609383	2025-05-14 10:41:16.609386	500	0
5466	\N	/favicon.ico	GET	2025-05-14 10:41:17.782556	2025-05-14 10:41:17.782558	500	0
5467	\N	/ads.txt	GET	2025-05-14 10:41:18.02277	2025-05-14 10:41:18.022773	500	0
5468	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-14 13:13:16.621414	2025-05-14 13:13:16.647317	200	25
5469	\N	/api/weather	GET	2025-05-14 13:13:17.014175	2025-05-14 13:13:17.092365	200	78
5470	\N	/api/restaurant	GET	2025-05-14 13:13:17.169309	2025-05-14 13:13:17.169384	200	0
5480	\N	/api/restaurant	GET	2025-05-14 13:14:04.007603	2025-05-14 13:14:04.007672	200	0
5481	\N	/api/restaurant	GET	2025-05-14 13:14:09.471567	2025-05-14 13:14:09.471622	200	0
5482	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-14 13:14:16.686909	2025-05-14 13:14:16.689173	200	2
5483	\N	/api/restaurant	GET	2025-05-14 13:14:21.485893	2025-05-14 13:14:22.678038	200	1192
5487	\N	/api/restaurant	GET	2025-05-14 13:14:32.428187	2025-05-14 13:14:32.428273	200	0
5493	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-14 14:22:17.604875	2025-05-14 14:22:17.606051	200	1
5494	\N	/	GET	2025-05-14 14:58:42.928734	2025-05-14 14:58:42.928739	500	0
5495	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 00:24:23.063225	2025-05-15 00:24:23.063315	200	0
5496	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 08:39:43.79659	2025-05-15 08:39:43.796682	200	0
5497	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 08:39:43.927734	2025-05-15 08:39:43.927787	200	0
5498	\N	/	GET	2025-05-15 09:58:33.394403	2025-05-15 09:58:33.394406	500	0
5499	\N	/favicon.ico	GET	2025-05-15 09:58:34.040551	2025-05-15 09:58:34.040554	500	0
5500	\N	/	GET	2025-05-15 09:59:00.608349	2025-05-15 09:59:00.608352	500	0
5501	\N	/ru	GET	2025-05-15 09:59:05.161652	2025-05-15 09:59:05.161654	500	0
5502	\N	/info	GET	2025-05-15 09:59:08.147026	2025-05-15 09:59:08.147028	500	0
5503	\N	/	GET	2025-05-15 09:59:11.738666	2025-05-15 09:59:11.738669	500	0
5504	\N	/api	GET	2025-05-15 09:59:15.726287	2025-05-15 09:59:15.726292	500	0
5505	\N	/	GET	2025-05-15 09:59:18.54261	2025-05-15 09:59:18.542613	500	0
5506	\N	/api	GET	2025-05-15 09:59:20.138852	2025-05-15 09:59:20.138858	500	0
5507	\N	/	GET	2025-05-15 09:59:20.609364	2025-05-15 09:59:20.609368	500	0
5508	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-15 10:00:30.229934	2025-05-15 10:00:30.268648	200	38
5509	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-15 10:00:30.67499	2025-05-15 10:00:30.712426	200	37
5510	\N	/api/restaurant	GET	2025-05-15 10:00:30.852373	2025-05-15 10:00:31.12117	200	268
5511	\N	/api/weather	GET	2025-05-15 10:00:30.852724	2025-05-15 10:00:32.054545	200	1201
5512	\N	/api/restaurant	GET	2025-05-15 10:00:32.429505	2025-05-15 10:00:32.42956	200	0
5513	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 10:01:28.341693	2025-05-15 10:01:28.341775	304	0
5514	\N	/api/restaurant	GET	2025-05-15 10:01:30.727314	2025-05-15 10:01:30.727389	200	0
5515	\N	/api/statistics/endpoints	GET	2025-05-15 10:02:05.033135	2025-05-15 10:02:05.09165	200	58
5516	\N	/favicon.ico	GET	2025-05-15 10:02:05.829754	2025-05-15 10:02:05.829756	500	0
5517	\N	/api/restaurant	GET	2025-05-15 10:02:10.934908	2025-05-15 10:02:10.934966	200	0
5518	\N	/api/restaurant	GET	2025-05-15 10:02:16.641341	2025-05-15 10:02:16.64141	200	0
5519	\N	/api/restaurant	GET	2025-05-15 10:02:17.98263	2025-05-15 10:02:17.982693	200	0
5520	\N	/api/restaurant	GET	2025-05-15 10:02:19.806282	2025-05-15 10:02:19.806353	200	0
5521	\N	/api/restaurant	GET	2025-05-15 10:02:33.349391	2025-05-15 10:02:33.349524	200	0
5522	\N	/api/restaurant	GET	2025-05-15 10:04:24.033904	2025-05-15 10:04:24.034006	200	0
5523	\N	/favicon.ico	GET	2025-05-15 10:04:24.40082	2025-05-15 10:04:24.400823	500	0
5524	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 10:49:20.075913	2025-05-15 10:49:20.076004	200	0
5525	\N	/api/restaurant	GET	2025-05-15 11:50:34.221736	2025-05-15 11:50:34.221821	200	0
5526	\N	/favicon.ico	GET	2025-05-15 11:50:34.609849	2025-05-15 11:50:34.609851	500	0
5527	\N	/api/restaurant	GET	2025-05-15 12:17:41.650804	2025-05-15 12:17:41.650896	200	0
5528	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-15 13:09:54.995983	2025-05-15 13:09:54.996084	200	0
5529	\N	/api/restaurant	GET	2025-05-15 13:19:03.29	2025-05-15 13:19:03.290088	200	0
5530	\N	/api/auth/login	POST	2025-05-15 14:28:57.911896	2025-05-15 14:28:58.001863	401	89
5531	\N	/	GET	2025-05-15 22:50:26.683799	2025-05-15 22:50:26.683801	500	0
5532	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 06:54:13.919947	2025-05-16 06:54:13.920023	200	0
5533	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 08:12:31.629095	2025-05-16 08:12:31.629262	200	0
5534	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 08:12:32.417167	2025-05-16 08:12:32.417234	200	0
5535	\N	/	GET	2025-05-16 11:18:20.770314	2025-05-16 11:18:20.770316	500	0
5536	\N	/	GET	2025-05-16 11:18:21.438905	2025-05-16 11:18:21.438908	500	0
5537	\N	/phpinfo.php	GET	2025-05-16 11:18:21.711009	2025-05-16 11:18:21.711013	500	0
5538	\N	/info.php	GET	2025-05-16 11:18:22.337538	2025-05-16 11:18:22.33754	500	0
5539	\N	/phpinfo	GET	2025-05-16 11:18:23.023062	2025-05-16 11:18:23.023064	500	0
5540	\N	/.env	GET	2025-05-16 11:18:23.650499	2025-05-16 11:18:23.650502	500	0
5541	\N	/config/.env	GET	2025-05-16 11:18:23.95833	2025-05-16 11:18:23.958333	500	0
5542	\N	/api/.env	GET	2025-05-16 11:18:24.767354	2025-05-16 11:18:24.767359	500	0
5543	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 14:14:17.833729	2025-05-16 14:14:17.833816	200	0
5544	\N	/api/auth/login	POST	2025-05-16 14:23:46.416771	2025-05-16 14:23:46.493187	200	76
5545	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:23:46.693969	2025-05-16 14:23:46.70369	200	9
5546	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:23:46.92304	2025-05-16 14:23:46.924878	200	1
5547	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:23:47.55439	2025-05-16 14:23:47.556208	200	1
5548	\N	/api/weather	GET	2025-05-16 14:23:47.554654	2025-05-16 14:23:47.784366	200	229
5549	\N	/api/restaurant	GET	2025-05-16 14:23:47.400436	2025-05-16 14:23:48.270524	200	870
5550	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-16 14:23:51.70471	2025-05-16 14:23:51.711419	200	6
5551	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:23:59.448175	2025-05-16 14:23:59.450061	200	1
5552	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-16 14:23:59.680385	2025-05-16 14:23:59.685548	200	5
5553	\N	/api/restaurant	GET	2025-05-16 14:24:08.944754	2025-05-16 14:24:08.944824	200	0
5554	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 14:31:18.719324	2025-05-16 14:31:18.719427	200	0
5555	\N	/status	GET	2025-05-16 14:31:19.87667	2025-05-16 14:31:19.876672	200	0
5556	\N	/status	GET	2025-05-16 14:31:20.06652	2025-05-16 14:31:20.066522	200	0
5557	\N	/api/statistics/global	GET	2025-05-16 14:31:20.225739	2025-05-16 14:31:20.233075	200	7
5558	\N	/api/statistics/global	GET	2025-05-16 14:31:20.425621	2025-05-16 14:31:20.427256	200	1
5559	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:20.494374	2025-05-16 14:31:20.499645	200	5
5560	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:21.216293	2025-05-16 14:31:21.221848	200	5
5561	\N	/status	GET	2025-05-16 14:31:24.804265	2025-05-16 14:31:24.804274	200	0
5562	\N	/status	GET	2025-05-16 14:31:25.001692	2025-05-16 14:31:25.001694	200	0
5563	\N	/api/statistics/global	GET	2025-05-16 14:31:25.001823	2025-05-16 14:31:25.00335	200	1
5564	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:25.231822	2025-05-16 14:31:25.236965	200	5
5565	\N	/status	GET	2025-05-16 14:31:29.808355	2025-05-16 14:31:29.808356	200	0
5566	\N	/status	GET	2025-05-16 14:31:30.011677	2025-05-16 14:31:30.011679	200	0
5567	\N	/api/statistics/global	GET	2025-05-16 14:31:30.011652	2025-05-16 14:31:30.013564	200	1
5568	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:30.223487	2025-05-16 14:31:30.228196	200	4
5569	\N	/status	GET	2025-05-16 14:31:34.806327	2025-05-16 14:31:34.806328	200	0
5570	\N	/status	GET	2025-05-16 14:31:35.001835	2025-05-16 14:31:35.001837	200	0
5571	\N	/api/statistics/global	GET	2025-05-16 14:31:35.012119	2025-05-16 14:31:35.014437	200	2
5572	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:35.235073	2025-05-16 14:31:35.240079	200	5
5573	\N	/status	GET	2025-05-16 14:31:39.804961	2025-05-16 14:31:39.804962	200	0
5574	\N	/status	GET	2025-05-16 14:31:39.997916	2025-05-16 14:31:39.997917	200	0
5575	\N	/api/statistics/global	GET	2025-05-16 14:31:40.03628	2025-05-16 14:31:40.037926	200	1
5576	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:40.236366	2025-05-16 14:31:40.242043	200	5
5577	\N	/status	GET	2025-05-16 14:31:44.810258	2025-05-16 14:31:44.810259	200	0
5578	\N	/status	GET	2025-05-16 14:31:45.033415	2025-05-16 14:31:45.033417	200	0
5585	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:35:55.934748	2025-05-16 14:35:55.936384	200	1
5586	\N	/api/weather	GET	2025-05-16 14:35:55.780073	2025-05-16 14:35:56.007967	200	227
5587	\N	/api/auth/register	POST	2025-05-16 14:36:37.811071	2025-05-16 14:36:37.897179	201	86
5588	\N	/api/auth/verify-account	POST	2025-05-16 14:37:00.866064	2025-05-16 14:37:01.145699	200	279
5589	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:01.540764	2025-05-16 14:37:01.542465	200	1
5590	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:01.730519	2025-05-16 14:37:01.731862	200	1
5591	\N	/api/weather	GET	2025-05-16 14:37:02.476168	2025-05-16 14:37:02.476192	200	0
5593	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:02.628681	2025-05-16 14:37:02.630827	200	2
5594	zephyr.dassouli@imt-atlantique.net	/api/newf/me	PATCH	2025-05-16 14:37:04.772769	2025-05-16 14:37:04.777311	200	4
5595	\N	/api/restaurant	GET	2025-05-16 14:37:12.905866	2025-05-16 14:37:12.90593	200	0
5596	\N	/api/traq/	GET	2025-05-16 14:37:15.521862	2025-05-16 14:37:15.528642	200	6
5597	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:20.023422	2025-05-16 14:37:20.024837	200	1
5598	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:21.557302	2025-05-16 14:37:21.558848	200	1
5599	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:23.605186	2025-05-16 14:37:23.606608	200	1
5600	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:54.528195	2025-05-16 14:37:54.529799	200	1
5601	\N	/api/weather	GET	2025-05-16 14:37:54.822728	2025-05-16 14:37:54.822754	200	0
5602	\N	/api/restaurant	GET	2025-05-16 14:37:54.974614	2025-05-16 14:37:54.974694	200	0
5617	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 14:44:42.540885	2025-05-16 14:44:42.540985	200	0
5618	\N	/status	GET	2025-05-16 14:44:45.277935	2025-05-16 14:44:45.277937	200	0
5619	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:44:53.183344	2025-05-16 14:44:53.185159	200	1
5620	\N	/status	GET	2025-05-16 14:44:55.829321	2025-05-16 14:44:55.829323	200	0
5621	\N	/api/statistics/global	GET	2025-05-16 14:44:56.037659	2025-05-16 14:44:56.039183	200	1
5622	\N	/api/statistics/endpoints	GET	2025-05-16 14:44:56.239706	2025-05-16 14:44:56.244587	200	4
5623	\N	/status	GET	2025-05-16 14:45:05.844158	2025-05-16 14:45:05.844159	200	0
5624	\N	/status	GET	2025-05-16 14:45:15.883466	2025-05-16 14:45:15.883468	200	0
5625	\N	/status	GET	2025-05-16 14:45:25.90868	2025-05-16 14:45:25.908682	200	0
5626	\N	/status	GET	2025-05-16 14:45:35.908028	2025-05-16 14:45:35.908029	200	0
5627	\N	/status	GET	2025-05-16 14:45:45.921903	2025-05-16 14:45:45.921904	200	0
5628	\N	/status	GET	2025-05-16 14:45:55.931696	2025-05-16 14:45:55.931697	200	0
5629	\N	/status	GET	2025-05-16 14:46:05.958165	2025-05-16 14:46:05.958166	200	0
5630	\N	/status	GET	2025-05-16 14:46:15.980532	2025-05-16 14:46:15.980534	200	0
5631	\N	/status	GET	2025-05-16 14:46:25.981774	2025-05-16 14:46:25.981776	200	0
5632	\N	/status	GET	2025-05-16 14:46:36.010842	2025-05-16 14:46:36.010843	200	0
5633	\N	/status	GET	2025-05-16 14:46:46.013654	2025-05-16 14:46:46.013655	200	0
5634	\N	/status	GET	2025-05-16 14:46:56.033725	2025-05-16 14:46:56.033727	200	0
5635	\N	/status	GET	2025-05-16 14:47:06.052751	2025-05-16 14:47:06.052752	200	0
5636	\N	/status	GET	2025-05-16 14:47:16.068186	2025-05-16 14:47:16.068188	200	0
5637	\N	/status	GET	2025-05-16 14:47:26.08906	2025-05-16 14:47:26.089061	200	0
5638	\N	/status	GET	2025-05-16 14:47:36.092423	2025-05-16 14:47:36.092425	200	0
5639	\N	/status	GET	2025-05-16 14:47:46.12144	2025-05-16 14:47:46.121441	200	0
5640	\N	/status	GET	2025-05-16 14:47:56.137453	2025-05-16 14:47:56.137454	200	0
5641	\N	/status	GET	2025-05-16 14:48:06.178951	2025-05-16 14:48:06.178952	200	0
5642	\N	/status	GET	2025-05-16 14:48:16.167539	2025-05-16 14:48:16.167541	200	0
5643	\N	/status	GET	2025-05-16 14:48:26.191047	2025-05-16 14:48:26.191048	200	0
5644	\N	/status	GET	2025-05-16 14:48:36.206876	2025-05-16 14:48:36.206878	200	0
5645	\N	/status	GET	2025-05-16 14:48:52.191354	2025-05-16 14:48:52.191355	200	0
5646	\N	/status	GET	2025-05-16 14:48:56.245446	2025-05-16 14:48:56.245447	200	0
5647	\N	/status	GET	2025-05-16 14:49:06.245082	2025-05-16 14:49:06.245083	200	0
5648	\N	/status	GET	2025-05-16 14:49:16.271575	2025-05-16 14:49:16.271576	200	0
5649	\N	/status	GET	2025-05-16 14:49:26.302247	2025-05-16 14:49:26.302249	200	0
5650	\N	/status	GET	2025-05-16 14:49:36.302541	2025-05-16 14:49:36.302542	200	0
5651	\N	/status	GET	2025-05-16 14:49:46.331534	2025-05-16 14:49:46.331536	200	0
5652	\N	/status	GET	2025-05-16 14:49:56.331651	2025-05-16 14:49:56.331653	200	0
5653	\N	/status	GET	2025-05-16 14:50:06.359671	2025-05-16 14:50:06.359673	200	0
5654	\N	/status	GET	2025-05-16 14:50:16.36822	2025-05-16 14:50:16.368221	200	0
5655	\N	/status	GET	2025-05-16 14:50:26.385747	2025-05-16 14:50:26.385749	200	0
5656	\N	/status	GET	2025-05-16 14:50:36.402709	2025-05-16 14:50:36.402711	200	0
5657	\N	/status	GET	2025-05-16 14:50:46.418703	2025-05-16 14:50:46.418705	200	0
5658	\N	/status	GET	2025-05-16 14:50:56.43096	2025-05-16 14:50:56.430961	200	0
5659	\N	/status	GET	2025-05-16 14:51:06.449734	2025-05-16 14:51:06.449735	200	0
5660	\N	/status	GET	2025-05-16 14:51:16.465713	2025-05-16 14:51:16.465715	200	0
5661	\N	/status	GET	2025-05-16 14:51:26.485042	2025-05-16 14:51:26.485043	200	0
5662	\N	/status	GET	2025-05-16 14:51:36.5028	2025-05-16 14:51:36.502801	200	0
5663	\N	/status	GET	2025-05-16 14:51:46.51812	2025-05-16 14:51:46.518121	200	0
5664	\N	/status	GET	2025-05-16 14:51:56.529053	2025-05-16 14:51:56.529054	200	0
5665	\N	/status	GET	2025-05-16 14:52:03.942739	2025-05-16 14:52:03.94274	200	0
5666	\N	/api/statistics/global	GET	2025-05-16 14:52:04.156753	2025-05-16 14:52:04.159767	200	3
5667	\N	/api/statistics/endpoints	GET	2025-05-16 14:52:04.406875	2025-05-16 14:52:04.411944	200	5
5668	\N	/status	GET	2025-05-16 14:52:07.926379	2025-05-16 14:52:07.92638	200	0
5669	\N	/api/statistics/global	GET	2025-05-16 14:52:08.13755	2025-05-16 14:52:08.139187	200	1
5670	\N	/api/statistics/endpoints	GET	2025-05-16 14:52:08.353994	2025-05-16 14:52:08.359825	200	5
5671	\N	/status	GET	2025-05-16 14:52:10.605332	2025-05-16 14:52:10.605333	200	0
5672	\N	/api/statistics/global	GET	2025-05-16 14:52:11.283405	2025-05-16 14:52:11.284925	200	1
5673	\N	/api/statistics/endpoints	GET	2025-05-16 14:52:11.506418	2025-05-16 14:52:11.512667	200	6
5674	\N	/status	GET	2025-05-16 14:52:13.95893	2025-05-16 14:52:13.958931	200	0
5675	\N	/status	GET	2025-05-16 14:52:23.990271	2025-05-16 14:52:23.990272	200	0
5676	\N	/status	GET	2025-05-16 14:52:34.012568	2025-05-16 14:52:34.01257	200	0
5677	\N	/status	GET	2025-05-16 14:52:44.022608	2025-05-16 14:52:44.02261	200	0
5678	\N	/status	GET	2025-05-16 14:52:54.036595	2025-05-16 14:52:54.036597	200	0
5679	\N	/status	GET	2025-05-16 14:53:04.053021	2025-05-16 14:53:04.053023	200	0
5680	\N	/status	GET	2025-05-16 14:53:14.079796	2025-05-16 14:53:14.079797	200	0
5681	\N	/status	GET	2025-05-16 14:53:24.092119	2025-05-16 14:53:24.09212	200	0
5682	\N	/status	GET	2025-05-16 14:53:34.122074	2025-05-16 14:53:34.122076	200	0
5683	\N	/status	GET	2025-05-16 14:53:44.655176	2025-05-16 14:53:44.655177	200	0
5684	\N	/status	GET	2025-05-16 14:54:09.033341	2025-05-16 14:54:09.033343	200	0
5685	\N	/status	GET	2025-05-16 14:54:09.189968	2025-05-16 14:54:09.189969	200	0
5686	\N	/status	GET	2025-05-16 14:54:15.642582	2025-05-16 14:54:15.642583	200	0
5687	\N	/status	GET	2025-05-16 14:54:34.996607	2025-05-16 14:54:34.996608	200	0
5688	\N	/status	GET	2025-05-16 14:54:42.950104	2025-05-16 14:54:42.950106	200	0
5689	\N	/status	GET	2025-05-16 14:54:49.519972	2025-05-16 14:54:49.519973	200	0
5690	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:54:53.627237	2025-05-16 14:54:53.628808	200	1
5691	\N	/api/restaurant	GET	2025-05-16 14:54:53.891111	2025-05-16 14:54:53.891165	200	0
5692	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:54:54.044912	2025-05-16 14:54:54.0464	200	1
5693	\N	/api/weather	GET	2025-05-16 14:54:54.043765	2025-05-16 14:54:54.274189	200	230
5694	\N	/status	GET	2025-05-16 14:54:58.945927	2025-05-16 14:54:58.945931	200	0
5695	\N	/status	GET	2025-05-16 14:55:08.961754	2025-05-16 14:55:08.961756	200	0
5696	\N	/status	GET	2025-05-16 14:55:18.966392	2025-05-16 14:55:18.966393	200	0
5697	\N	/status	GET	2025-05-16 14:55:28.990003	2025-05-16 14:55:28.990004	200	0
5698	\N	/status	GET	2025-05-16 14:55:39.004267	2025-05-16 14:55:39.004269	200	0
5699	\N	/status	GET	2025-05-16 14:55:49.010625	2025-05-16 14:55:49.010627	200	0
5579	\N	/api/statistics/global	GET	2025-05-16 14:31:45.034462	2025-05-16 14:31:45.036048	200	1
5580	\N	/api/statistics/endpoints	GET	2025-05-16 14:31:45.243367	2025-05-16 14:31:45.248654	200	5
5581	\N	/api/auth/login	POST	2025-05-16 14:35:54.749618	2025-05-16 14:35:54.824853	200	75
5582	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:35:55.089839	2025-05-16 14:35:55.091657	200	1
5583	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:35:55.297533	2025-05-16 14:35:55.298975	200	1
5584	\N	/api/restaurant	GET	2025-05-16 14:35:55.933582	2025-05-16 14:35:55.933674	200	0
5592	\N	/api/restaurant	GET	2025-05-16 14:37:02.629077	2025-05-16 14:37:02.629203	200	0
5603	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:54.974764	2025-05-16 14:37:54.975997	200	1
5604	zephyr.dassouli@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:37:56.174684	2025-05-16 14:37:56.176181	200	1
5605	\N	/api/restaurant	GET	2025-05-16 14:37:57.497756	2025-05-16 14:37:57.497808	200	0
5606	\N	/api/restaurant	GET	2025-05-16 14:37:58.594343	2025-05-16 14:37:58.594409	200	0
5607	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 14:42:40.725968	2025-05-16 14:42:40.727832	200	1
5608	\N	/status	GET	2025-05-16 14:43:35.173225	2025-05-16 14:43:35.173226	200	0
5609	\N	/api/statistics/global	GET	2025-05-16 14:43:35.445209	2025-05-16 14:43:35.447436	200	2
5610	\N	/api/statistics/endpoints	GET	2025-05-16 14:43:35.658698	2025-05-16 14:43:35.664493	200	5
5611	\N	/status	GET	2025-05-16 14:43:45.177796	2025-05-16 14:43:45.177797	200	0
5612	\N	/status	GET	2025-05-16 14:43:55.204518	2025-05-16 14:43:55.20452	200	0
5613	\N	/status	GET	2025-05-16 14:44:05.22314	2025-05-16 14:44:05.223142	200	0
5614	\N	/status	GET	2025-05-16 14:44:15.231087	2025-05-16 14:44:15.231088	200	0
5615	\N	/status	GET	2025-05-16 14:44:25.238757	2025-05-16 14:44:25.238759	200	0
5616	\N	/status	GET	2025-05-16 14:44:35.259873	2025-05-16 14:44:35.259875	200	0
8671	\N	/api/restaurant	GET	2025-05-21 20:43:08.033041	2025-05-21 20:43:08.033216	200	0
8672	\N	/api/restaurant	GET	2025-05-21 20:45:28.992635	2025-05-21 20:45:28.99276	200	0
8673	\N	/api/restaurant	GET	2025-05-21 20:46:10.2909	2025-05-21 20:46:10.291033	200	0
8674	\N	/api/restaurant	GET	2025-05-21 20:46:54.925757	2025-05-21 20:46:54.925894	200	0
8675	\N	/api/restaurant	GET	2025-05-21 20:47:04.66888	2025-05-21 20:47:04.668951	200	0
8676	\N	/api/restaurant	GET	2025-05-21 20:47:16.836669	2025-05-21 20:47:16.836753	200	0
8677	\N	/api/restaurant	GET	2025-05-21 20:47:30.982107	2025-05-21 20:47:30.982319	200	0
8678	\N	/api/restaurant	GET	2025-05-21 20:47:43.55006	2025-05-21 20:47:43.550136	200	0
8679	\N	/api/restaurant	GET	2025-05-21 20:47:48.35505	2025-05-21 20:47:48.355136	200	0
8680	\N	/api/restaurant	GET	2025-05-21 20:47:53.624834	2025-05-21 20:47:53.624908	200	0
8681	\N	/api/restaurant	GET	2025-05-21 20:48:12.786278	2025-05-21 20:48:12.786412	200	0
8682	\N	/api/restaurant	GET	2025-05-21 20:48:54.337141	2025-05-21 20:48:54.337226	200	0
8683	\N	/api/restaurant	GET	2025-05-21 20:49:38.507192	2025-05-21 20:49:38.507324	200	0
8684	\N	/api/restaurant	GET	2025-05-21 20:51:10.479589	2025-05-21 20:51:10.479732	200	0
8685	\N	/api/restaurant	GET	2025-05-21 20:51:20.826557	2025-05-21 20:51:20.826631	200	0
8686	\N	/api/restaurant	GET	2025-05-21 20:51:22.040324	2025-05-21 20:51:22.040434	200	0
8687	\N	/api/restaurant	GET	2025-05-21 20:51:23.824411	2025-05-21 20:51:23.824503	200	0
8688	\N	/api/restaurant	GET	2025-05-21 20:51:40.404747	2025-05-21 20:51:40.404826	200	0
8689	\N	/api/restaurant	GET	2025-05-21 20:51:50.528706	2025-05-21 20:51:50.528788	200	0
8690	\N	/api/restaurant	GET	2025-05-21 20:51:52.979694	2025-05-21 20:51:52.979781	200	0
8691	\N	/api/restaurant	GET	2025-05-21 20:52:17.295508	2025-05-21 20:52:17.295707	200	0
8692	\N	/api/restaurant	GET	2025-05-21 20:52:19.687414	2025-05-21 20:52:19.687515	200	0
8693	\N	/api/restaurant	GET	2025-05-21 20:52:33.656443	2025-05-21 20:52:33.656527	200	0
8694	\N	/api/restaurant	GET	2025-05-21 20:52:36.573594	2025-05-21 20:52:36.573706	200	0
8695	\N	/api/restaurant	GET	2025-05-21 20:53:01.251072	2025-05-21 20:53:01.251199	200	0
8696	\N	/api/restaurant	GET	2025-05-21 20:53:08.430194	2025-05-21 20:53:08.430277	200	0
8697	\N	/api/restaurant	GET	2025-05-21 20:53:18.284674	2025-05-21 20:53:18.284755	200	0
8698	\N	/api/restaurant	GET	2025-05-21 20:53:21.7079	2025-05-21 20:53:21.707994	200	0
8699	\N	/api/restaurant	GET	2025-05-21 20:53:24.102839	2025-05-21 20:53:24.102925	200	0
8700	\N	/api/restaurant	GET	2025-05-21 20:53:27.599697	2025-05-21 20:53:27.599844	200	0
8701	\N	/api/restaurant	GET	2025-05-21 20:53:35.63666	2025-05-21 20:53:35.636737	200	0
8702	\N	/api/restaurant	GET	2025-05-21 20:53:39.053613	2025-05-21 20:53:39.053704	200	0
8703	\N	/api/restaurant	GET	2025-05-21 20:54:22.140092	2025-05-21 20:54:22.140222	200	0
8704	\N	/api/weather	GET	2025-05-21 20:54:22.145987	2025-05-21 20:54:22.194906	200	48
8705	\N	/api/restaurant	GET	2025-05-21 20:54:51.093671	2025-05-21 20:54:51.093768	200	0
8706	\N	/api/restaurant	GET	2025-05-21 20:54:58.552768	2025-05-21 20:54:58.552847	200	0
8707	\N	/api/restaurant	GET	2025-05-21 20:55:00.162158	2025-05-21 20:55:00.162235	200	0
8708	\N	/api/washingmachines	GET	2025-05-21 20:55:29.188497	2025-05-21 20:55:29.299005	200	110
8709	\N	/api/restaurant	GET	2025-05-21 21:00:53.464354	2025-05-21 21:00:53.464478	200	0
8710	\N	/api/restaurant	GET	2025-05-21 21:01:01.996253	2025-05-21 21:01:01.996416	200	0
8711	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 21:01:24.384849	2025-05-21 21:01:24.386679	200	1
8712	\N	/api/washingmachines	GET	2025-05-21 21:01:24.568899	2025-05-21 21:01:24.642449	200	73
8717	\N	/api/weather	GET	2025-05-21 21:02:03.80659	2025-05-21 21:02:03.806619	200	0
8718	\N	/api/washingmachines	GET	2025-05-21 21:02:03.806205	2025-05-21 21:02:03.836316	200	30
8719	\N	/api/restaurant	GET	2025-05-21 21:02:06.720867	2025-05-21 21:02:06.721106	200	0
8720	\N	/api/restaurant	GET	2025-05-21 21:02:09.644396	2025-05-21 21:02:09.644495	200	0
8721	\N	/api/restaurant	GET	2025-05-21 21:02:14.961918	2025-05-21 21:02:14.962013	200	0
8722	\N	/api/restaurant	GET	2025-05-21 21:02:19.667135	2025-05-21 21:02:19.667214	200	0
8723	\N	/api/restaurant	GET	2025-05-21 21:04:15.449787	2025-05-21 21:04:15.449927	200	0
8724	\N	/api/restaurant	GET	2025-05-21 21:04:21.656109	2025-05-21 21:04:21.656218	200	0
8725	\N	/api/restaurant	GET	2025-05-21 21:04:32.957957	2025-05-21 21:04:32.95806	200	0
8726	\N	/api/restaurant	GET	2025-05-21 21:04:35.459487	2025-05-21 21:04:35.459565	200	0
8727	\N	/api/restaurant	GET	2025-05-21 21:05:38.847081	2025-05-21 21:05:38.84722	200	0
8728	\N	/api/restaurant	GET	2025-05-21 21:09:41.386386	2025-05-21 21:09:41.386494	200	0
8729	\N	/api/restaurant	GET	2025-05-21 21:12:54.735696	2025-05-21 21:12:54.735801	200	0
8730	\N	/api/weather	GET	2025-05-21 21:12:54.73544	2025-05-21 21:12:54.799126	200	63
8731	\N	/api/weather	GET	2025-05-21 21:13:01.95283	2025-05-21 21:13:01.952855	200	0
8732	\N	/api/restaurant	GET	2025-05-21 21:14:04.04646	2025-05-21 21:14:04.046531	200	0
8733	\N	/api/restaurant	GET	2025-05-21 21:14:48.266158	2025-05-21 21:14:48.266329	200	0
8734	\N	/api/restaurant	GET	2025-05-21 21:15:46.410311	2025-05-21 21:15:46.410395	200	0
8735	\N	/api/restaurant	GET	2025-05-21 21:15:53.478567	2025-05-21 21:15:53.479077	200	0
8736	\N	/api/restaurant	GET	2025-05-21 21:16:21.051346	2025-05-21 21:16:21.051428	200	0
8737	\N	/api/restaurant	GET	2025-05-21 21:16:36.369983	2025-05-21 21:16:36.370064	200	0
8738	\N	/api/restaurant	GET	2025-05-21 21:16:42.768227	2025-05-21 21:16:42.768328	200	0
8739	\N	/api/restaurant	GET	2025-05-21 21:18:00.418095	2025-05-21 21:18:00.418229	200	0
8740	\N	/api/restaurant	GET	2025-05-21 21:18:06.770912	2025-05-21 21:18:06.771019	200	0
8741	\N	/api/restaurant	GET	2025-05-21 21:20:26.719534	2025-05-21 21:20:26.719647	200	0
8742	\N	/api/restaurant	GET	2025-05-21 21:20:44.542363	2025-05-21 21:20:44.542471	200	0
8743	\N	/api/restaurant	GET	2025-05-21 21:21:28.581601	2025-05-21 21:21:28.581691	200	0
8744	\N	/api/restaurant	GET	2025-05-21 21:21:54.44098	2025-05-21 21:21:54.44113	200	0
8745	\N	/api/restaurant	GET	2025-05-21 21:22:44.08674	2025-05-21 21:22:44.086825	200	0
8746	\N	/api/restaurant	GET	2025-05-21 21:22:46.288658	2025-05-21 21:22:46.28877	200	0
8747	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 21:23:03.07728	2025-05-21 21:23:03.079244	200	1
8748	test@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-21 21:23:05.521856	2025-05-21 21:23:05.52414	200	2
8749	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 21:23:09.439804	2025-05-21 21:23:09.445731	200	5
8750	\N	/api/washingmachines	GET	2025-05-21 21:23:42.030042	2025-05-21 21:23:42.13339	200	103
8751	\N	/api/washingmachines	GET	2025-05-21 21:23:58.289837	2025-05-21 21:23:58.327487	200	37
8752	\N	/api/washingmachines	GET	2025-05-21 21:24:11.534028	2025-05-21 21:24:11.571903	200	37
8753	\N	/api/washingmachines	GET	2025-05-21 21:24:30.047495	2025-05-21 21:24:30.078904	200	31
8754	\N	/api/washingmachines	GET	2025-05-21 21:24:40.509858	2025-05-21 21:24:40.550389	200	40
8755	\N	/api/washingmachines	GET	2025-05-21 21:24:59.451321	2025-05-21 21:24:59.481332	200	30
5701	\N	/status	GET	2025-05-16 14:56:09.055221	2025-05-16 14:56:09.055222	200	0
5702	\N	/status	GET	2025-05-16 14:56:19.079785	2025-05-16 14:56:19.079787	200	0
5703	\N	/status	GET	2025-05-16 14:56:29.078623	2025-05-16 14:56:29.078624	200	0
5704	\N	/status	GET	2025-05-16 14:56:39.103969	2025-05-16 14:56:39.103971	200	0
5705	\N	/status	GET	2025-05-16 14:56:49.117122	2025-05-16 14:56:49.117123	200	0
5706	\N	/status	GET	2025-05-16 14:56:59.131218	2025-05-16 14:56:59.13122	200	0
5707	\N	/status	GET	2025-05-16 14:57:09.145947	2025-05-16 14:57:09.145949	200	0
5708	\N	/status	GET	2025-05-16 14:57:19.164713	2025-05-16 14:57:19.164715	200	0
5709	\N	/status	GET	2025-05-16 14:57:29.173671	2025-05-16 14:57:29.173673	200	0
5710	\N	/status	GET	2025-05-16 14:57:31.80463	2025-05-16 14:57:31.804632	200	0
5711	\N	/api/statistics/global	GET	2025-05-16 14:57:32.028027	2025-05-16 14:57:32.029705	200	1
5712	\N	/api/statistics/top-users	GET	2025-05-16 14:57:32.242472	2025-05-16 14:57:32.242478	500	0
5713	\N	/status	GET	2025-05-16 14:57:41.829342	2025-05-16 14:57:41.829344	200	0
5714	\N	/status	GET	2025-05-16 14:57:51.843988	2025-05-16 14:57:51.84399	200	0
5715	\N	/status	GET	2025-05-16 14:58:01.87221	2025-05-16 14:58:01.872212	200	0
5716	\N	/status	GET	2025-05-16 14:58:11.880694	2025-05-16 14:58:11.880696	200	0
5717	\N	/status	GET	2025-05-16 14:58:21.921203	2025-05-16 14:58:21.921206	200	0
5718	\N	/status	GET	2025-05-16 14:58:31.921358	2025-05-16 14:58:31.92136	200	0
5719	\N	/status	GET	2025-05-16 14:58:41.947972	2025-05-16 14:58:41.947973	200	0
5720	\N	/status	GET	2025-05-16 14:58:52.027152	2025-05-16 14:58:52.027153	200	0
5721	\N	/status	GET	2025-05-16 14:59:01.981638	2025-05-16 14:59:01.981639	200	0
5722	\N	/status	GET	2025-05-16 14:59:12.000781	2025-05-16 14:59:12.000783	200	0
5723	\N	/status	GET	2025-05-16 14:59:22.008844	2025-05-16 14:59:22.008846	200	0
5724	\N	/status	GET	2025-05-16 14:59:32.050958	2025-05-16 14:59:32.050959	200	0
5725	\N	/status	GET	2025-05-16 14:59:42.057078	2025-05-16 14:59:42.057079	200	0
5726	\N	/status	GET	2025-05-16 14:59:52.075014	2025-05-16 14:59:52.075016	200	0
5727	\N	/status	GET	2025-05-16 15:00:02.076655	2025-05-16 15:00:02.076657	200	0
5728	\N	/status	GET	2025-05-16 15:00:12.099836	2025-05-16 15:00:12.099838	200	0
5729	\N	/status	GET	2025-05-16 15:00:22.139323	2025-05-16 15:00:22.139324	200	0
5730	\N	/status	GET	2025-05-16 15:00:32.140787	2025-05-16 15:00:32.140788	200	0
5731	\N	/status	GET	2025-05-16 15:00:42.159719	2025-05-16 15:00:42.15972	200	0
5732	\N	/status	GET	2025-05-16 15:00:52.160388	2025-05-16 15:00:52.16039	200	0
5733	\N	/status	GET	2025-05-16 15:01:02.201878	2025-05-16 15:01:02.20188	200	0
5734	\N	/status	GET	2025-05-16 15:01:12.191623	2025-05-16 15:01:12.191625	200	0
5735	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 15:01:14.901294	2025-05-16 15:01:14.903125	200	1
5736	\N	/status	GET	2025-05-16 15:01:20.187423	2025-05-16 15:01:20.187424	200	0
5737	\N	/api/statistics/global	GET	2025-05-16 15:01:20.389285	2025-05-16 15:01:20.391101	200	1
5738	\N	/api/statistics/top-users	GET	2025-05-16 15:01:20.589109	2025-05-16 15:01:20.589116	500	0
5739	\N	/api/statistics/top-users	GET	2025-05-16 17:01:51.270779	2025-05-16 17:01:51.45917	200	188
5740	\N	/status	GET	2025-05-16 15:02:50.878689	2025-05-16 15:02:50.87869	200	0
5741	\N	/api/statistics/global	GET	2025-05-16 15:02:51.256051	2025-05-16 15:02:51.258237	200	2
5742	\N	/api/statistics/top-users	GET	2025-05-16 15:02:51.496372	2025-05-16 15:02:51.496378	500	0
5743	\N	/status	GET	2025-05-16 15:03:00.879564	2025-05-16 15:03:00.879566	200	0
5744	\N	/status	GET	2025-05-16 15:03:10.90612	2025-05-16 15:03:10.906122	200	0
5745	\N	/status	GET	2025-05-16 15:03:20.91893	2025-05-16 15:03:20.918931	200	0
5746	\N	/status	GET	2025-05-16 15:03:41.262816	2025-05-16 15:03:41.262819	200	0
5747	\N	/status	GET	2025-05-16 15:03:50.964497	2025-05-16 15:03:50.964498	200	0
5748	\N	/status	GET	2025-05-16 15:04:00.981963	2025-05-16 15:04:00.981965	200	0
5749	\N	/status	GET	2025-05-16 15:04:10.998135	2025-05-16 15:04:10.998136	200	0
5750	\N	/status	GET	2025-05-16 15:04:12.895893	2025-05-16 15:04:12.895895	200	0
5751	\N	/api/statistics/global	GET	2025-05-16 15:04:13.108844	2025-05-16 15:04:13.110959	200	2
5752	\N	/api/statistics/top-users	GET	2025-05-16 15:04:13.325664	2025-05-16 15:04:13.327132	200	1
5753	\N	/status	GET	2025-05-16 15:04:21.022128	2025-05-16 15:04:21.022129	200	0
5754	\N	/status	GET	2025-05-16 15:04:31.039703	2025-05-16 15:04:31.039704	200	0
5755	\N	/status	GET	2025-05-16 15:04:41.044773	2025-05-16 15:04:41.044774	200	0
5756	\N	/status	GET	2025-05-16 15:04:51.07466	2025-05-16 15:04:51.074662	200	0
5757	\N	/status	GET	2025-05-16 15:05:01.080225	2025-05-16 15:05:01.080226	200	0
5758	\N	/status	GET	2025-05-16 15:05:11.097645	2025-05-16 15:05:11.097658	200	0
5759	\N	/status	GET	2025-05-16 15:05:21.120979	2025-05-16 15:05:21.12098	200	0
5760	\N	/status	GET	2025-05-16 15:05:31.128313	2025-05-16 15:05:31.128314	200	0
5761	\N	/status	GET	2025-05-16 15:05:41.15994	2025-05-16 15:05:41.159942	200	0
5762	\N	/status	GET	2025-05-16 15:05:51.170383	2025-05-16 15:05:51.170385	200	0
5763	\N	/status	GET	2025-05-16 15:06:01.188873	2025-05-16 15:06:01.188875	200	0
5764	\N	/status	GET	2025-05-16 15:06:11.211532	2025-05-16 15:06:11.211533	200	0
5765	\N	/status	GET	2025-05-16 15:06:21.221045	2025-05-16 15:06:21.221047	200	0
5766	\N	/status	GET	2025-05-16 15:06:31.231914	2025-05-16 15:06:31.231916	200	0
5767	\N	/status	GET	2025-05-16 15:06:41.260516	2025-05-16 15:06:41.260518	200	0
5768	\N	/status	GET	2025-05-16 15:06:51.259285	2025-05-16 15:06:51.259287	200	0
5769	\N	/status	GET	2025-05-16 15:07:01.290675	2025-05-16 15:07:01.290676	200	0
5770	\N	/status	GET	2025-05-16 15:07:11.31424	2025-05-16 15:07:11.314241	200	0
5771	\N	/status	GET	2025-05-16 15:07:21.312492	2025-05-16 15:07:21.312493	200	0
5772	\N	/status	GET	2025-05-16 15:07:31.328016	2025-05-16 15:07:31.328017	200	0
5773	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-16 15:07:39.514214	2025-05-16 15:07:39.515316	200	1
5774	\N	/status	GET	2025-05-16 15:07:41.357063	2025-05-16 15:07:41.357065	200	0
5775	\N	/status	GET	2025-05-16 15:07:51.40033	2025-05-16 15:07:51.400332	200	0
5776	\N	/status	GET	2025-05-16 15:08:01.37717	2025-05-16 15:08:01.377172	200	0
5777	\N	/status	GET	2025-05-16 15:08:11.394977	2025-05-16 15:08:11.394979	200	0
5778	\N	/status	GET	2025-05-16 15:08:21.418786	2025-05-16 15:08:21.418788	200	0
5779	\N	/status	GET	2025-05-16 15:08:31.457238	2025-05-16 15:08:31.45724	200	0
5780	\N	/status	GET	2025-05-16 15:08:41.445237	2025-05-16 15:08:41.445239	200	0
5781	\N	/status	GET	2025-05-16 15:08:51.462945	2025-05-16 15:08:51.462946	200	0
5782	\N	/status	GET	2025-05-16 15:09:01.479934	2025-05-16 15:09:01.479935	200	0
5783	\N	/status	GET	2025-05-16 15:09:11.503321	2025-05-16 15:09:11.503322	200	0
5784	\N	/status	GET	2025-05-16 15:09:21.510142	2025-05-16 15:09:21.510143	200	0
5785	\N	/status	GET	2025-05-16 15:09:31.526798	2025-05-16 15:09:31.526799	200	0
5786	\N	/status	GET	2025-05-16 15:09:41.565454	2025-05-16 15:09:41.565456	200	0
5787	\N	/status	GET	2025-05-16 15:09:51.568391	2025-05-16 15:09:51.568392	200	0
5788	\N	/status	GET	2025-05-16 15:09:56.903384	2025-05-16 15:09:56.903386	200	0
5789	\N	/api/statistics/global	GET	2025-05-16 15:09:57.130372	2025-05-16 15:09:57.13276	200	2
5790	\N	/api/statistics/top-users	GET	2025-05-16 15:09:57.34728	2025-05-16 15:09:57.349033	200	1
5791	\N	/status	GET	2025-05-16 15:10:06.910767	2025-05-16 15:10:06.910768	200	0
5792	\N	/status	GET	2025-05-16 15:10:16.930236	2025-05-16 15:10:16.930237	200	0
5793	\N	/status	GET	2025-05-16 15:10:26.953736	2025-05-16 15:10:26.953737	200	0
5794	\N	/status	GET	2025-05-16 15:10:36.984806	2025-05-16 15:10:36.984808	200	0
5795	\N	/status	GET	2025-05-16 15:10:46.987448	2025-05-16 15:10:46.98745	200	0
5796	\N	/status	GET	2025-05-16 15:10:57.000929	2025-05-16 15:10:57.000931	200	0
5797	\N	/status	GET	2025-05-16 15:11:07.015964	2025-05-16 15:11:07.015966	200	0
5798	\N	/status	GET	2025-05-16 15:11:17.029088	2025-05-16 15:11:17.02909	200	0
5799	\N	/status	GET	2025-05-16 15:11:27.045805	2025-05-16 15:11:27.045807	200	0
5800	\N	/status	GET	2025-05-16 15:11:37.057214	2025-05-16 15:11:37.057216	200	0
5801	\N	/status	GET	2025-05-16 15:11:47.079152	2025-05-16 15:11:47.079154	200	0
5802	\N	/status	GET	2025-05-16 15:11:57.101699	2025-05-16 15:11:57.101701	200	0
5803	\N	/status	GET	2025-05-16 15:12:07.151845	2025-05-16 15:12:07.151847	200	0
5804	\N	/status	GET	2025-05-16 15:12:17.137435	2025-05-16 15:12:17.137437	200	0
5805	\N	/status	GET	2025-05-16 15:12:27.152798	2025-05-16 15:12:27.152799	200	0
5806	\N	/status	GET	2025-05-16 15:12:37.168618	2025-05-16 15:12:37.16862	200	0
5807	\N	/status	GET	2025-05-16 15:12:47.178257	2025-05-16 15:12:47.178259	200	0
5808	\N	/status	GET	2025-05-16 15:12:57.193614	2025-05-16 15:12:57.193615	200	0
5809	\N	/status	GET	2025-05-16 15:13:07.219383	2025-05-16 15:13:07.219384	200	0
5810	\N	/status	GET	2025-05-16 15:13:17.24653	2025-05-16 15:13:17.246531	200	0
5811	\N	/status	GET	2025-05-16 15:13:27.246046	2025-05-16 15:13:27.246048	200	0
5812	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 15:13:30.136089	2025-05-16 15:13:30.138531	200	2
5813	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 15:13:30.562184	2025-05-16 15:13:30.563684	200	1
5814	\N	/api/weather	GET	2025-05-16 15:13:30.400043	2025-05-16 15:13:30.634665	200	234
5815	\N	/api/restaurant	GET	2025-05-16 15:13:30.562412	2025-05-16 15:13:31.346517	200	784
5816	\N	/status	GET	2025-05-16 15:13:37.283115	2025-05-16 15:13:37.283117	200	0
5817	\N	/status	GET	2025-05-16 15:13:47.283354	2025-05-16 15:13:47.283356	200	0
5818	\N	/status	GET	2025-05-16 15:13:57.304622	2025-05-16 15:13:57.304623	200	0
5819	\N	/status	GET	2025-05-16 15:14:07.318189	2025-05-16 15:14:07.31819	200	0
5820	\N	/status	GET	2025-05-16 15:14:17.331541	2025-05-16 15:14:17.331542	200	0
5821	\N	/status	GET	2025-05-16 15:14:27.344033	2025-05-16 15:14:27.344035	200	0
5822	\N	/status	GET	2025-05-16 15:14:37.368575	2025-05-16 15:14:37.368577	200	0
5823	\N	/status	GET	2025-05-16 15:14:47.380973	2025-05-16 15:14:47.380975	200	0
5824	\N	/status	GET	2025-05-16 15:14:57.405338	2025-05-16 15:14:57.40534	200	0
5825	\N	/status	GET	2025-05-16 15:15:07.415409	2025-05-16 15:15:07.415411	200	0
5826	\N	/status	GET	2025-05-16 15:15:17.429559	2025-05-16 15:15:17.429561	200	0
5827	\N	/status	GET	2025-05-16 15:15:27.452251	2025-05-16 15:15:27.452252	200	0
5828	\N	/status	GET	2025-05-16 15:15:37.461558	2025-05-16 15:15:37.46156	200	0
5829	\N	/status	GET	2025-05-16 15:15:47.489182	2025-05-16 15:15:47.489183	200	0
5830	\N	/status	GET	2025-05-16 15:15:57.495121	2025-05-16 15:15:57.495122	200	0
5831	\N	/status	GET	2025-05-16 15:16:07.511396	2025-05-16 15:16:07.511398	200	0
5832	\N	/status	GET	2025-05-16 15:16:17.556707	2025-05-16 15:16:17.556709	200	0
5833	\N	/status	GET	2025-05-16 15:16:27.541203	2025-05-16 15:16:27.541204	200	0
5834	\N	/status	GET	2025-05-16 15:16:37.567032	2025-05-16 15:16:37.567034	200	0
5835	\N	/status	GET	2025-05-16 15:16:47.584988	2025-05-16 15:16:47.58499	200	0
5836	\N	/status	GET	2025-05-16 15:16:57.607518	2025-05-16 15:16:57.60752	200	0
5837	\N	/status	GET	2025-05-16 15:17:07.61584	2025-05-16 15:17:07.615842	200	0
5838	\N	/status	GET	2025-05-16 15:17:17.633191	2025-05-16 15:17:17.633192	200	0
5839	\N	/status	GET	2025-05-16 15:17:27.642155	2025-05-16 15:17:27.642157	200	0
5840	\N	/status	GET	2025-05-16 15:17:37.688172	2025-05-16 15:17:37.688174	200	0
5841	\N	/status	GET	2025-05-16 15:17:47.681906	2025-05-16 15:17:47.681907	200	0
5842	\N	/status	GET	2025-05-16 15:17:57.708104	2025-05-16 15:17:57.708105	200	0
5843	\N	/status	GET	2025-05-16 15:18:07.729816	2025-05-16 15:18:07.729817	200	0
5844	\N	/status	GET	2025-05-16 15:18:17.731541	2025-05-16 15:18:17.731542	200	0
5845	\N	/status	GET	2025-05-16 15:18:27.763927	2025-05-16 15:18:27.763928	200	0
5846	\N	/status	GET	2025-05-16 15:18:37.770176	2025-05-16 15:18:37.770178	200	0
5847	\N	/status	GET	2025-05-16 15:18:47.783054	2025-05-16 15:18:47.783056	200	0
5848	\N	/status	GET	2025-05-16 15:18:57.801668	2025-05-16 15:18:57.801669	200	0
5849	\N	/status	GET	2025-05-16 15:19:07.818446	2025-05-16 15:19:07.818447	200	0
5850	\N	/status	GET	2025-05-16 15:19:17.830872	2025-05-16 15:19:17.830873	200	0
5851	\N	/status	GET	2025-05-16 15:19:27.851256	2025-05-16 15:19:27.851257	200	0
5852	\N	/status	GET	2025-05-16 15:19:37.878469	2025-05-16 15:19:37.878471	200	0
5853	\N	/status	GET	2025-05-16 15:19:47.88666	2025-05-16 15:19:47.886662	200	0
5854	\N	/status	GET	2025-05-16 15:19:57.90459	2025-05-16 15:19:57.904591	200	0
5855	\N	/status	GET	2025-05-16 15:20:07.927624	2025-05-16 15:20:07.927625	200	0
5856	\N	/status	GET	2025-05-16 15:20:17.934039	2025-05-16 15:20:17.934041	200	0
5857	\N	/status	GET	2025-05-16 15:20:27.948464	2025-05-16 15:20:27.948466	200	0
5858	\N	/status	GET	2025-05-16 15:20:37.965889	2025-05-16 15:20:37.96589	200	0
5859	\N	/status	GET	2025-05-16 15:20:47.983424	2025-05-16 15:20:47.983426	200	0
5860	\N	/status	GET	2025-05-16 15:20:58.001459	2025-05-16 15:20:58.001461	200	0
5861	\N	/status	GET	2025-05-16 15:21:08.017675	2025-05-16 15:21:08.017676	200	0
5862	\N	/status	GET	2025-05-16 15:21:18.029139	2025-05-16 15:21:18.029141	200	0
5863	\N	/status	GET	2025-05-16 15:21:28.048004	2025-05-16 15:21:28.048005	200	0
5864	\N	/status	GET	2025-05-16 15:21:38.069885	2025-05-16 15:21:38.069886	200	0
5865	\N	/status	GET	2025-05-16 15:21:48.092729	2025-05-16 15:21:48.09273	200	0
5866	\N	/status	GET	2025-05-16 15:21:58.096878	2025-05-16 15:21:58.09688	200	0
5867	\N	/status	GET	2025-05-16 15:22:08.122124	2025-05-16 15:22:08.122126	200	0
5868	\N	/status	GET	2025-05-16 15:22:18.128129	2025-05-16 15:22:18.12813	200	0
5869	\N	/status	GET	2025-05-16 15:22:28.152771	2025-05-16 15:22:28.152773	200	0
5870	\N	/status	GET	2025-05-16 15:22:38.167689	2025-05-16 15:22:38.167691	200	0
5871	\N	/status	GET	2025-05-16 15:22:48.177108	2025-05-16 15:22:48.177109	200	0
5872	\N	/status	GET	2025-05-16 15:22:58.234897	2025-05-16 15:22:58.234898	200	0
5873	\N	/status	GET	2025-05-16 15:23:08.222531	2025-05-16 15:23:08.222533	200	0
5874	\N	/status	GET	2025-05-16 15:23:18.241616	2025-05-16 15:23:18.241618	200	0
5875	\N	/status	GET	2025-05-16 15:23:28.252012	2025-05-16 15:23:28.252014	200	0
5876	\N	/status	GET	2025-05-16 15:23:38.282403	2025-05-16 15:23:38.282404	200	0
5877	\N	/status	GET	2025-05-16 15:23:48.280856	2025-05-16 15:23:48.280857	200	0
5878	\N	/status	GET	2025-05-16 15:23:58.296203	2025-05-16 15:23:58.296204	200	0
5879	\N	/status	GET	2025-05-16 15:24:08.324138	2025-05-16 15:24:08.32414	200	0
5880	\N	/status	GET	2025-05-16 15:24:18.330013	2025-05-16 15:24:18.330015	200	0
5881	\N	/status	GET	2025-05-16 15:24:28.346234	2025-05-16 15:24:28.346236	200	0
5882	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 15:24:40.270816	2025-05-16 15:24:40.272493	200	1
5883	\N	/status	GET	2025-05-16 15:24:41.967154	2025-05-16 15:24:41.967155	200	0
5884	\N	/api/statistics/global	GET	2025-05-16 15:24:42.17747	2025-05-16 15:24:42.179325	200	1
5885	\N	/api/statistics/top-users	GET	2025-05-16 15:24:42.390301	2025-05-16 15:24:42.391547	200	1
5886	\N	/status	GET	2025-05-16 15:24:51.977969	2025-05-16 15:24:51.97797	200	0
5887	\N	/status	GET	2025-05-16 15:25:01.993188	2025-05-16 15:25:01.99319	200	0
5888	\N	/status	GET	2025-05-16 15:25:12.016181	2025-05-16 15:25:12.016183	200	0
5889	\N	/status	GET	2025-05-16 15:25:22.032875	2025-05-16 15:25:22.032877	200	0
5890	\N	/status	GET	2025-05-16 15:25:32.051656	2025-05-16 15:25:32.051658	200	0
5891	\N	/status	GET	2025-05-16 15:25:42.071674	2025-05-16 15:25:42.071676	200	0
5892	\N	/status	GET	2025-05-16 15:25:52.096379	2025-05-16 15:25:52.096381	200	0
5893	\N	/status	GET	2025-05-16 15:26:02.098881	2025-05-16 15:26:02.098882	200	0
5894	\N	/status	GET	2025-05-16 15:26:12.128022	2025-05-16 15:26:12.128024	200	0
5895	\N	/status	GET	2025-05-16 15:26:22.169122	2025-05-16 15:26:22.169123	200	0
5896	\N	/status	GET	2025-05-16 15:26:32.150543	2025-05-16 15:26:32.150544	200	0
5897	\N	/status	GET	2025-05-16 15:26:42.160508	2025-05-16 15:26:42.16051	200	0
5898	\N	/status	GET	2025-05-16 15:26:52.180914	2025-05-16 15:26:52.180916	200	0
5899	\N	/status	GET	2025-05-16 15:27:02.198501	2025-05-16 15:27:02.198503	200	0
5900	\N	/status	GET	2025-05-16 15:27:12.211453	2025-05-16 15:27:12.211455	200	0
5901	\N	/status	GET	2025-05-16 15:27:22.236664	2025-05-16 15:27:22.236666	200	0
5902	\N	/status	GET	2025-05-16 15:27:32.249473	2025-05-16 15:27:32.249475	200	0
5903	\N	/status	GET	2025-05-16 15:27:42.267551	2025-05-16 15:27:42.267553	200	0
5904	\N	/status	GET	2025-05-16 15:27:52.281359	2025-05-16 15:27:52.281361	200	0
5905	\N	/status	GET	2025-05-16 15:28:02.297993	2025-05-16 15:28:02.297995	200	0
5906	\N	/status	GET	2025-05-16 15:28:12.328943	2025-05-16 15:28:12.328945	200	0
5907	\N	/status	GET	2025-05-16 15:28:22.325822	2025-05-16 15:28:22.325823	200	0
5908	\N	/status	GET	2025-05-16 15:28:32.347832	2025-05-16 15:28:32.347834	200	0
5909	\N	/status	GET	2025-05-16 15:28:42.369832	2025-05-16 15:28:42.369834	200	0
5910	\N	/status	GET	2025-05-16 15:28:52.386577	2025-05-16 15:28:52.386579	200	0
5911	\N	/status	GET	2025-05-16 15:29:02.402533	2025-05-16 15:29:02.402534	200	0
5912	\N	/status	GET	2025-05-16 15:29:12.424857	2025-05-16 15:29:12.424859	200	0
5913	\N	/status	GET	2025-05-16 15:29:22.430181	2025-05-16 15:29:22.430182	200	0
5914	\N	/status	GET	2025-05-16 15:29:32.453238	2025-05-16 15:29:32.453239	200	0
5915	\N	/status	GET	2025-05-16 15:29:42.464681	2025-05-16 15:29:42.464683	200	0
5916	\N	/status	GET	2025-05-16 15:29:52.516271	2025-05-16 15:29:52.516273	200	0
5917	\N	/status	GET	2025-05-16 15:30:02.497936	2025-05-16 15:30:02.497937	200	0
5918	\N	/status	GET	2025-05-16 15:37:40.091993	2025-05-16 15:37:40.091998	200	0
5919	\N	/status	GET	2025-05-16 15:37:41.30705	2025-05-16 15:37:41.307052	200	0
5920	\N	/api/statistics/global	GET	2025-05-16 15:37:41.462209	2025-05-16 15:37:42.198901	200	736
5921	\N	/api/statistics/endpoints	GET	2025-05-16 15:37:42.679068	2025-05-16 15:37:43.107523	200	428
5922	\N	/	GET	2025-05-16 17:08:56.152604	2025-05-16 17:08:56.15261	500	0
5923	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 17:49:58.756899	2025-05-16 17:49:58.763316	200	6
5924	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-16 17:49:59.066513	2025-05-16 17:49:59.068645	200	2
5925	\N	/api/weather	GET	2025-05-16 17:49:59.061156	2025-05-16 17:49:59.120819	200	59
5926	\N	/api/restaurant	GET	2025-05-16 17:49:59.066433	2025-05-16 17:49:59.662494	200	596
5927	\N	/status	GET	2025-05-16 17:50:09.744045	2025-05-16 17:50:09.744056	200	0
5928	\N	/api/statistics/global	GET	2025-05-16 17:50:09.840436	2025-05-16 17:50:09.843373	200	2
5929	\N	/api/statistics/top-users	GET	2025-05-16 17:50:09.96395	2025-05-16 17:50:09.970267	200	6
5930	\N	/status	GET	2025-05-16 17:50:19.741132	2025-05-16 17:50:19.741134	200	0
5931	\N	/status	GET	2025-05-16 17:50:29.775161	2025-05-16 17:50:29.775162	200	0
5932	\N	/status	GET	2025-05-16 17:50:39.794177	2025-05-16 17:50:39.794179	200	0
5933	\N	/status	GET	2025-05-16 17:50:49.809519	2025-05-16 17:50:49.809521	200	0
5934	\N	/status	GET	2025-05-16 17:50:59.820712	2025-05-16 17:50:59.820714	200	0
5935	\N	/status	GET	2025-05-16 17:51:09.911835	2025-05-16 17:51:09.911837	200	0
5936	\N	/status	GET	2025-05-16 17:51:19.859629	2025-05-16 17:51:19.85963	200	0
5937	\N	/status	GET	2025-05-16 17:51:29.86349	2025-05-16 17:51:29.863492	200	0
5938	\N	/status	GET	2025-05-16 17:51:39.878208	2025-05-16 17:51:39.87821	200	0
5939	\N	/status	GET	2025-05-16 17:51:49.889416	2025-05-16 17:51:49.889418	200	0
5940	\N	/status	GET	2025-05-16 17:51:59.907408	2025-05-16 17:51:59.90741	200	0
5941	\N	/status	GET	2025-05-16 17:52:09.934434	2025-05-16 17:52:09.934436	200	0
5942	\N	/status	GET	2025-05-16 17:52:19.974132	2025-05-16 17:52:19.974134	200	0
5943	\N	/status	GET	2025-05-16 17:52:29.957756	2025-05-16 17:52:29.957758	200	0
5944	\N	/status	GET	2025-05-16 17:52:39.989687	2025-05-16 17:52:39.989689	200	0
5945	\N	/status	GET	2025-05-16 17:52:49.9994	2025-05-16 17:52:49.999402	200	0
5946	\N	/status	GET	2025-05-16 17:53:00.007629	2025-05-16 17:53:00.007631	200	0
5947	\N	/status	GET	2025-05-16 17:53:10.040137	2025-05-16 17:53:10.040139	200	0
5948	\N	/status	GET	2025-05-16 17:53:20.040329	2025-05-16 17:53:20.04033	200	0
5949	\N	/status	GET	2025-05-16 17:53:30.05566	2025-05-16 17:53:30.055662	200	0
5950	\N	/status	GET	2025-05-16 17:53:40.108579	2025-05-16 17:53:40.108581	200	0
5951	\N	/status	GET	2025-05-16 17:53:50.109696	2025-05-16 17:53:50.109698	200	0
5952	\N	/status	GET	2025-05-16 17:54:00.164914	2025-05-16 17:54:00.164916	200	0
5953	\N	/status	GET	2025-05-16 17:54:10.157799	2025-05-16 17:54:10.157802	200	0
5954	\N	/status	GET	2025-05-16 17:54:20.142376	2025-05-16 17:54:20.142378	200	0
5955	\N	/status	GET	2025-05-16 17:54:30.153123	2025-05-16 17:54:30.153125	200	0
5956	\N	/status	GET	2025-05-16 17:54:40.164842	2025-05-16 17:54:40.164844	200	0
5957	\N	/status	GET	2025-05-16 17:54:50.184399	2025-05-16 17:54:50.184401	200	0
5958	\N	/status	GET	2025-05-16 17:55:00.213606	2025-05-16 17:55:00.213608	200	0
5959	\N	/status	GET	2025-05-16 17:55:10.248821	2025-05-16 17:55:10.248823	200	0
5960	\N	/status	GET	2025-05-16 17:55:20.235453	2025-05-16 17:55:20.235454	200	0
5961	\N	/status	GET	2025-05-16 17:55:30.259369	2025-05-16 17:55:30.259371	200	0
5962	\N	/status	GET	2025-05-16 17:55:40.294188	2025-05-16 17:55:40.294189	200	0
5963	\N	/status	GET	2025-05-16 17:55:50.320865	2025-05-16 17:55:50.320867	200	0
5964	\N	/status	GET	2025-05-16 17:56:00.3033	2025-05-16 17:56:00.303302	200	0
5965	\N	/status	GET	2025-05-16 17:56:10.323416	2025-05-16 17:56:10.323417	200	0
5966	\N	/status	GET	2025-05-16 17:56:20.355236	2025-05-16 17:56:20.355238	200	0
5967	\N	/status	GET	2025-05-16 17:56:30.364084	2025-05-16 17:56:30.364086	200	0
5968	\N	/status	GET	2025-05-16 17:56:40.420294	2025-05-16 17:56:40.420295	200	0
5969	\N	/status	GET	2025-05-16 17:56:50.411666	2025-05-16 17:56:50.411667	200	0
5970	\N	/status	GET	2025-05-16 17:57:00.408261	2025-05-16 17:57:00.408263	200	0
5971	\N	/status	GET	2025-05-16 17:57:10.438097	2025-05-16 17:57:10.438099	200	0
5972	\N	/status	GET	2025-05-16 17:57:15.385785	2025-05-16 17:57:15.385787	200	0
5973	\N	/api/statistics/global	GET	2025-05-16 17:57:16.487489	2025-05-16 17:57:16.489945	200	2
5974	\N	/api/statistics/top-users	GET	2025-05-16 17:57:16.723353	2025-05-16 17:57:16.724942	200	1
5975	\N	/status	GET	2025-05-16 17:57:25.408041	2025-05-16 17:57:25.408042	200	0
5976	\N	/status	GET	2025-05-16 17:57:34.485133	2025-05-16 17:57:34.485135	200	0
5977	\N	/api/statistics/global	GET	2025-05-16 17:57:34.607387	2025-05-16 17:57:34.609439	200	2
5978	\N	/api/statistics/top-users	GET	2025-05-16 17:57:34.734715	2025-05-16 17:57:34.736545	200	1
5979	\N	/status	GET	2025-05-16 17:57:39.066345	2025-05-16 17:57:39.066347	200	0
5980	\N	/api/statistics/global	GET	2025-05-16 17:57:39.211349	2025-05-16 17:57:39.213303	200	1
5981	\N	/api/statistics/top-users	GET	2025-05-16 17:57:39.312951	2025-05-16 17:57:39.314273	200	1
5982	\N	/status	GET	2025-05-16 17:57:46.782082	2025-05-16 17:57:46.782084	200	0
5983	\N	/api/statistics/global	GET	2025-05-16 17:57:46.867262	2025-05-16 17:57:46.86969	200	2
5984	\N	/api/statistics/top-users	GET	2025-05-16 17:57:46.962735	2025-05-16 17:57:46.964065	200	1
5985	\N	/status	GET	2025-05-16 17:57:49.065673	2025-05-16 17:57:49.065675	200	0
5986	\N	/status	GET	2025-05-16 17:57:59.091397	2025-05-16 17:57:59.091398	200	0
5987	\N	/status	GET	2025-05-16 17:58:09.122392	2025-05-16 17:58:09.122394	200	0
5988	\N	/status	GET	2025-05-16 17:58:19.14159	2025-05-16 17:58:19.141592	200	0
5989	\N	/status	GET	2025-05-16 17:58:29.162977	2025-05-16 17:58:29.162979	200	0
5990	\N	/status	GET	2025-05-16 17:58:39.173979	2025-05-16 17:58:39.173981	200	0
5991	\N	/status	GET	2025-05-16 17:58:49.17523	2025-05-16 17:58:49.175232	200	0
5992	\N	/status	GET	2025-05-16 17:58:59.195242	2025-05-16 17:58:59.195243	200	0
5993	\N	/status	GET	2025-05-16 17:59:09.206076	2025-05-16 17:59:09.206078	200	0
5994	\N	/status	GET	2025-05-16 17:59:19.235394	2025-05-16 17:59:19.235396	200	0
5995	\N	/status	GET	2025-05-16 17:59:29.259706	2025-05-16 17:59:29.259708	200	0
5996	\N	/status	GET	2025-05-16 17:59:39.271339	2025-05-16 17:59:39.27134	200	0
5997	\N	/status	GET	2025-05-16 17:59:49.280312	2025-05-16 17:59:49.280313	200	0
5998	\N	/status	GET	2025-05-16 17:59:59.297595	2025-05-16 17:59:59.297597	200	0
5999	\N	/status	GET	2025-05-16 18:00:09.314463	2025-05-16 18:00:09.314464	200	0
6000	\N	/status	GET	2025-05-16 18:00:19.348935	2025-05-16 18:00:19.34895	200	0
6001	\N	/status	GET	2025-05-16 18:00:29.354714	2025-05-16 18:00:29.354716	200	0
6002	\N	/status	GET	2025-05-16 18:00:39.344754	2025-05-16 18:00:39.344756	200	0
6003	\N	/status	GET	2025-05-16 18:00:49.53664	2025-05-16 18:00:49.536642	200	0
6004	\N	/status	GET	2025-05-16 18:00:59.403907	2025-05-16 18:00:59.403908	200	0
6005	\N	/status	GET	2025-05-16 18:01:09.395638	2025-05-16 18:01:09.39564	200	0
6006	\N	/status	GET	2025-05-16 18:01:19.415665	2025-05-16 18:01:19.415667	200	0
6007	\N	/status	GET	2025-05-16 18:01:29.446045	2025-05-16 18:01:29.446047	200	0
6008	\N	/status	GET	2025-05-16 18:01:39.46997	2025-05-16 18:01:39.469972	200	0
6009	\N	/status	GET	2025-05-16 18:01:49.48415	2025-05-16 18:01:49.484152	200	0
6010	\N	/status	GET	2025-05-16 18:01:59.506166	2025-05-16 18:01:59.506168	200	0
6011	\N	/status	GET	2025-05-16 18:02:09.511411	2025-05-16 18:02:09.511413	200	0
6012	\N	/status	GET	2025-05-16 18:02:19.523413	2025-05-16 18:02:19.523415	200	0
6013	\N	/status	GET	2025-05-16 18:02:29.546527	2025-05-16 18:02:29.546529	200	0
6014	\N	/status	GET	2025-05-16 18:02:39.582709	2025-05-16 18:02:39.58271	200	0
6015	\N	/status	GET	2025-05-16 18:02:49.570509	2025-05-16 18:02:49.57051	200	0
6016	\N	/status	GET	2025-05-16 18:02:59.61166	2025-05-16 18:02:59.611662	200	0
6017	\N	/status	GET	2025-05-16 18:03:09.608149	2025-05-16 18:03:09.608151	200	0
6018	\N	/status	GET	2025-05-16 18:03:19.62062	2025-05-16 18:03:19.620622	200	0
6019	\N	/status	GET	2025-05-16 18:03:29.632449	2025-05-16 18:03:29.63245	200	0
6020	\N	/status	GET	2025-05-16 18:03:39.673117	2025-05-16 18:03:39.673119	200	0
6021	\N	/status	GET	2025-05-16 18:03:49.664142	2025-05-16 18:03:49.664144	200	0
6022	\N	/status	GET	2025-05-16 18:03:59.686357	2025-05-16 18:03:59.686359	200	0
6023	\N	/status	GET	2025-05-16 18:04:09.726145	2025-05-16 18:04:09.726146	200	0
6024	\N	/status	GET	2025-05-16 18:04:19.720063	2025-05-16 18:04:19.720065	200	0
6025	\N	/status	GET	2025-05-16 18:04:29.746548	2025-05-16 18:04:29.74655	200	0
6026	\N	/status	GET	2025-05-16 18:04:39.772372	2025-05-16 18:04:39.772373	200	0
6027	\N	/status	GET	2025-05-16 18:04:49.773109	2025-05-16 18:04:49.77311	200	0
6028	\N	/status	GET	2025-05-16 18:04:59.799913	2025-05-16 18:04:59.799915	200	0
6029	\N	/status	GET	2025-05-16 18:05:08.430924	2025-05-16 18:05:08.430926	200	0
6030	\N	/api/statistics/global	GET	2025-05-16 18:05:08.532564	2025-05-16 18:05:08.534634	200	2
6031	\N	/api/statistics/top-users	GET	2025-05-16 18:05:08.67644	2025-05-16 18:05:08.67797	200	1
6032	\N	/status	GET	2025-05-16 18:05:09.135744	2025-05-16 18:05:09.135747	200	0
6033	\N	/api/statistics/global	GET	2025-05-16 18:05:09.232074	2025-05-16 18:05:09.234153	200	2
6034	\N	/api/statistics/top-users	GET	2025-05-16 18:05:09.316172	2025-05-16 18:05:09.317569	200	1
6035	\N	/status	GET	2025-05-16 18:05:09.820389	2025-05-16 18:05:09.820391	200	0
6036	\N	/status	GET	2025-05-16 18:05:19.819708	2025-05-16 18:05:19.81971	200	0
6037	\N	/status	GET	2025-05-16 18:05:29.84652	2025-05-16 18:05:29.846521	200	0
6038	\N	/status	GET	2025-05-16 18:05:39.873713	2025-05-16 18:05:39.873715	200	0
6039	\N	/status	GET	2025-05-16 18:05:49.891649	2025-05-16 18:05:49.891651	200	0
6040	\N	/status	GET	2025-05-16 18:05:57.659808	2025-05-16 18:05:57.65981	200	0
6041	\N	/api/statistics/global	GET	2025-05-16 18:05:57.80311	2025-05-16 18:05:57.80559	200	2
6042	\N	/api/statistics/top-users	GET	2025-05-16 18:05:57.902133	2025-05-16 18:05:57.903853	200	1
6043	\N	/status	GET	2025-05-16 18:06:04.564912	2025-05-16 18:06:04.564913	200	0
6044	\N	/api/statistics/global	GET	2025-05-16 18:06:04.702294	2025-05-16 18:06:04.704417	200	2
6045	\N	/api/statistics/top-users	GET	2025-05-16 18:06:04.813334	2025-05-16 18:06:04.815299	200	1
6046	\N	/status	GET	2025-05-16 18:06:14.580656	2025-05-16 18:06:14.580658	200	0
6047	\N	/status	GET	2025-05-16 18:06:15.421702	2025-05-16 18:06:15.421704	200	0
6048	\N	/api/statistics/global	GET	2025-05-16 18:06:15.516061	2025-05-16 18:06:15.518449	200	2
6049	\N	/api/statistics/top-users	GET	2025-05-16 18:06:15.62248	2025-05-16 18:06:15.624231	200	1
6050	\N	/status	GET	2025-05-16 18:06:25.440853	2025-05-16 18:06:25.440855	200	0
6051	\N	/status	GET	2025-05-16 18:06:35.452413	2025-05-16 18:06:35.452415	200	0
6052	\N	/status	GET	2025-05-16 18:06:45.465654	2025-05-16 18:06:45.465656	200	0
6053	\N	/status	GET	2025-05-16 18:06:55.470483	2025-05-16 18:06:55.470484	200	0
6054	\N	/status	GET	2025-05-16 18:07:05.503983	2025-05-16 18:07:05.503985	200	0
6055	\N	/status	GET	2025-05-16 18:07:15.528523	2025-05-16 18:07:15.528525	200	0
6056	\N	/status	GET	2025-05-16 18:07:15.714779	2025-05-16 18:07:15.714781	200	0
6057	\N	/api/statistics/global	GET	2025-05-16 18:07:15.822024	2025-05-16 18:07:15.824316	200	2
6058	\N	/api/statistics/top-users	GET	2025-05-16 18:07:15.959034	2025-05-16 18:07:15.961417	200	2
6059	\N	/status	GET	2025-05-16 18:07:25.692148	2025-05-16 18:07:25.69215	200	0
6060	\N	/status	GET	2025-05-16 18:07:29.197461	2025-05-16 18:07:29.197463	200	0
6061	\N	/api/statistics/global	GET	2025-05-16 18:07:29.407803	2025-05-16 18:07:29.409861	200	2
6062	\N	/api/statistics/top-users	GET	2025-05-16 18:07:29.526023	2025-05-16 18:07:29.527554	200	1
6063	\N	/api/restaurant	GET	2025-05-16 18:07:33.536972	2025-05-16 18:07:33.537028	200	0
6064	\N	/status	GET	2025-05-16 18:07:39.210001	2025-05-16 18:07:39.210003	200	0
6065	\N	/status	GET	2025-05-16 18:07:49.256613	2025-05-16 18:07:49.256615	200	0
6066	\N	/status	GET	2025-05-16 18:07:59.248397	2025-05-16 18:07:59.248399	200	0
6067	\N	/status	GET	2025-05-16 18:08:09.255908	2025-05-16 18:08:09.25591	200	0
6068	\N	/status	GET	2025-05-16 18:08:11.390984	2025-05-16 18:08:11.390986	200	0
6069	\N	/api/statistics/global	GET	2025-05-16 18:08:11.474494	2025-05-16 18:08:11.476663	200	2
6070	\N	/api/statistics/top-users	GET	2025-05-16 18:08:11.592343	2025-05-16 18:08:11.593924	200	1
6071	\N	/status	GET	2025-05-16 18:08:19.274174	2025-05-16 18:08:19.274176	200	0
6072	\N	/status	GET	2025-05-16 18:08:29.287402	2025-05-16 18:08:29.287404	200	0
6073	\N	/status	GET	2025-05-16 18:08:35.881966	2025-05-16 18:08:35.881967	200	0
6074	\N	/api/statistics/global	GET	2025-05-16 18:08:35.963446	2025-05-16 18:08:35.965281	200	1
6075	\N	/api/statistics/top-users	GET	2025-05-16 18:08:36.053533	2025-05-16 18:08:36.054889	200	1
6076	\N	/status	GET	2025-05-16 18:08:45.924186	2025-05-16 18:08:45.924188	200	0
6077	\N	/status	GET	2025-05-16 18:08:55.937521	2025-05-16 18:08:55.937523	200	0
6078	\N	/status	GET	2025-05-16 18:09:06.47952	2025-05-16 18:09:06.479523	200	0
6079	\N	/status	GET	2025-05-16 18:09:16.053411	2025-05-16 18:09:16.053412	200	0
6080	\N	/status	GET	2025-05-16 18:09:25.974303	2025-05-16 18:09:25.974304	200	0
6081	\N	/status	GET	2025-05-16 18:09:35.993057	2025-05-16 18:09:35.993059	200	0
6082	\N	/status	GET	2025-05-16 18:09:46.065779	2025-05-16 18:09:46.065781	200	0
6083	\N	/status	GET	2025-05-16 18:09:56.020644	2025-05-16 18:09:56.020646	200	0
6084	\N	/status	GET	2025-05-16 18:10:06.040982	2025-05-16 18:10:06.040984	200	0
6085	\N	/status	GET	2025-05-16 18:10:16.056007	2025-05-16 18:10:16.056009	200	0
6086	\N	/status	GET	2025-05-16 18:10:26.096346	2025-05-16 18:10:26.096347	200	0
6087	\N	/status	GET	2025-05-16 18:10:36.128769	2025-05-16 18:10:36.128771	200	0
6088	\N	/status	GET	2025-05-16 18:10:46.112724	2025-05-16 18:10:46.112726	200	0
6089	\N	/status	GET	2025-05-16 18:10:56.116831	2025-05-16 18:10:56.116833	200	0
6090	\N	/status	GET	2025-05-16 18:11:06.159897	2025-05-16 18:11:06.159899	200	0
6091	\N	/status	GET	2025-05-16 18:11:16.191218	2025-05-16 18:11:16.19122	200	0
6092	\N	/status	GET	2025-05-16 18:11:26.176896	2025-05-16 18:11:26.176898	200	0
6093	\N	/status	GET	2025-05-16 18:11:36.18802	2025-05-16 18:11:36.188022	200	0
6094	\N	/status	GET	2025-05-16 18:11:46.204309	2025-05-16 18:11:46.204311	200	0
6095	\N	/status	GET	2025-05-16 18:11:56.224068	2025-05-16 18:11:56.224069	200	0
6096	\N	/status	GET	2025-05-16 18:12:06.305691	2025-05-16 18:12:06.305693	200	0
6097	\N	/status	GET	2025-05-16 18:12:16.300872	2025-05-16 18:12:16.300874	200	0
6098	\N	/status	GET	2025-05-16 18:12:27.275798	2025-05-16 18:12:27.2758	200	0
6099	\N	/status	GET	2025-05-16 18:12:36.308389	2025-05-16 18:12:36.308391	200	0
6100	\N	/status	GET	2025-05-16 18:12:47.307845	2025-05-16 18:12:47.307846	200	0
6101	\N	/status	GET	2025-05-16 18:12:56.34708	2025-05-16 18:12:56.347082	200	0
6102	\N	/status	GET	2025-05-16 18:13:06.352135	2025-05-16 18:13:06.352136	200	0
6103	\N	/status	GET	2025-05-16 18:13:16.349427	2025-05-16 18:13:16.349429	200	0
6104	\N	/status	GET	2025-05-16 18:13:26.3765	2025-05-16 18:13:26.376502	200	0
6105	\N	/status	GET	2025-05-16 18:13:36.413963	2025-05-16 18:13:36.413964	200	0
6106	\N	/status	GET	2025-05-16 18:13:46.420551	2025-05-16 18:13:46.420552	200	0
6107	\N	/status	GET	2025-05-16 18:13:56.432114	2025-05-16 18:13:56.432116	200	0
6108	\N	/status	GET	2025-05-16 18:14:06.448515	2025-05-16 18:14:06.448517	200	0
6109	\N	/status	GET	2025-05-16 18:14:16.476937	2025-05-16 18:14:16.476938	200	0
6110	\N	/status	GET	2025-05-16 18:14:26.479256	2025-05-16 18:14:26.479257	200	0
6111	\N	/status	GET	2025-05-16 18:14:36.500433	2025-05-16 18:14:36.500434	200	0
6112	\N	/status	GET	2025-05-16 18:14:46.512073	2025-05-16 18:14:46.512075	200	0
6113	\N	/status	GET	2025-05-16 18:14:56.544187	2025-05-16 18:14:56.544189	200	0
6114	\N	/status	GET	2025-05-16 18:15:06.559487	2025-05-16 18:15:06.559488	200	0
6115	\N	/status	GET	2025-05-16 18:15:16.56879	2025-05-16 18:15:16.568792	200	0
6116	\N	/status	GET	2025-05-16 18:15:26.607371	2025-05-16 18:15:26.607373	200	0
6117	\N	/status	GET	2025-05-16 18:15:36.628171	2025-05-16 18:15:36.628172	200	0
6118	\N	/status	GET	2025-05-16 18:15:46.6097	2025-05-16 18:15:46.609701	200	0
6119	\N	/status	GET	2025-05-16 18:15:57.718328	2025-05-16 18:15:57.71833	200	0
6120	\N	/status	GET	2025-05-16 18:16:06.683327	2025-05-16 18:16:06.683329	200	0
6121	\N	/status	GET	2025-05-16 18:16:16.680685	2025-05-16 18:16:16.680687	200	0
6122	\N	/status	GET	2025-05-16 18:16:26.720194	2025-05-16 18:16:26.720196	200	0
6123	\N	/status	GET	2025-05-16 18:16:36.728059	2025-05-16 18:16:36.72806	200	0
6124	\N	/status	GET	2025-05-16 18:16:46.726789	2025-05-16 18:16:46.726791	200	0
6125	\N	/status	GET	2025-05-16 18:16:56.738136	2025-05-16 18:16:56.738138	200	0
6126	\N	/status	GET	2025-05-16 18:17:06.744039	2025-05-16 18:17:06.74404	200	0
6127	\N	/status	GET	2025-05-16 18:17:16.754384	2025-05-16 18:17:16.754386	200	0
6128	\N	/status	GET	2025-05-16 18:17:26.817625	2025-05-16 18:17:26.817627	200	0
6129	\N	/status	GET	2025-05-16 18:17:36.810441	2025-05-16 18:17:36.810443	200	0
6130	\N	/status	GET	2025-05-16 18:17:46.805921	2025-05-16 18:17:46.805922	200	0
6131	\N	/status	GET	2025-05-16 18:17:56.826048	2025-05-16 18:17:56.826051	200	0
6132	\N	/status	GET	2025-05-16 18:18:06.868982	2025-05-16 18:18:06.868984	200	0
6133	\N	/status	GET	2025-05-16 18:18:16.874655	2025-05-16 18:18:16.874657	200	0
6134	\N	/status	GET	2025-05-16 18:18:26.897273	2025-05-16 18:18:26.897275	200	0
6135	\N	/status	GET	2025-05-16 18:18:36.905376	2025-05-16 18:18:36.905378	200	0
6136	\N	/status	GET	2025-05-16 18:18:46.91069	2025-05-16 18:18:46.910692	200	0
6137	\N	/api/statistics/top-users	GET	2025-05-16 20:50:34.179736	2025-05-16 20:50:34.181549	200	1
6138	\N	/api/auth/login	POST	2025-05-16 21:23:23.324287	2025-05-16 21:23:23.408479	200	84
6139	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-16 21:23:23.849914	2025-05-16 21:23:23.852218	200	2
6140	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-16 21:23:24.282211	2025-05-16 21:23:24.284274	200	2
6141	\N	/api/weather	GET	2025-05-16 21:23:25.002481	2025-05-16 21:23:25.067778	200	65
6142	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-16 21:23:25.172211	2025-05-16 21:23:25.17387	200	1
6143	\N	/api/restaurant	GET	2025-05-16 21:23:25.171866	2025-05-16 21:23:25.271355	200	99
6146	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-16 21:23:29.193598	2025-05-16 21:23:29.19679	200	3
6151	\N	/api/restaurant	GET	2025-05-16 22:37:06.434732	2025-05-16 22:37:06.434956	200	0
6158	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-17 03:44:53.510451	2025-05-17 03:44:53.510619	200	0
6162	\N	/favicon.ico	GET	2025-05-17 09:57:17.994738	2025-05-17 09:57:17.994743	500	0
6168	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-17 11:15:25.929896	2025-05-17 11:15:25.932139	200	2
6169	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-17 11:15:29.464206	2025-05-17 11:15:29.472064	200	7
6170	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-17 11:15:30.906881	2025-05-17 11:15:30.908648	200	1
6174	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 00:38:14.497323	2025-05-18 00:38:14.497427	200	0
6175	\N	/sftp-config.json	GET	2025-05-18 01:20:32.055439	2025-05-18 01:20:32.055446	500	0
6176	\N	/.vscode/sftp.json	GET	2025-05-18 01:20:32.422067	2025-05-18 01:20:32.42207	500	0
6187	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 10:34:20.420363	2025-05-18 10:34:20.420462	200	0
6144	lucie.delestre@imt-atlantique.net	/api/newf/me	PATCH	2025-05-16 21:23:26.908734	2025-05-16 21:23:26.917617	200	8
6145	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-16 21:23:28.709006	2025-05-16 21:23:28.714136	200	5
6147	\N	/api/auth/register	POST	2025-05-16 22:36:42.992904	2025-05-16 22:36:43.083668	201	90
6148	\N	/api/auth/verify-account	POST	2025-05-16 22:37:05.889964	2025-05-16 22:37:05.902785	200	12
6149	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-16 22:37:05.982155	2025-05-16 22:37:05.984647	200	2
6150	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-16 22:37:06.061115	2025-05-16 22:37:06.063086	200	1
6152	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-16 22:37:06.434747	2025-05-16 22:37:06.437434	200	2
6153	\N	/api/weather	GET	2025-05-16 22:37:06.429896	2025-05-16 22:37:06.497343	200	67
6154	matis.byar@imt-atlantique.net	/api/newf/me	PATCH	2025-05-16 22:37:12.061715	2025-05-16 22:37:12.067133	200	5
6155	\N	/api/restaurant	GET	2025-05-16 22:39:39.511754	2025-05-16 22:39:39.511872	200	0
6156	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-16 22:39:45.138414	2025-05-16 22:39:45.140486	200	2
6157	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-16 22:40:13.818388	2025-05-16 22:40:13.821298	200	2
6159	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-17 08:16:52.614751	2025-05-17 08:16:52.61484	200	0
6160	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-17 08:16:52.624748	2025-05-17 08:16:52.624792	200	0
6161	\N	/api/restaurant	GET	2025-05-17 09:57:17.655487	2025-05-17 09:57:17.760554	200	105
6163	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-17 10:50:47.704998	2025-05-17 10:50:47.705117	200	0
6164	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-17 11:15:21.393226	2025-05-17 11:15:21.39559	200	2
6165	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-17 11:15:21.530164	2025-05-17 11:15:21.532691	200	2
6166	\N	/api/weather	GET	2025-05-17 11:15:21.523608	2025-05-17 11:15:21.572899	200	49
6167	\N	/api/restaurant	GET	2025-05-17 11:15:21.530137	2025-05-17 11:15:22.167773	200	637
6171	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-17 14:04:02.184278	2025-05-17 14:04:02.184509	200	0
6172	\N	/robots.txt	GET	2025-05-17 15:47:43.705552	2025-05-17 15:47:43.705556	500	0
6173	\N	/api/statistics/top-users	GET	2025-05-17 17:33:35.891058	2025-05-17 17:33:35.893667	200	2
6177	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 06:42:49.653199	2025-05-18 06:42:49.653283	200	0
6178	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 06:42:49.682702	2025-05-18 06:42:49.682756	200	0
6179	\N	/	GET	2025-05-18 08:17:24.555112	2025-05-18 08:17:24.555115	500	0
6180	\N	/	GET	2025-05-18 08:17:24.669819	2025-05-18 08:17:24.669821	500	0
6181	\N	/phpinfo.php	GET	2025-05-18 08:17:24.905905	2025-05-18 08:17:24.905909	500	0
6182	\N	/info.php	GET	2025-05-18 08:17:25.041565	2025-05-18 08:17:25.041568	500	0
6183	\N	/phpinfo	GET	2025-05-18 08:17:25.454053	2025-05-18 08:17:25.454057	500	0
6184	\N	/.env	GET	2025-05-18 08:17:25.754176	2025-05-18 08:17:25.754179	500	0
6185	\N	/config/.env	GET	2025-05-18 08:17:26.00445	2025-05-18 08:17:26.004454	500	0
6186	\N	/api/.env	GET	2025-05-18 08:17:26.265155	2025-05-18 08:17:26.265161	500	0
6188	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 12:50:55.071263	2025-05-18 12:50:55.071346	200	0
6194	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:06:17.288774	2025-05-18 16:06:17.290763	200	1
6195	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:07:51.911632	2025-05-18 16:07:51.913727	200	2
6196	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:07:56.163355	2025-05-18 16:07:56.165655	200	2
6197	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:08:08.969352	2025-05-18 16:08:08.971274	200	1
6198	\N	/api/restaurant	GET	2025-05-18 16:08:25.082536	2025-05-18 16:08:25.082607	200	0
6199	\N	/api/traq/	GET	2025-05-18 16:08:26.525331	2025-05-18 16:08:26.533191	200	7
6200	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:09:11.396392	2025-05-18 16:09:11.398934	200	2
6201	\N	/api/restaurant	GET	2025-05-18 16:09:12.447197	2025-05-18 16:09:12.447256	200	0
6202	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:09:12.452384	2025-05-18 16:09:12.454062	200	1
6203	\N	/api/weather	GET	2025-05-18 16:09:12.452405	2025-05-18 16:09:12.513573	200	61
6204	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:09:14.19026	2025-05-18 16:09:14.192741	200	2
6205	\N	/api/restaurant	GET	2025-05-18 16:10:05.323283	2025-05-18 16:10:05.323339	200	0
6206	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:10:13.915752	2025-05-18 16:10:13.917891	200	2
6207	\N	/api/restaurant	GET	2025-05-18 16:10:24.377859	2025-05-18 16:10:24.377954	200	0
6208	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:10:48.491583	2025-05-18 16:10:48.493683	200	2
6209	\N	/api/weather	GET	2025-05-18 16:10:48.669663	2025-05-18 16:10:48.727585	200	57
6210	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:10:49.722483	2025-05-18 16:10:49.724507	200	2
6226	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 00:27:25.473433	2025-05-19 00:27:25.47352	200	0
6228	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 07:22:26.29873	2025-05-19 07:22:26.298791	200	0
6229	\N	/api/newf/me	GET	2025-05-19 08:42:45.045037	2025-05-19 08:42:45.045157	401	0
6230	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-19 08:42:47.594785	2025-05-19 08:42:47.596777	200	1
6231	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-19 08:42:47.79523	2025-05-19 08:42:47.797004	200	1
6232	\N	/api/weather	GET	2025-05-19 08:42:47.789664	2025-05-19 08:42:47.854665	200	65
6233	\N	/api/restaurant	GET	2025-05-19 08:42:47.795227	2025-05-19 08:42:47.924617	200	129
6239	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 09:11:26.72831	2025-05-19 09:11:26.729725	200	1
6240	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 09:12:01.216429	2025-05-19 09:12:01.218322	200	1
6241	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-19 09:12:02.018126	2025-05-19 09:12:02.025002	200	6
6242	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-19 09:12:07.883083	2025-05-19 09:12:07.89257	200	9
6243	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-19 09:12:09.168511	2025-05-19 09:12:09.172098	200	3
8715	\N	/api/restaurant	GET	2025-05-21 21:01:24.645819	2025-05-21 21:01:24.646026	200	0
6189	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-18 15:47:36.028891	2025-05-18 15:47:36.028998	200	0
6190	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:06:15.286893	2025-05-18 16:06:15.288441	200	1
6191	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:06:15.477001	2025-05-18 16:06:15.481107	200	4
6192	\N	/api/weather	GET	2025-05-18 16:06:15.471781	2025-05-18 16:06:15.544599	200	72
6193	\N	/api/restaurant	GET	2025-05-18 16:06:15.477015	2025-05-18 16:06:16.10134	200	624
6211	\N	/api/restaurant	GET	2025-05-18 16:10:49.727338	2025-05-18 16:10:49.727431	200	0
6212	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:11:16.023574	2025-05-18 16:11:16.025955	200	2
6213	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:11:18.809732	2025-05-18 16:11:18.814401	200	4
6214	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-18 16:11:22.99093	2025-05-18 16:11:22.992659	200	1
6215	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-18 16:11:33.322241	2025-05-18 16:11:33.324448	200	2
6216	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-18 16:11:39.861912	2025-05-18 16:11:39.863861	200	1
6217	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:41.093233	2025-05-18 16:11:41.099354	200	6
6218	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:42.841047	2025-05-18 16:11:42.845795	200	4
6219	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-18 16:11:48.011136	2025-05-18 16:11:48.013358	200	2
6220	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:48.971702	2025-05-18 16:11:48.980419	200	8
6221	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:49.234449	2025-05-18 16:11:49.240806	200	6
6222	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:50.871638	2025-05-18 16:11:50.877454	200	5
6223	matis.byar@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-18 16:11:51.661758	2025-05-18 16:11:51.668411	200	6
6224	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:12:51.800248	2025-05-18 16:12:51.802096	200	1
6225	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-18 16:13:00.820758	2025-05-18 16:13:00.822395	200	1
6227	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 07:22:25.632973	2025-05-19 07:22:25.633066	200	0
6234	\N	/api/restaurant	GET	2025-05-19 08:42:58.06956	2025-05-19 08:42:58.069628	200	0
6235	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 09:11:19.031798	2025-05-19 09:11:19.034123	200	2
6236	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 09:11:19.30398	2025-05-19 09:11:19.306349	200	2
6237	\N	/api/weather	GET	2025-05-19 09:11:19.304924	2025-05-19 09:11:19.34842	200	43
6238	\N	/api/restaurant	GET	2025-05-19 09:11:19.296735	2025-05-19 09:11:20.124965	200	828
6244	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:46:48.108432	2025-05-19 09:46:48.11004	304	1
6245	\N	/api/restaurant	GET	2025-05-19 09:46:50.283911	2025-05-19 09:46:50.377514	200	93
6246	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:48:01.540094	2025-05-19 09:48:01.540196	200	0
6247	\N	/api/restaurant	GET	2025-05-19 09:48:02.144341	2025-05-19 09:48:02.144433	200	0
6248	\N	/api/restaurant	GET	2025-05-19 09:48:02.608552	2025-05-19 09:48:02.608621	200	0
6249	\N	/status	GET	2025-05-19 09:48:14.876122	2025-05-19 09:48:14.876126	200	0
6250	\N	/api/statistics/global	GET	2025-05-19 09:48:14.925877	2025-05-19 09:48:14.92842	200	2
6251	\N	/api/statistics/endpoints	GET	2025-05-19 09:48:14.976281	2025-05-19 09:48:14.98485	200	8
6252	\N	/api/washingmachines	GET	2025-05-19 09:48:20.120291	2025-05-19 09:48:20.12031	500	0
6253	\N	/status	GET	2025-05-19 09:48:20.120307	2025-05-19 09:48:20.120315	200	0
6254	\N	/status	GET	2025-05-19 09:48:20.172064	2025-05-19 09:48:20.172066	200	0
6255	\N	/api/statistics/global	GET	2025-05-19 09:48:20.172078	2025-05-19 09:48:20.174376	200	2
6256	\N	/api/statistics/endpoints	GET	2025-05-19 09:48:20.240101	2025-05-19 09:48:20.24735	200	7
6257	\N	/api/restaurant	GET	2025-05-19 09:48:20.892558	2025-05-19 09:48:20.892682	200	0
6258	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:48:20.898067	2025-05-19 09:48:20.898165	200	0
6259	\N	/api/washingmachines	GET	2025-05-19 09:48:26.9683	2025-05-19 09:48:26.968308	500	0
6260	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 09:48:46.035512	2025-05-19 09:48:46.038672	200	3
6261	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 09:48:46.307072	2025-05-19 09:48:46.308847	200	1
6262	\N	/api/weather	GET	2025-05-19 09:48:46.297481	2025-05-19 09:48:46.359874	200	62
6263	\N	/api/restaurant	GET	2025-05-19 09:48:46.29201	2025-05-19 09:48:47.133024	200	841
6264	\N	/api/washingmachines	GET	2025-05-19 09:49:02.189761	2025-05-19 09:49:02.189772	500	0
6265	\N	/api/washingmachines	GET	2025-05-19 09:51:51.124856	2025-05-19 09:51:51.124862	500	0
6266	\N	/api/washingmachines	GET	2025-05-19 09:52:01.358198	2025-05-19 09:52:01.358207	500	0
6267	\N	/favicon.ico	GET	2025-05-19 09:52:01.71709	2025-05-19 09:52:01.717094	500	0
6268	\N	/api/washingmachines	GET	2025-05-19 09:53:20.207169	2025-05-19 09:53:20.207177	500	0
6269	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:54:09.418007	2025-05-19 09:54:09.418107	200	0
6270	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:56:08.403911	2025-05-19 09:56:08.40447	200	0
6271	\N	/api/restaurant	GET	2025-05-19 09:56:08.890696	2025-05-19 09:56:08.972692	200	81
6272	\N	/api/washingmachines	GET	2025-05-19 09:56:25.687273	2025-05-19 09:56:25.687282	500	0
6273	\N	/api/washingmachines	GET	2025-05-19 11:56:51.069117	2025-05-19 11:56:51.069145	500	0
6274	\N	/api/washingmachines	GET	2025-05-19 09:56:51.47006	2025-05-19 09:56:51.470065	500	0
6275	\N	/favicon.ico	GET	2025-05-19 11:56:51.383646	2025-05-19 11:56:51.38365	500	0
6276	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:56:58.020528	2025-05-19 09:56:58.020616	200	0
6277	\N	/api/restaurant	GET	2025-05-19 09:56:58.4796	2025-05-19 09:56:58.47968	200	0
6278	\N	/api/washingmachines/	GET	2025-05-19 11:57:19.682135	2025-05-19 11:57:19.821768	200	139
6279	\N	/api/washingmachines/	GET	2025-05-19 11:57:39.218685	2025-05-19 11:57:39.311365	200	92
6280	\N	/api/washingmachines/	GET	2025-05-19 11:57:41.611938	2025-05-19 11:57:41.655644	200	43
6281	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 09:57:42.145728	2025-05-19 09:57:42.145919	200	0
6282	\N	/api/restaurant	GET	2025-05-19 09:57:42.206801	2025-05-19 09:57:42.20686	200	0
6283	\N	/api/washingmachines/	GET	2025-05-19 11:57:45.813632	2025-05-19 11:57:45.852128	200	38
6284	\N	/api/washingmachines	GET	2025-05-19 09:58:20.155568	2025-05-19 09:58:20.155577	500	0
6285	\N	/api/washingmachines/	GET	2025-05-19 11:58:27.468038	2025-05-19 11:58:27.550536	200	82
6286	\N	/api/restaurant	GET	2025-05-19 09:58:30.807751	2025-05-19 09:58:30.807816	200	0
6287	\N	/api/washingmachines	GET	2025-05-19 09:58:40.683642	2025-05-19 09:58:40.683651	500	0
6288	\N	/status	GET	2025-05-19 09:59:03.929094	2025-05-19 09:59:03.929096	200	0
6289	\N	/api/statistics/global	GET	2025-05-19 09:59:04.010459	2025-05-19 09:59:04.013155	200	2
6290	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:04.082406	2025-05-19 09:59:04.089103	200	6
6291	\N	/api/restaurant	GET	2025-05-19 09:59:04.687707	2025-05-19 09:59:04.687809	200	0
6292	\N	/api/washingmachines	GET	2025-05-19 09:59:06.329016	2025-05-19 09:59:06.329023	500	0
6293	\N	/api/washingmachines	GET	2025-05-19 09:59:08.450877	2025-05-19 09:59:08.450888	500	0
6294	\N	/status	GET	2025-05-19 09:59:08.939342	2025-05-19 09:59:08.939344	200	0
6295	\N	/status	GET	2025-05-19 09:59:09.027149	2025-05-19 09:59:09.02715	200	0
6296	\N	/api/statistics/global	GET	2025-05-19 09:59:09.032555	2025-05-19 09:59:09.034295	200	1
6297	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:09.128211	2025-05-19 09:59:09.13376	200	5
6298	\N	/status	GET	2025-05-19 09:59:13.945895	2025-05-19 09:59:13.945904	200	0
6299	\N	/status	GET	2025-05-19 09:59:14.0637	2025-05-19 09:59:14.063704	200	0
6300	\N	/api/statistics/global	GET	2025-05-19 09:59:14.065598	2025-05-19 09:59:14.068524	200	2
6301	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:14.178207	2025-05-19 09:59:14.185296	200	7
6302	\N	/status	GET	2025-05-19 09:59:14.445951	2025-05-19 09:59:14.445953	200	0
6303	\N	/api/statistics/global	GET	2025-05-19 09:59:14.562278	2025-05-19 09:59:14.566897	200	4
6304	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:14.657091	2025-05-19 09:59:14.661902	200	4
6305	\N	/status	GET	2025-05-19 09:59:18.963844	2025-05-19 09:59:18.963845	200	0
6306	\N	/status	GET	2025-05-19 09:59:19.061448	2025-05-19 09:59:19.06145	200	0
6307	\N	/api/statistics/global	GET	2025-05-19 09:59:19.062222	2025-05-19 09:59:19.064365	200	2
6308	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:19.145557	2025-05-19 09:59:19.151198	200	5
6309	\N	/status	GET	2025-05-19 09:59:23.961511	2025-05-19 09:59:23.961513	200	0
6310	\N	/status	GET	2025-05-19 09:59:24.042793	2025-05-19 09:59:24.042795	200	0
6311	\N	/api/statistics/global	GET	2025-05-19 09:59:24.081768	2025-05-19 09:59:24.084142	200	2
6312	\N	/api/statistics/endpoints	GET	2025-05-19 09:59:24.167266	2025-05-19 09:59:24.182357	200	15
8756	\N	/api/washingmachines	GET	2025-05-21 21:27:17.066875	2025-05-21 21:27:17.169187	200	102
8757	\N	/api/washingmachines	GET	2025-05-21 21:27:20.239497	2025-05-21 21:27:20.269507	200	30
8758	\N	/api/washingmachines	GET	2025-05-21 21:27:21.887411	2025-05-21 21:27:21.917985	200	30
8759	\N	/api/washingmachines	GET	2025-05-21 21:27:29.187801	2025-05-21 21:27:29.218181	200	30
8760	\N	/api/washingmachines	GET	2025-05-21 21:27:32.503079	2025-05-21 21:27:32.5326	200	29
8761	\N	/api/washingmachines	GET	2025-05-21 21:27:44.264149	2025-05-21 21:27:44.293607	200	29
8762	\N	/api/weather	GET	2025-05-21 21:28:36.104154	2025-05-21 21:28:36.169177	200	65
8763	\N	/api/weather	GET	2025-05-21 21:30:42.963835	2025-05-21 21:30:42.963862	200	0
8764	\N	/api/weather	GET	2025-05-21 21:31:00.780679	2025-05-21 21:31:00.780711	200	0
8765	\N	/api/weather	GET	2025-05-21 21:31:16.681221	2025-05-21 21:31:16.681258	200	0
8766	\N	/api/weather	GET	2025-05-21 21:31:34.367681	2025-05-21 21:31:34.367707	200	0
8767	\N	/api/weather	GET	2025-05-21 21:31:51.89121	2025-05-21 21:31:51.891238	200	0
8768	\N	/api/weather	GET	2025-05-21 21:33:49.231401	2025-05-21 21:33:49.289598	200	58
8769	\N	/api/weather	GET	2025-05-21 21:33:57.195275	2025-05-21 21:33:57.195307	200	0
8770	\N	/api/weather	GET	2025-05-21 21:34:05.920052	2025-05-21 21:34:05.920078	200	0
8771	\N	/api/weather	GET	2025-05-21 21:34:15.626623	2025-05-21 21:34:15.626649	200	0
8772	\N	/api/weather	GET	2025-05-21 21:34:23.319892	2025-05-21 21:34:23.31994	200	0
8773	\N	/api/weather	GET	2025-05-21 21:39:27.235784	2025-05-21 21:39:27.297136	200	61
8774	\N	/api/weather	GET	2025-05-21 21:39:33.2656	2025-05-21 21:39:33.265624	200	0
8775	\N	/api/weather	GET	2025-05-21 21:39:39.838838	2025-05-21 21:39:39.838865	200	0
8776	\N	/api/weather	GET	2025-05-21 21:39:46.18408	2025-05-21 21:39:46.184109	200	0
8777	\N	/api/weather	GET	2025-05-21 21:39:53.500529	2025-05-21 21:39:53.500553	200	0
8778	\N	/api/weather	GET	2025-05-21 21:40:01.329648	2025-05-21 21:40:01.329674	200	0
8779	\N	/api/weather	GET	2025-05-21 21:40:08.288114	2025-05-21 21:40:08.288154	200	0
8780	\N	/api/weather	GET	2025-05-21 21:40:14.814528	2025-05-21 21:40:14.814553	200	0
8781	\N	/api/weather	GET	2025-05-21 21:40:21.633297	2025-05-21 21:40:21.63332	200	0
8782	\N	/api/weather	GET	2025-05-21 21:40:30.563707	2025-05-21 21:40:30.56375	200	0
8783	\N	/api/weather	GET	2025-05-21 21:40:37.193677	2025-05-21 21:40:37.193708	200	0
8784	\N	/api/weather	GET	2025-05-21 21:40:42.833529	2025-05-21 21:40:42.833558	200	0
8785	\N	/api/weather	GET	2025-05-21 21:40:49.250364	2025-05-21 21:40:49.250389	200	0
8786	\N	/api/weather	GET	2025-05-21 21:40:56.48285	2025-05-21 21:40:56.482879	200	0
8787	\N	/api/weather	GET	2025-05-21 21:41:03.052338	2025-05-21 21:41:03.052364	200	0
8788	\N	/api/weather	GET	2025-05-21 21:41:32.249248	2025-05-21 21:41:32.249272	200	0
8789	\N	/api/weather	GET	2025-05-21 21:42:07.161044	2025-05-21 21:42:07.161067	200	0
8793	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 22:03:00.116854	2025-05-21 22:03:00.118695	200	1
8794	\N	/api/weather	GET	2025-05-21 22:03:00.116838	2025-05-21 22:03:00.168121	200	51
8795	\N	/api/washingmachines	GET	2025-05-21 22:03:00.095138	2025-05-21 22:03:00.188196	200	93
8796	\N	/api/auth/login	POST	2025-05-21 22:03:20.664287	2025-05-21 22:03:20.73752	401	73
8797	\N	/api/auth/login	POST	2025-05-21 22:03:32.459263	2025-05-21 22:03:32.530958	401	71
8798	\N	/api/auth/login	POST	2025-05-21 22:03:44.51991	2025-05-21 22:03:44.591273	401	71
8799	\N	/api/auth/login	POST	2025-05-21 22:03:47.837985	2025-05-21 22:03:47.910611	401	72
8800	\N	/api/auth/login	POST	2025-05-21 22:04:34.019991	2025-05-21 22:04:34.09335	401	73
8801	\N	/api/auth/login	POST	2025-05-21 22:04:49.218081	2025-05-21 22:04:49.292186	401	74
8802	\N	/api/auth/login	POST	2025-05-21 22:05:36.461574	2025-05-21 22:05:36.533312	401	71
8803	\N	/api/auth/login	POST	2025-05-21 22:05:47.44057	2025-05-21 22:05:47.512991	401	72
8804	\N	/api/auth/login	POST	2025-05-21 22:05:58.107878	2025-05-21 22:05:58.184201	401	76
8805	\N	/api/auth/verification-code	POST	2025-05-21 22:06:03.762584	2025-05-21 22:06:03.793872	200	31
8806	\N	/api/auth/login	POST	2025-05-21 22:06:14.208623	2025-05-21 22:06:14.283408	401	74
8807	\N	/api/auth/login	POST	2025-05-21 22:06:25.838635	2025-05-21 22:06:25.912124	401	73
8808	\N	/api/auth/login	POST	2025-05-21 22:06:42.217001	2025-05-21 22:06:42.29197	401	74
8809	\N	/api/auth/login	POST	2025-05-21 22:06:53.370882	2025-05-21 22:06:53.444713	401	73
8810	\N	/api/auth/login	POST	2025-05-21 22:07:07.769074	2025-05-21 22:07:07.840981	401	71
8811	\N	/api/auth/login	POST	2025-05-21 22:07:25.310528	2025-05-21 22:07:25.384855	401	74
8812	\N	/api/auth/login	POST	2025-05-21 22:07:37.760803	2025-05-21 22:07:37.832135	401	71
8813	\N	/api/auth/login	POST	2025-05-21 22:07:51.144534	2025-05-21 22:07:51.218269	401	73
8814	\N	/api/auth/login	POST	2025-05-21 22:07:55.067975	2025-05-21 22:07:55.141606	401	73
8815	\N	/api/auth/login	POST	2025-05-21 22:08:45.314961	2025-05-21 22:08:45.387277	401	72
8816	\N	/api/auth/login	POST	2025-05-21 22:09:01.368595	2025-05-21 22:09:01.443851	401	75
8817	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 22:11:45.891108	2025-05-21 22:11:45.893773	200	2
8818	\N	/api/washingmachines	GET	2025-05-21 22:11:46.354842	2025-05-21 22:11:46.444119	200	89
8820	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 22:11:46.545472	2025-05-21 22:11:46.547589	200	2
8821	\N	/api/weather	GET	2025-05-21 22:11:46.545604	2025-05-21 22:11:46.59668	200	51
8823	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 22:26:10.851527	2025-05-21 22:26:10.851595	200	0
6313	\N	/api/washingmachines/	GET	2025-05-19 12:00:37.479177	2025-05-19 12:00:37.570392	200	91
6314	\N	/api/washingmachines	GET	2025-05-19 10:00:40.838028	2025-05-19 10:00:40.927993	200	89
6315	\N	/api/washingmachines	GET	2025-05-19 10:01:22.886893	2025-05-19 10:01:22.914043	200	27
6316	\N	/api/washingmachines	GET	2025-05-19 10:03:20.271847	2025-05-19 10:03:20.371353	200	99
6317	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 10:04:05.056391	2025-05-19 10:04:05.056947	200	0
6318	\N	/api/restaurant	GET	2025-05-19 10:04:05.378864	2025-05-19 10:04:05.477003	200	98
6319	\N	/api/washingmachines	GET	2025-05-19 10:04:08.389448	2025-05-19 10:04:08.418247	200	28
6320	\N	/api/washingmachines	GET	2025-05-19 10:05:40.896914	2025-05-19 10:05:40.994468	200	97
6321	\N	/api/washingmachines	GET	2025-05-19 10:08:20.085496	2025-05-19 10:08:20.180053	200	94
6322	\N	/api/washingmachines	GET	2025-05-19 10:09:08.195052	2025-05-19 10:09:08.226419	200	31
6323	\N	/api/washingmachines	GET	2025-05-19 10:09:11.143482	2025-05-19 10:09:11.173592	200	30
6324	\N	/api/washingmachines	GET	2025-05-19 10:10:40.654123	2025-05-19 10:10:40.691438	200	37
6325	\N	/api/washingmachines	GET	2025-05-19 10:11:22.812297	2025-05-19 10:11:22.842637	200	30
6326	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:11:35.042259	2025-05-19 10:11:35.045731	200	3
6327	\N	/api/washingmachines	GET	2025-05-19 10:11:35.905853	2025-05-19 10:11:35.951879	200	46
6328	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:11:36.464271	2025-05-19 10:11:36.466048	200	1
6329	\N	/api/weather	GET	2025-05-19 10:11:36.439502	2025-05-19 10:11:36.488223	200	48
6330	\N	/api/restaurant	GET	2025-05-19 10:11:36.44844	2025-05-19 10:11:37.43642	200	987
6331	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:11:39.163115	2025-05-19 10:11:39.164874	200	1
6332	\N	/api/washingmachines	GET	2025-05-19 10:11:40.285854	2025-05-19 10:11:40.316108	200	30
6333	\N	/api/weather	GET	2025-05-19 10:11:40.89232	2025-05-19 10:11:40.892349	200	0
6334	\N	/api/restaurant	GET	2025-05-19 10:11:40.914243	2025-05-19 10:11:40.914301	200	0
6335	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:11:40.925466	2025-05-19 10:11:40.928822	200	3
6336	\N	/api/washingmachines	GET	2025-05-19 10:11:42.842758	2025-05-19 10:11:42.890172	200	47
6337	\N	/api/washingmachines	GET	2025-05-19 10:12:07.55339	2025-05-19 10:12:07.585334	200	31
6338	\N	/api/washingmachines	GET	2025-05-19 10:12:42.385372	2025-05-19 10:12:42.416488	200	31
6339	\N	/api/washingmachines	GET	2025-05-19 10:12:43.769428	2025-05-19 10:12:43.80469	200	35
6340	\N	/api/washingmachines	GET	2025-05-19 10:12:44.205896	2025-05-19 10:12:44.23753	200	31
6341	\N	/api/washingmachines	GET	2025-05-19 10:12:44.669944	2025-05-19 10:12:44.707607	200	37
6342	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:12:50.934079	2025-05-19 10:12:50.936037	200	1
6343	\N	/api/washingmachines	GET	2025-05-19 10:12:51.151591	2025-05-19 10:12:51.18155	200	29
6344	\N	/api/weather	GET	2025-05-19 10:12:51.221393	2025-05-19 10:12:51.221422	200	0
6345	\N	/api/restaurant	GET	2025-05-19 10:12:51.226252	2025-05-19 10:12:51.226473	200	0
6346	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 10:12:51.226939	2025-05-19 10:12:51.379479	200	152
6347	\N	/api/washingmachines	GET	2025-05-19 10:12:53.318486	2025-05-19 10:12:53.349109	200	30
6348	\N	/api/washingmachines	GET	2025-05-19 10:12:54.821549	2025-05-19 10:12:54.857373	200	35
6349	\N	/api/washingmachines	GET	2025-05-19 10:12:55.274965	2025-05-19 10:12:55.304904	200	29
6350	\N	/api/washingmachines	GET	2025-05-19 10:12:55.800987	2025-05-19 10:12:55.838287	200	37
6351	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 10:13:06.88744	2025-05-19 10:13:06.887555	200	0
6352	\N	/api/restaurant	GET	2025-05-19 10:13:07.46912	2025-05-19 10:13:07.469186	200	0
6353	\N	/api/washingmachines	GET	2025-05-19 10:13:14.049654	2025-05-19 10:13:14.079589	200	29
6354	\N	/api/washingmachines	GET	2025-05-19 10:13:20.220389	2025-05-19 10:13:20.251138	200	30
6355	\N	/api/washingmachines	GET	2025-05-19 10:14:08.190351	2025-05-19 10:14:08.221448	200	31
6356	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 10:23:18.652403	2025-05-19 10:23:18.652484	200	0
6357	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 10:23:18.995523	2025-05-19 10:23:18.995592	200	0
6358	\N	/api/restaurant	GET	2025-05-19 10:23:19.131904	2025-05-19 10:23:19.13202	200	0
6359	\N	/api/washingmachines	GET	2025-05-19 10:54:45.030316	2025-05-19 10:54:45.124703	200	94
6360	\N	/api/washingmachines	GET	2025-05-19 10:56:49.504183	2025-05-19 10:56:49.598234	200	94
6361	\N	/api/washingmachines	GET	2025-05-19 10:57:37.793387	2025-05-19 10:57:37.825641	200	32
6362	\N	/api/washingmachines	GET	2025-05-19 10:59:45.316863	2025-05-19 10:59:45.39042	200	73
6363	\N	/api/washingmachines	GET	2025-05-19 11:01:49.392836	2025-05-19 11:01:49.468923	200	76
6364	\N	/api/washingmachines	GET	2025-05-19 11:02:37.712311	2025-05-19 11:02:37.743838	200	31
6365	\N	/api/washingmachines	GET	2025-05-19 11:04:45.120655	2025-05-19 11:04:45.210875	200	90
6366	\N	/api/washingmachines	GET	2025-05-19 11:06:49.393808	2025-05-19 11:06:49.450733	200	56
6367	\N	/api/washingmachines	GET	2025-05-19 11:07:37.713763	2025-05-19 11:07:37.741985	200	28
6368	\N	/api/washingmachines	GET	2025-05-19 11:09:45.276395	2025-05-19 11:09:45.368431	200	92
6369	\N	/api/washingmachines	GET	2025-05-19 11:11:49.506438	2025-05-19 11:11:49.601985	200	95
6370	\N	/api/washingmachines	GET	2025-05-19 11:12:38.024305	2025-05-19 11:12:38.054886	200	30
6371	\N	/api/washingmachines	GET	2025-05-19 11:27:47.08997	2025-05-19 11:27:47.160052	200	70
6372	\N	/api/washingmachines	GET	2025-05-19 11:33:10.237197	2025-05-19 11:33:10.303541	200	66
6373	\N	/api/washingmachines	GET	2025-05-19 11:38:09.890568	2025-05-19 11:38:09.982939	200	92
6374	\N	/api/washingmachines	GET	2025-05-19 11:43:10.419512	2025-05-19 11:43:10.498378	200	78
6375	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 11:43:30.384834	2025-05-19 11:43:30.384939	200	0
6376	\N	/api/restaurant	GET	2025-05-19 11:43:30.671401	2025-05-19 11:43:30.671506	200	0
6377	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 11:43:45.510679	2025-05-19 11:43:45.510798	200	0
6378	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 11:43:45.73695	2025-05-19 11:43:45.737023	200	0
6379	\N	/api/restaurant	GET	2025-05-19 11:43:45.742233	2025-05-19 11:43:45.742301	200	0
6380	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 11:46:14.217757	2025-05-19 11:46:14.219787	200	2
6381	\N	/api/restaurant	GET	2025-05-19 11:46:23.468953	2025-05-19 11:46:23.469016	200	0
6382	\N	/api/restaurant	GET	2025-05-19 11:46:39.884553	2025-05-19 11:46:39.884629	200	0
6383	\N	/api/washingmachines	GET	2025-05-19 11:48:10.154304	2025-05-19 11:48:10.25076	200	96
6384	\N	/api/washingmachines	GET	2025-05-19 11:53:10.120711	2025-05-19 11:53:10.191079	200	70
6385	\N	/api/restaurant	GET	2025-05-19 11:55:31.936287	2025-05-19 11:55:31.936377	200	0
6386	\N	/api/restaurant	GET	2025-05-19 11:55:48.738614	2025-05-19 11:55:48.73867	200	0
6387	\N	/api/washingmachines	GET	2025-05-19 11:58:10.454592	2025-05-19 11:58:10.551571	200	96
6388	\N	/api/washingmachines	GET	2025-05-19 11:59:48.832688	2025-05-19 11:59:48.937746	200	105
6389	\N	/api/washingmachines	GET	2025-05-19 12:01:08.954222	2025-05-19 12:01:09.003961	200	49
6390	\N	/api/washingmachines	GET	2025-05-19 12:01:17.255856	2025-05-19 12:01:17.28709	200	31
6391	\N	/api/washingmachines	GET	2025-05-19 12:01:18.84284	2025-05-19 12:01:18.875375	200	32
6392	\N	/api/washingmachines	GET	2025-05-19 12:03:10.034048	2025-05-19 12:03:10.104634	200	70
6393	\N	/api/washingmachines	GET	2025-05-19 12:03:55.974957	2025-05-19 12:03:56.005096	200	30
6394	\N	/api/washingmachines	GET	2025-05-19 12:03:56.432555	2025-05-19 12:03:56.461639	200	29
6395	\N	/api/washingmachines	GET	2025-05-19 12:03:56.920213	2025-05-19 12:03:56.951681	200	31
6396	\N	/api/washingmachines	GET	2025-05-19 12:04:17.316224	2025-05-19 12:04:17.345767	200	29
6397	\N	/api/washingmachines	GET	2025-05-19 12:04:23.286382	2025-05-19 12:04:23.316717	200	30
6398	\N	/api/washingmachines	GET	2025-05-19 12:04:30.671742	2025-05-19 12:04:30.700745	200	29
6399	\N	/api/restaurant	GET	2025-05-19 12:05:13.432853	2025-05-19 12:05:13.433096	200	0
6400	\N	/api/washingmachines	GET	2025-05-19 12:05:13.424761	2025-05-19 12:05:13.455997	200	31
6401	\N	/api/weather	GET	2025-05-19 12:05:13.436932	2025-05-19 12:05:13.491486	200	54
6402	\N	/api/restaurant	GET	2025-05-19 12:05:23.744507	2025-05-19 12:05:23.744591	200	0
6403	\N	/api/restaurant	GET	2025-05-19 12:05:27.66165	2025-05-19 12:05:27.661794	200	0
6404	\N	/api/weather	GET	2025-05-19 12:05:27.66156	2025-05-19 12:05:27.673861	200	12
6405	\N	/api/washingmachines	GET	2025-05-19 12:05:27.661344	2025-05-19 12:05:27.702532	200	41
6406	\N	/api/washingmachines	GET	2025-05-19 12:08:10.320438	2025-05-19 12:08:10.395497	200	75
6407	\N	/api/restaurant	GET	2025-05-19 12:08:31.953103	2025-05-19 12:08:31.95322	200	0
6408	\N	/api/washingmachines	GET	2025-05-19 12:08:46.399114	2025-05-19 12:08:46.480359	200	81
6409	\N	/api/restaurant	GET	2025-05-19 12:12:45.057486	2025-05-19 12:12:45.05761	200	0
6410	\N	/api/restaurant	GET	2025-05-19 12:12:45.32159	2025-05-19 12:12:45.32165	200	0
6411	\N	/api/restaurant	GET	2025-05-19 12:12:45.605464	2025-05-19 12:12:45.605533	200	0
6412	\N	/api/restaurant	GET	2025-05-19 12:12:45.890398	2025-05-19 12:12:45.890472	200	0
6413	\N	/api/restaurant	GET	2025-05-19 12:12:46.973508	2025-05-19 12:12:46.973563	200	0
6414	\N	/api/restaurant	GET	2025-05-19 12:12:47.254617	2025-05-19 12:12:47.254675	200	0
6415	\N	/api/restaurant	GET	2025-05-19 12:13:01.866632	2025-05-19 12:13:01.866744	200	0
6416	\N	/api/restaurant	GET	2025-05-19 12:13:04.358934	2025-05-19 12:13:04.359012	200	0
6417	\N	/api/restaurant	GET	2025-05-19 12:13:04.6411	2025-05-19 12:13:04.641154	200	0
6418	\N	/api/restaurant	GET	2025-05-19 12:13:08.720631	2025-05-19 12:13:08.720702	200	0
6419	\N	/api/washingmachines	GET	2025-05-19 12:13:10.202578	2025-05-19 12:13:10.272881	200	70
6420	\N	/api/washingmachines	GET	2025-05-19 12:18:10.368704	2025-05-19 12:18:10.464762	200	96
6421	\N	/api/washingmachines	GET	2025-05-19 12:23:09.896126	2025-05-19 12:23:09.98449	200	88
6422	\N	/api/washingmachines	GET	2025-05-19 12:25:37.500331	2025-05-19 12:25:37.574221	200	73
6423	\N	/api/washingmachines	GET	2025-05-19 12:26:05.482255	2025-05-19 12:26:05.513832	200	31
6424	\N	/api/washingmachines	GET	2025-05-19 12:28:09.899213	2025-05-19 12:28:09.995888	200	96
6425	\N	/api/traq/	GET	2025-05-19 12:30:32.56334	2025-05-19 12:30:32.564957	200	1
6426	\N	/api/traq/	GET	2025-05-19 12:31:04.727063	2025-05-19 12:31:04.728517	200	1
6427	\N	/api/washingmachines	GET	2025-05-19 12:31:05.579835	2025-05-19 12:31:05.647946	200	68
6428	\N	/api/washingmachines	GET	2025-05-19 12:32:10.363966	2025-05-19 12:32:10.391505	200	27
6429	\N	/api/washingmachines	GET	2025-05-19 12:32:22.4942	2025-05-19 12:32:22.533157	200	38
6430	\N	/api/washingmachines	GET	2025-05-19 12:32:31.029554	2025-05-19 12:32:31.056977	200	27
6431	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:32:36.921576	2025-05-19 12:32:36.923986	200	2
6432	\N	/api/washingmachines	GET	2025-05-19 12:33:10.380961	2025-05-19 12:33:10.422168	200	41
6433	yohann.chavanel@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 12:35:05.303953	2025-05-19 12:35:05.305989	200	2
6434	\N	/api/restaurant	GET	2025-05-19 12:36:12.307463	2025-05-19 12:36:12.307528	200	0
6435	\N	/api/restaurant	GET	2025-05-19 12:36:23.976563	2025-05-19 12:36:23.976661	200	0
6436	\N	/api/washingmachines	GET	2025-05-19 12:36:23.963357	2025-05-19 12:36:24.032106	200	68
6437	\N	/api/weather	GET	2025-05-19 12:36:23.976812	2025-05-19 12:36:24.042896	200	66
6438	\N	/api/washingmachines	GET	2025-05-19 12:36:25.827581	2025-05-19 12:36:25.85602	200	28
6439	\N	/api/washingmachines	GET	2025-05-19 12:36:28.75295	2025-05-19 12:36:28.789988	200	37
6440	\N	/api/restaurant	GET	2025-05-19 12:36:29.007381	2025-05-19 12:36:29.007435	200	0
6441	\N	/api/washingmachines	GET	2025-05-19 12:36:31.62545	2025-05-19 12:36:31.654581	200	29
6442	\N	/api/traq/	GET	2025-05-19 12:36:35.982131	2025-05-19 12:36:35.985529	200	3
6443	\N	/api/traq/	GET	2025-05-19 12:36:38.416645	2025-05-19 12:36:38.417807	200	1
6444	\N	/api/traq/	GET	2025-05-19 12:36:38.696101	2025-05-19 12:36:38.697358	200	1
6445	\N	/api/traq/	GET	2025-05-19 12:36:38.962359	2025-05-19 12:36:38.963432	200	1
6446	\N	/api/traq/	GET	2025-05-19 12:36:39.379516	2025-05-19 12:36:39.380741	200	1
6447	\N	/api/traq/	GET	2025-05-19 12:37:06.435238	2025-05-19 12:37:06.436458	200	1
6448	\N	/api/traq/	GET	2025-05-19 12:37:17.460375	2025-05-19 12:37:17.461513	200	1
6449	\N	/api/traq/	GET	2025-05-19 12:37:34.92897	2025-05-19 12:37:34.930181	200	1
6450	\N	/api/traq/	GET	2025-05-19 12:37:50.481446	2025-05-19 12:37:50.482634	200	1
6451	\N	/api/traq/	GET	2025-05-19 12:37:50.795681	2025-05-19 12:37:50.796995	200	1
6452	\N	/api/washingmachines	GET	2025-05-19 12:38:10.284836	2025-05-19 12:38:10.351954	200	67
6453	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:38:12.097209	2025-05-19 12:38:12.09906	200	1
6454	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:38:14.07569	2025-05-19 12:38:14.077335	200	1
6455	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:38:14.523854	2025-05-19 12:38:14.525827	200	1
6456	\N	/api/traq/	GET	2025-05-19 12:38:27.669253	2025-05-19 12:38:27.670439	200	1
6457	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:39:06.582716	2025-05-19 12:39:06.584636	200	1
6458	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:39:06.9479	2025-05-19 12:39:06.950046	200	2
6459	\N	/api/washingmachines	GET	2025-05-19 12:41:28.605255	2025-05-19 12:41:28.681567	200	76
6460	\N	/api/auth/login	POST	2025-05-19 12:42:32.020865	2025-05-19 12:42:32.10039	200	79
6461	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:32.21156	2025-05-19 12:42:32.213455	200	1
6462	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:32.295097	2025-05-19 12:42:32.296775	200	1
6463	\N	/api/washingmachines	GET	2025-05-19 12:42:32.4273	2025-05-19 12:42:32.456527	200	29
6464	\N	/api/weather	GET	2025-05-19 12:42:33.156238	2025-05-19 12:42:33.212291	200	56
6465	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:33.236132	2025-05-19 12:42:33.237792	200	1
6471	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-19 12:42:43.813991	2025-05-19 12:42:43.814103	200	0
6466	\N	/api/restaurant	GET	2025-05-19 12:42:33.241792	2025-05-19 12:42:33.241861	200	0
6467	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 12:42:35.137905	2025-05-19 12:42:35.142121	200	4
6468	\N	/api/restaurant	GET	2025-05-19 12:42:38.332087	2025-05-19 12:42:38.332166	200	0
6469	\N	/api/restaurant	GET	2025-05-19 12:42:41.502539	2025-05-19 12:42:41.502604	200	0
6470	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:43.500294	2025-05-19 12:42:43.50237	200	2
6472	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:45.339074	2025-05-19 12:42:45.34113	200	2
6473	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 12:42:48.483006	2025-05-19 12:42:48.484976	200	1
6474	\N	/api/washingmachines	GET	2025-05-19 12:42:50.546048	2025-05-19 12:42:50.57436	200	28
6475	\N	/api/restaurant	GET	2025-05-19 12:42:53.237177	2025-05-19 12:42:53.237257	200	0
6476	\N	/api/traq/	GET	2025-05-19 12:42:59.709129	2025-05-19 12:42:59.713846	200	4
6477	\N	/api/washingmachines	GET	2025-05-19 12:43:09.892402	2025-05-19 12:43:09.920502	200	28
6478	\N	/status	GET	2025-05-19 12:43:10.084875	2025-05-19 12:43:10.084884	200	0
6479	\N	/api/statistics/global	GET	2025-05-19 12:43:10.141966	2025-05-19 12:43:10.144557	200	2
6480	\N	/api/statistics/top-users	GET	2025-05-19 12:43:10.200792	2025-05-19 12:43:10.202529	200	1
6481	\N	/status	GET	2025-05-19 12:43:20.09169	2025-05-19 12:43:20.091691	200	0
6482	\N	/status	GET	2025-05-19 12:43:30.1017	2025-05-19 12:43:30.101702	200	0
6483	\N	/status	GET	2025-05-19 12:43:40.104328	2025-05-19 12:43:40.10433	200	0
6484	\N	/status	GET	2025-05-19 12:43:50.117951	2025-05-19 12:43:50.117952	200	0
6485	\N	/status	GET	2025-05-19 12:43:53.471336	2025-05-19 12:43:53.471338	200	0
6486	\N	/api/statistics/global	GET	2025-05-19 12:43:53.559642	2025-05-19 12:43:53.561777	200	2
6487	\N	/api/statistics/top-users	GET	2025-05-19 12:43:53.664378	2025-05-19 12:43:53.669154	200	4
6488	\N	/status	GET	2025-05-19 12:43:55.908364	2025-05-19 12:43:55.908366	200	0
6489	\N	/api/statistics/global	GET	2025-05-19 12:43:56.013123	2025-05-19 12:43:56.015163	200	2
6490	\N	/api/statistics/top-users	GET	2025-05-19 12:43:56.11506	2025-05-19 12:43:56.116251	200	1
6491	\N	/status	GET	2025-05-19 12:44:00.119324	2025-05-19 12:44:00.119325	200	0
6492	\N	/status	GET	2025-05-19 12:44:10.127398	2025-05-19 12:44:10.1274	200	0
6493	\N	/api/restaurant	GET	2025-05-19 12:56:57.373265	2025-05-19 12:56:57.37338	200	0
6494	\N	/api/restaurant	GET	2025-05-19 12:56:58.090746	2025-05-19 12:56:58.090833	200	0
6495	\N	/api/washingmachines	GET	2025-05-19 12:56:59.025464	2025-05-19 12:56:59.123403	200	97
6496	\N	/api/washingmachines	GET	2025-05-19 13:01:03.427677	2025-05-19 13:01:03.526375	200	98
6497	\N	/api/washingmachines	GET	2025-05-19 13:02:44.145767	2025-05-19 13:02:44.249125	200	103
6498	\N	/api/newf/me	GET	2025-05-19 13:03:14.828252	2025-05-19 13:03:14.828418	401	0
6499	\N	/api/auth/login	POST	2025-05-19 13:04:36.8129	2025-05-19 13:04:36.889616	200	76
6500	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:04:36.961976	2025-05-19 13:04:36.963643	200	1
6501	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:04:37.017873	2025-05-19 13:04:37.019526	200	1
6502	\N	/api/restaurant	GET	2025-05-19 13:04:37.249211	2025-05-19 13:04:37.249297	200	0
6503	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:04:37.249196	2025-05-19 13:04:37.252968	200	3
6504	\N	/api/washingmachines	GET	2025-05-19 13:04:37.166453	2025-05-19 13:04:37.265281	200	98
6505	\N	/api/weather	GET	2025-05-19 13:04:37.249253	2025-05-19 13:04:37.317089	200	67
6506	\N	/api/restaurant	GET	2025-05-19 13:04:43.560451	2025-05-19 13:04:43.560553	200	0
6507	\N	/api/washingmachines	GET	2025-05-19 13:06:02.743664	2025-05-19 13:06:02.773206	200	29
6508	\N	/api/washingmachines	GET	2025-05-19 13:07:44.14446	2025-05-19 13:07:44.231642	200	87
6509	\N	/api/restaurant	GET	2025-05-19 13:10:43.208771	2025-05-19 13:10:43.208876	200	0
6510	\N	/api/weather	GET	2025-05-19 13:10:43.203904	2025-05-19 13:10:43.251836	200	47
6511	\N	/api/washingmachines	GET	2025-05-19 13:10:43.338566	2025-05-19 13:10:43.413314	200	74
6512	\N	/api/washingmachines	GET	2025-05-19 13:11:02.744587	2025-05-19 13:11:02.776552	200	31
6513	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:12:05.272132	2025-05-19 13:12:05.274792	200	2
6514	\N	/api/washingmachines	GET	2025-05-19 13:12:05.485756	2025-05-19 13:12:05.533153	200	47
6515	\N	/api/weather	GET	2025-05-19 13:12:05.549843	2025-05-19 13:12:05.549878	200	0
6516	\N	/api/restaurant	GET	2025-05-19 13:12:05.550671	2025-05-19 13:12:05.550713	200	0
6517	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:12:05.558266	2025-05-19 13:12:05.559647	200	1
6518	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 13:12:11.613807	2025-05-19 13:12:11.613907	200	0
6519	\N	/api/weather	GET	2025-05-19 13:12:35.151977	2025-05-19 13:12:35.151999	200	0
6520	\N	/api/restaurant	GET	2025-05-19 13:12:35.152233	2025-05-19 13:12:35.152301	200	0
6521	\N	/api/washingmachines	GET	2025-05-19 13:12:35.135744	2025-05-19 13:12:35.170822	200	35
6522	\N	/api/washingmachines	GET	2025-05-19 13:12:44.147339	2025-05-19 13:12:44.189312	200	41
6523	\N	/api/restaurant	GET	2025-05-19 13:12:52.171083	2025-05-19 13:12:52.171144	200	0
6524	\N	/api/restaurant	GET	2025-05-19 13:14:27.483565	2025-05-19 13:14:27.483668	200	0
6525	\N	/api/restaurant	GET	2025-05-19 13:15:05.147383	2025-05-19 13:15:05.147451	200	0
6526	\N	/api/weather	GET	2025-05-19 13:15:05.147128	2025-05-19 13:15:05.147158	200	0
6527	\N	/api/washingmachines	GET	2025-05-19 13:15:05.111472	2025-05-19 13:15:05.185072	200	73
6528	\N	/api/restaurant	GET	2025-05-19 13:15:12.540963	2025-05-19 13:15:12.54103	200	0
6529	\N	/api/restaurant	GET	2025-05-19 13:15:24.613632	2025-05-19 13:15:24.613697	200	0
6530	\N	/api/washingmachines	GET	2025-05-19 13:16:02.740732	2025-05-19 13:16:02.770797	200	30
6531	\N	/api/restaurant	GET	2025-05-19 13:16:28.013021	2025-05-19 13:16:28.013132	200	0
6532	\N	/api/restaurant	GET	2025-05-19 13:16:35.509797	2025-05-19 13:16:35.509874	200	0
6533	\N	/api/restaurant	GET	2025-05-19 13:16:47.27345	2025-05-19 13:16:47.273518	200	0
6534	\N	/api/washingmachines	GET	2025-05-19 13:17:44.145287	2025-05-19 13:17:44.217961	200	72
6535	\N	/api/washingmachines	GET	2025-05-19 13:21:02.74609	2025-05-19 13:21:02.821321	200	75
6536	\N	/api/washingmachines	GET	2025-05-19 13:22:44.144214	2025-05-19 13:22:44.213192	200	68
6537	\N	/api/washingmachines	GET	2025-05-19 13:26:02.736989	2025-05-19 13:26:02.790251	200	53
6538	\N	/api/washingmachines	GET	2025-05-19 13:27:44.144441	2025-05-19 13:27:44.234892	200	90
6539	\N	/api/washingmachines	GET	2025-05-19 13:31:02.747256	2025-05-19 13:31:02.804312	200	57
6540	\N	/api/washingmachines	GET	2025-05-19 13:32:44.153424	2025-05-19 13:32:44.222073	200	68
6541	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:34:14.340271	2025-05-19 13:34:14.34192	200	1
6542	\N	/api/washingmachines	GET	2025-05-19 13:34:14.53865	2025-05-19 13:34:14.61341	200	74
6543	\N	/api/restaurant	GET	2025-05-19 13:34:14.61454	2025-05-19 13:34:14.614638	200	0
6544	\N	/api/weather	GET	2025-05-19 13:34:14.613981	2025-05-19 13:34:14.66199	200	48
6545	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:34:14.616775	2025-05-19 13:34:14.766074	200	149
6546	\N	/api/washingmachines	GET	2025-05-19 13:34:20.081144	2025-05-19 13:34:20.108814	200	27
6547	\N	/api/restaurant	GET	2025-05-19 13:34:24.380637	2025-05-19 13:34:24.380703	200	0
6548	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:35:17.286964	2025-05-19 13:35:17.288947	200	1
6549	\N	/api/washingmachines	GET	2025-05-19 13:35:17.469199	2025-05-19 13:35:17.499707	200	30
6550	\N	/api/weather	GET	2025-05-19 13:35:17.535166	2025-05-19 13:35:17.535192	200	0
6551	\N	/api/restaurant	GET	2025-05-19 13:35:17.535724	2025-05-19 13:35:17.535803	200	0
6552	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:35:17.537857	2025-05-19 13:35:17.688369	200	150
6553	\N	/api/restaurant	GET	2025-05-19 13:35:18.535456	2025-05-19 13:35:18.535513	200	0
6554	\N	/api/washingmachines	GET	2025-05-19 13:36:02.74826	2025-05-19 13:36:02.781987	200	33
6555	\N	/api/washingmachines	GET	2025-05-19 13:37:44.147882	2025-05-19 13:37:44.226656	200	78
6556	\N	/api/restaurant	GET	2025-05-19 13:39:21.380524	2025-05-19 13:39:21.380651	200	0
6557	\N	/api/restaurant	GET	2025-05-19 13:40:34.652282	2025-05-19 13:40:34.652341	200	0
6558	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:40:42.269701	2025-05-19 13:40:42.271736	200	2
6559	\N	/api/washingmachines	GET	2025-05-19 13:40:42.455727	2025-05-19 13:40:42.523735	200	68
6560	\N	/api/restaurant	GET	2025-05-19 13:40:42.532349	2025-05-19 13:40:42.532452	200	0
6561	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 13:40:42.532598	2025-05-19 13:40:42.534407	200	1
6562	\N	/api/weather	GET	2025-05-19 13:40:42.532188	2025-05-19 13:40:42.609335	200	77
6563	\N	/api/restaurant	GET	2025-05-19 13:40:51.47416	2025-05-19 13:40:51.474222	200	0
6564	\N	/api/restaurant	GET	2025-05-19 13:40:53.115576	2025-05-19 13:40:53.115666	200	0
6565	\N	/api/restaurant	GET	2025-05-19 13:41:00.570122	2025-05-19 13:41:00.570194	200	0
6566	\N	/api/washingmachines	GET	2025-05-19 13:41:02.773009	2025-05-19 13:41:02.803184	200	30
6567	\N	/api/washingmachines	GET	2025-05-19 13:42:44.147482	2025-05-19 13:42:44.23843	200	90
6568	\N	/api/washingmachines	GET	2025-05-19 13:46:02.753466	2025-05-19 13:46:02.85011	200	96
6569	\N	/api/restaurant	GET	2025-05-19 13:46:55.475531	2025-05-19 13:46:55.475606	200	0
6570	\N	/api/washingmachines	GET	2025-05-19 13:47:44.146977	2025-05-19 13:47:44.218145	200	71
6571	\N	/api/restaurant	GET	2025-05-19 13:48:05.222676	2025-05-19 13:48:05.222772	200	0
6572	\N	/.git/config	GET	2025-05-19 13:48:55.182867	2025-05-19 13:48:55.182875	500	0
6573	\N	/api/washingmachines	GET	2025-05-19 13:51:02.756311	2025-05-19 13:51:02.843777	200	87
6574	\N	/api/washingmachines	GET	2025-05-19 13:52:44.15565	2025-05-19 13:52:44.224713	200	69
6575	\N	/api/restaurant	GET	2025-05-19 13:55:41.364819	2025-05-19 13:55:41.364921	200	0
6616	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:05:18.540789	2025-05-19 14:05:18.542776	200	1
6625	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:09:05.074506	2025-05-19 14:09:05.076679	200	2
6626	\N	/api/weather	GET	2025-05-19 14:09:05.067645	2025-05-19 14:09:05.138906	200	71
6627	\N	/api/washingmachines	GET	2025-05-19 14:09:06.028969	2025-05-19 14:09:06.05973	200	30
6628	\N	/api/restaurant	GET	2025-05-19 14:09:22.711725	2025-05-19 14:09:22.711784	200	0
6629	\N	/api/restaurant	GET	2025-05-19 14:09:53.098658	2025-05-19 14:09:53.098737	200	0
6630	\N	/api/washingmachines	GET	2025-05-19 14:11:02.906329	2025-05-19 14:11:02.994371	200	88
6631	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:12:16.013122	2025-05-19 14:12:16.015026	200	1
6632	\N	/api/washingmachines	GET	2025-05-19 14:12:16.20871	2025-05-19 14:12:16.250031	200	41
6633	\N	/api/weather	GET	2025-05-19 14:12:16.274675	2025-05-19 14:12:16.274707	200	0
6634	\N	/api/restaurant	GET	2025-05-19 14:12:16.280679	2025-05-19 14:12:16.280746	200	0
6648	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:16:29.411893	2025-05-19 14:16:29.414738	200	2
6649	\N	/api/restaurant	GET	2025-05-19 14:16:31.149673	2025-05-19 14:16:31.149728	200	0
6650	\N	/api/restaurant	GET	2025-05-19 14:17:07.560765	2025-05-19 14:17:07.560834	200	0
6651	\N	/api/washingmachines	GET	2025-05-19 14:17:44.155663	2025-05-19 14:17:44.187679	200	32
6652	\N	/api/washingmachines	GET	2025-05-19 14:21:02.757381	2025-05-19 14:21:02.850926	200	93
6653	\N	/api/washingmachines	GET	2025-05-19 14:22:44.151079	2025-05-19 14:22:44.234117	200	83
6654	\N	/api/washingmachines	GET	2025-05-19 14:26:02.751341	2025-05-19 14:26:02.828536	200	77
6655	\N	/api/washingmachines	GET	2025-05-19 14:27:44.153314	2025-05-19 14:27:44.246341	200	93
6656	\N	/api/washingmachines	GET	2025-05-19 14:31:02.755623	2025-05-19 14:31:02.834879	200	79
6657	\N	/api/restaurant	GET	2025-05-19 14:32:31.549021	2025-05-19 14:32:31.54915	200	0
6658	\N	/api/restaurant	GET	2025-05-19 14:32:31.664885	2025-05-19 14:32:31.664972	200	0
6659	\N	/api/washingmachines	GET	2025-05-19 14:32:44.161856	2025-05-19 14:32:44.234581	200	72
6660	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:33:32.954922	2025-05-19 14:33:32.956381	200	1
6661	\N	/api/washingmachines	GET	2025-05-19 14:33:33.199938	2025-05-19 14:33:33.229513	200	29
6662	\N	/api/restaurant	GET	2025-05-19 14:33:33.317956	2025-05-19 14:33:33.318042	200	0
6666	\N	/api/restaurant	GET	2025-05-19 14:33:41.291693	2025-05-19 14:33:41.291732	200	0
6667	\N	/api/washingmachines	GET	2025-05-19 14:33:41.268516	2025-05-19 14:33:41.300176	200	31
6668	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:33:50.64531	2025-05-19 14:33:50.647974	200	2
6669	\N	/api/washingmachines	GET	2025-05-19 14:33:50.822554	2025-05-19 14:33:50.861685	200	39
6670	\N	/api/weather	GET	2025-05-19 14:33:50.889439	2025-05-19 14:33:50.889464	200	0
6684	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:40:18.523439	2025-05-19 14:40:18.525198	200	1
6685	\N	/api/weather	GET	2025-05-19 14:40:18.521995	2025-05-19 14:40:18.578558	200	56
6687	\N	/api/restaurant	GET	2025-05-19 14:40:20.096188	2025-05-19 14:40:20.096278	200	0
6688	\N	/api/restaurant	GET	2025-05-19 14:40:36.291274	2025-05-19 14:40:36.291392	200	0
6689	\N	/api/restaurant	GET	2025-05-19 14:40:46.590088	2025-05-19 14:40:46.590165	200	0
6690	\N	/api/restaurant	GET	2025-05-19 14:40:46.974531	2025-05-19 14:40:46.974583	200	0
6691	\N	/api/washingmachines	GET	2025-05-19 14:41:02.837086	2025-05-19 14:41:02.866297	200	29
6692	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:41:07.831198	2025-05-19 14:41:07.832964	200	1
6693	\N	/api/washingmachines	GET	2025-05-19 14:41:08.008206	2025-05-19 14:41:08.038332	200	30
6694	\N	/api/weather	GET	2025-05-19 14:41:08.074644	2025-05-19 14:41:08.074671	200	0
6695	\N	/api/restaurant	GET	2025-05-19 14:41:08.081724	2025-05-19 14:41:08.081805	200	0
6708	\N	/api/restaurant	GET	2025-05-19 14:58:37.36765	2025-05-19 14:58:37.367743	200	0
6736	\N	/api/restaurant	GET	2025-05-19 15:14:18.053556	2025-05-19 15:14:18.053647	200	0
6763	\N	/api/restaurant	GET	2025-05-19 15:27:10.577886	2025-05-19 15:27:10.578001	200	0
6765	\N	/api/weather	GET	2025-05-19 15:27:10.57219	2025-05-19 15:27:10.624179	200	51
6766	enzo.morvan@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 15:27:12.460039	2025-05-19 15:27:12.467707	200	7
6767	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:27:17.585792	2025-05-19 15:27:17.587915	200	2
6768	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-19 15:27:17.796873	2025-05-19 15:27:17.796992	200	0
6769	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:27:21.091405	2025-05-19 15:27:21.093183	200	1
6770	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-19 15:27:21.265591	2025-05-19 15:27:21.26566	200	0
6771	\N	/api/washingmachines	GET	2025-05-19 15:27:44.156001	2025-05-19 15:27:44.261279	200	105
6772	\N	/api/restaurant	GET	2025-05-19 15:27:47.731364	2025-05-19 15:27:47.731423	200	0
6773	\N	/api/washingmachines	GET	2025-05-19 15:31:03.030589	2025-05-19 15:31:03.120999	200	90
6774	\N	/api/washingmachines	GET	2025-05-19 15:32:44.497571	2025-05-19 15:32:44.568821	200	71
6775	\N	/api/washingmachines	GET	2025-05-19 15:36:02.979911	2025-05-19 15:36:03.05235	200	72
6776	\N	/api/washingmachines	GET	2025-05-19 15:37:44.302432	2025-05-19 15:37:44.375063	200	72
6778	\N	/api/restaurant	GET	2025-05-19 15:39:47.351559	2025-05-19 15:39:47.351658	200	0
6779	\N	/api/washingmachines	GET	2025-05-19 15:39:50.813736	2025-05-19 15:39:50.912168	200	98
6780	\N	/status	GET	2025-05-19 15:40:14.400126	2025-05-19 15:40:14.40013	200	0
6781	\N	/api/statistics/global	GET	2025-05-19 15:40:14.448201	2025-05-19 15:40:14.450501	200	2
6782	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:14.498113	2025-05-19 15:40:14.505283	200	7
6783	\N	/status	GET	2025-05-19 15:40:16.493949	2025-05-19 15:40:16.493951	200	0
6784	\N	/api/statistics/global	GET	2025-05-19 15:40:16.547166	2025-05-19 15:40:16.549118	200	1
6785	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:16.599057	2025-05-19 15:40:16.606405	200	7
6786	\N	/status	GET	2025-05-19 15:40:17.459837	2025-05-19 15:40:17.459839	200	0
6787	\N	/api/statistics/global	GET	2025-05-19 15:40:17.507362	2025-05-19 15:40:17.509223	200	1
6788	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:17.560553	2025-05-19 15:40:17.565698	200	5
6789	\N	/status	GET	2025-05-19 15:40:17.998241	2025-05-19 15:40:17.998242	200	0
6790	\N	/api/statistics/global	GET	2025-05-19 15:40:18.054133	2025-05-19 15:40:18.056271	200	2
6791	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:18.11355	2025-05-19 15:40:18.119229	200	5
6792	\N	/status	GET	2025-05-19 15:40:18.393529	2025-05-19 15:40:18.393531	200	0
6793	\N	/api/statistics/global	GET	2025-05-19 15:40:18.447322	2025-05-19 15:40:18.449607	200	2
6794	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:18.50823	2025-05-19 15:40:18.514423	200	6
6795	\N	/status	GET	2025-05-19 15:40:18.759157	2025-05-19 15:40:18.759159	200	0
6796	\N	/api/statistics/global	GET	2025-05-19 15:40:18.819951	2025-05-19 15:40:18.821855	200	1
6797	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:18.873913	2025-05-19 15:40:18.87881	200	4
6798	\N	/status	GET	2025-05-19 15:40:18.967061	2025-05-19 15:40:18.967062	200	0
6799	\N	/api/statistics/global	GET	2025-05-19 15:40:19.019933	2025-05-19 15:40:19.022534	200	2
6800	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:19.071576	2025-05-19 15:40:19.079066	200	7
6801	\N	/status	GET	2025-05-19 15:40:19.138758	2025-05-19 15:40:19.13876	200	0
6802	\N	/api/statistics/global	GET	2025-05-19 15:40:19.188315	2025-05-19 15:40:19.191042	200	2
6803	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:19.257929	2025-05-19 15:40:19.266138	200	8
6804	\N	/status	GET	2025-05-19 15:40:19.322962	2025-05-19 15:40:19.322963	200	0
6805	\N	/status	GET	2025-05-19 15:40:19.385435	2025-05-19 15:40:19.385437	200	0
6576	\N	/api/restaurant	GET	2025-05-19 13:55:41.365781	2025-05-19 13:55:41.365869	200	0
6577	\N	/api/restaurant	GET	2025-05-19 13:55:41.415399	2025-05-19 13:55:41.415456	200	0
6578	\N	/api/restaurant	GET	2025-05-19 13:55:42.200613	2025-05-19 13:55:42.200677	200	0
6579	\N	/api/restaurant	GET	2025-05-19 13:55:44.533263	2025-05-19 13:55:44.533316	200	0
6580	\N	/api/washingmachines	GET	2025-05-19 13:56:02.750498	2025-05-19 13:56:02.823415	200	72
6581	\N	/api/restaurant	GET	2025-05-19 13:56:33.267374	2025-05-19 13:56:33.267449	200	0
6582	\N	/api/restaurant	GET	2025-05-19 13:56:33.550662	2025-05-19 13:56:33.550722	200	0
6583	\N	/api/restaurant	GET	2025-05-19 13:56:33.837381	2025-05-19 13:56:33.837455	200	0
6584	\N	/api/restaurant	GET	2025-05-19 13:56:34.366451	2025-05-19 13:56:34.366522	200	0
6585	\N	/api/restaurant	GET	2025-05-19 13:57:04.005211	2025-05-19 13:57:04.005262	200	0
6586	\N	/api/restaurant	GET	2025-05-19 13:57:07.888093	2025-05-19 13:57:07.888219	200	0
6587	\N	/api/restaurant	GET	2025-05-19 13:57:08.168137	2025-05-19 13:57:08.168208	200	0
6588	\N	/api/restaurant	GET	2025-05-19 13:57:08.450751	2025-05-19 13:57:08.450808	200	0
6589	\N	/api/restaurant	GET	2025-05-19 13:57:11.367702	2025-05-19 13:57:11.367764	200	0
6590	\N	/api/restaurant	GET	2025-05-19 13:57:13.203822	2025-05-19 13:57:13.203891	200	0
6591	\N	/api/restaurant	GET	2025-05-19 13:57:15.406029	2025-05-19 13:57:15.406081	200	0
6592	\N	/api/washingmachines	GET	2025-05-19 13:57:44.150964	2025-05-19 13:57:44.234097	200	83
6593	\N	/api/restaurant	GET	2025-05-19 13:57:56.944247	2025-05-19 13:57:56.944316	200	0
6594	\N	/api/restaurant	GET	2025-05-19 13:57:57.218021	2025-05-19 13:57:57.218081	200	0
6595	\N	/api/restaurant	GET	2025-05-19 13:57:57.504031	2025-05-19 13:57:57.504089	200	0
6596	\N	/api/restaurant	GET	2025-05-19 13:57:57.782855	2025-05-19 13:57:57.782933	200	0
6597	\N	/api/restaurant	GET	2025-05-19 13:57:58.158483	2025-05-19 13:57:58.158556	200	0
6598	\N	/api/washingmachines	GET	2025-05-19 13:58:37.042286	2025-05-19 13:58:37.075554	200	33
6599	\N	/api/washingmachines	GET	2025-05-19 14:01:02.749397	2025-05-19 14:01:02.827475	200	78
6600	\N	/api/restaurant	GET	2025-05-19 14:01:05.66398	2025-05-19 14:01:05.66408	200	0
6601	\N	/api/washingmachines	GET	2025-05-19 14:01:05.631682	2025-05-19 14:01:05.679469	200	47
6602	\N	/api/weather	GET	2025-05-19 14:01:05.661625	2025-05-19 14:01:05.702369	200	40
6603	\N	/api/restaurant	GET	2025-05-19 14:01:26.642781	2025-05-19 14:01:26.64284	200	0
6604	\N	/api/restaurant	GET	2025-05-19 14:01:42.050854	2025-05-19 14:01:42.050908	200	0
6605	\N	/api/restaurant	GET	2025-05-19 14:01:42.782908	2025-05-19 14:01:42.782968	200	0
6606	\N	/api/restaurant	GET	2025-05-19 14:01:45.169628	2025-05-19 14:01:45.16968	200	0
6607	\N	/api/restaurant	GET	2025-05-19 14:01:48.432713	2025-05-19 14:01:48.432783	200	0
6608	\N	/api/restaurant	GET	2025-05-19 14:01:59.245349	2025-05-19 14:01:59.245421	200	0
6609	\N	/api/washingmachines	GET	2025-05-19 14:02:44.149934	2025-05-19 14:02:44.237766	200	87
6610	\N	/api/restaurant	GET	2025-05-19 14:02:58.137517	2025-05-19 14:02:58.137588	200	0
6611	\N	/api/restaurant	GET	2025-05-19 14:04:23.589372	2025-05-19 14:04:23.589461	200	0
6612	\N	/api/restaurant	GET	2025-05-19 14:04:56.910999	2025-05-19 14:04:56.911064	200	0
6613	\N	/api/restaurant	GET	2025-05-19 14:04:57.184327	2025-05-19 14:04:57.184405	200	0
6614	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:05:18.247818	2025-05-19 14:05:18.249873	200	2
6615	\N	/api/restaurant	GET	2025-05-19 14:05:18.540474	2025-05-19 14:05:18.540564	200	0
6617	\N	/api/washingmachines	GET	2025-05-19 14:05:18.451143	2025-05-19 14:05:18.545428	200	94
6618	\N	/api/weather	GET	2025-05-19 14:05:18.540347	2025-05-19 14:05:18.618229	200	77
6619	\N	/api/restaurant	GET	2025-05-19 14:05:23.931648	2025-05-19 14:05:23.931706	200	0
6620	\N	/api/washingmachines	GET	2025-05-19 14:06:03.059197	2025-05-19 14:06:03.089854	200	30
6621	\N	/api/washingmachines	GET	2025-05-19 14:07:44.415365	2025-05-19 14:07:44.479296	200	63
6622	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:09:04.786373	2025-05-19 14:09:04.788362	200	1
6623	\N	/api/washingmachines	GET	2025-05-19 14:09:04.993563	2025-05-19 14:09:05.024303	200	30
6624	\N	/api/restaurant	GET	2025-05-19 14:09:05.073699	2025-05-19 14:09:05.073791	200	0
6635	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:12:16.280871	2025-05-19 14:12:16.282389	200	1
6636	\N	/api/restaurant	GET	2025-05-19 14:12:17.319087	2025-05-19 14:12:17.319148	200	0
6637	\N	/api/restaurant	GET	2025-05-19 14:12:44.038768	2025-05-19 14:12:44.038848	200	0
6638	\N	/api/washingmachines	GET	2025-05-19 14:12:44.145173	2025-05-19 14:12:44.192397	200	47
6639	\N	/api/restaurant	GET	2025-05-19 14:14:07.522754	2025-05-19 14:14:07.522865	200	0
6640	\N	/api/restaurant	GET	2025-05-19 14:14:47.465232	2025-05-19 14:14:47.46532	200	0
6641	\N	/api/weather	GET	2025-05-19 14:14:47.464433	2025-05-19 14:14:47.521282	200	56
6642	\N	/api/washingmachines	GET	2025-05-19 14:14:47.428697	2025-05-19 14:14:47.525793	200	97
6643	\N	/api/washingmachines	GET	2025-05-19 14:16:02.76273	2025-05-19 14:16:02.793888	200	31
6644	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:16:29.105819	2025-05-19 14:16:29.10781	200	1
6645	\N	/api/washingmachines	GET	2025-05-19 14:16:29.328523	2025-05-19 14:16:29.35934	200	30
6646	\N	/api/weather	GET	2025-05-19 14:16:29.404377	2025-05-19 14:16:29.404414	200	0
6647	\N	/api/restaurant	GET	2025-05-19 14:16:29.41095	2025-05-19 14:16:29.411034	200	0
6663	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:33:33.319276	2025-05-19 14:33:33.321175	200	1
6664	\N	/api/weather	GET	2025-05-19 14:33:33.311888	2025-05-19 14:33:33.39214	200	80
6665	\N	/api/weather	GET	2025-05-19 14:33:41.29148	2025-05-19 14:33:41.291506	200	0
6671	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:33:50.896491	2025-05-19 14:33:50.898439	200	1
6672	\N	/api/restaurant	GET	2025-05-19 14:33:50.896467	2025-05-19 14:33:50.896559	200	0
6673	\N	/api/restaurant	GET	2025-05-19 14:33:52.135432	2025-05-19 14:33:52.135492	200	0
6674	\N	/api/restaurant	GET	2025-05-19 14:34:26.772291	2025-05-19 14:34:26.772334	200	0
6675	\N	/api/restaurant	GET	2025-05-19 14:35:47.610976	2025-05-19 14:35:47.61104	200	0
6676	\N	/api/restaurant	GET	2025-05-19 14:35:56.363182	2025-05-19 14:35:56.363316	200	0
6677	\N	/api/washingmachines	GET	2025-05-19 14:36:02.868239	2025-05-19 14:36:02.940738	200	72
6678	\N	/api/restaurant	GET	2025-05-19 14:36:23.456512	2025-05-19 14:36:23.456582	200	0
6679	\N	/api/restaurant	GET	2025-05-19 14:36:23.742404	2025-05-19 14:36:23.742476	200	0
6680	\N	/api/restaurant	GET	2025-05-19 14:36:24.020846	2025-05-19 14:36:24.020902	200	0
6681	\N	/api/washingmachines	GET	2025-05-19 14:37:44.309556	2025-05-19 14:37:44.39911	200	89
6682	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:40:18.230544	2025-05-19 14:40:18.232476	200	1
6683	\N	/api/restaurant	GET	2025-05-19 14:40:18.522425	2025-05-19 14:40:18.522535	200	0
6686	\N	/api/washingmachines	GET	2025-05-19 14:40:18.450686	2025-05-19 14:40:18.524351	200	73
6696	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:41:08.081941	2025-05-19 14:41:08.084136	200	2
6697	\N	/api/restaurant	GET	2025-05-19 14:41:10.024335	2025-05-19 14:41:10.024438	200	0
6698	\N	/api/restaurant	GET	2025-05-19 14:41:31.276461	2025-05-19 14:41:31.276543	200	0
6699	\N	/api/washingmachines	GET	2025-05-19 14:42:44.153941	2025-05-19 14:42:44.21038	200	56
6700	\N	/api/restaurant	GET	2025-05-19 14:44:17.431082	2025-05-19 14:44:17.431178	200	0
6701	\N	/api/washingmachines	GET	2025-05-19 14:46:02.75694	2025-05-19 14:46:02.850121	200	93
6702	\N	/api/washingmachines	GET	2025-05-19 14:47:44.286745	2025-05-19 14:47:44.359991	200	73
6703	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:48:47.225058	2025-05-19 14:48:47.22697	200	1
6704	\N	/api/washingmachines	GET	2025-05-19 14:51:02.759956	2025-05-19 14:51:02.833495	200	73
6705	\N	/api/washingmachines	GET	2025-05-19 14:52:44.154118	2025-05-19 14:52:44.221894	200	67
6706	\N	/api/washingmachines	GET	2025-05-19 14:56:02.756365	2025-05-19 14:56:02.828548	200	72
6707	\N	/api/washingmachines	GET	2025-05-19 14:57:44.156178	2025-05-19 14:57:44.24717	200	90
6709	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 14:58:37.367605	2025-05-19 14:58:37.369639	200	2
6710	\N	/api/washingmachines	GET	2025-05-19 14:58:37.357316	2025-05-19 14:58:37.387167	200	29
6711	\N	/api/weather	GET	2025-05-19 14:58:37.357282	2025-05-19 14:58:37.427448	200	70
6712	\N	/api/restaurant	GET	2025-05-19 14:58:47.549727	2025-05-19 14:58:47.549787	200	0
6713	\N	/api/washingmachines	GET	2025-05-19 14:58:51.809612	2025-05-19 14:58:51.838475	200	28
6714	\N	/api/washingmachines	GET	2025-05-19 14:59:38.930107	2025-05-19 14:59:38.960277	200	30
6715	\N	/api/washingmachines	GET	2025-05-19 14:59:43.703523	2025-05-19 14:59:43.734323	200	30
6716	\N	/api/washingmachines	GET	2025-05-19 15:00:02.678423	2025-05-19 15:00:02.709262	200	30
6717	\N	/api/washingmachines	GET	2025-05-19 15:00:33.141794	2025-05-19 15:00:33.17219	200	30
6718	\N	/api/washingmachines	GET	2025-05-19 15:00:40.488892	2025-05-19 15:00:40.521407	200	32
6719	\N	/api/washingmachines	GET	2025-05-19 15:00:55.604643	2025-05-19 15:00:55.634449	200	29
6720	\N	/api/washingmachines	GET	2025-05-19 15:01:02.776634	2025-05-19 15:01:02.807417	200	30
6721	\N	/api/washingmachines	GET	2025-05-19 15:02:44.157154	2025-05-19 15:02:44.249167	200	92
6722	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:03:49.360996	2025-05-19 15:03:49.363216	200	2
6723	\N	/api/washingmachines	GET	2025-05-19 15:06:02.756463	2025-05-19 15:06:02.83616	200	79
6724	\N	/api/restaurant	GET	2025-05-19 15:07:20.524646	2025-05-19 15:07:20.524778	200	0
6725	\N	/api/restaurant	GET	2025-05-19 15:07:21.827203	2025-05-19 15:07:21.827272	200	0
6726	\N	/api/restaurant	GET	2025-05-19 15:07:27.686658	2025-05-19 15:07:27.686733	200	0
6727	\N	/api/washingmachines	GET	2025-05-19 15:07:44.163327	2025-05-19 15:07:44.235397	200	72
6728	\N	/api/washingmachines	GET	2025-05-19 15:11:02.766707	2025-05-19 15:11:02.856482	200	89
6729	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:11:19.659341	2025-05-19 15:11:19.660886	200	1
6730	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:11:41.977845	2025-05-19 15:11:41.979427	200	1
6731	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 15:12:26.520905	2025-05-19 15:12:26.525922	200	5
6732	\N	/api/washingmachines	GET	2025-05-19 15:12:44.15533	2025-05-19 15:12:44.225323	200	69
6733	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:13:12.46121	2025-05-19 15:13:12.462807	200	1
6734	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:13:24.643116	2025-05-19 15:13:24.646488	200	3
6735	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:14:17.886883	2025-05-19 15:14:17.888486	200	1
6737	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:14:18.052844	2025-05-19 15:14:18.054395	200	1
6738	\N	/api/washingmachines	GET	2025-05-19 15:14:18.003033	2025-05-19 15:14:18.100186	200	97
6739	\N	/api/weather	GET	2025-05-19 15:14:18.051606	2025-05-19 15:14:18.120486	200	68
6740	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:14:18.491546	2025-05-19 15:14:18.494578	200	3
6741	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:14:19.297131	2025-05-19 15:14:19.299803	200	2
6742	\N	/api/washingmachines	GET	2025-05-19 15:16:02.771156	2025-05-19 15:16:02.849488	200	78
6743	\N	/api/washingmachines	GET	2025-05-19 15:17:44.163001	2025-05-19 15:17:44.24499	200	81
6744	\N	/status	GET	2025-05-19 15:17:56.190136	2025-05-19 15:17:56.19014	200	0
6745	\N	/api/statistics/global	GET	2025-05-19 15:17:56.250247	2025-05-19 15:17:56.252667	200	2
6746	\N	/api/statistics/top-users	GET	2025-05-19 15:17:56.319574	2025-05-19 15:17:56.321533	200	1
6747	\N	/status	GET	2025-05-19 15:18:01.883434	2025-05-19 15:18:01.883436	200	0
6748	\N	/api/statistics/global	GET	2025-05-19 15:18:01.940431	2025-05-19 15:18:01.944834	200	4
6749	\N	/api/statistics/top-users	GET	2025-05-19 15:18:02.000564	2025-05-19 15:18:02.002087	200	1
6750	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:18:20.476464	2025-05-19 15:18:20.482708	200	6
6751	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:18:20.743476	2025-05-19 15:18:20.745308	200	1
6752	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:18:21.025767	2025-05-19 15:18:21.027683	200	1
6753	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:18:21.329322	2025-05-19 15:18:21.331251	200	1
6754	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:18:21.8993	2025-05-19 15:18:21.901487	200	2
6755	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:20:58.770523	2025-05-19 15:20:58.772327	200	1
6756	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:20:59.070656	2025-05-19 15:20:59.072508	200	1
6757	\N	/api/washingmachines	GET	2025-05-19 15:21:02.764253	2025-05-19 15:21:02.832642	200	68
6758	\N	/api/washingmachines	GET	2025-05-19 15:22:44.156222	2025-05-19 15:22:44.22394	200	67
6759	\N	/api/washingmachines	GET	2025-05-19 15:26:02.758775	2025-05-19 15:26:02.855446	200	96
6760	\N	/api/auth/login	POST	2025-05-19 15:27:10.048871	2025-05-19 15:27:10.122529	200	73
6761	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:27:10.173968	2025-05-19 15:27:10.176391	200	2
6762	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:27:10.231504	2025-05-19 15:27:10.23338	200	1
6764	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 15:27:10.577927	2025-05-19 15:27:10.579524	200	1
6777	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 15:39:47.164396	2025-05-19 15:39:47.16451	200	0
6808	\N	/api/statistics/global	GET	2025-05-19 15:40:19.456252	2025-05-19 15:40:19.458169	200	1
6809	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:19.462205	2025-05-19 15:40:19.467414	200	5
6810	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:19.516081	2025-05-19 15:40:19.527126	200	11
6811	\N	/status	GET	2025-05-19 15:40:19.709763	2025-05-19 15:40:19.70977	200	0
6812	\N	/api/statistics/global	GET	2025-05-19 15:40:19.762379	2025-05-19 15:40:19.765662	200	3
6813	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:19.814816	2025-05-19 15:40:19.819984	200	5
6814	\N	/status	GET	2025-05-19 15:40:19.908587	2025-05-19 15:40:19.908589	200	0
6815	\N	/api/statistics/global	GET	2025-05-19 15:40:19.961228	2025-05-19 15:40:19.963167	200	1
6816	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:20.014974	2025-05-19 15:40:20.019577	200	4
6817	\N	/status	GET	2025-05-19 15:40:20.106716	2025-05-19 15:40:20.106718	200	0
6818	\N	/api/statistics/global	GET	2025-05-19 15:40:20.159727	2025-05-19 15:40:20.161797	200	2
6819	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:20.210952	2025-05-19 15:40:20.215837	200	4
6820	\N	/status	GET	2025-05-19 15:40:20.324534	2025-05-19 15:40:20.324536	200	0
6821	\N	/api/statistics/global	GET	2025-05-19 15:40:20.376929	2025-05-19 15:40:20.379339	200	2
6822	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:20.434396	2025-05-19 15:40:20.44105	200	6
6823	\N	/status	GET	2025-05-19 15:40:20.545662	2025-05-19 15:40:20.545663	200	0
6824	\N	/api/statistics/global	GET	2025-05-19 15:40:20.598341	2025-05-19 15:40:20.600452	200	2
6825	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:20.660452	2025-05-19 15:40:20.665467	200	5
6826	\N	/status	GET	2025-05-19 15:40:20.757989	2025-05-19 15:40:20.757991	200	0
6827	\N	/api/statistics/global	GET	2025-05-19 15:40:20.809454	2025-05-19 15:40:20.81174	200	2
6828	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:20.862814	2025-05-19 15:40:20.868553	200	5
6829	\N	/status	GET	2025-05-19 15:40:20.957088	2025-05-19 15:40:20.957089	200	0
6830	\N	/api/statistics/global	GET	2025-05-19 15:40:21.010845	2025-05-19 15:40:21.013049	200	2
6831	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:21.074631	2025-05-19 15:40:21.08101	200	6
6832	\N	/status	GET	2025-05-19 15:40:21.156642	2025-05-19 15:40:21.156644	200	0
6833	\N	/api/statistics/global	GET	2025-05-19 15:40:21.205267	2025-05-19 15:40:21.207203	200	1
6834	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:21.259064	2025-05-19 15:40:21.265237	200	6
6835	\N	/status	GET	2025-05-19 15:40:21.356588	2025-05-19 15:40:21.35659	200	0
6836	\N	/api/statistics/global	GET	2025-05-19 15:40:21.406322	2025-05-19 15:40:21.408217	200	1
6837	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:21.464943	2025-05-19 15:40:21.470858	200	5
6838	\N	/status	GET	2025-05-19 15:40:21.572394	2025-05-19 15:40:21.572395	200	0
6839	\N	/api/statistics/global	GET	2025-05-19 15:40:21.626108	2025-05-19 15:40:21.6283	200	2
6840	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:21.678287	2025-05-19 15:40:21.683355	200	5
6841	\N	/status	GET	2025-05-19 15:40:21.771758	2025-05-19 15:40:21.77176	200	0
6842	\N	/api/statistics/global	GET	2025-05-19 15:40:21.825292	2025-05-19 15:40:21.827389	200	2
6843	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:21.878118	2025-05-19 15:40:21.885054	200	6
6844	\N	/status	GET	2025-05-19 15:40:21.992109	2025-05-19 15:40:21.992111	200	0
6845	\N	/api/statistics/global	GET	2025-05-19 15:40:22.041153	2025-05-19 15:40:22.043226	200	2
6846	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:22.095773	2025-05-19 15:40:22.10103	200	5
6847	\N	/status	GET	2025-05-19 15:40:22.325969	2025-05-19 15:40:22.325971	200	0
6848	\N	/api/statistics/global	GET	2025-05-19 15:40:22.377095	2025-05-19 15:40:22.379255	200	2
6849	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:22.432544	2025-05-19 15:40:22.438593	200	6
6850	\N	/status	GET	2025-05-19 15:40:22.744112	2025-05-19 15:40:22.744114	200	0
6851	\N	/api/statistics/global	GET	2025-05-19 15:40:22.794514	2025-05-19 15:40:22.79692	200	2
6852	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:22.846129	2025-05-19 15:40:22.851449	200	5
6853	\N	/status	GET	2025-05-19 15:40:23.628651	2025-05-19 15:40:23.628653	200	0
6854	\N	/api/statistics/global	GET	2025-05-19 15:40:23.682446	2025-05-19 15:40:23.684384	200	1
7620	\N	/status	GET	2025-05-20 08:49:46.449469	2025-05-20 08:49:46.44947	200	0
6806	\N	/api/statistics/global	GET	2025-05-19 15:40:19.39254	2025-05-19 15:40:19.395232	200	2
6807	\N	/status	GET	2025-05-19 15:40:19.456179	2025-05-19 15:40:19.456181	200	0
6858	\N	/api/statistics/global	GET	2025-05-19 15:40:24.436558	2025-05-19 15:40:24.438955	200	2
6859	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:24.505155	2025-05-19 15:40:24.510225	200	5
6860	\N	/status	GET	2025-05-19 15:40:29.373794	2025-05-19 15:40:29.373797	200	0
6861	\N	/status	GET	2025-05-19 15:40:29.422655	2025-05-19 15:40:29.422657	200	0
6866	\N	/api/statistics/global	GET	2025-05-19 15:40:34.441649	2025-05-19 15:40:34.443436	200	1
6867	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:34.496014	2025-05-19 15:40:34.504428	200	8
6868	\N	/status	GET	2025-05-19 15:40:39.387819	2025-05-19 15:40:39.387821	200	0
6869	\N	/status	GET	2025-05-19 15:40:39.448752	2025-05-19 15:40:39.448754	200	0
6872	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 15:40:44.364919	2025-05-19 15:40:44.365021	200	0
6875	\N	/api/statistics/global	GET	2025-05-19 15:40:47.081198	2025-05-19 15:40:47.083265	200	2
6876	\N	/api/statistics/global	GET	2025-05-19 15:40:47.133512	2025-05-19 15:40:47.135506	200	1
6877	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:47.133627	2025-05-19 15:40:47.139995	200	6
6878	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:47.185915	2025-05-19 15:40:47.191101	200	5
6879	\N	/status	GET	2025-05-19 15:40:51.974857	2025-05-19 15:40:51.974859	200	0
6880	\N	/status	GET	2025-05-19 15:40:52.01983	2025-05-19 15:40:52.019832	200	0
6881	\N	/api/statistics/global	GET	2025-05-19 15:40:52.05956	2025-05-19 15:40:52.06168	200	2
6882	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:52.10878	2025-05-19 15:40:52.11681	200	8
6883	\N	/status	GET	2025-05-19 15:40:56.964672	2025-05-19 15:40:56.964674	200	0
6884	\N	/status	GET	2025-05-19 15:40:57.011206	2025-05-19 15:40:57.011208	200	0
6885	\N	/api/statistics/global	GET	2025-05-19 15:40:57.069935	2025-05-19 15:40:57.071946	200	2
6886	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:57.138788	2025-05-19 15:40:57.143955	200	5
6887	\N	/status	GET	2025-05-19 15:41:01.964349	2025-05-19 15:41:01.964351	200	0
6888	\N	/status	GET	2025-05-19 15:41:02.009208	2025-05-19 15:41:02.00921	200	0
6889	\N	/api/statistics/global	GET	2025-05-19 15:41:02.053079	2025-05-19 15:41:02.055716	200	2
6890	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:02.098955	2025-05-19 15:41:02.104264	200	5
6891	\N	/api/washingmachines	GET	2025-05-19 15:41:03.008542	2025-05-19 15:41:03.04058	200	32
6892	\N	/status	GET	2025-05-19 15:41:06.966444	2025-05-19 15:41:06.966445	200	0
6893	\N	/status	GET	2025-05-19 15:41:07.011118	2025-05-19 15:41:07.01112	200	0
6894	\N	/api/statistics/global	GET	2025-05-19 15:41:07.052323	2025-05-19 15:41:07.054813	200	2
6895	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:07.098635	2025-05-19 15:41:07.104135	200	5
6955	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 16:45:11.738203	2025-05-19 16:45:11.738286	200	0
6975	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 16:57:36.922442	2025-05-19 16:57:36.927596	200	5
6976	\N	/api/restaurant	GET	2025-05-19 16:57:38.496014	2025-05-19 16:57:38.496072	200	0
6977	\N	/api/restaurant	GET	2025-05-19 16:57:38.93469	2025-05-19 16:57:38.934747	200	0
6978	\N	/api/auth/register	POST	2025-05-19 16:59:08.181808	2025-05-19 16:59:08.260897	201	79
6979	\N	/api/auth/verify-account	POST	2025-05-19 17:00:29.064503	2025-05-19 17:00:29.072954	200	8
6980	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:00:29.223379	2025-05-19 17:00:29.225716	200	2
6981	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:00:29.357604	2025-05-19 17:00:29.359187	200	1
6982	\N	/api/washingmachines	GET	2025-05-19 17:00:29.497098	2025-05-19 17:00:29.586577	200	89
6983	\N	/api/restaurant	GET	2025-05-19 17:00:29.777578	2025-05-19 17:00:29.777658	200	0
7020	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:10:15.342414	2025-05-19 17:10:15.344216	200	1
7021	\N	/api/weather	GET	2025-05-19 17:10:15.337624	2025-05-19 17:10:15.388393	200	50
7022	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:10:29.526434	2025-05-19 17:10:29.528047	200	1
7023	\N	/api/restaurant	GET	2025-05-19 17:10:37.952218	2025-05-19 17:10:37.952289	200	0
7024	\N	/status	GET	2025-05-19 17:21:59.427441	2025-05-19 17:21:59.427444	200	0
7025	\N	/status	GET	2025-05-19 17:22:09.290299	2025-05-19 17:22:09.290301	200	0
7026	\N	/status	GET	2025-05-19 17:22:19.277073	2025-05-19 17:22:19.277074	200	0
7027	\N	/status	GET	2025-05-19 17:22:29.256118	2025-05-19 17:22:29.25612	200	0
7028	\N	/status	GET	2025-05-19 17:22:39.280807	2025-05-19 17:22:39.280809	200	0
7029	\N	/api/auth/register	POST	2025-05-19 17:24:03.078814	2025-05-19 17:24:03.158737	201	79
7030	\N	/api/auth/verify-account	POST	2025-05-19 17:24:29.47473	2025-05-19 17:24:29.482976	200	8
7031	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:24:29.593584	2025-05-19 17:24:29.596102	200	2
7032	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:24:29.678465	2025-05-19 17:24:29.680141	200	1
7033	\N	/api/washingmachines	GET	2025-05-19 17:24:29.785217	2025-05-19 17:24:29.878177	200	92
7034	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:24:30.074485	2025-05-19 17:24:30.076271	200	1
7035	\N	/api/weather	GET	2025-05-19 17:24:30.059289	2025-05-19 17:24:30.117074	200	57
7036	\N	/api/restaurant	GET	2025-05-19 17:24:30.073971	2025-05-19 17:24:31.424715	200	1350
7059	\N	/api/restaurant	GET	2025-05-19 17:44:33.048792	2025-05-19 17:44:33.04888	200	0
7065	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-19 17:44:36.428892	2025-05-19 17:44:36.42898	200	0
7066	\N	/api/restaurant	GET	2025-05-19 17:44:39.664764	2025-05-19 17:44:39.664869	200	0
7067	\N	/api/traq/	GET	2025-05-19 17:44:45.848911	2025-05-19 17:44:45.850128	200	1
7068	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:44:53.086114	2025-05-19 17:44:53.087705	200	1
7069	\N	/api/restaurant	GET	2025-05-19 17:45:03.496837	2025-05-19 17:45:03.496895	200	0
7070	\N	/api/washingmachines	GET	2025-05-19 17:45:32.431412	2025-05-19 17:45:32.462406	200	30
7071	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:45:50.298096	2025-05-19 17:45:50.300317	200	2
7072	\N	/api/washingmachines	GET	2025-05-19 17:45:55.519568	2025-05-19 17:45:55.553321	200	33
7074	\N	/status	GET	2025-05-19 17:48:00.3352	2025-05-19 17:48:00.335202	200	0
7075	\N	/status	GET	2025-05-19 17:48:04.92815	2025-05-19 17:48:04.928151	200	0
7076	\N	/api/statistics/global	GET	2025-05-19 17:48:05.989464	2025-05-19 17:48:05.991901	200	2
7077	\N	/api/statistics/top-users	GET	2025-05-19 17:48:06.069448	2025-05-19 17:48:06.070902	200	1
7078	\N	/status	GET	2025-05-19 17:48:07.313674	2025-05-19 17:48:07.313675	200	0
7079	\N	/status	GET	2025-05-19 17:48:17.634567	2025-05-19 17:48:17.634569	200	0
7080	\N	/status	GET	2025-05-19 17:48:27.677437	2025-05-19 17:48:27.677439	200	0
7081	\N	/api/traq/	GET	2025-05-19 17:48:27.93282	2025-05-19 17:48:27.933978	200	1
7082	\N	/status	GET	2025-05-19 17:48:37.327912	2025-05-19 17:48:37.327914	200	0
7083	\N	/api/washingmachines	GET	2025-05-19 17:48:45.58326	2025-05-19 17:48:45.654753	200	71
7084	\N	/status	GET	2025-05-19 17:48:47.322734	2025-05-19 17:48:47.322735	200	0
7085	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:49:01.690546	2025-05-19 17:49:01.692122	200	1
7086	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:49:24.09154	2025-05-19 17:49:24.095172	200	3
7087	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:50:19.20262	2025-05-19 17:50:19.204306	200	1
7088	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:51:17.440101	2025-05-19 17:51:17.444438	200	4
7089	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:51:54.354945	2025-05-19 17:51:54.356715	200	1
7090	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:53:23.456367	2025-05-19 17:53:23.458055	200	1
7091	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:53:29.397514	2025-05-19 17:53:29.399486	200	1
7092	aurelien.moignet@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 17:53:32.29131	2025-05-19 17:53:32.294358	200	3
7093	aurelien.moignet@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-19 17:53:32.903761	2025-05-19 17:53:32.908683	200	4
7094	aurelien.moignet@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-19 17:53:33.594139	2025-05-19 17:53:33.600361	200	6
7095	\N	/status	GET	2025-05-19 17:53:42.744967	2025-05-19 17:53:42.744973	200	0
7096	\N	/api/statistics/global	GET	2025-05-19 17:53:42.856821	2025-05-19 17:53:42.859076	200	2
7097	\N	/api/statistics/top-users	GET	2025-05-19 17:53:42.970058	2025-05-19 17:53:42.97177	200	1
7102	\N	/status	GET	2025-05-19 18:03:48.467907	2025-05-19 18:03:48.467909	200	0
6855	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:23.729594	2025-05-19 15:40:23.734688	200	5
6856	\N	/status	GET	2025-05-19 15:40:24.383719	2025-05-19 15:40:24.383721	200	0
6857	\N	/status	GET	2025-05-19 15:40:24.435565	2025-05-19 15:40:24.435567	200	0
6862	\N	/api/statistics/global	GET	2025-05-19 15:40:29.422794	2025-05-19 15:40:29.425103	200	2
6863	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:29.469547	2025-05-19 15:40:29.475872	200	6
6864	\N	/status	GET	2025-05-19 15:40:34.388664	2025-05-19 15:40:34.388666	200	0
6865	\N	/status	GET	2025-05-19 15:40:34.438472	2025-05-19 15:40:34.438475	200	0
6870	\N	/api/statistics/global	GET	2025-05-19 15:40:39.45104	2025-05-19 15:40:39.453538	200	2
6871	\N	/api/statistics/endpoints	GET	2025-05-19 15:40:39.511074	2025-05-19 15:40:39.517731	200	6
6873	\N	/status	GET	2025-05-19 15:40:47.035307	2025-05-19 15:40:47.03531	200	0
6874	\N	/status	GET	2025-05-19 15:40:47.080215	2025-05-19 15:40:47.080217	200	0
6896	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 15:41:09.415156	2025-05-19 15:41:09.415256	200	0
6897	\N	/status	GET	2025-05-19 15:41:11.968003	2025-05-19 15:41:11.968004	200	0
6898	\N	/status	GET	2025-05-19 15:41:12.020698	2025-05-19 15:41:12.0207	200	0
6899	\N	/api/statistics/global	GET	2025-05-19 15:41:12.06279	2025-05-19 15:41:12.065206	200	2
6900	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:12.110752	2025-05-19 15:41:12.116664	200	5
6901	\N	/status	GET	2025-05-19 15:41:16.967717	2025-05-19 15:41:16.967719	200	0
6902	\N	/status	GET	2025-05-19 15:41:17.012819	2025-05-19 15:41:17.012821	200	0
6903	\N	/api/statistics/global	GET	2025-05-19 15:41:17.056459	2025-05-19 15:41:17.058723	200	2
6904	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:17.109439	2025-05-19 15:41:17.116325	200	6
6905	\N	/status	GET	2025-05-19 15:41:21.968552	2025-05-19 15:41:21.968553	200	0
6906	\N	/status	GET	2025-05-19 15:41:22.013108	2025-05-19 15:41:22.01311	200	0
6907	\N	/api/statistics/global	GET	2025-05-19 15:41:22.05608	2025-05-19 15:41:22.057964	200	1
6908	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:22.104075	2025-05-19 15:41:22.110709	200	6
6909	\N	/api/restaurant	GET	2025-05-19 15:41:26.442217	2025-05-19 15:41:26.442294	200	0
6910	\N	/status	GET	2025-05-19 15:41:26.972232	2025-05-19 15:41:26.972233	200	0
6911	\N	/status	GET	2025-05-19 15:41:27.024909	2025-05-19 15:41:27.024911	200	0
6912	\N	/api/statistics/global	GET	2025-05-19 15:41:27.060049	2025-05-19 15:41:27.062312	200	2
6913	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:27.109891	2025-05-19 15:41:27.11578	200	5
6914	\N	/status	GET	2025-05-19 15:41:31.96681	2025-05-19 15:41:31.966812	200	0
6915	\N	/status	GET	2025-05-19 15:41:32.011747	2025-05-19 15:41:32.011748	200	0
6916	\N	/api/statistics/global	GET	2025-05-19 15:41:32.068224	2025-05-19 15:41:32.070326	200	2
6917	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:32.115021	2025-05-19 15:41:32.122536	200	7
6918	\N	/status	GET	2025-05-19 15:41:36.963145	2025-05-19 15:41:36.963147	200	0
6919	\N	/status	GET	2025-05-19 15:41:37.006586	2025-05-19 15:41:37.006588	200	0
6920	\N	/api/statistics/global	GET	2025-05-19 15:41:37.065001	2025-05-19 15:41:37.067221	200	2
6921	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:37.120194	2025-05-19 15:41:37.126107	200	5
6922	\N	/status	GET	2025-05-19 15:41:41.964564	2025-05-19 15:41:41.964566	200	0
6923	\N	/status	GET	2025-05-19 15:41:42.00824	2025-05-19 15:41:42.008242	200	0
6924	\N	/api/statistics/global	GET	2025-05-19 15:41:42.048173	2025-05-19 15:41:42.050732	200	2
6925	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:42.095461	2025-05-19 15:41:42.101086	200	5
6926	\N	/status	GET	2025-05-19 15:41:46.971762	2025-05-19 15:41:46.971764	200	0
6927	\N	/status	GET	2025-05-19 15:41:47.016604	2025-05-19 15:41:47.016605	200	0
6928	\N	/api/statistics/global	GET	2025-05-19 15:41:47.061045	2025-05-19 15:41:47.063067	200	2
6929	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:47.106432	2025-05-19 15:41:47.111653	200	5
6930	\N	/status	GET	2025-05-19 15:41:51.967637	2025-05-19 15:41:51.967639	200	0
6931	\N	/status	GET	2025-05-19 15:41:52.018696	2025-05-19 15:41:52.018698	200	0
6932	\N	/api/statistics/global	GET	2025-05-19 15:41:52.055215	2025-05-19 15:41:52.057431	200	2
6933	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:52.10297	2025-05-19 15:41:52.109204	200	6
6934	\N	/status	GET	2025-05-19 15:41:56.967515	2025-05-19 15:41:56.967516	200	0
6935	\N	/status	GET	2025-05-19 15:41:57.014948	2025-05-19 15:41:57.01495	200	0
6936	\N	/api/statistics/global	GET	2025-05-19 15:41:57.064867	2025-05-19 15:41:57.067439	200	2
6937	\N	/api/statistics/endpoints	GET	2025-05-19 15:41:57.117802	2025-05-19 15:41:57.123299	200	5
6938	\N	/api/washingmachines	GET	2025-05-19 15:42:44.272811	2025-05-19 15:42:44.342578	200	69
6939	\N	/api/washingmachines	GET	2025-05-19 15:52:03.92812	2025-05-19 15:52:03.995581	200	67
6940	\N	/api/auth/login	POST	2025-05-19 15:54:59.992223	2025-05-19 15:55:00.067031	401	74
6941	\N	/api/auth/login	POST	2025-05-19 15:55:11.996065	2025-05-19 15:55:12.069774	401	73
6942	\N	/api/auth/login	POST	2025-05-19 15:55:24.648473	2025-05-19 15:55:24.721689	401	73
6943	\N	/api/auth/login	POST	2025-05-19 15:55:37.251234	2025-05-19 15:55:37.327547	401	76
6944	\N	/api/auth/login	POST	2025-05-19 15:55:52.534815	2025-05-19 15:55:52.614294	401	79
6945	\N	/api/auth/login	POST	2025-05-19 15:55:54.529366	2025-05-19 15:55:54.611413	401	82
6946	\N	/api/auth/login	POST	2025-05-19 15:56:39.342929	2025-05-19 15:56:39.41902	401	76
6947	\N	/api/auth/login	POST	2025-05-19 15:56:53.300783	2025-05-19 15:56:53.372468	401	71
6948	\N	/api/auth/login	POST	2025-05-19 15:57:05.357	2025-05-19 15:57:05.430935	401	73
6949	\N	/api/washingmachines	GET	2025-05-19 15:57:06.227098	2025-05-19 15:57:06.301628	200	74
6950	\N	/api/auth/login	POST	2025-05-19 15:57:17.992599	2025-05-19 15:57:18.068135	401	75
6951	\N	/api/auth/login	POST	2025-05-19 15:58:09.500886	2025-05-19 15:58:09.574061	401	73
6952	\N	/api/washingmachines	GET	2025-05-19 16:02:03.703526	2025-05-19 16:02:03.791606	200	88
6953	\N	/api/washingmachines	GET	2025-05-19 16:07:03.621997	2025-05-19 16:07:03.699619	200	77
6954	\N	/api/washingmachines	GET	2025-05-19 16:12:03.670766	2025-05-19 16:12:03.750454	200	79
6956	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 16:45:11.79924	2025-05-19 16:45:11.799307	200	0
6957	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 16:55:11.776468	2025-05-19 16:55:11.778348	200	1
6958	\N	/status	GET	2025-05-19 16:55:13.034323	2025-05-19 16:55:13.034333	200	0
6959	\N	/api/statistics/global	GET	2025-05-19 16:55:13.214622	2025-05-19 16:55:13.216573	200	1
6960	\N	/api/statistics/top-users	GET	2025-05-19 16:55:13.292412	2025-05-19 16:55:13.294809	200	2
6961	\N	/status	GET	2025-05-19 16:55:23.012832	2025-05-19 16:55:23.012834	200	0
6962	\N	/status	GET	2025-05-19 16:55:33.039321	2025-05-19 16:55:33.039322	200	0
6963	\N	/status	GET	2025-05-19 16:55:43.080389	2025-05-19 16:55:43.080391	200	0
6964	\N	/status	GET	2025-05-19 16:55:53.019389	2025-05-19 16:55:53.01939	200	0
6965	\N	/status	GET	2025-05-19 16:56:03.036359	2025-05-19 16:56:03.036361	200	0
6966	\N	/status	GET	2025-05-19 16:56:13.025288	2025-05-19 16:56:13.025289	200	0
6967	\N	/status	GET	2025-05-19 16:56:23.041897	2025-05-19 16:56:23.041898	200	0
6968	\N	/status	GET	2025-05-19 16:56:28.059753	2025-05-19 16:56:28.059755	200	0
6969	\N	/api/statistics/global	GET	2025-05-19 16:56:28.129704	2025-05-19 16:56:28.13158	200	1
6970	\N	/api/statistics/top-users	GET	2025-05-19 16:56:28.207032	2025-05-19 16:56:28.208467	200	1
6971	\N	/status	GET	2025-05-19 16:57:16.339909	2025-05-19 16:57:16.33991	200	0
6972	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 16:57:17.57441	2025-05-19 16:57:17.576286	200	1
6973	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 16:57:21.61788	2025-05-19 16:57:21.62751	200	9
6974	\N	/api/restaurant	GET	2025-05-19 16:57:24.287103	2025-05-19 16:57:25.614539	200	1327
6984	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:00:29.777562	2025-05-19 17:00:29.779455	200	1
6985	\N	/api/weather	GET	2025-05-19 17:00:29.771887	2025-05-19 17:00:29.830181	200	58
6986	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:00:33.035001	2025-05-19 17:00:33.057524	200	22
6987	\N	/api/washingmachines	GET	2025-05-19 17:00:35.747748	2025-05-19 17:00:35.776649	200	28
6988	\N	/api/washingmachines	GET	2025-05-19 17:00:46.879569	2025-05-19 17:00:46.922768	200	43
6989	\N	/api/washingmachines	GET	2025-05-19 17:00:55.871823	2025-05-19 17:00:55.912048	200	40
6990	\N	/api/restaurant	GET	2025-05-19 17:00:58.727832	2025-05-19 17:00:58.727887	200	0
6991	\N	/api/traq/	GET	2025-05-19 17:01:13.950548	2025-05-19 17:01:13.951796	200	1
6992	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:01:49.510281	2025-05-19 17:01:49.512086	200	1
6993	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:01:53.579967	2025-05-19 17:01:53.581699	200	1
7623	\N	/status	GET	2025-05-20 08:49:46.759014	2025-05-20 08:49:46.759016	200	0
6994	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:01:58.659043	2025-05-19 17:01:58.664282	200	5
6995	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:02:32.912405	2025-05-19 17:02:32.92176	200	9
6996	\N	/api/support	GET	2025-05-19 17:05:05.92411	2025-05-19 17:05:05.924118	500	0
6997	\N	/api/support	GET	2025-05-19 17:05:07.03452	2025-05-19 17:05:07.034528	500	0
6998	\N	/api/support	GET	2025-05-19 17:05:09.160366	2025-05-19 17:05:09.160374	500	0
6999	\N	/api/support	GET	2025-05-19 17:05:13.318004	2025-05-19 17:05:13.318012	500	0
7000	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:05:14.527838	2025-05-19 17:05:14.530618	200	2
7001	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:05:18.298616	2025-05-19 17:05:18.303085	200	4
7002	aurelien.moignet@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 17:05:20.640327	2025-05-19 17:05:20.642328	200	2
7003	\N	/api/support	GET	2025-05-19 17:05:33.424028	2025-05-19 17:05:33.424038	500	0
7004	\N	/api/support	GET	2025-05-19 17:05:34.556817	2025-05-19 17:05:34.556825	500	0
7005	\N	/api/support	GET	2025-05-19 17:05:36.694443	2025-05-19 17:05:36.69445	500	0
7006	\N	/api/support	GET	2025-05-19 17:05:40.835097	2025-05-19 17:05:40.835103	500	0
7007	\N	/api/support	GET	2025-05-19 17:05:45.200002	2025-05-19 17:05:45.20001	500	0
7008	\N	/api/support	GET	2025-05-19 17:05:46.326988	2025-05-19 17:05:46.326995	500	0
7009	\N	/status	GET	2025-05-19 17:05:50.421794	2025-05-19 17:05:50.421796	200	0
7010	\N	/api/statistics/global	GET	2025-05-19 17:05:50.495943	2025-05-19 17:05:50.499803	200	3
7011	\N	/api/statistics/top-users	GET	2025-05-19 17:05:50.583684	2025-05-19 17:05:50.586001	200	2
7012	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:05:56.054249	2025-05-19 17:05:56.113225	200	58
7013	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:05:58.174113	2025-05-19 17:05:58.180053	200	5
7014	\N	/status	GET	2025-05-19 17:06:00.298269	2025-05-19 17:06:00.298271	200	0
7015	\N	/api/statistics/global	GET	2025-05-19 17:06:00.38678	2025-05-19 17:06:00.38915	200	2
7016	\N	/api/statistics/top-users	GET	2025-05-19 17:06:00.467075	2025-05-19 17:06:00.468589	200	1
7017	\N	/status	GET	2025-05-19 17:06:10.330125	2025-05-19 17:06:10.330126	200	0
7018	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:10:15.128026	2025-05-19 17:10:15.129982	200	1
7019	\N	/api/restaurant	GET	2025-05-19 17:10:15.342409	2025-05-19 17:10:15.342511	200	0
7037	nathan.marie2@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:24:32.804806	2025-05-19 17:24:32.809172	200	4
7038	\N	/api/washingmachines	GET	2025-05-19 17:24:41.818578	2025-05-19 17:24:41.847663	200	29
7039	\N	/api/traq/	GET	2025-05-19 17:25:06.072422	2025-05-19 17:25:06.073571	200	1
7040	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:25:12.638457	2025-05-19 17:25:12.640162	200	1
7041	\N	/status	GET	2025-05-19 17:33:07.799001	2025-05-19 17:33:07.799005	200	0
7042	\N	/api/statistics/global	GET	2025-05-19 17:33:07.905243	2025-05-19 17:33:07.90734	200	2
7043	\N	/api/statistics/top-users	GET	2025-05-19 17:33:08.00512	2025-05-19 17:33:08.006615	200	1
7044	\N	/status	GET	2025-05-19 17:33:15.088069	2025-05-19 17:33:15.088071	200	0
7045	\N	/api/statistics/global	GET	2025-05-19 17:33:15.184998	2025-05-19 17:33:15.187096	200	2
7046	\N	/api/statistics/top-users	GET	2025-05-19 17:33:15.273042	2025-05-19 17:33:15.274496	200	1
7047	\N	/status	GET	2025-05-19 17:33:17.226748	2025-05-19 17:33:17.22675	200	0
7048	\N	/status	GET	2025-05-19 17:33:27.223641	2025-05-19 17:33:27.223643	200	0
7049	\N	/status	GET	2025-05-19 17:33:37.226269	2025-05-19 17:33:37.22627	200	0
7050	\N	/status	GET	2025-05-19 17:33:47.227486	2025-05-19 17:33:47.227488	200	0
7051	\N	/status	GET	2025-05-19 17:34:09.838309	2025-05-19 17:34:09.838312	200	0
7052	\N	/api/statistics/global	GET	2025-05-19 17:34:09.940989	2025-05-19 17:34:09.94328	200	2
7053	\N	/api/statistics/top-users	GET	2025-05-19 17:34:10.052844	2025-05-19 17:34:10.054075	200	1
7054	\N	/status	GET	2025-05-19 17:43:54.745179	2025-05-19 17:43:54.745182	200	0
7055	\N	/status	GET	2025-05-19 17:43:55.813406	2025-05-19 17:43:55.813407	200	0
7056	\N	/api/statistics/global	GET	2025-05-19 17:43:55.922439	2025-05-19 17:43:55.929562	200	7
7057	\N	/api/statistics/top-users	GET	2025-05-19 17:43:56.014556	2025-05-19 17:43:56.01591	200	1
7058	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:44:32.753096	2025-05-19 17:44:32.754711	200	1
7060	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:44:33.048768	2025-05-19 17:44:33.052388	200	3
7061	\N	/api/washingmachines	GET	2025-05-19 17:44:32.997116	2025-05-19 17:44:33.097183	200	100
7062	\N	/api/weather	GET	2025-05-19 17:44:33.048769	2025-05-19 17:44:33.102135	200	53
7063	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 17:44:34.580534	2025-05-19 17:44:34.585246	200	4
7064	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 17:44:36.001097	2025-05-19 17:44:36.003006	200	1
7073	\N	/status	GET	2025-05-19 17:48:00.335222	2025-05-19 17:48:00.335232	200	0
7098	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 18:03:34.82045	2025-05-19 18:03:34.820528	200	0
7099	\N	/api/restaurant	GET	2025-05-19 18:03:34.915006	2025-05-19 18:03:34.915059	200	0
7100	\N	/api/restaurant	GET	2025-05-19 18:03:35.129739	2025-05-19 18:03:35.129805	200	0
7101	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 18:03:35.135384	2025-05-19 18:03:35.135435	200	0
7106	\N	/status	GET	2025-05-19 18:03:53.494987	2025-05-19 18:03:53.494989	200	0
7107	\N	/api/statistics/global	GET	2025-05-19 18:03:53.556053	2025-05-19 18:03:53.558187	200	2
7108	\N	/api/statistics/endpoints	GET	2025-05-19 18:03:53.599145	2025-05-19 18:03:53.604722	200	5
7109	\N	/api/restaurant	GET	2025-05-19 18:03:56.427073	2025-05-19 18:03:56.42715	200	0
7110	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:14:57.255789	2025-05-19 18:14:57.257917	200	2
7112	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:14:57.402541	2025-05-19 18:14:57.404703	200	2
7113	\N	/api/weather	GET	2025-05-19 18:14:57.402546	2025-05-19 18:14:57.454829	200	52
7123	\N	/api/restaurant	GET	2025-05-19 18:18:25.548729	2025-05-19 18:18:25.548821	200	0
7103	\N	/api/statistics/global	GET	2025-05-19 18:03:48.507206	2025-05-19 18:03:48.509515	200	2
7104	\N	/api/statistics/endpoints	GET	2025-05-19 18:03:48.724881	2025-05-19 18:03:48.732025	200	7
7105	\N	/status	GET	2025-05-19 18:03:53.495007	2025-05-19 18:03:53.495009	200	0
7111	\N	/api/restaurant	GET	2025-05-19 18:14:57.402539	2025-05-19 18:14:57.40271	200	0
7114	\N	/api/washingmachines	GET	2025-05-19 18:14:57.366355	2025-05-19 18:14:57.455149	200	88
7115	\N	/api/washingmachines	GET	2025-05-19 18:15:11.374752	2025-05-19 18:15:11.402422	200	27
7116	\N	/api/washingmachines	GET	2025-05-19 18:15:24.146195	2025-05-19 18:15:24.179591	200	33
7117	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:15:47.157936	2025-05-19 18:15:47.160072	200	2
7118	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:16:14.545009	2025-05-19 18:16:14.546969	200	1
7119	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:16:16.754537	2025-05-19 18:16:16.756495	200	1
7120	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:16:26.104936	2025-05-19 18:16:26.107547	200	2
7121	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:17:09.176759	2025-05-19 18:17:09.178761	200	2
7122	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:18:24.477479	2025-05-19 18:18:24.479068	200	1
7124	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:18:25.54869	2025-05-19 18:18:25.55079	200	2
7125	\N	/api/weather	GET	2025-05-19 18:18:25.548689	2025-05-19 18:18:25.596679	200	47
7126	\N	/api/washingmachines	GET	2025-05-19 18:18:25.543729	2025-05-19 18:18:25.613697	200	69
7127	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:18:30.493821	2025-05-19 18:18:30.495355	200	1
7128	\N	/api/washingmachines	GET	2025-05-19 18:19:02.908676	2025-05-19 18:19:02.939207	200	30
7129	\N	/api/washingmachines	GET	2025-05-19 18:19:04.479345	2025-05-19 18:19:04.508279	200	28
7130	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:19:18.136251	2025-05-19 18:19:18.138255	200	2
7131	aurelien.moignet@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 18:19:22.840442	2025-05-19 18:19:22.842037	200	1
7132	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:19:27.264694	2025-05-19 18:19:27.266582	200	1
7133	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:19:29.503856	2025-05-19 18:19:29.505896	200	2
7134	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 18:19:30.372669	2025-05-19 18:19:30.377066	200	4
7135	\N	/status	GET	2025-05-19 18:19:31.295093	2025-05-19 18:19:31.295096	200	0
7136	\N	/api/statistics/global	GET	2025-05-19 18:19:32.394032	2025-05-19 18:19:32.396692	200	2
7137	\N	/api/statistics/top-users	GET	2025-05-19 18:19:33.549743	2025-05-19 18:19:33.551744	200	2
7138	\N	/status	GET	2025-05-19 18:19:41.287684	2025-05-19 18:19:41.287686	200	0
7139	\N	/status	GET	2025-05-19 18:19:52.293672	2025-05-19 18:19:52.293673	200	0
7140	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:19:54.852124	2025-05-19 18:19:54.853865	200	1
7141	\N	/status	GET	2025-05-19 18:19:55.779539	2025-05-19 18:19:55.779543	200	0
7142	\N	/api/statistics/global	GET	2025-05-19 18:19:55.858145	2025-05-19 18:19:55.86034	200	2
7143	\N	/api/statistics/top-users	GET	2025-05-19 18:19:55.92685	2025-05-19 18:19:55.928501	200	1
7144	\N	/status	GET	2025-05-19 18:20:01.306372	2025-05-19 18:20:01.306374	200	0
7145	\N	/status	GET	2025-05-19 18:20:05.767501	2025-05-19 18:20:05.767503	200	0
7146	\N	/status	GET	2025-05-19 18:20:11.315352	2025-05-19 18:20:11.315354	200	0
7147	\N	/status	GET	2025-05-19 18:20:34.549529	2025-05-19 18:20:34.549533	200	0
7148	\N	/status	GET	2025-05-19 18:20:44.256921	2025-05-19 18:20:44.256922	200	0
7149	\N	/status	GET	2025-05-19 18:20:54.267241	2025-05-19 18:20:54.267243	200	0
7150	\N	/status	GET	2025-05-19 18:20:59.776245	2025-05-19 18:20:59.776248	200	0
7151	\N	/status	GET	2025-05-19 18:21:04.269559	2025-05-19 18:21:04.269561	200	0
7152	\N	/status	GET	2025-05-19 18:21:08.746924	2025-05-19 18:21:08.746926	200	0
7153	\N	/status	GET	2025-05-19 18:21:14.278518	2025-05-19 18:21:14.27852	200	0
7154	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:21:16.434788	2025-05-19 18:21:16.436662	200	1
7155	\N	/status	GET	2025-05-19 18:21:16.538523	2025-05-19 18:21:16.538524	200	0
7156	\N	/api/statistics/global	GET	2025-05-19 18:21:16.644762	2025-05-19 18:21:16.646929	200	2
7157	\N	/api/statistics/top-users	GET	2025-05-19 18:21:16.714708	2025-05-19 18:21:16.716351	200	1
7158	\N	/status	GET	2025-05-19 18:21:19.448989	2025-05-19 18:21:19.448991	200	0
7159	\N	/api/statistics/global	GET	2025-05-19 18:21:19.52488	2025-05-19 18:21:19.527311	200	2
7160	\N	/api/statistics/top-users	GET	2025-05-19 18:21:19.589428	2025-05-19 18:21:19.591298	200	1
7161	\N	/status	GET	2025-05-19 18:21:24.294829	2025-05-19 18:21:24.29483	200	0
7162	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:21:25.68868	2025-05-19 18:21:25.690587	200	1
7163	\N	/status	GET	2025-05-19 18:21:27.681408	2025-05-19 18:21:27.681409	200	0
7164	\N	/api/statistics/global	GET	2025-05-19 18:21:27.78608	2025-05-19 18:21:27.788176	200	2
7165	\N	/api/statistics/top-users	GET	2025-05-19 18:21:27.890938	2025-05-19 18:21:27.892525	200	1
7166	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 18:21:30.067427	2025-05-19 18:21:30.071986	200	4
7167	\N	/status	GET	2025-05-19 18:21:30.107791	2025-05-19 18:21:30.107792	200	0
7168	\N	/api/statistics/global	GET	2025-05-19 18:21:30.192174	2025-05-19 18:21:30.194233	200	2
7169	\N	/api/statistics/top-users	GET	2025-05-19 18:21:30.248859	2025-05-19 18:21:30.250721	200	1
7170	\N	/status	GET	2025-05-19 18:21:31.498814	2025-05-19 18:21:31.498816	200	0
7171	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:21:31.538497	2025-05-19 18:21:31.540011	200	1
7172	\N	/api/statistics/global	GET	2025-05-19 18:21:31.568554	2025-05-19 18:21:31.570586	200	2
7173	\N	/api/statistics/top-users	GET	2025-05-19 18:21:31.648685	2025-05-19 18:21:31.65076	200	2
7174	\N	/status	GET	2025-05-19 18:21:34.288872	2025-05-19 18:21:34.288874	200	0
7175	\N	/status	GET	2025-05-19 18:21:37.686238	2025-05-19 18:21:37.686239	200	0
7176	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-19 18:21:39.255573	2025-05-19 18:21:39.260391	200	4
7177	\N	/status	GET	2025-05-19 18:21:44.299572	2025-05-19 18:21:44.299573	200	0
7178	\N	/api/washingmachines	GET	2025-05-19 18:21:46.698806	2025-05-19 18:21:46.800359	200	101
7179	\N	/status	GET	2025-05-19 18:21:47.687379	2025-05-19 18:21:47.68738	200	0
7180	\N	/api/washingmachines	GET	2025-05-19 18:21:57.499298	2025-05-19 18:21:57.529977	200	30
7181	\N	/api/washingmachines	GET	2025-05-19 18:22:46.147109	2025-05-19 18:22:46.177202	200	30
7182	\N	/api/washingmachines	GET	2025-05-19 18:22:51.084725	2025-05-19 18:22:51.114244	200	29
7183	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:25:24.968534	2025-05-19 18:25:24.970797	200	2
7184	\N	/api/restaurant	GET	2025-05-19 18:25:25.218805	2025-05-19 18:25:25.218895	200	0
7185	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:25:25.218429	2025-05-19 18:25:25.221398	200	2
7186	\N	/api/washingmachines	GET	2025-05-19 18:25:25.156355	2025-05-19 18:25:25.253092	200	96
7187	\N	/api/weather	GET	2025-05-19 18:25:25.218406	2025-05-19 18:25:25.282891	200	64
7188	\N	/api/washingmachines	GET	2025-05-19 18:25:27.020497	2025-05-19 18:25:27.054053	200	33
7189	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:25:39.767997	2025-05-19 18:25:39.769966	200	1
7190	\N	/api/washingmachines	GET	2025-05-19 18:25:39.963873	2025-05-19 18:25:39.99509	200	31
7191	\N	/api/weather	GET	2025-05-19 18:25:40.020875	2025-05-19 18:25:40.020902	200	0
7192	\N	/api/restaurant	GET	2025-05-19 18:25:40.027212	2025-05-19 18:25:40.027277	200	0
7193	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:25:40.030028	2025-05-19 18:25:40.031859	200	1
7194	\N	/api/washingmachines	GET	2025-05-19 18:25:42.170385	2025-05-19 18:25:42.204803	200	34
7195	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:26:23.559516	2025-05-19 18:26:23.561082	200	1
7196	\N	/api/washingmachines	GET	2025-05-19 18:26:23.732723	2025-05-19 18:26:23.762804	200	30
7197	\N	/api/weather	GET	2025-05-19 18:26:23.786163	2025-05-19 18:26:23.786183	200	0
7198	\N	/api/restaurant	GET	2025-05-19 18:26:23.79143	2025-05-19 18:26:23.791518	200	0
7199	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:26:23.792324	2025-05-19 18:26:23.950176	200	157
7200	\N	/api/washingmachines	GET	2025-05-19 18:26:25.286566	2025-05-19 18:26:25.334944	200	48
7201	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:26:44.420792	2025-05-19 18:26:44.423353	200	2
7202	\N	/api/washingmachines	GET	2025-05-19 18:26:44.609246	2025-05-19 18:26:44.640929	200	31
7203	\N	/api/weather	GET	2025-05-19 18:26:44.651533	2025-05-19 18:26:44.651556	200	0
7211	\N	/api/auth/login	POST	2025-05-19 18:35:41.997675	2025-05-19 18:35:42.071954	200	74
7214	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.204496	2025-05-19 18:35:42.214242	200	9
7215	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.204607	2025-05-19 18:35:42.220478	200	15
7217	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.227739	2025-05-19 18:35:42.229456	200	1
7221	\N	/api/auth/login	POST	2025-05-19 18:35:41.998109	2025-05-19 18:35:42.232373	200	234
7204	\N	/api/restaurant	GET	2025-05-19 18:26:44.651533	2025-05-19 18:26:44.651626	200	0
7205	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:26:44.654152	2025-05-19 18:26:44.803723	200	149
7206	\N	/api/washingmachines	GET	2025-05-19 18:26:46.057769	2025-05-19 18:26:46.088104	200	30
7207	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:28:11.466332	2025-05-19 18:28:11.468571	200	2
7208	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:30:24.693319	2025-05-19 18:30:24.695361	200	2
7209	\N	/api/washingmachines	GET	2025-05-19 18:34:50.186913	2025-05-19 18:34:50.275459	200	88
7210	\N	/api/auth/login	POST	2025-05-19 18:35:41.997904	2025-05-19 18:35:42.071196	200	73
7212	\N	/api/auth/login	POST	2025-05-19 18:35:42.004363	2025-05-19 18:35:42.078436	200	74
7213	\N	/api/auth/login	POST	2025-05-19 18:35:42.045281	2025-05-19 18:35:42.120896	200	75
7216	\N	/api/auth/login	POST	2025-05-19 18:35:41.99829	2025-05-19 18:35:42.221399	200	223
7218	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.227747	2025-05-19 18:35:42.229408	200	1
7219	\N	/api/auth/login	POST	2025-05-19 18:35:41.998137	2025-05-19 18:35:42.221475	200	223
7220	\N	/api/auth/login	POST	2025-05-19 18:35:41.997944	2025-05-19 18:35:42.221961	200	224
7222	\N	/api/auth/login	POST	2025-05-19 18:35:41.998075	2025-05-19 18:35:42.231799	200	233
7223	\N	/api/auth/login	POST	2025-05-19 18:35:42.204354	2025-05-19 18:35:42.282043	200	77
7224	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.33426	2025-05-19 18:35:42.337322	200	3
7225	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.337922	2025-05-19 18:35:42.339722	200	1
7226	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.338702	2025-05-19 18:35:42.348614	200	9
7227	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.37868	2025-05-19 18:35:42.380542	200	1
7228	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.378666	2025-05-19 18:35:42.380598	200	1
7229	\N	/api/auth/login	POST	2025-05-19 18:35:42.380262	2025-05-19 18:35:42.380272	429	0
7230	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.386087	2025-05-19 18:35:42.389056	200	2
7231	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.38169	2025-05-19 18:35:42.391989	200	10
7232	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.381669	2025-05-19 18:35:42.393543	200	11
7233	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.381379	2025-05-19 18:35:42.394069	200	12
7234	\N	/api/washingmachines	GET	2025-05-19 18:35:42.440047	2025-05-19 18:35:42.470086	200	30
7235	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.758395	2025-05-19 18:35:42.761294	200	2
7236	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.761525	2025-05-19 18:35:42.7638	200	2
7237	\N	/api/restaurant	GET	2025-05-19 18:35:42.761553	2025-05-19 18:35:42.761607	200	0
7238	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.790629	2025-05-19 18:35:42.792698	200	2
7239	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.794177	2025-05-19 18:35:42.79591	200	1
7240	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.794209	2025-05-19 18:35:42.803971	200	9
7241	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.803644	2025-05-19 18:35:42.804855	200	1
7242	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:42.803511	2025-05-19 18:35:42.80463	200	1
7243	\N	/api/weather	GET	2025-05-19 18:35:42.758391	2025-05-19 18:35:42.810502	200	52
7244	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:43.018428	2025-05-19 18:35:43.020336	200	1
7245	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:53.227121	2025-05-19 18:35:53.232099	200	4
7246	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:35:57.398229	2025-05-19 18:35:57.400017	200	1
7247	\N	/status	GET	2025-05-19 18:37:02.936834	2025-05-19 18:37:02.936837	200	0
7248	\N	/api/statistics/global	GET	2025-05-19 18:37:03.054427	2025-05-19 18:37:03.056476	200	2
7249	\N	/api/statistics/top-users	GET	2025-05-19 18:37:03.166636	2025-05-19 18:37:03.168194	200	1
7250	\N	/status	GET	2025-05-19 18:37:12.920293	2025-05-19 18:37:12.920294	200	0
7251	\N	/api/washingmachines	GET	2025-05-19 18:38:21.298577	2025-05-19 18:38:21.381743	200	83
7252	\N	/api/washingmachines	GET	2025-05-19 18:43:20.935327	2025-05-19 18:43:21.007621	200	72
7253	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:47:20.216417	2025-05-19 18:47:20.218416	200	1
7254	\N	/api/restaurant	GET	2025-05-19 18:47:20.449361	2025-05-19 18:47:20.449555	200	0
7255	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-19 18:47:20.44936	2025-05-19 18:47:20.451689	200	2
7256	\N	/api/weather	GET	2025-05-19 18:47:20.449375	2025-05-19 18:47:20.500375	200	51
7257	\N	/api/washingmachines	GET	2025-05-19 18:47:20.393453	2025-05-19 18:47:20.50653	200	113
7258	\N	/api/traq/	GET	2025-05-19 18:47:27.212184	2025-05-19 18:47:27.213332	200	1
7259	\N	/api/washingmachines	GET	2025-05-19 18:48:21.054989	2025-05-19 18:48:21.106293	200	51
7260	\N	/api/washingmachines	GET	2025-05-19 18:53:21.16616	2025-05-19 18:53:21.238839	200	72
7261	\N	/api/washingmachines	GET	2025-05-19 18:58:21.0515	2025-05-19 18:58:21.139621	200	88
7262	\N	/api/washingmachines	GET	2025-05-19 19:00:35.270661	2025-05-19 19:00:35.361895	200	91
7263	\N	/api/washingmachines	GET	2025-05-19 19:03:20.93373	2025-05-19 19:03:21.033418	200	99
7264	\N	/api/washingmachines	GET	2025-05-19 19:08:20.932698	2025-05-19 19:08:21.02658	200	93
7265	\N	/api/washingmachines	GET	2025-05-19 19:13:20.935103	2025-05-19 19:13:21.009863	200	74
7266	\N	/api/washingmachines	GET	2025-05-19 19:18:20.930221	2025-05-19 19:18:21.047655	200	117
7267	\N	/api/washingmachines	GET	2025-05-19 19:23:20.929294	2025-05-19 19:23:21.029528	200	100
7268	\N	/api/washingmachines	GET	2025-05-19 19:28:20.935312	2025-05-19 19:28:21.003591	200	68
7269	\N	/api/washingmachines	GET	2025-05-19 19:33:20.934116	2025-05-19 19:33:21.023043	200	88
7270	\N	/api/washingmachines	GET	2025-05-19 19:38:20.933218	2025-05-19 19:38:21.03791	200	104
7271	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:39:15.265614	2025-05-19 19:39:15.267327	200	1
7272	\N	/api/washingmachines	GET	2025-05-19 19:39:15.726947	2025-05-19 19:39:15.761877	200	34
7273	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:39:15.925522	2025-05-19 19:39:15.927355	200	1
7274	\N	/api/restaurant	GET	2025-05-19 19:39:15.925475	2025-05-19 19:39:15.92559	200	0
7275	\N	/api/weather	GET	2025-05-19 19:39:15.925451	2025-05-19 19:39:15.987033	200	61
7276	\N	/api/restaurant	GET	2025-05-19 19:39:19.408046	2025-05-19 19:39:19.408122	200	0
7277	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:39:27.163142	2025-05-19 19:39:27.165123	200	1
7278	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-19 19:39:27.589557	2025-05-19 19:39:27.589673	200	0
7279	\N	/api/washingmachines	GET	2025-05-19 19:43:20.930927	2025-05-19 19:43:21.012589	200	81
7280	\N	/api/washingmachines	GET	2025-05-19 19:48:21.216733	2025-05-19 19:48:21.287574	200	70
7281	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:50:57.98059	2025-05-19 19:50:57.982618	200	2
7282	\N	/api/restaurant	GET	2025-05-19 19:50:58.198991	2025-05-19 19:50:58.199076	200	0
7283	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:50:58.199014	2025-05-19 19:50:58.201655	200	2
7284	\N	/api/washingmachines	GET	2025-05-19 19:50:58.111391	2025-05-19 19:50:58.206504	200	95
7285	\N	/api/weather	GET	2025-05-19 19:50:58.173342	2025-05-19 19:50:58.248557	200	75
7286	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:51:02.137472	2025-05-19 19:51:02.139277	200	1
7287	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-19 19:51:22.764332	2025-05-19 19:51:22.76608	200	1
7288	\N	/status	GET	2025-05-19 19:51:26.919758	2025-05-19 19:51:26.919768	200	0
7289	\N	/api/statistics/global	GET	2025-05-19 19:51:27.162399	2025-05-19 19:51:27.165666	200	3
7290	\N	/api/statistics/top-users	GET	2025-05-19 19:51:27.252894	2025-05-19 19:51:27.255446	200	2
7291	\N	/status	GET	2025-05-19 19:51:36.938318	2025-05-19 19:51:36.93832	200	0
7292	\N	/api/washingmachines	GET	2025-05-19 19:53:20.936302	2025-05-19 19:53:21.037057	200	100
7293	\N	/api/washingmachines	GET	2025-05-19 19:58:20.934521	2025-05-19 19:58:21.013733	200	79
7294	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 20:01:42.000582	2025-05-19 20:01:42.002555	200	1
7295	\N	/api/restaurant	GET	2025-05-19 20:01:42.22137	2025-05-19 20:01:42.221456	200	0
7319	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-19 21:35:29.813407	2025-05-19 21:35:29.813486	200	0
7320	\N	/api/washingmachines	GET	2025-05-19 21:43:25.600199	2025-05-19 21:43:25.677672	200	77
7321	\N	/api/washingmachines	GET	2025-05-19 21:48:15.357665	2025-05-19 21:48:15.448665	200	90
7322	\N	/api/washingmachines	GET	2025-05-19 21:53:15.372148	2025-05-19 21:53:15.452725	200	80
7323	\N	/api/washingmachines	GET	2025-05-19 22:41:48.060835	2025-05-19 22:41:48.15509	200	94
7324	\N	/robots.txt	GET	2025-05-20 02:05:18.609658	2025-05-20 02:05:18.609662	500	0
7325	\N	/robots.txt	GET	2025-05-20 02:17:26.473445	2025-05-20 02:17:26.473449	500	0
7326	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:02:07.194298	2025-05-20 05:02:07.196052	200	1
7327	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:02:07.377133	2025-05-20 05:02:07.379051	200	1
7328	\N	/api/washingmachines	GET	2025-05-20 05:02:07.312518	2025-05-20 05:02:07.420367	200	107
7329	\N	/api/weather	GET	2025-05-20 05:02:07.376965	2025-05-20 05:02:07.448276	200	71
7330	\N	/api/restaurant	GET	2025-05-20 05:02:07.376969	2025-05-20 05:02:07.503637	200	126
7338	\N	/api/restaurant	GET	2025-05-20 05:24:21.939747	2025-05-20 05:24:21.939845	200	0
7348	\N	/api/restaurant	GET	2025-05-20 06:21:44.911855	2025-05-20 06:21:44.911964	200	0
7351	\N	/api/washingmachines	GET	2025-05-20 06:21:44.8767	2025-05-20 06:21:44.970941	200	94
7352	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:21:47.155976	2025-05-20 06:21:47.158227	200	2
7353	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-20 06:21:47.223013	2025-05-20 06:21:47.223133	200	0
7354	\N	/api/restaurant	GET	2025-05-20 06:21:50.838255	2025-05-20 06:21:50.838341	200	0
7355	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:21:54.118786	2025-05-20 06:21:54.120494	200	1
7356	\N	/api/washingmachines	GET	2025-05-20 06:21:54.291244	2025-05-20 06:21:54.321699	200	30
7357	\N	/api/weather	GET	2025-05-20 06:21:54.345508	2025-05-20 06:21:54.345533	200	0
7359	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:21:54.350621	2025-05-20 06:21:54.352725	200	2
7360	\N	/api/restaurant	GET	2025-05-20 06:21:56.308234	2025-05-20 06:21:56.308292	200	0
7361	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:28:08.259408	2025-05-20 06:28:08.261239	200	1
7362	\N	/api/washingmachines	GET	2025-05-20 06:28:08.462172	2025-05-20 06:28:08.56404	200	101
7364	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:28:08.569779	2025-05-20 06:28:08.571067	200	1
7365	\N	/api/weather	GET	2025-05-20 06:28:08.566873	2025-05-20 06:28:08.620155	200	53
7366	\N	/api/restaurant	GET	2025-05-20 06:28:12.193191	2025-05-20 06:28:12.193245	200	0
7367	\N	/api/washingmachines	GET	2025-05-20 06:30:04.894171	2025-05-20 06:30:04.98407	200	89
7368	\N	/api/restaurant	GET	2025-05-20 06:30:05.025587	2025-05-20 06:30:05.025644	200	0
7369	\N	/api/weather	GET	2025-05-20 06:30:05.025517	2025-05-20 06:30:05.069821	200	44
7370	\N	/api/weather	GET	2025-05-20 06:30:22.744135	2025-05-20 06:30:22.744169	200	0
7371	\N	/api/washingmachines	GET	2025-05-20 06:30:22.724271	2025-05-20 06:30:22.753112	200	28
7372	\N	/api/restaurant	GET	2025-05-20 06:30:22.769955	2025-05-20 06:30:22.770022	200	0
7373	\N	/api/restaurant	GET	2025-05-20 06:30:31.140221	2025-05-20 06:30:31.140311	200	0
7375	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:31:27.500311	2025-05-20 06:31:27.501797	200	1
7376	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:31:30.709367	2025-05-20 06:31:30.716186	200	6
7377	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:31:31.109237	2025-05-20 06:31:31.111556	200	2
7378	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:31:31.561699	2025-05-20 06:31:31.564005	200	2
7379	\N	/api/restaurant	GET	2025-05-20 06:34:10.007368	2025-05-20 06:34:10.007452	200	0
7380	\N	/api/washingmachines	GET	2025-05-20 06:34:09.977375	2025-05-20 06:34:10.029994	200	52
7381	\N	/api/weather	GET	2025-05-20 06:34:10.007084	2025-05-20 06:34:10.083593	200	76
7382	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:34:17.1879	2025-05-20 06:34:17.189636	200	1
7383	\N	/api/washingmachines	GET	2025-05-20 06:34:17.379073	2025-05-20 06:34:17.415506	200	36
7384	\N	/api/weather	GET	2025-05-20 06:34:17.467273	2025-05-20 06:34:17.467299	200	0
7385	\N	/api/restaurant	GET	2025-05-20 06:34:17.473208	2025-05-20 06:34:17.473302	200	0
7387	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 06:39:44.539688	2025-05-20 06:39:44.539781	200	0
7388	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 06:39:44.872347	2025-05-20 06:39:44.872408	200	0
7389	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:06:01.380076	2025-05-20 07:06:01.382603	200	2
7390	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:06:01.804088	2025-05-20 07:06:01.806393	200	2
7391	\N	/status	GET	2025-05-20 07:06:03.586123	2025-05-20 07:06:03.586138	200	0
7392	\N	/api/statistics/global	GET	2025-05-20 07:06:03.63882	2025-05-20 07:06:03.640907	200	2
7393	\N	/api/statistics/top-users	GET	2025-05-20 07:06:03.80174	2025-05-20 07:06:03.803174	200	1
7394	\N	/status	GET	2025-05-20 07:06:09.3549	2025-05-20 07:06:09.354901	200	0
7395	\N	/api/statistics/global	GET	2025-05-20 07:06:09.435662	2025-05-20 07:06:09.43815	200	2
7396	\N	/api/statistics/top-users	GET	2025-05-20 07:06:09.530912	2025-05-20 07:06:09.532589	200	1
7397	\N	/status	GET	2025-05-20 07:06:13.636009	2025-05-20 07:06:13.636011	200	0
7398	\N	/api/washingmachines	GET	2025-05-20 07:06:17.059102	2025-05-20 07:06:17.140609	200	81
7399	\N	/api/restaurant	GET	2025-05-20 07:06:19.480511	2025-05-20 07:06:19.480618	200	0
7400	\N	/status	GET	2025-05-20 07:06:23.628308	2025-05-20 07:06:23.628309	200	0
7401	\N	/api/restaurant	GET	2025-05-20 07:06:23.871584	2025-05-20 07:06:23.871677	200	0
7402	\N	/api/auth/register	POST	2025-05-20 07:07:04.338246	2025-05-20 07:07:04.429676	201	91
7403	\N	/api/auth/verify-account	POST	2025-05-20 07:07:43.135846	2025-05-20 07:07:43.144823	200	8
7404	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:07:43.400023	2025-05-20 07:07:43.402414	200	2
7405	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:07:43.654724	2025-05-20 07:07:43.656362	200	1
7406	\N	/api/washingmachines	GET	2025-05-20 07:07:44.050814	2025-05-20 07:07:44.087233	200	36
7408	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:07:45.213342	2025-05-20 07:07:45.215152	200	1
7409	\N	/api/weather	GET	2025-05-20 07:07:45.164131	2025-05-20 07:07:45.222213	200	58
7410	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:07:47.626555	2025-05-20 07:07:47.630569	200	4
7411	\N	/api/traq/	GET	2025-05-20 07:08:08.394451	2025-05-20 07:08:08.396896	200	2
7412	\N	/api/restaurant	GET	2025-05-20 07:08:31.568765	2025-05-20 07:08:31.56882	200	0
7413	\N	/api/washingmachines	GET	2025-05-20 07:08:41.073677	2025-05-20 07:08:41.106203	200	32
7414	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:08:56.240799	2025-05-20 07:08:56.243021	200	2
7415	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:09:01.208978	2025-05-20 07:09:01.210959	200	1
7416	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:09:03.645042	2025-05-20 07:09:03.646897	200	1
7417	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:10:36.659495	2025-05-20 07:10:36.661886	200	2
7418	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:10:36.830895	2025-05-20 07:10:36.837495	200	6
7419	\N	/api/washingmachines	GET	2025-05-20 07:10:36.784814	2025-05-20 07:10:36.859908	200	75
7420	\N	/api/weather	GET	2025-05-20 07:10:36.830889	2025-05-20 07:10:36.877232	200	46
7421	\N	/api/restaurant	GET	2025-05-20 07:10:36.830936	2025-05-20 07:10:37.487787	200	656
7432	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-20 07:14:21.958225	2025-05-20 07:14:21.958343	200	0
7442	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:15:23.960944	2025-05-20 07:15:23.962777	200	1
7443	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:15:35.341853	2025-05-20 07:15:35.346683	200	4
7444	\N	/api/traq/	GET	2025-05-20 07:15:59.241417	2025-05-20 07:15:59.24281	200	1
7445	\N	/api/washingmachines	GET	2025-05-20 07:16:01.567163	2025-05-20 07:16:01.59569	200	28
7446	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:31:43.856101	2025-05-20 07:31:43.857812	200	1
7447	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:31:47.584795	2025-05-20 07:31:47.590003	200	5
7448	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:31:54.046165	2025-05-20 07:31:54.054613	200	8
7296	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 20:01:42.222168	2025-05-19 20:01:42.224327	200	2
7297	\N	/api/weather	GET	2025-05-19 20:01:42.214511	2025-05-19 20:01:42.272103	200	57
7298	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 20:01:42.772219	2025-05-19 20:01:42.774415	200	2
7299	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-19 20:01:42.848576	2025-05-19 20:01:42.848672	200	0
7300	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 20:01:44.043505	2025-05-19 20:01:44.045605	200	2
7301	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-19 20:01:44.113329	2025-05-19 20:01:44.113406	200	0
7302	enzo.morvan@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-19 20:01:45.245264	2025-05-19 20:01:45.247098	200	1
7303	\N	/api/traq/	GET	2025-05-19 20:01:52.003067	2025-05-19 20:01:52.004221	200	1
7304	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-19 20:02:19.94911	2025-05-19 20:02:19.951374	200	2
7305	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-19 20:02:20.0281	2025-05-19 20:02:20.028235	200	0
7306	\N	/api/washingmachines	GET	2025-05-19 20:03:20.932497	2025-05-19 20:03:21.019421	200	86
7307	\N	/api/washingmachines	GET	2025-05-19 20:08:20.930702	2025-05-19 20:08:21.024858	200	94
7308	\N	/api/washingmachines	GET	2025-05-19 20:13:20.936599	2025-05-19 20:13:21.038144	200	101
7309	\N	/api/washingmachines	GET	2025-05-19 20:18:20.934265	2025-05-19 20:18:21.009227	200	74
7310	\N	/api/washingmachines	GET	2025-05-19 20:23:20.929699	2025-05-19 20:23:21.004365	200	74
7311	\N	/api/washingmachines	GET	2025-05-19 20:28:20.933591	2025-05-19 20:28:21.019725	200	86
7312	\N	/api/washingmachines	GET	2025-05-19 20:33:20.933791	2025-05-19 20:33:21.027935	200	94
7313	\N	/api/washingmachines	GET	2025-05-19 20:38:20.933193	2025-05-19 20:38:21.026879	200	93
7314	\N	/api/washingmachines	GET	2025-05-19 20:43:20.944797	2025-05-19 20:43:21.025027	200	80
7315	\N	/api/washingmachines	GET	2025-05-19 20:48:20.934626	2025-05-19 20:48:21.040168	200	105
7316	\N	/api/washingmachines	GET	2025-05-19 20:53:21.0811	2025-05-19 20:53:21.15005	200	68
7317	\N	/api/washingmachines	GET	2025-05-19 20:58:20.932393	2025-05-19 20:58:21.004491	200	72
7318	\N	/api/washingmachines	GET	2025-05-19 21:10:33.744619	2025-05-19 21:10:33.833103	200	88
7331	\N	/api/restaurant	GET	2025-05-20 05:02:12.241258	2025-05-20 05:02:12.241318	200	0
7332	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:02:17.789762	2025-05-20 05:02:17.79131	200	1
7333	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-20 05:02:17.865267	2025-05-20 05:02:17.865366	200	0
7334	\N	/api/washingmachines	GET	2025-05-20 05:02:22.352874	2025-05-20 05:02:22.388174	200	35
7335	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:02:28.092279	2025-05-20 05:02:28.094007	200	1
7336	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-20 05:02:28.171765	2025-05-20 05:02:28.171852	200	0
7337	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:24:21.699531	2025-05-20 05:24:21.70169	200	2
7339	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:24:21.939534	2025-05-20 05:24:21.94164	200	2
7340	\N	/api/weather	GET	2025-05-20 05:24:21.939573	2025-05-20 05:24:22.008126	200	68
7341	\N	/api/washingmachines	GET	2025-05-20 05:24:21.933764	2025-05-20 05:24:22.029255	200	95
7342	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:24:22.927447	2025-05-20 05:24:22.930126	200	2
7343	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 05:24:25.908487	2025-05-20 05:24:25.910333	200	1
7344	\N	/status	GET	2025-05-20 05:24:28.383253	2025-05-20 05:24:28.383268	200	0
7345	\N	/api/statistics/global	GET	2025-05-20 05:24:28.610325	2025-05-20 05:24:28.612475	200	2
7346	\N	/api/statistics/top-users	GET	2025-05-20 05:24:28.819013	2025-05-20 05:24:28.820846	200	1
7347	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:21:44.750923	2025-05-20 06:21:44.752674	200	1
7349	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:21:44.911696	2025-05-20 06:21:44.913839	200	2
7350	\N	/api/weather	GET	2025-05-20 06:21:44.911829	2025-05-20 06:21:44.967772	200	55
7358	\N	/api/restaurant	GET	2025-05-20 06:21:54.350675	2025-05-20 06:21:54.350782	200	0
7363	\N	/api/restaurant	GET	2025-05-20 06:28:08.567759	2025-05-20 06:28:08.567838	200	0
7374	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 06:30:31.140347	2025-05-20 06:30:31.140471	200	0
7386	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 06:34:17.47333	2025-05-20 06:34:17.474896	200	1
7407	\N	/api/restaurant	GET	2025-05-20 07:07:45.213358	2025-05-20 07:07:45.21345	200	0
7422	\N	/api/washingmachines	GET	2025-05-20 07:10:47.21369	2025-05-20 07:10:47.244394	200	30
7423	\N	/api/traq/	GET	2025-05-20 07:10:50.138524	2025-05-20 07:10:50.139709	200	1
7424	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:11:11.237871	2025-05-20 07:11:11.242583	200	4
7425	\N	/api/auth/login	POST	2025-05-20 07:14:18.626224	2025-05-20 07:14:18.701586	200	75
7426	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:14:18.772878	2025-05-20 07:14:18.774975	200	2
7427	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:14:18.83134	2025-05-20 07:14:18.832839	200	1
7428	\N	/api/washingmachines	GET	2025-05-20 07:14:19.007756	2025-05-20 07:14:19.075374	200	67
7429	\N	/api/weather	GET	2025-05-20 07:14:19.093952	2025-05-20 07:14:19.093978	200	0
7430	\N	/api/restaurant	GET	2025-05-20 07:14:19.103486	2025-05-20 07:14:19.103523	200	0
7431	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:14:19.116659	2025-05-20 07:14:19.119048	200	2
7433	\N	/status	GET	2025-05-20 07:15:02.582371	2025-05-20 07:15:02.582374	200	0
7434	\N	/api/statistics/global	GET	2025-05-20 07:15:02.683085	2025-05-20 07:15:02.685874	200	2
7435	\N	/api/statistics/top-users	GET	2025-05-20 07:15:02.790822	2025-05-20 07:15:02.792502	200	1
7436	\N	/status	GET	2025-05-20 07:15:12.60557	2025-05-20 07:15:12.605571	200	0
7437	\N	/status	GET	2025-05-20 07:15:22.610163	2025-05-20 07:15:22.610165	200	0
7438	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:15:23.529387	2025-05-20 07:15:23.531544	200	2
7439	\N	/api/washingmachines	GET	2025-05-20 07:15:23.873478	2025-05-20 07:15:23.902089	200	28
7440	\N	/api/weather	GET	2025-05-20 07:15:23.95346	2025-05-20 07:15:23.953487	200	0
7441	\N	/api/restaurant	GET	2025-05-20 07:15:23.96088	2025-05-20 07:15:23.960974	200	0
7479	\N	/api/restaurant	GET	2025-05-20 07:42:34.758513	2025-05-20 07:42:35.307235	200	548
7484	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:43:22.091968	2025-05-20 07:43:22.093274	200	1
7485	\N	/status	GET	2025-05-20 07:44:23.1698	2025-05-20 07:44:23.169803	200	0
7486	\N	/status	GET	2025-05-20 07:44:24.848731	2025-05-20 07:44:24.848733	200	0
7487	\N	/api/statistics/global	GET	2025-05-20 07:44:25.063988	2025-05-20 07:44:25.066181	200	2
7488	\N	/api/statistics/top-users	GET	2025-05-20 07:44:25.164121	2025-05-20 07:44:25.16556	200	1
7489	\N	/status	GET	2025-05-20 07:44:32.622791	2025-05-20 07:44:32.622792	200	0
7490	\N	/status	GET	2025-05-20 07:44:42.597844	2025-05-20 07:44:42.597846	200	0
7491	\N	/status	GET	2025-05-20 07:44:51.457495	2025-05-20 07:44:51.457497	200	0
7492	\N	/api/statistics/global	GET	2025-05-20 07:44:51.569576	2025-05-20 07:44:51.571561	200	1
7493	\N	/api/statistics/top-users	GET	2025-05-20 07:44:51.728781	2025-05-20 07:44:51.730634	200	1
7494	\N	/status	GET	2025-05-20 07:44:52.599482	2025-05-20 07:44:52.599484	200	0
7495	\N	/api/auth/login	POST	2025-05-20 07:45:07.809637	2025-05-20 07:45:07.811557	401	1
7496	\N	/api/auth/verification-code	POST	2025-05-20 07:45:21.225676	2025-05-20 07:45:21.227396	400	1
7497	\N	/api/auth/register	POST	2025-05-20 07:45:50.921481	2025-05-20 07:45:51.002975	201	81
7498	\N	/api/auth/verify-account	POST	2025-05-20 07:46:23.895095	2025-05-20 07:46:23.902651	200	7
7499	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:46:24.122269	2025-05-20 07:46:24.124532	200	2
7500	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:46:24.283405	2025-05-20 07:46:24.285149	200	1
7501	\N	/api/washingmachines	GET	2025-05-20 07:46:24.53772	2025-05-20 07:46:24.62697	200	89
7502	\N	/api/restaurant	GET	2025-05-20 07:46:26.245897	2025-05-20 07:46:26.245989	200	0
7503	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:46:26.280676	2025-05-20 07:46:26.282412	200	1
7504	\N	/api/weather	GET	2025-05-20 07:46:26.298426	2025-05-20 07:46:26.364604	200	66
7505	\N	/api/restaurant	GET	2025-05-20 07:46:26.38409	2025-05-20 07:46:26.384149	200	0
7506	marina.carbone@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:46:49.87716	2025-05-20 07:46:49.881571	200	4
7449	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:31:57.295186	2025-05-20 07:31:57.300226	200	5
7450	\N	/status	GET	2025-05-20 07:32:34.948665	2025-05-20 07:32:34.948669	200	0
7451	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:32:37.695876	2025-05-20 07:32:37.700554	200	4
7452	nathaniel.guitton@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-20 07:32:41.519268	2025-05-20 07:32:41.521546	200	2
7453	nathaniel.guitton@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-20 07:32:48.869186	2025-05-20 07:32:48.870749	200	1
7454	nathaniel.guitton@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-20 07:32:49.596878	2025-05-20 07:32:49.602449	200	5
7455	\N	/status	GET	2025-05-20 07:32:56.134764	2025-05-20 07:32:56.134765	200	0
7456	\N	/api/statistics/global	GET	2025-05-20 07:32:56.241143	2025-05-20 07:32:56.244157	200	3
7457	\N	/api/statistics/top-users	GET	2025-05-20 07:32:56.340029	2025-05-20 07:32:56.341812	200	1
7458	\N	/status	GET	2025-05-20 07:33:06.194658	2025-05-20 07:33:06.19466	200	0
7459	\N	/status	GET	2025-05-20 07:33:16.02864	2025-05-20 07:33:16.028641	200	0
7460	\N	/status	GET	2025-05-20 07:33:26.014439	2025-05-20 07:33:26.014441	200	0
7461	\N	/api/washingmachines	GET	2025-05-20 07:33:27.612749	2025-05-20 07:33:27.712657	200	99
7462	\N	/api/restaurant	GET	2025-05-20 07:33:35.153915	2025-05-20 07:33:35.153994	200	0
7463	\N	/status	GET	2025-05-20 07:33:36.009605	2025-05-20 07:33:36.009607	200	0
7464	\N	/api/restaurant	GET	2025-05-20 07:33:39.241733	2025-05-20 07:33:39.241799	200	0
7465	\N	/api/restaurant	GET	2025-05-20 07:33:40.450675	2025-05-20 07:33:40.450733	200	0
7466	\N	/api/restaurant	GET	2025-05-20 07:33:41.611463	2025-05-20 07:33:41.611527	200	0
7467	\N	/api/washingmachines	GET	2025-05-20 07:33:45.081604	2025-05-20 07:33:45.112389	200	30
7468	\N	/status	GET	2025-05-20 07:33:46.025899	2025-05-20 07:33:46.025901	200	0
7469	\N	/api/washingmachines	GET	2025-05-20 07:33:51.640629	2025-05-20 07:33:51.688659	200	48
7470	\N	/status	GET	2025-05-20 07:33:56.023881	2025-05-20 07:33:56.023882	200	0
7471	\N	/status	GET	2025-05-20 07:34:08.16872	2025-05-20 07:34:08.168722	200	0
7472	\N	/status	GET	2025-05-20 07:34:17.852909	2025-05-20 07:34:17.852911	200	0
7473	\N	/status	GET	2025-05-20 07:35:30.37524	2025-05-20 07:35:30.375242	200	0
7474	\N	/status	GET	2025-05-20 07:39:13.807892	2025-05-20 07:39:13.807896	200	0
7475	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:42:34.497429	2025-05-20 07:42:34.499074	200	1
7476	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:42:34.758786	2025-05-20 07:42:34.760647	200	1
7477	\N	/api/washingmachines	GET	2025-05-20 07:42:34.680716	2025-05-20 07:42:34.78137	200	100
7478	\N	/api/weather	GET	2025-05-20 07:42:34.758795	2025-05-20 07:42:34.806747	200	47
7480	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:43:21.802326	2025-05-20 07:43:21.804453	200	2
7481	\N	/api/washingmachines	GET	2025-05-20 07:43:22.00261	2025-05-20 07:43:22.033322	200	30
7482	\N	/api/restaurant	GET	2025-05-20 07:43:22.085853	2025-05-20 07:43:22.085907	200	0
7483	\N	/api/weather	GET	2025-05-20 07:43:22.091152	2025-05-20 07:43:22.091169	200	0
7514	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:48:26.992387	2025-05-20 07:48:26.994584	200	2
7515	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:48:27.058966	2025-05-20 07:48:27.061325	200	2
7516	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:48:27.131547	2025-05-20 07:48:27.134537	200	2
7517	\N	/api/washingmachines	GET	2025-05-20 07:48:27.27921	2025-05-20 07:48:27.308659	200	29
7518	\N	/api/restaurant	GET	2025-05-20 07:48:27.3618	2025-05-20 07:48:27.36188	200	0
7526	\N	/api/restaurant	GET	2025-05-20 08:05:53.03553	2025-05-20 08:05:53.035655	200	0
7527	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:05:53.035535	2025-05-20 08:05:53.040195	200	4
7528	\N	/api/washingmachines	GET	2025-05-20 08:05:52.995562	2025-05-20 08:05:53.070407	200	74
7529	\N	/api/weather	GET	2025-05-20 08:05:53.035531	2025-05-20 08:05:53.098824	200	63
7530	\N	/api/restaurant	GET	2025-05-20 08:05:55.993382	2025-05-20 08:05:55.993448	200	0
7531	\N	/api/restaurant	GET	2025-05-20 08:05:57.075254	2025-05-20 08:05:57.075326	200	0
7532	\N	/api/restaurant	GET	2025-05-20 08:05:58.307408	2025-05-20 08:05:58.307471	200	0
7533	\N	/api/restaurant	GET	2025-05-20 08:06:01.777869	2025-05-20 08:06:01.777928	200	0
7534	\N	/api/restaurant	GET	2025-05-20 08:06:03.461415	2025-05-20 08:06:03.461469	200	0
7535	\N	/api/restaurant	GET	2025-05-20 08:06:05.048652	2025-05-20 08:06:05.048717	200	0
7537	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:17:47.144144	2025-05-20 08:17:47.144221	200	0
7538	\N	/api/restaurant	GET	2025-05-20 08:17:47.254925	2025-05-20 08:17:47.25503	200	0
7539	\N	/api/newf/me	GET	2025-05-20 08:22:54.067701	2025-05-20 08:22:54.067839	401	0
7540	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:22:57.931554	2025-05-20 08:22:57.933733	200	2
7541	\N	/api/washingmachines	GET	2025-05-20 08:22:58.134376	2025-05-20 08:22:58.206623	200	72
7542	\N	/api/restaurant	GET	2025-05-20 08:22:58.271263	2025-05-20 08:22:58.271354	200	0
7507	\N	/api/washingmachines	GET	2025-05-20 07:47:59.925034	2025-05-20 07:48:00.001837	200	76
7508	\N	/status	GET	2025-05-20 07:48:01.503567	2025-05-20 07:48:01.503572	200	0
7509	\N	/api/restaurant	GET	2025-05-20 07:48:03.451431	2025-05-20 07:48:03.451498	200	0
7510	\N	/api/traq/	GET	2025-05-20 07:48:08.161116	2025-05-20 07:48:08.162292	200	1
7511	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:48:19.350385	2025-05-20 07:48:19.352271	200	1
7512	\N	/api/auth/login	POST	2025-05-20 07:48:20.680897	2025-05-20 07:48:20.761293	401	80
7513	\N	/api/auth/login	POST	2025-05-20 07:48:26.917642	2025-05-20 07:48:26.992199	200	74
7519	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 07:48:27.362715	2025-05-20 07:48:27.365754	200	3
7520	\N	/api/weather	GET	2025-05-20 07:48:27.356616	2025-05-20 07:48:27.417658	200	61
7521	marina.carbone@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 07:53:27.523264	2025-05-20 07:53:27.528745	200	5
7522	\N	/api/traq/	GET	2025-05-20 07:54:01.89776	2025-05-20 07:54:01.899201	200	1
7523	\N	/api/washingmachines	GET	2025-05-20 07:54:10.732867	2025-05-20 07:54:10.833149	200	100
7524	\N	/status	GET	2025-05-20 08:03:30.87613	2025-05-20 08:03:30.876133	200	0
7525	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:05:52.871248	2025-05-20 08:05:52.87332	200	2
7536	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:17:45.628988	2025-05-20 08:17:45.629085	200	0
7543	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:22:58.27129	2025-05-20 08:22:58.273254	200	1
7544	\N	/api/weather	GET	2025-05-20 08:22:58.231626	2025-05-20 08:22:58.284802	200	53
7545	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:27:06.742527	2025-05-20 10:27:06.82513	200	82
7546	\N	/api/weather	GET	2025-05-20 08:37:18.491438	2025-05-20 08:37:18.56025	200	68
7547	\N	/api/weather	GET	2025-05-20 08:37:20.935889	2025-05-20 08:37:20.935949	200	0
7548	\N	/api/weather	GET	2025-05-20 08:37:21.972148	2025-05-20 08:37:21.972169	200	0
7549	\N	/api/weather	GET	2025-05-20 08:43:17.097859	2025-05-20 08:43:17.143694	200	45
7550	\N	/health	GET	2025-05-20 08:43:20.800985	2025-05-20 08:43:20.800996	200	0
7551	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:43:22.889957	2025-05-20 08:43:22.890394	200	0
7552	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:43:23.119548	2025-05-20 08:43:23.119624	200	0
7553	\N	/health	GET	2025-05-20 08:43:23.257772	2025-05-20 08:43:23.257777	200	0
7554	\N	/api/restaurant	GET	2025-05-20 08:43:23.212201	2025-05-20 08:43:23.327244	200	115
7555	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:43:44.402793	2025-05-20 08:43:44.40289	200	0
7556	\N	/api/restaurant	GET	2025-05-20 08:43:44.512387	2025-05-20 08:43:44.512469	200	0
7557	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:43:44.976288	2025-05-20 08:43:44.976358	200	0
7558	\N	/api/restaurant	GET	2025-05-20 08:43:44.989714	2025-05-20 08:43:44.989804	200	0
7559	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:43:49.334882	2025-05-20 08:43:49.334964	200	0
7560	\N	/api/restaurant	GET	2025-05-20 08:43:49.62511	2025-05-20 08:43:49.625207	200	0
7561	\N	/api/restaurant	GET	2025-05-20 08:44:09.872435	2025-05-20 08:44:09.872514	200	0
7562	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:44:20.703435	2025-05-20 08:44:20.703569	200	0
7563	\N	/api/restaurant	GET	2025-05-20 08:44:21.015905	2025-05-20 08:44:21.016012	200	0
7564	\N	/health	GET	2025-05-20 08:44:26.09828	2025-05-20 08:44:26.098283	200	0
7565	\N	/api/restaurant	GET	2025-05-20 08:44:41.225829	2025-05-20 08:44:41.309913	200	84
7566	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:45:41.056744	2025-05-20 08:45:41.0572	200	0
7567	\N	/api/restaurant	GET	2025-05-20 08:45:44.734735	2025-05-20 08:45:44.734843	200	0
7568	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:46:12.889835	2025-05-20 08:46:12.892309	200	2
7569	\N	/api/restaurant	GET	2025-05-20 08:46:13.051874	2025-05-20 08:46:13.052021	200	0
7570	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:46:13.051884	2025-05-20 08:46:13.054199	200	2
7571	\N	/api/weather	GET	2025-05-20 08:46:13.051869	2025-05-20 08:46:13.098891	200	47
7572	\N	/api/washingmachines	GET	2025-05-20 08:46:13.00958	2025-05-20 08:46:13.102819	200	93
7573	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:46:15.195995	2025-05-20 08:46:15.199114	200	3
7574	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:46:16.367594	2025-05-20 08:46:16.369593	200	1
7575	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:46:17.859	2025-05-20 08:46:17.85913	200	0
7576	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 08:46:18.844067	2025-05-20 08:46:18.849649	200	5
7577	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 08:46:21.66045	2025-05-20 08:46:21.664933	200	4
7578	\N	/api/restaurant	GET	2025-05-20 08:46:23.338274	2025-05-20 08:46:24.391953	200	1053
7579	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:04.23378	2025-05-20 08:47:04.235836	200	2
7580	\N	/api/washingmachines	GET	2025-05-20 08:47:04.592673	2025-05-20 08:47:04.635807	200	43
7581	\N	/api/weather	GET	2025-05-20 08:47:04.872679	2025-05-20 08:47:04.872702	200	0
7582	\N	/api/restaurant	GET	2025-05-20 08:47:04.913444	2025-05-20 08:47:04.913581	200	0
7583	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:04.919127	2025-05-20 08:47:04.921147	200	2
7584	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:06.79151	2025-05-20 08:47:06.793387	200	1
7585	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:08.031367	2025-05-20 08:47:08.032859	200	1
7586	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 08:47:10.405013	2025-05-20 08:47:10.409594	200	4
7587	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:15.643267	2025-05-20 08:47:15.645058	200	1
7588	\N	/api/washingmachines	GET	2025-05-20 08:47:16.044358	2025-05-20 08:47:16.073623	200	29
7589	\N	/api/weather	GET	2025-05-20 08:47:16.197763	2025-05-20 08:47:16.213749	200	15
7590	\N	/api/restaurant	GET	2025-05-20 08:47:16.57147	2025-05-20 08:47:16.571549	200	0
7591	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:47:16.576942	2025-05-20 08:47:16.578552	200	1
7592	\N	/api/restaurant	GET	2025-05-20 08:47:35.125937	2025-05-20 08:47:35.126009	200	0
7593	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:48:32.082995	2025-05-20 08:48:32.083104	200	0
7594	\N	/api/restaurant	GET	2025-05-20 08:48:32.32125	2025-05-20 08:48:32.321343	200	0
7595	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:48:36.660852	2025-05-20 08:48:36.660924	200	0
7596	\N	/api/restaurant	GET	2025-05-20 08:48:36.918631	2025-05-20 08:48:36.918742	200	0
7597	\N	/api/restaurant	GET	2025-05-20 08:48:55.791597	2025-05-20 08:48:55.791689	200	0
7598	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:48:55.791667	2025-05-20 08:48:55.791801	200	0
7599	\N	/api/washingmachines	GET	2025-05-20 08:49:08.719814	2025-05-20 08:49:08.812591	200	92
7600	\N	/api/washingmachines	GET	2025-05-20 08:49:11.737681	2025-05-20 08:49:11.768146	200	30
7601	\N	/status	GET	2025-05-20 08:49:31.751945	2025-05-20 08:49:31.751947	200	0
7602	\N	/api/statistics/global	GET	2025-05-20 08:49:31.798494	2025-05-20 08:49:31.801111	200	2
7603	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:31.849235	2025-05-20 08:49:31.858416	200	9
7604	\N	/status	GET	2025-05-20 08:49:36.442809	2025-05-20 08:49:36.442811	200	0
7605	\N	/api/statistics/global	GET	2025-05-20 08:49:36.500845	2025-05-20 08:49:36.503039	200	2
7606	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:36.547758	2025-05-20 08:49:36.553698	200	5
7608	\N	/status	GET	2025-05-20 08:49:36.757386	2025-05-20 08:49:36.75739	200	0
7607	\N	/status	GET	2025-05-20 08:49:36.757369	2025-05-20 08:49:36.757371	200	0
7609	\N	/api/statistics/global	GET	2025-05-20 08:49:36.810998	2025-05-20 08:49:36.813448	200	2
7610	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:36.861696	2025-05-20 08:49:36.867405	200	5
7611	\N	/status	GET	2025-05-20 08:49:41.478219	2025-05-20 08:49:41.478221	200	0
7612	\N	/status	GET	2025-05-20 08:49:41.483433	2025-05-20 08:49:41.483436	200	0
7613	\N	/api/statistics/global	GET	2025-05-20 08:49:41.531228	2025-05-20 08:49:41.533902	200	2
7614	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:41.580126	2025-05-20 08:49:41.588347	200	8
7615	\N	/status	GET	2025-05-20 08:49:41.755058	2025-05-20 08:49:41.755059	200	0
7616	\N	/status	GET	2025-05-20 08:49:41.75571	2025-05-20 08:49:41.755711	200	0
7617	\N	/api/statistics/global	GET	2025-05-20 08:49:41.80435	2025-05-20 08:49:41.806593	200	2
7618	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:41.855831	2025-05-20 08:49:41.862459	200	6
7619	\N	/status	GET	2025-05-20 08:49:46.448671	2025-05-20 08:49:46.448672	200	0
7621	\N	/api/statistics/global	GET	2025-05-20 08:49:46.497243	2025-05-20 08:49:46.49973	200	2
7622	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:46.548489	2025-05-20 08:49:46.555169	200	6
7624	\N	/status	GET	2025-05-20 08:49:46.759016	2025-05-20 08:49:46.759016	200	0
7625	\N	/api/statistics/global	GET	2025-05-20 08:49:46.812152	2025-05-20 08:49:46.814425	200	2
7626	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:46.863993	2025-05-20 08:49:46.870467	200	6
7627	\N	/status	GET	2025-05-20 08:49:51.452591	2025-05-20 08:49:51.452592	200	0
7631	\N	/status	GET	2025-05-20 08:49:51.759305	2025-05-20 08:49:51.759306	200	0
7635	\N	/status	GET	2025-05-20 08:49:56.586019	2025-05-20 08:49:56.58602	200	0
7639	\N	/status	GET	2025-05-20 08:49:56.761174	2025-05-20 08:49:56.761175	200	0
7640	\N	/api/statistics/global	GET	2025-05-20 08:49:56.812306	2025-05-20 08:49:56.814463	200	2
7641	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:56.86177	2025-05-20 08:49:56.867657	200	5
7643	\N	/api/restaurant	GET	2025-05-20 08:49:57.253211	2025-05-20 08:49:57.253301	200	0
7644	\N	/status	GET	2025-05-20 08:49:58.472195	2025-05-20 08:49:58.472197	200	0
7645	\N	/api/statistics/global	GET	2025-05-20 08:49:58.520238	2025-05-20 08:49:58.52245	200	2
7646	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:58.576414	2025-05-20 08:49:58.582322	200	5
7647	\N	/status	GET	2025-05-20 08:50:01.922407	2025-05-20 08:50:01.922408	200	0
7651	\N	/status	GET	2025-05-20 08:50:03.474909	2025-05-20 08:50:03.47491	200	0
7656	\N	/status	GET	2025-05-20 08:50:06.872226	2025-05-20 08:50:06.872227	200	0
7657	\N	/api/statistics/global	GET	2025-05-20 08:50:06.962008	2025-05-20 08:50:06.964631	200	2
7658	\N	/api/statistics/endpoints	GET	2025-05-20 08:50:07.165997	2025-05-20 08:50:07.171846	200	5
7659	\N	/status	GET	2025-05-20 08:50:08.506696	2025-05-20 08:50:08.506698	200	0
7663	\N	/status	GET	2025-05-20 08:50:13.583047	2025-05-20 08:50:13.583047	200	0
7670	\N	/api/restaurant	GET	2025-05-20 08:52:34.331347	2025-05-20 08:52:34.331452	200	0
7672	\N	/api/washingmachines	GET	2025-05-20 08:52:45.488038	2025-05-20 08:52:45.587937	200	99
7673	\N	/status	GET	2025-05-20 08:52:55.143766	2025-05-20 08:52:55.143768	200	0
7674	\N	/api/statistics/global	GET	2025-05-20 08:52:55.205735	2025-05-20 08:52:55.207948	200	2
7675	\N	/api/statistics/endpoints	GET	2025-05-20 08:52:55.260588	2025-05-20 08:52:55.268214	200	7
7676	\N	/api/washingmachines	GET	2025-05-20 08:52:58.942117	2025-05-20 08:52:58.971995	200	29
7677	\N	/status	GET	2025-05-20 08:53:07.796955	2025-05-20 08:53:07.796957	200	0
7678	\N	/api/statistics/global	GET	2025-05-20 08:53:07.855745	2025-05-20 08:53:07.858495	200	2
7679	\N	/api/statistics/endpoints	GET	2025-05-20 08:53:07.921377	2025-05-20 08:53:07.928262	200	6
7680	\N	/api/washingmachines	GET	2025-05-20 08:53:08.505228	2025-05-20 08:53:08.53418	200	28
7681	\N	/api/restaurant	GET	2025-05-20 08:53:09.059698	2025-05-20 08:53:09.059801	200	0
7683	\N	/status	GET	2025-05-20 08:53:29.24412	2025-05-20 08:53:29.244124	200	0
7684	\N	/api/statistics/global	GET	2025-05-20 08:53:29.293527	2025-05-20 08:53:29.29565	200	2
7685	\N	/api/statistics/endpoints	GET	2025-05-20 08:53:29.341175	2025-05-20 08:53:29.346865	200	5
7628	\N	/status	GET	2025-05-20 08:49:51.45259	2025-05-20 08:49:51.452592	200	0
7629	\N	/api/statistics/global	GET	2025-05-20 08:49:51.503844	2025-05-20 08:49:51.506111	200	2
7630	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:51.553765	2025-05-20 08:49:51.559889	200	6
7632	\N	/status	GET	2025-05-20 08:49:51.759263	2025-05-20 08:49:51.759264	200	0
7633	\N	/api/statistics/global	GET	2025-05-20 08:49:51.809658	2025-05-20 08:49:51.811768	200	2
7634	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:51.865029	2025-05-20 08:49:51.87018	200	5
7636	\N	/status	GET	2025-05-20 08:49:56.585829	2025-05-20 08:49:56.58583	200	0
7637	\N	/api/statistics/global	GET	2025-05-20 08:49:56.692716	2025-05-20 08:49:56.695341	200	2
7638	\N	/status	GET	2025-05-20 08:49:56.761247	2025-05-20 08:49:56.761248	200	0
7642	\N	/api/statistics/endpoints	GET	2025-05-20 08:49:56.755819	2025-05-20 08:49:56.76206	200	6
7648	\N	/status	GET	2025-05-20 08:50:01.922657	2025-05-20 08:50:01.922658	200	0
7649	\N	/api/statistics/global	GET	2025-05-20 08:50:02.018629	2025-05-20 08:50:02.020463	200	1
7650	\N	/api/statistics/endpoints	GET	2025-05-20 08:50:02.107881	2025-05-20 08:50:02.113574	200	5
7652	\N	/status	GET	2025-05-20 08:50:03.474909	2025-05-20 08:50:03.47491	200	0
7653	\N	/api/statistics/global	GET	2025-05-20 08:50:03.545938	2025-05-20 08:50:03.547937	200	1
7654	\N	/api/statistics/endpoints	GET	2025-05-20 08:50:03.597548	2025-05-20 08:50:03.603036	200	5
7655	\N	/status	GET	2025-05-20 08:50:06.872017	2025-05-20 08:50:06.872018	200	0
7660	\N	/status	GET	2025-05-20 08:50:08.506637	2025-05-20 08:50:08.506639	200	0
7661	\N	/api/statistics/global	GET	2025-05-20 08:50:08.55547	2025-05-20 08:50:08.557529	200	2
7662	\N	/api/statistics/endpoints	GET	2025-05-20 08:50:08.604541	2025-05-20 08:50:08.609866	200	5
7664	\N	/status	GET	2025-05-20 08:50:13.58304	2025-05-20 08:50:13.583042	200	0
7665	\N	/api/statistics/global	GET	2025-05-20 08:50:13.634151	2025-05-20 08:50:13.636067	200	1
7666	\N	/api/statistics/endpoints	GET	2025-05-20 08:50:13.685166	2025-05-20 08:50:13.691937	200	6
7667	\N	/api/upload/	POST	2025-05-20 10:51:38.326129	2025-05-20 10:51:38.326174	500	0
7668	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:52:33.922036	2025-05-20 08:52:33.922144	200	0
7669	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 08:52:34.115961	2025-05-20 08:52:34.116047	200	0
7671	yohann.chavanel@imt-atlantique.net	/api/upload	POST	2025-05-20 10:52:38.749063	2025-05-20 10:52:38.840891	200	91
7682	\N	/status	GET	2025-05-20 08:53:29.244109	2025-05-20 08:53:29.244112	200	0
7686	\N	/health	GET	2025-05-20 08:53:56.011432	2025-05-20 08:53:56.01144	200	0
7687	\N	/api/restaurant	GET	2025-05-20 08:54:56.181791	2025-05-20 08:54:56.293713	200	111
7688	\N	/api/washingmachines	GET	2025-05-20 08:54:57.749497	2025-05-20 08:54:57.83079	200	81
7689	\N	/api/restaurant	GET	2025-05-20 08:55:05.462519	2025-05-20 08:55:06.583521	200	1121
7690	\N	/api/restaurant	GET	2025-05-20 08:55:05.655173	2025-05-20 08:55:06.583694	200	928
7691	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:07.046484	2025-05-20 08:55:07.048924	200	2
7692	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:13.334647	2025-05-20 08:55:13.342776	200	8
7693	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:13.495874	2025-05-20 08:55:13.498059	200	2
7694	\N	/api/washingmachines	GET	2025-05-20 08:55:13.465267	2025-05-20 08:55:13.51407	200	48
7695	\N	/api/weather	GET	2025-05-20 08:55:13.496029	2025-05-20 08:55:13.569393	200	73
7696	\N	/api/restaurant	GET	2025-05-20 08:55:13.495883	2025-05-20 08:55:13.970776	200	474
7697	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:19.518822	2025-05-20 08:55:19.521928	200	3
7698	\N	/api/auth/login	POST	2025-05-20 08:55:25.923168	2025-05-20 08:55:25.998266	200	75
7699	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:26.173317	2025-05-20 08:55:26.175472	200	2
7700	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:26.249669	2025-05-20 08:55:26.251441	200	1
7701	\N	/api/washingmachines	GET	2025-05-20 08:55:26.361194	2025-05-20 08:55:26.395295	200	34
7702	\N	/api/weather	GET	2025-05-20 08:55:26.753712	2025-05-20 08:55:26.7771	200	23
7703	\N	/api/restaurant	GET	2025-05-20 08:55:26.799463	2025-05-20 08:55:26.799533	200	0
7704	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:26.804813	2025-05-20 08:55:26.809653	200	4
7705	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:30.753562	2025-05-20 08:55:30.755142	200	1
7706	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:32.860536	2025-05-20 08:55:32.863375	200	2
7707	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:34.29829	2025-05-20 08:55:34.300235	200	1
7708	test@imt-atlantique.net	/api/upload	POST	2025-05-20 08:55:55.819702	2025-05-20 08:55:55.826415	200	6
7709	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 08:55:55.882671	2025-05-20 08:55:55.888694	200	6
7710	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:56.118175	2025-05-20 08:55:56.120177	200	2
7711	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:55:57.872774	2025-05-20 08:55:57.874304	200	1
7712	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:08.384074	2025-05-20 08:56:08.386106	200	2
7713	\N	/api/washingmachines	GET	2025-05-20 08:56:08.503635	2025-05-20 08:56:08.536604	200	32
7714	\N	/api/weather	GET	2025-05-20 08:56:08.551177	2025-05-20 08:56:08.570029	200	18
7715	\N	/api/restaurant	GET	2025-05-20 08:56:08.588326	2025-05-20 08:56:08.588421	200	0
7716	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:08.593597	2025-05-20 08:56:08.595213	200	1
7717	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:09.787287	2025-05-20 08:56:09.789164	200	1
7718	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:10.985874	2025-05-20 08:56:10.987432	200	1
7719	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:11.727575	2025-05-20 08:56:11.730145	200	2
7720	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:12.4933	2025-05-20 08:56:12.494963	200	1
7721	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:13.366547	2025-05-20 08:56:13.370501	200	3
7722	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:15.18652	2025-05-20 08:56:15.189441	200	2
7723	test@imt-atlantique.net	/api/upload	POST	2025-05-20 08:56:19.601731	2025-05-20 08:56:19.607803	200	6
7724	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 08:56:19.674596	2025-05-20 08:56:19.681403	200	6
7725	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:19.862109	2025-05-20 08:56:19.864142	200	2
7726	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:21.632964	2025-05-20 08:56:21.63479	200	1
7727	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:22.508844	2025-05-20 08:56:22.510653	200	1
7728	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:25.432311	2025-05-20 08:56:25.434248	200	1
7729	\N	/api/washingmachines	GET	2025-05-20 08:56:25.560121	2025-05-20 08:56:25.594375	200	34
7730	\N	/api/weather	GET	2025-05-20 08:56:25.609963	2025-05-20 08:56:25.609986	200	0
7731	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:25.643744	2025-05-20 08:56:25.645965	200	2
7732	\N	/api/restaurant	GET	2025-05-20 08:56:25.649151	2025-05-20 08:56:25.64921	200	0
7733	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:26.267383	2025-05-20 08:56:26.26914	200	1
7734	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:27.23584	2025-05-20 08:56:27.238313	200	2
7735	\N	/api/auth/login	POST	2025-05-20 08:56:43.276087	2025-05-20 08:56:43.349646	200	73
7736	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:43.509644	2025-05-20 08:56:43.512441	200	2
7737	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:43.577473	2025-05-20 08:56:43.579451	200	1
7738	\N	/api/washingmachines	GET	2025-05-20 08:56:43.692148	2025-05-20 08:56:43.741591	200	49
7739	\N	/api/weather	GET	2025-05-20 08:56:44.076174	2025-05-20 08:56:44.076196	200	0
7740	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:44.130293	2025-05-20 08:56:44.132774	200	2
7741	\N	/api/restaurant	GET	2025-05-20 08:56:44.135667	2025-05-20 08:56:44.135726	200	0
7742	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 08:56:45.13812	2025-05-20 08:56:45.140285	200	2
7743	\N	/status	GET	2025-05-20 08:58:29.702521	2025-05-20 08:58:29.702524	200	0
7744	\N	/status	GET	2025-05-20 08:58:29.702531	2025-05-20 08:58:29.702534	200	0
7745	\N	/api/statistics/global	GET	2025-05-20 08:58:29.74752	2025-05-20 08:58:29.750518	200	2
7746	\N	/api/statistics/endpoints	GET	2025-05-20 08:58:29.794518	2025-05-20 08:58:29.800296	200	5
7747	\N	/api/auth/login	POST	2025-05-20 09:00:20.762405	2025-05-20 09:00:20.764625	401	2
7748	\N	/api/auth/register	POST	2025-05-20 09:00:54.658986	2025-05-20 09:00:54.744424	201	85
7749	\N	/api/auth/verify-account	POST	2025-05-20 09:01:23.265438	2025-05-20 09:01:23.274187	200	8
7750	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:01:23.352846	2025-05-20 09:01:23.354866	200	2
7751	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:01:23.434163	2025-05-20 09:01:23.436158	200	1
7752	\N	/api/washingmachines	GET	2025-05-20 09:01:23.536527	2025-05-20 09:01:23.644159	200	107
7754	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:01:23.91255	2025-05-20 09:01:23.914273	200	1
7755	\N	/api/weather	GET	2025-05-20 09:01:23.906661	2025-05-20 09:01:23.953201	200	46
7756	aurelien.pautet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 09:01:26.326699	2025-05-20 09:01:26.336597	200	9
7757	\N	/api/washingmachines	GET	2025-05-20 09:01:39.198052	2025-05-20 09:01:39.226375	200	28
7758	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:01:47.797041	2025-05-20 09:01:47.799126	200	2
7759	\N	/api/traq/	GET	2025-05-20 09:01:53.738884	2025-05-20 09:01:53.740392	200	1
7760	\N	/api/restaurant	GET	2025-05-20 09:02:07.844247	2025-05-20 09:02:07.844323	200	0
7764	\N	/api/restaurant	GET	2025-05-20 09:02:58.511234	2025-05-20 09:02:58.51136	200	0
7765	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:05:00.953475	2025-05-20 09:05:00.95738	200	3
7766	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:05:01.120091	2025-05-20 09:05:01.122353	200	2
7767	test@imt-atlantique.net	/api/upload	POST	2025-05-20 09:05:06.140084	2025-05-20 09:05:06.150237	200	10
7768	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 09:05:06.202853	2025-05-20 09:05:06.23925	200	36
7769	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:05:06.329872	2025-05-20 09:05:06.333417	200	3
7772	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:10:40.43196	2025-05-20 09:10:40.432057	200	0
7773	\N	/api/restaurant	GET	2025-05-20 09:10:42.887035	2025-05-20 09:10:42.887179	200	0
7753	\N	/api/restaurant	GET	2025-05-20 09:01:23.912552	2025-05-20 09:01:23.912805	200	0
7761	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:02:33.591568	2025-05-20 09:02:33.591822	200	0
7762	\N	/api/restaurant	GET	2025-05-20 09:02:34.452264	2025-05-20 09:02:34.452405	200	0
7763	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:02:58.049984	2025-05-20 09:02:58.050081	200	0
7770	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:09:55.649172	2025-05-20 09:09:55.649295	200	0
7771	\N	/api/restaurant	GET	2025-05-20 09:09:55.863872	2025-05-20 09:09:55.864009	200	0
7774	\N	/api/user	PUT	2025-05-20 11:13:26.723038	2025-05-20 11:13:26.72306	500	0
7775	\N	/health	GET	2025-05-20 09:17:23.249259	2025-05-20 09:17:23.249263	200	0
7776	test@imt-atlantique.net	/api/upload	POST	2025-05-20 09:17:37.288103	2025-05-20 09:17:37.293558	200	5
7777	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 09:17:37.346241	2025-05-20 09:17:37.352491	200	6
7778	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:37.422296	2025-05-20 09:17:37.425617	200	3
7779	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:39.239327	2025-05-20 09:17:39.241073	200	1
7780	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:39.498482	2025-05-20 09:17:39.500149	200	1
7781	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:39.767578	2025-05-20 09:17:39.76919	200	1
7782	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:40.032004	2025-05-20 09:17:40.033827	200	1
7783	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:17:40.322316	2025-05-20 09:17:40.324214	200	1
7784	\N	/api/data/467bdeed-c69f-4074-bb6e-765100a569f5_85fe4e125190b7b4.jpeg	GET	2025-05-20 09:17:41.812327	2025-05-20 09:17:41.812582	200	0
7785	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:18:48.846031	2025-05-20 09:18:48.849671	200	3
7786	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:18:49.196799	2025-05-20 09:18:49.198683	200	1
7787	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:18:49.584538	2025-05-20 09:18:49.586255	200	1
7788	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:18:50.34219	2025-05-20 09:18:50.344259	200	2
7789	yohann.chavanel@imt-atlantique.net	/api/files	GET	2025-05-20 11:24:17.160713	2025-05-20 11:24:17.263752	200	103
7790	yohann.chavanel@imt-atlantique.net	/api/all-files	GET	2025-05-20 11:24:25.312144	2025-05-20 11:24:25.31279	200	0
7791	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:24:58.055705	2025-05-20 09:24:58.057803	200	2
7792	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:24:58.063708	2025-05-20 09:24:58.066365	200	2
7793	\N	/api/washingmachines	GET	2025-05-20 09:25:02.026721	2025-05-20 09:25:02.103379	200	76
7794	\N	/api/weather	GET	2025-05-20 09:25:02.035209	2025-05-20 09:25:02.111532	200	76
7795	\N	/api/restaurant	GET	2025-05-20 09:25:02.02864	2025-05-20 09:25:02.947992	200	919
7796	\N	/health	GET	2025-05-20 09:27:46.309174	2025-05-20 09:27:46.309178	200	0
7797	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:02.340552	2025-05-20 09:29:02.343688	200	3
7798	\N	/api/washingmachines	GET	2025-05-20 09:29:02.697589	2025-05-20 09:29:02.774596	200	77
7799	\N	/api/weather	GET	2025-05-20 09:29:02.83884	2025-05-20 09:29:02.899886	200	61
7800	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:02.904469	2025-05-20 09:29:03.094336	200	189
7801	\N	/api/restaurant	GET	2025-05-20 09:29:02.904266	2025-05-20 09:29:03.627276	200	723
7802	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:24.041686	2025-05-20 09:29:24.043597	200	1
7803	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:24.905101	2025-05-20 09:29:24.906745	200	1
7804	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:26.755882	2025-05-20 09:29:26.758036	200	2
7805	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:28.528356	2025-05-20 09:29:28.530428	200	2
7806	nathaniel.guitton@imt-atlantique.net	/api/upload	POST	2025-05-20 09:29:57.730217	2025-05-20 09:29:57.736821	200	6
7807	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 09:29:57.860016	2025-05-20 09:29:57.865624	200	5
7808	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:29:58.132863	2025-05-20 09:29:58.134869	200	2
7809	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:29:58.806896	2025-05-20 09:29:58.807215	200	0
7810	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:30:36.488121	2025-05-20 09:30:36.488226	200	0
7811	\N	/favicon.ico	GET	2025-05-20 09:30:36.836723	2025-05-20 09:30:36.836728	500	0
7812	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:52.21249	2025-05-20 09:30:52.214645	200	2
7813	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:52.354509	2025-05-20 09:30:52.356372	200	1
7814	\N	/api/washingmachines	GET	2025-05-20 09:30:52.319362	2025-05-20 09:30:52.419142	200	99
7815	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:30:52.42089	2025-05-20 09:30:52.42098	200	0
7816	\N	/api/weather	GET	2025-05-20 09:30:52.353654	2025-05-20 09:30:52.438228	200	84
7817	\N	/api/restaurant	GET	2025-05-20 09:30:52.354621	2025-05-20 09:30:52.458281	200	103
7818	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:53.211656	2025-05-20 09:30:53.213476	200	1
7819	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-20 09:30:53.277098	2025-05-20 09:30:53.27727	200	0
7820	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:30:53.83871	2025-05-20 09:30:53.838802	200	0
7821	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:55.027128	2025-05-20 09:30:55.029058	200	1
7822	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-20 09:30:55.096728	2025-05-20 09:30:55.096804	200	0
7823	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:56.099372	2025-05-20 09:30:56.101131	200	1
7824	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-20 09:30:56.177686	2025-05-20 09:30:56.177759	200	0
7825	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:56.208549	2025-05-20 09:30:56.210546	200	1
7826	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:30:57.784306	2025-05-20 09:30:57.786129	200	1
7827	\N	/api/data/29daf123-19cd-4d46-89cf-a06090f86da4.jpeg	GET	2025-05-20 09:30:57.854543	2025-05-20 09:30:57.85462	200	0
7828	\N	/health	GET	2025-05-20 09:36:18.588678	2025-05-20 09:36:18.588681	200	0
7829	\N	/status	GET	2025-05-20 09:45:18.652446	2025-05-20 09:45:18.652452	200	0
7830	\N	/api/statistics/global	GET	2025-05-20 09:45:18.706477	2025-05-20 09:45:18.709727	200	3
7831	\N	/api/statistics/endpoints	GET	2025-05-20 09:45:18.759189	2025-05-20 09:45:18.766744	200	7
7832	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:45:18.652416	2025-05-20 09:45:18.652925	200	0
7833	\N	/status	GET	2025-05-20 09:45:22.619317	2025-05-20 09:45:22.61932	200	0
7834	\N	/status	GET	2025-05-20 09:45:22.621315	2025-05-20 09:45:22.621316	200	0
7835	\N	/api/statistics/global	GET	2025-05-20 09:45:22.665047	2025-05-20 09:45:22.66817	200	3
7836	\N	/api/statistics/endpoints	GET	2025-05-20 09:45:22.713632	2025-05-20 09:45:22.719606	200	5
7837	\N	/status	GET	2025-05-20 09:45:27.633443	2025-05-20 09:45:27.633445	200	0
7838	\N	/status	GET	2025-05-20 09:45:27.633668	2025-05-20 09:45:27.63367	200	0
7839	\N	/api/statistics/global	GET	2025-05-20 09:45:27.684172	2025-05-20 09:45:27.687105	200	2
7840	\N	/api/statistics/endpoints	GET	2025-05-20 09:45:27.736138	2025-05-20 09:45:27.744156	200	8
7841	\N	/status	GET	2025-05-20 09:45:32.644075	2025-05-20 09:45:32.644077	200	0
7842	\N	/status	GET	2025-05-20 09:45:32.643971	2025-05-20 09:45:32.643973	200	0
7843	\N	/api/statistics/global	GET	2025-05-20 09:45:32.69664	2025-05-20 09:45:32.698913	200	2
7844	\N	/api/statistics/endpoints	GET	2025-05-20 09:45:32.748822	2025-05-20 09:45:32.755066	200	6
7845	\N	/status	GET	2025-05-20 09:46:01.257293	2025-05-20 09:46:01.257297	200	0
7846	\N	/status	GET	2025-05-20 09:46:01.257352	2025-05-20 09:46:01.257354	200	0
7847	\N	/api/statistics/global	GET	2025-05-20 09:46:01.340413	2025-05-20 09:46:01.342882	200	2
7848	\N	/api/statistics/endpoints	GET	2025-05-20 09:46:01.509225	2025-05-20 09:46:01.515608	200	6
7849	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:46:45.363167	2025-05-20 09:46:45.366564	200	3
7850	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 09:46:45.661178	2025-05-20 09:46:45.66305	200	1
7851	\N	/api/washingmachines	GET	2025-05-20 09:46:45.626953	2025-05-20 09:46:45.697137	200	70
7852	\N	/api/weather	GET	2025-05-20 09:46:45.661183	2025-05-20 09:46:45.731721	200	70
7853	\N	/api/restaurant	GET	2025-05-20 09:46:45.661181	2025-05-20 09:46:45.755272	200	94
7854	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:47:10.742894	2025-05-20 09:47:10.74301	200	0
7855	\N	/api/restaurant	GET	2025-05-20 09:47:11.738092	2025-05-20 09:47:11.738222	200	0
7856	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:48:04.962946	2025-05-20 09:48:04.963098	200	0
7857	\N	/api/data/b8aa4b8d-ce3c-4528-9936-71e3aee47ced_569efa10a6cfed29.jpeg	GET	2025-05-20 09:48:07.032788	2025-05-20 09:48:07.032857	200	0
7858	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 09:57:39.149653	2025-05-20 09:57:39.149781	200	0
7859	\N	/api/restaurant	GET	2025-05-20 10:03:00.641663	2025-05-20 10:03:00.641801	200	0
7860	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:04:20.003024	2025-05-20 10:04:20.005605	200	2
7861	\N	/api/washingmachines	GET	2025-05-20 10:04:20.355787	2025-05-20 10:04:20.460096	200	104
7862	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:04:20.576683	2025-05-20 10:04:20.578931	200	2
7863	\N	/api/weather	GET	2025-05-20 10:04:20.527368	2025-05-20 10:04:20.596469	200	69
7864	\N	/api/restaurant	GET	2025-05-20 10:04:20.576621	2025-05-20 10:04:21.563359	200	986
7865	\N	/api/washingmachines	GET	2025-05-20 10:04:57.684212	2025-05-20 10:04:57.721145	200	36
7866	\N	/api/washingmachines	GET	2025-05-20 10:05:08.610215	2025-05-20 10:05:08.643797	200	33
7867	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:05:26.200152	2025-05-20 10:05:26.202084	200	1
7868	\N	/api/weather	GET	2025-05-20 10:05:26.380096	2025-05-20 10:05:26.38012	200	0
7885	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 10:06:49.588263	2025-05-20 10:06:49.588372	200	0
7886	\N	/api/restaurant	GET	2025-05-20 10:06:49.763596	2025-05-20 10:06:49.763723	200	0
7887	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:07:42.122975	2025-05-20 10:07:42.12467	200	1
7902	\N	/api/restaurant	GET	2025-05-20 10:17:24.961497	2025-05-20 10:17:24.96159	200	0
7904	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 10:37:52.798127	2025-05-20 10:37:52.798264	200	0
7905	\N	/api/washingmachines	GET	2025-05-20 10:37:52.984433	2025-05-20 10:37:53.075441	200	91
7906	\N	/robots.txt	GET	2025-05-20 10:56:39.644976	2025-05-20 10:56:39.644984	500	0
7907	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:00:53.034487	2025-05-20 11:00:53.037401	200	2
7908	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:01:02.677407	2025-05-20 11:01:02.679746	200	2
7909	\N	/status	GET	2025-05-20 11:01:07.443664	2025-05-20 11:01:07.443665	200	0
7910	\N	/api/statistics/global	GET	2025-05-20 11:01:07.511261	2025-05-20 11:01:07.513768	200	2
7911	\N	/api/statistics/top-users	GET	2025-05-20 11:01:07.581557	2025-05-20 11:01:07.583826	200	2
7913	\N	/api/auth/register	POST	2025-05-20 11:01:15.722734	2025-05-20 11:01:15.807561	409	84
7914	\N	/status	GET	2025-05-20 11:01:17.462463	2025-05-20 11:01:17.462465	200	0
7915	\N	/status	GET	2025-05-20 11:01:27.451027	2025-05-20 11:01:27.451028	200	0
7916	\N	/status	GET	2025-05-20 11:01:37.465069	2025-05-20 11:01:37.46507	200	0
7917	\N	/status	GET	2025-05-20 11:01:47.456913	2025-05-20 11:01:47.456915	200	0
7918	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:01:56.100763	2025-05-20 11:01:56.102968	200	2
7919	\N	/api/auth/verify-account	POST	2025-05-20 11:02:00.707479	2025-05-20 11:02:00.714938	200	7
7920	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:00.852407	2025-05-20 11:02:00.854677	200	2
7921	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:01.007433	2025-05-20 11:02:01.00984	200	2
7922	\N	/api/washingmachines	GET	2025-05-20 11:02:01.203112	2025-05-20 11:02:01.300792	200	97
7923	\N	/api/restaurant	GET	2025-05-20 11:02:01.556507	2025-05-20 11:02:01.556637	200	0
7935	\N	/api/data/2F773180-2F5B-41EF-8DA2-1109D921BBE0_4d1b0b7d141e7f9e.jpg	GET	2025-05-20 11:03:20.392736	2025-05-20 11:03:20.392835	200	0
7936	\N	/status	GET	2025-05-20 11:03:22.792569	2025-05-20 11:03:22.792571	200	0
7937	\N	/api/statistics/global	GET	2025-05-20 11:03:22.849971	2025-05-20 11:03:22.852376	200	2
7938	\N	/api/statistics/top-users	GET	2025-05-20 11:03:22.901025	2025-05-20 11:03:22.902568	200	1
7939	\N	/status	GET	2025-05-20 11:03:32.799237	2025-05-20 11:03:32.799238	200	0
7940	\N	/status	GET	2025-05-20 11:03:42.802737	2025-05-20 11:03:42.802738	200	0
7941	\N	/status	GET	2025-05-20 11:03:52.804006	2025-05-20 11:03:52.804008	200	0
7942	\N	/status	GET	2025-05-20 11:04:02.346142	2025-05-20 11:04:02.346144	200	0
7943	\N	/api/statistics/global	GET	2025-05-20 11:04:02.413859	2025-05-20 11:04:02.416042	200	2
7944	\N	/api/statistics/top-users	GET	2025-05-20 11:04:02.478799	2025-05-20 11:04:02.480882	200	2
7945	\N	/status	GET	2025-05-20 11:04:02.796169	2025-05-20 11:04:02.796171	200	0
7946	\N	/status	GET	2025-05-20 11:04:03.791785	2025-05-20 11:04:03.791787	200	0
7947	\N	/api/statistics/global	GET	2025-05-20 11:04:03.867568	2025-05-20 11:04:03.869831	200	2
7948	\N	/api/statistics/top-users	GET	2025-05-20 11:04:03.932976	2025-05-20 11:04:03.93472	200	1
7949	\N	/status	GET	2025-05-20 11:04:07.443867	2025-05-20 11:04:07.443868	200	0
7950	\N	/api/statistics/global	GET	2025-05-20 11:04:07.509959	2025-05-20 11:04:07.512095	200	2
7951	\N	/api/statistics/top-users	GET	2025-05-20 11:04:07.574071	2025-05-20 11:04:07.575658	200	1
7952	\N	/status	GET	2025-05-20 11:04:12.801185	2025-05-20 11:04:12.801187	200	0
7869	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:05:26.380095	2025-05-20 10:05:26.381614	200	1
7870	\N	/api/washingmachines	GET	2025-05-20 10:05:26.373679	2025-05-20 10:05:26.407084	200	33
7871	\N	/api/restaurant	GET	2025-05-20 10:05:26.380355	2025-05-20 10:05:26.380452	200	0
7872	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:05:28.771637	2025-05-20 10:05:28.773535	200	1
7873	\N	/api/washingmachines	GET	2025-05-20 10:05:29.121312	2025-05-20 10:05:29.158413	200	37
7874	\N	/api/weather	GET	2025-05-20 10:05:29.246137	2025-05-20 10:05:29.246167	200	0
7875	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:05:29.28508	2025-05-20 10:05:29.287117	200	2
7876	\N	/api/restaurant	GET	2025-05-20 10:05:29.29538	2025-05-20 10:05:29.295478	200	0
7877	\N	/api/washingmachines	GET	2025-05-20 10:05:30.215161	2025-05-20 10:05:30.248312	200	33
7878	\N	/api/washingmachines	GET	2025-05-20 10:05:31.221697	2025-05-20 10:05:31.254501	200	32
7879	\N	/api/washingmachines	GET	2025-05-20 10:05:31.910964	2025-05-20 10:05:31.945324	200	34
7880	\N	/api/washingmachines	GET	2025-05-20 10:05:32.839066	2025-05-20 10:05:32.871677	200	32
7881	\N	/api/washingmachines	GET	2025-05-20 10:05:33.896643	2025-05-20 10:05:33.930346	200	33
7882	\N	/api/washingmachines	GET	2025-05-20 10:05:34.324651	2025-05-20 10:05:34.359008	200	34
7883	\N	/api/washingmachines	GET	2025-05-20 10:05:35.732997	2025-05-20 10:05:35.779344	200	46
7884	\N	/api/washingmachines	GET	2025-05-20 10:05:36.398787	2025-05-20 10:05:36.436363	200	37
7888	\N	/api/data/467bdeed-c69f-4074-bb6e-765100a569f5_85fe4e125190b7b4.jpeg	GET	2025-05-20 10:07:42.570093	2025-05-20 10:07:42.570224	200	0
7889	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:07:53.762505	2025-05-20 10:07:53.764571	200	2
7890	\N	/status	GET	2025-05-20 10:07:56.183593	2025-05-20 10:07:56.183598	200	0
7891	\N	/api/statistics/global	GET	2025-05-20 10:07:56.381386	2025-05-20 10:07:56.38372	200	2
7892	\N	/api/statistics/top-users	GET	2025-05-20 10:07:56.649037	2025-05-20 10:07:56.650947	200	1
7893	\N	/status	GET	2025-05-20 10:08:06.197609	2025-05-20 10:08:06.197611	200	0
7894	\N	/status	GET	2025-05-20 10:08:16.206431	2025-05-20 10:08:16.206433	200	0
7895	\N	/api/restaurant	GET	2025-05-20 10:08:17.145747	2025-05-20 10:08:17.14583	200	0
7896	\N	/status	GET	2025-05-20 10:16:51.736569	2025-05-20 10:16:51.736574	200	0
7897	\N	/api/restaurant	GET	2025-05-20 10:16:52.08815	2025-05-20 10:16:52.08823	200	0
7898	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 10:16:58.466599	2025-05-20 10:16:58.468269	200	1
7899	test@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-20 10:17:03.226384	2025-05-20 10:17:03.22822	200	1
7900	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 10:17:10.796291	2025-05-20 10:17:10.801963	200	5
7901	\N	/api/restaurant	GET	2025-05-20 10:17:13.721301	2025-05-20 10:17:13.72138	200	0
7903	\N	/api/restaurant	GET	2025-05-20 10:17:24.961427	2025-05-20 10:17:24.961534	200	0
7912	\N	/api/auth/register	POST	2025-05-20 11:01:15.721984	2025-05-20 11:01:15.807628	201	85
7924	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:01.556753	2025-05-20 11:02:01.558536	200	1
7925	\N	/api/weather	GET	2025-05-20 11:02:01.536807	2025-05-20 11:02:01.585168	200	48
7926	ninon.basle-blin@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 11:02:03.911594	2025-05-20 11:02:03.915977	200	4
7927	\N	/api/traq/	GET	2025-05-20 11:02:11.722049	2025-05-20 11:02:11.723676	200	1
7928	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:24.518755	2025-05-20 11:02:24.522144	200	3
7929	aurelien.pautet@imt-atlantique.net	/api/upload	POST	2025-05-20 11:02:24.535228	2025-05-20 11:02:24.541155	200	5
7930	aurelien.pautet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 11:02:24.600587	2025-05-20 11:02:24.605304	200	4
7931	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:24.731893	2025-05-20 11:02:24.733635	200	1
7932	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:02:30.486015	2025-05-20 11:02:30.487726	200	1
7933	aurelien.pautet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 11:03:13.838608	2025-05-20 11:03:13.843396	200	4
7934	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-20 11:03:20.304193	2025-05-20 11:03:20.30732	200	3
7953	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 14:24:53.506455	2025-05-20 14:24:53.506566	200	0
7954	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 15:30:21.728201	2025-05-20 15:30:21.728291	200	0
7955	\N	/api/restaurant	GET	2025-05-20 15:30:21.905871	2025-05-20 15:30:21.905976	200	0
7956	\N	/api/washingmachines	GET	2025-05-20 15:30:29.970949	2025-05-20 15:30:30.066409	200	95
7957	\N	/api/washingmachines	GET	2025-05-20 15:36:08.722228	2025-05-20 15:36:08.815648	200	93
7958	\N	/api/washingmachines	GET	2025-05-20 15:41:08.780579	2025-05-20 15:41:08.884051	200	103
7959	\N	/api/restaurant	GET	2025-05-20 16:24:23.231563	2025-05-20 16:24:23.231704	200	0
7960	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 16:24:23.231507	2025-05-20 16:24:23.231778	200	0
7961	\N	/api/restaurant	GET	2025-05-20 16:32:41.181171	2025-05-20 16:32:41.181292	200	0
7962	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 16:32:41.181174	2025-05-20 16:32:41.181276	200	0
7963	\N	/api/restaurant	GET	2025-05-20 16:34:17.781636	2025-05-20 16:34:17.781774	200	0
7964	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 16:34:17.781708	2025-05-20 16:34:17.781801	200	0
7965	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:37:22.613571	2025-05-20 16:37:22.615499	200	1
7966	\N	/api/washingmachines	GET	2025-05-20 16:37:22.921181	2025-05-20 16:37:23.01934	200	98
7967	\N	/api/restaurant	GET	2025-05-20 16:37:23.159656	2025-05-20 16:37:23.159755	200	0
7968	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:37:23.160161	2025-05-20 16:37:23.161573	200	1
7969	\N	/api/weather	GET	2025-05-20 16:37:23.118171	2025-05-20 16:37:23.175694	200	57
7970	\N	/api/washingmachines	GET	2025-05-20 16:37:32.669061	2025-05-20 16:37:32.697764	200	28
7971	\N	/api/washingmachines	GET	2025-05-20 16:37:38.3012	2025-05-20 16:37:38.330622	200	29
7972	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:39:46.927999	2025-05-20 16:39:46.930001	200	2
7973	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:39:48.130534	2025-05-20 16:39:48.132192	200	1
7974	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:40:14.762976	2025-05-20 16:40:14.764834	200	1
7975	\N	/api/washingmachines	GET	2025-05-20 16:40:15.047567	2025-05-20 16:40:15.13798	200	90
7976	\N	/api/weather	GET	2025-05-20 16:40:15.186065	2025-05-20 16:40:15.186087	200	0
7977	\N	/api/restaurant	GET	2025-05-20 16:40:15.218232	2025-05-20 16:40:15.218319	200	0
7978	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:40:15.224292	2025-05-20 16:40:15.226174	200	1
7979	\N	/api/washingmachines	GET	2025-05-20 16:40:22.827385	2025-05-20 16:40:22.856649	200	29
7980	\N	/api/washingmachines	GET	2025-05-20 16:40:24.966041	2025-05-20 16:40:24.993331	200	27
7981	\N	/api/washingmachines	GET	2025-05-20 16:40:26.434224	2025-05-20 16:40:26.46233	200	28
7982	\N	/api/washingmachines	GET	2025-05-20 16:40:27.662291	2025-05-20 16:40:27.691134	200	28
7983	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:41:06.015478	2025-05-20 16:41:06.017886	200	2
7984	\N	/api/washingmachines	GET	2025-05-20 16:41:06.128883	2025-05-20 16:41:06.175936	200	47
7985	\N	/api/restaurant	GET	2025-05-20 16:41:06.183698	2025-05-20 16:41:06.183878	200	0
7986	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 16:41:06.18377	2025-05-20 16:41:06.185559	200	1
7987	\N	/api/weather	GET	2025-05-20 16:41:06.183694	2025-05-20 16:41:06.249486	200	65
7988	\N	/api/washingmachines	GET	2025-05-20 16:41:08.404514	2025-05-20 16:41:08.432774	200	28
7989	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:43:29.216791	2025-05-20 17:43:29.218306	200	1
7990	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:43:29.379009	2025-05-20 17:43:29.380887	200	1
7991	\N	/api/weather	GET	2025-05-20 17:43:29.379074	2025-05-20 17:43:29.446146	200	67
7992	\N	/api/washingmachines	GET	2025-05-20 17:43:29.34342	2025-05-20 17:43:29.449677	200	106
7993	\N	/api/restaurant	GET	2025-05-20 17:43:29.379078	2025-05-20 17:43:30.485133	200	1106
7994	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:43:31.98326	2025-05-20 17:43:31.985674	200	2
7995	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:43:50.803067	2025-05-20 17:43:50.805264	200	2
7996	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:43:52.203409	2025-05-20 17:43:52.205871	200	2
7997	nathan.marie2@imt-atlantique.net	/api/upload	POST	2025-05-20 17:44:09.005446	2025-05-20 17:44:09.010817	200	5
7998	nathan.marie2@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 17:44:09.063967	2025-05-20 17:44:09.068859	200	4
7999	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:44:09.148228	2025-05-20 17:44:09.150378	200	2
8000	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-20 17:44:21.063114	2025-05-20 17:44:21.06558	200	2
8001	\N	/api/data/EDEB72E5-0F88-4926-8E9A-F2336C16ED38_53fa02e12cbd1f70.jpg	GET	2025-05-20 17:44:21.147876	2025-05-20 17:44:21.148084	200	0
8018	\N	/api/restaurant	GET	2025-05-20 19:50:41.791355	2025-05-20 19:50:41.791456	200	0
8031	\N	/api/restaurant	GET	2025-05-20 19:51:18.953383	2025-05-20 19:51:18.953513	200	0
8032	\N	/api/washingmachines	GET	2025-05-20 19:51:20.327989	2025-05-20 19:51:20.355832	200	27
8033	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:52:19.657203	2025-05-20 19:52:19.659306	200	2
8034	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:52:22.385878	2025-05-20 19:52:22.387902	200	2
8035	\N	/status	GET	2025-05-20 19:52:28.803804	2025-05-20 19:52:28.803811	200	0
8036	\N	/api/statistics/global	GET	2025-05-20 19:52:28.890385	2025-05-20 19:52:28.89297	200	2
8037	\N	/api/statistics/top-users	GET	2025-05-20 19:52:28.976845	2025-05-20 19:52:28.978591	200	1
8038	\N	/status	GET	2025-05-20 19:52:38.83786	2025-05-20 19:52:38.837862	200	0
8039	\N	/status	GET	2025-05-20 19:52:48.836733	2025-05-20 19:52:48.836734	200	0
8040	\N	/status	GET	2025-05-20 19:52:49.473892	2025-05-20 19:52:49.473894	200	0
8041	\N	/api/statistics/global	GET	2025-05-20 19:52:49.594216	2025-05-20 19:52:49.59648	200	2
8042	\N	/api/statistics/top-users	GET	2025-05-20 19:52:49.696981	2025-05-20 19:52:49.699063	200	2
8043	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:54:30.02713	2025-05-20 19:54:30.031732	200	4
8044	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:55:18.015382	2025-05-20 19:55:18.021691	200	6
8045	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:55:32.482082	2025-05-20 19:55:32.487215	200	5
8046	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:55:37.56685	2025-05-20 19:55:37.571627	200	4
8047	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:55:42.48714	2025-05-20 19:55:42.495059	200	7
8048	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-20 19:55:48.20959	2025-05-20 19:55:48.21444	200	4
8049	\N	/api/auth/register	POST	2025-05-20 19:58:25.45918	2025-05-20 19:58:25.532798	409	73
8050	\N	/api/auth/register	POST	2025-05-20 19:58:34.300311	2025-05-20 19:58:34.380076	201	79
8052	\N	/robots.txt	GET	2025-05-20 19:59:55.396791	2025-05-20 19:59:55.396796	500	0
8053	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 20:01:35.251338	2025-05-20 20:01:35.253782	200	2
8054	\N	/robots.txt	GET	2025-05-20 20:03:54.708108	2025-05-20 20:03:54.708111	500	0
8055	\N	/robots.txt	GET	2025-05-20 20:07:53.658155	2025-05-20 20:07:53.65816	500	0
8056	\N	/api/washingmachines	GET	2025-05-20 20:08:15.418809	2025-05-20 20:08:15.518966	200	100
8057	\N	/robots.txt	GET	2025-05-20 20:12:00.85233	2025-05-20 20:12:00.852333	500	0
8058	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 21:42:42.152921	2025-05-20 21:42:42.154686	200	1
8059	\N	/api/restaurant	GET	2025-05-20 21:42:42.42271	2025-05-20 21:42:42.422846	200	0
8063	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 03:19:56.624798	2025-05-21 03:19:56.624914	200	0
8064	nathan.marie2@imt-atlantique.net	/api/newf/me	GET	2025-05-21 06:49:59.482645	2025-05-21 06:49:59.484631	200	1
8065	\N	/api/data/EDEB72E5-0F88-4926-8E9A-F2336C16ED38_53fa02e12cbd1f70.jpg	GET	2025-05-21 06:49:59.560354	2025-05-21 06:49:59.560485	200	0
8066	\N	/api/restaurant	GET	2025-05-21 06:51:56.258739	2025-05-21 06:51:57.047185	200	788
8072	\N	/api/restaurant	GET	2025-05-21 07:27:55.329505	2025-05-21 07:27:56.372725	200	1043
8084	matheo.vallee@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 07:32:32.019745	2025-05-21 07:32:32.041872	200	22
8085	\N	/api/restaurant	GET	2025-05-21 07:32:35.643194	2025-05-21 07:32:35.643283	200	0
8086	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:32:47.308861	2025-05-21 07:32:47.310803	200	1
8087	\N	/api/traq/	GET	2025-05-21 07:33:04.962276	2025-05-21 07:33:04.963621	200	1
8088	\N	/api/restaurant	GET	2025-05-21 07:33:25.718157	2025-05-21 07:33:25.718244	200	0
8089	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:35:08.943675	2025-05-21 07:35:08.947568	200	3
8090	\N	/api/washingmachines	GET	2025-05-21 07:35:09.482153	2025-05-21 07:35:09.578479	200	96
8091	\N	/api/weather	GET	2025-05-21 07:35:09.642643	2025-05-21 07:35:09.642668	200	0
8092	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:35:09.697465	2025-05-21 07:35:09.699417	200	1
8093	\N	/api/restaurant	GET	2025-05-21 07:35:09.703449	2025-05-21 07:35:09.703555	200	0
8094	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:47:33.991448	2025-05-21 07:47:33.994241	200	2
8095	aurelien.moignet@imt-atlantique.net	/api/upload	POST	2025-05-21 07:47:44.949187	2025-05-21 07:47:44.95589	200	6
8096	aurelien.moignet@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 07:47:45.017505	2025-05-21 07:47:45.021567	200	4
8097	aurelien.moignet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:47:45.083653	2025-05-21 07:47:45.086097	200	2
8098	\N	/api/restaurant	GET	2025-05-21 07:51:47.451513	2025-05-21 07:51:47.451625	200	0
8099	\N	/api/washingmachines	GET	2025-05-21 07:51:48.869628	2025-05-21 07:51:48.96606	200	96
8100	\N	/api/restaurant	GET	2025-05-21 07:52:04.314082	2025-05-21 07:52:04.314186	200	0
8101	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:52:08.012212	2025-05-21 07:52:08.014449	200	2
8102	\N	/api/restaurant	GET	2025-05-21 07:52:28.411413	2025-05-21 07:52:28.411492	200	0
8103	\N	/api/washingmachines	GET	2025-05-21 07:54:28.159691	2025-05-21 07:54:28.233072	200	73
8104	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:00:05.584005	2025-05-21 08:00:05.585819	200	1
8105	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:00:15.296119	2025-05-21 08:00:15.297968	200	1
8106	marina.carbone@imt-atlantique.net	/api/upload	POST	2025-05-21 08:00:48.041597	2025-05-21 08:00:48.047763	200	6
8107	marina.carbone@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 08:00:48.208425	2025-05-21 08:00:48.213298	200	4
8108	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:00:48.637779	2025-05-21 08:00:48.639605	200	1
8114	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:01:08.352052	2025-05-21 08:01:08.354379	200	2
8115	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:01:10.208512	2025-05-21 08:01:10.210485	200	1
8116	\N	/api/restaurant	GET	2025-05-21 08:01:15.691413	2025-05-21 08:01:15.6915	200	0
8117	\N	/api/restaurant	GET	2025-05-21 08:01:25.40976	2025-05-21 08:01:25.409835	200	0
8118	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:01:30.73649	2025-05-21 08:01:30.739015	200	2
8119	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 08:01:33.126904	2025-05-21 08:01:33.131412	200	4
8120	\N	/api/restaurant	GET	2025-05-21 08:01:35.121996	2025-05-21 08:01:35.122138	200	0
8121	\N	/api/weather	GET	2025-05-21 08:01:35.122116	2025-05-21 08:01:35.138355	200	16
8122	\N	/api/washingmachines	GET	2025-05-21 08:01:35.114919	2025-05-21 08:01:35.144283	200	29
8124	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 08:14:28.4581	2025-05-21 08:14:28.458167	200	0
8126	\N	/api/restaurant	GET	2025-05-21 08:47:20.800732	2025-05-21 08:47:20.801023	200	0
8132	\N	/api/restaurant	GET	2025-05-21 08:53:16.345514	2025-05-21 08:53:16.345757	200	0
8135	\N	/api/washingmachines	GET	2025-05-21 08:53:16.308832	2025-05-21 08:53:16.397597	200	88
8136	\N	/api/restaurant	GET	2025-05-21 08:53:20.333639	2025-05-21 08:53:20.333771	200	0
8144	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:13:35.121095	2025-05-21 09:13:35.122666	200	1
8145	\N	/api/weather	GET	2025-05-21 09:13:35.115473	2025-05-21 09:13:35.17178	200	56
8146	maxime.bodin@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 09:13:37.23668	2025-05-21 09:13:37.39336	200	156
8147	\N	/api/restaurant	GET	2025-05-21 09:13:40.435287	2025-05-21 09:13:40.43536	200	0
8148	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:13:51.677378	2025-05-21 09:13:51.679374	200	1
8149	\N	/api/traq/	GET	2025-05-21 09:13:54.356862	2025-05-21 09:13:54.35849	200	1
8150	\N	/api/traq/	GET	2025-05-21 09:13:59.812485	2025-05-21 09:13:59.81411	200	1
8151	\N	/api/restaurant	GET	2025-05-21 09:14:01.957537	2025-05-21 09:14:01.957614	200	0
8152	\N	/api/washingmachines	GET	2025-05-21 09:14:04.705728	2025-05-21 09:14:04.734293	200	28
8153	\N	/api/auth/register	POST	2025-05-21 09:15:54.812441	2025-05-21 09:15:54.89493	201	82
8002	\N	/api/restaurant	GET	2025-05-20 17:59:02.410197	2025-05-20 17:59:02.410309	200	0
8003	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 17:59:02.410158	2025-05-20 17:59:02.410307	200	0
8004	\N	/api/washingmachines	GET	2025-05-20 17:59:05.849265	2025-05-20 17:59:05.949637	200	100
8005	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:24.537495	2025-05-20 19:48:24.539501	200	2
8006	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:24.753893	2025-05-20 19:48:24.756021	200	2
8007	\N	/api/restaurant	GET	2025-05-20 19:48:24.771395	2025-05-20 19:48:24.7715	200	0
8008	\N	/api/washingmachines	GET	2025-05-20 19:48:24.706988	2025-05-20 19:48:24.801304	200	94
8009	\N	/api/weather	GET	2025-05-20 19:48:24.764982	2025-05-20 19:48:24.825442	200	60
8010	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:25.573094	2025-05-20 19:48:25.575017	200	1
8011	test@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:27.665672	2025-05-20 19:48:27.667722	200	2
8012	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:45.682177	2025-05-20 19:48:45.68399	200	1
8013	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-20 19:48:45.773863	2025-05-20 19:48:45.774004	200	0
8014	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:46.964059	2025-05-20 19:48:46.966141	200	2
8015	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:48:54.970605	2025-05-20 19:48:54.972122	200	1
8016	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:50:41.295404	2025-05-20 19:50:41.297305	200	1
8017	\N	/api/washingmachines	GET	2025-05-20 19:50:41.599385	2025-05-20 19:50:41.712677	200	113
8019	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:50:41.791404	2025-05-20 19:50:41.793452	200	2
8020	\N	/api/weather	GET	2025-05-20 19:50:41.737905	2025-05-20 19:50:41.823835	200	85
8021	\N	/api/washingmachines	GET	2025-05-20 19:50:48.790446	2025-05-20 19:50:48.818042	200	27
8022	\N	/api/washingmachines	GET	2025-05-20 19:50:55.770624	2025-05-20 19:50:55.798984	200	28
8023	\N	/api/washingmachines	GET	2025-05-20 19:50:56.85741	2025-05-20 19:50:56.890989	200	33
8024	\N	/api/washingmachines	GET	2025-05-20 19:50:57.928619	2025-05-20 19:50:57.971884	200	43
8025	\N	/api/washingmachines	GET	2025-05-20 19:50:59.970486	2025-05-20 19:50:59.997804	200	27
8026	\N	/api/washingmachines	GET	2025-05-20 19:51:05.770436	2025-05-20 19:51:05.798706	200	28
8027	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:51:18.435266	2025-05-20 19:51:18.436866	200	1
8028	\N	/api/washingmachines	GET	2025-05-20 19:51:18.751548	2025-05-20 19:51:18.780631	200	29
8029	\N	/api/weather	GET	2025-05-20 19:51:18.928906	2025-05-20 19:51:18.92893	200	0
8030	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-20 19:51:18.946844	2025-05-20 19:51:18.94873	200	1
8051	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-20 19:59:53.928243	2025-05-20 19:59:53.928367	206	0
8060	ninon.basle-blin@imt-atlantique.net	/api/newf/me	GET	2025-05-20 21:42:42.425511	2025-05-20 21:42:42.427284	200	1
8061	\N	/api/weather	GET	2025-05-20 21:42:42.422686	2025-05-20 21:42:42.486154	200	63
8062	\N	/api/washingmachines	GET	2025-05-20 21:42:42.413135	2025-05-20 21:42:42.516633	200	103
8067	\N	/	GET	2025-05-21 07:14:14.099726	2025-05-21 07:14:14.099735	500	0
8068	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:27:54.635973	2025-05-21 07:27:54.63781	200	1
8069	\N	/api/washingmachines	GET	2025-05-21 07:27:54.991905	2025-05-21 07:27:55.084639	200	92
8070	\N	/api/weather	GET	2025-05-21 07:27:55.176167	2025-05-21 07:27:55.24954	200	73
8071	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:27:55.33585	2025-05-21 07:27:55.33778	200	1
8073	\N	/api/washingmachines	GET	2025-05-21 07:27:58.020136	2025-05-21 07:27:58.06614	200	46
8074	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:28:01.511295	2025-05-21 07:28:01.514086	200	2
8075	\N	/api/auth/register	POST	2025-05-21 07:32:01.177312	2025-05-21 07:32:01.258688	201	81
8076	\N	/api/auth/verify-account	POST	2025-05-21 07:32:28.349421	2025-05-21 07:32:28.37272	200	23
8077	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:32:28.574758	2025-05-21 07:32:28.576773	200	2
8078	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:32:28.728928	2025-05-21 07:32:28.730749	200	1
8079	\N	/api/washingmachines	GET	2025-05-21 07:32:28.932998	2025-05-21 07:32:29.027046	200	94
8080	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 07:32:29.337795	2025-05-21 07:32:29.341157	200	3
8081	\N	/api/auth/verify-account	POST	2025-05-21 07:32:29.377761	2025-05-21 07:32:29.383302	400	5
8082	\N	/api/weather	GET	2025-05-21 07:32:29.337791	2025-05-21 07:32:29.394505	200	56
8083	\N	/api/restaurant	GET	2025-05-21 07:32:29.298008	2025-05-21 07:32:29.414261	200	116
8109	\N	/api/data/31714441-fe62-48c8-bf57-36dd24469c1f_73c0a62bf1d69f33.jpeg	GET	2025-05-21 08:00:49.395043	2025-05-21 08:00:49.395136	200	0
8110	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:01:07.903534	2025-05-21 08:01:07.905618	200	2
8111	\N	/api/washingmachines	GET	2025-05-21 08:01:08.173221	2025-05-21 08:01:08.249801	200	76
8112	\N	/api/weather	GET	2025-05-21 08:01:08.273177	2025-05-21 08:01:08.336438	200	63
8113	\N	/api/restaurant	GET	2025-05-21 08:01:08.348535	2025-05-21 08:01:08.34864	200	0
8123	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 08:14:27.942962	2025-05-21 08:14:27.943061	200	0
8125	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:47:20.651119	2025-05-21 08:47:20.652989	200	1
8127	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:47:20.800749	2025-05-21 08:47:20.802622	200	1
8128	\N	/api/weather	GET	2025-05-21 08:47:20.800169	2025-05-21 08:47:20.868629	200	68
8129	\N	/api/washingmachines	GET	2025-05-21 08:47:20.769265	2025-05-21 08:47:20.874433	200	105
8130	\N	/api/restaurant	GET	2025-05-21 08:47:25.130558	2025-05-21 08:47:25.130693	200	0
8131	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:53:16.183216	2025-05-21 08:53:16.185341	200	2
8133	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 08:53:16.345507	2025-05-21 08:53:16.347168	200	1
8134	\N	/api/weather	GET	2025-05-21 08:53:16.345506	2025-05-21 08:53:16.395494	200	49
8137	\N	/api/restaurant	GET	2025-05-21 09:08:07.470339	2025-05-21 09:08:07.470461	200	0
8138	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 09:08:07.47028	2025-05-21 09:08:07.470409	200	0
8139	\N	/api/auth/login	POST	2025-05-21 09:13:34.518451	2025-05-21 09:13:34.591553	200	73
8140	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:13:34.679154	2025-05-21 09:13:34.680722	200	1
8141	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:13:34.74843	2025-05-21 09:13:34.750858	200	2
8142	\N	/api/washingmachines	GET	2025-05-21 09:13:34.859416	2025-05-21 09:13:34.948248	200	88
8143	\N	/api/restaurant	GET	2025-05-21 09:13:35.121068	2025-05-21 09:13:35.121177	200	0
8158	\N	/api/restaurant	GET	2025-05-21 09:16:41.891727	2025-05-21 09:16:41.892098	200	0
8170	\N	/api/data/8e13f62f-9789-481a-b93c-d4c2cef0ccec_db66b49d7de717b3.jpeg	GET	2025-05-21 09:19:35.520429	2025-05-21 09:19:35.52052	200	0
8171	pacome.cailleteau@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 09:19:40.740504	2025-05-21 09:19:40.746402	200	5
8172	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:19:45.303072	2025-05-21 09:19:45.30477	200	1
8173	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:19:57.113439	2025-05-21 09:19:57.115475	200	2
8175	\N	/api/washingmachines	GET	2025-05-21 09:19:57.371854	2025-05-21 09:19:57.447936	200	76
8177	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:19:57.478808	2025-05-21 09:19:57.480594	200	1
8178	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:20:02.492953	2025-05-21 09:20:02.49504	200	2
8179	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:20:17.834205	2025-05-21 09:20:17.836544	200	2
8180	pacome.cailleteau@imt-atlantique.net	/api/newf/notifications/subscriptions	GET	2025-05-21 09:20:21.004747	2025-05-21 09:20:21.006726	200	1
8181	pacome.cailleteau@imt-atlantique.net	/api/newf/notifications/subscriptions	POST	2025-05-21 09:20:28.411699	2025-05-21 09:20:28.416824	200	5
8182	\N	/api/support	GET	2025-05-21 09:21:36.696209	2025-05-21 09:21:36.696218	500	0
8183	\N	/api/support	GET	2025-05-21 09:21:37.816875	2025-05-21 09:21:37.816885	500	0
8184	\N	/status	GET	2025-05-21 09:21:42.771542	2025-05-21 09:21:42.771544	200	0
8185	\N	/api/statistics/global	GET	2025-05-21 09:21:42.842063	2025-05-21 09:21:42.844764	200	2
8186	\N	/api/statistics/top-users	GET	2025-05-21 09:21:42.930385	2025-05-21 09:21:42.932465	200	2
8154	\N	/api/auth/verify-account	POST	2025-05-21 09:16:40.805893	2025-05-21 09:16:40.813995	200	8
8155	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:16:40.985461	2025-05-21 09:16:40.987796	200	2
8156	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:16:41.076507	2025-05-21 09:16:41.078539	200	2
8157	\N	/api/washingmachines	GET	2025-05-21 09:16:41.258938	2025-05-21 09:16:41.337849	200	78
8159	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:16:41.891751	2025-05-21 09:16:41.893781	200	2
8160	\N	/api/weather	GET	2025-05-21 09:16:41.851309	2025-05-21 09:16:41.898531	200	47
8161	pacome.cailleteau@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 09:16:43.69474	2025-05-21 09:16:43.698862	200	4
8162	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:16:56.789126	2025-05-21 09:16:56.79128	200	2
8163	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:17:03.116919	2025-05-21 09:17:03.118764	200	1
8164	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:17:06.08902	2025-05-21 09:17:06.090971	200	1
8165	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:17:07.450945	2025-05-21 09:17:07.453171	200	2
8166	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:17:08.014137	2025-05-21 09:17:08.01679	200	2
8167	pacome.cailleteau@imt-atlantique.net	/api/upload	POST	2025-05-21 09:19:34.391284	2025-05-21 09:19:34.396462	200	5
8168	pacome.cailleteau@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 09:19:34.477152	2025-05-21 09:19:34.482704	200	5
8169	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:19:34.797082	2025-05-21 09:19:34.799299	200	2
8174	\N	/api/weather	GET	2025-05-21 09:19:57.448297	2025-05-21 09:19:57.448315	200	0
8176	\N	/api/restaurant	GET	2025-05-21 09:19:57.478825	2025-05-21 09:19:57.478893	200	0
8195	\N	/api/restaurant	GET	2025-05-21 09:22:12.298457	2025-05-21 09:22:12.298562	200	0
8221	\N	/api/restaurant	GET	2025-05-21 09:44:03.991984	2025-05-21 09:44:03.992047	200	0
8222	\N	/status	GET	2025-05-21 09:48:02.537204	2025-05-21 09:48:02.537208	200	0
8231	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 10:49:55.099206	2025-05-21 10:49:55.108206	200	8
8232	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 10:49:57.755416	2025-05-21 10:49:57.758421	200	3
8238	\N	/api/data/8e13f62f-9789-481a-b93c-d4c2cef0ccec_db66b49d7de717b3.jpeg	GET	2025-05-21 10:51:20.9387	2025-05-21 10:51:20.938788	200	0
8187	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:21:57.173582	2025-05-21 09:21:57.175394	200	1
8188	\N	/status	GET	2025-05-21 09:21:59.054628	2025-05-21 09:21:59.054629	200	0
8189	\N	/api/statistics/global	GET	2025-05-21 09:21:59.138803	2025-05-21 09:21:59.141198	200	2
8190	\N	/api/statistics/top-users	GET	2025-05-21 09:21:59.204677	2025-05-21 09:21:59.20628	200	1
8191	\N	/status	GET	2025-05-21 09:22:09.065871	2025-05-21 09:22:09.065873	200	0
8192	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:22:11.875523	2025-05-21 09:22:11.877302	200	1
8193	\N	/api/washingmachines	GET	2025-05-21 09:22:12.097383	2025-05-21 09:22:12.196351	200	98
8194	\N	/api/weather	GET	2025-05-21 09:22:12.236844	2025-05-21 09:22:12.287327	200	50
8196	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:22:12.298452	2025-05-21 09:22:12.30023	200	1
8197	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:22:15.31621	2025-05-21 09:22:15.318027	200	1
8198	\N	/status	GET	2025-05-21 09:22:19.081318	2025-05-21 09:22:19.081332	200	0
8199	\N	/status	GET	2025-05-21 09:22:29.078275	2025-05-21 09:22:29.078277	200	0
8200	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:22:29.205526	2025-05-21 09:22:29.207379	200	1
8201	\N	/status	GET	2025-05-21 09:22:38.964764	2025-05-21 09:22:38.964766	200	0
8202	\N	/api/statistics/global	GET	2025-05-21 09:22:39.049452	2025-05-21 09:22:39.052159	200	2
8203	\N	/status	GET	2025-05-21 09:22:39.095997	2025-05-21 09:22:39.095998	200	0
8204	\N	/api/statistics/top-users	GET	2025-05-21 09:22:39.142278	2025-05-21 09:22:39.143971	200	1
8205	\N	/api/washingmachines	GET	2025-05-21 09:22:44.262121	2025-05-21 09:22:44.293514	200	31
8206	\N	/status	GET	2025-05-21 09:22:49.087408	2025-05-21 09:22:49.08741	200	0
8207	\N	/api/restaurant	GET	2025-05-21 09:22:56.802068	2025-05-21 09:22:56.802136	200	0
8208	\N	/status	GET	2025-05-21 09:22:59.089771	2025-05-21 09:22:59.089773	200	0
8209	\N	/api/washingmachines	GET	2025-05-21 09:23:02.994797	2025-05-21 09:23:03.026205	200	31
8210	\N	/api/restaurant	GET	2025-05-21 09:23:05.177164	2025-05-21 09:23:05.177256	200	0
8211	\N	/api/traq/	GET	2025-05-21 09:23:07.903713	2025-05-21 09:23:07.905494	200	1
8212	\N	/status	GET	2025-05-21 09:23:09.071669	2025-05-21 09:23:09.071671	200	0
8213	\N	/status	GET	2025-05-21 09:23:19.091146	2025-05-21 09:23:19.091147	200	0
8214	\N	/api/newf/me	GET	2025-05-21 09:30:37.661213	2025-05-21 09:30:37.66134	401	0
8215	\N	/status	GET	2025-05-21 09:32:18.610157	2025-05-21 09:32:18.610165	200	0
8216	\N	/api/newf/me	GET	2025-05-21 09:32:22.834653	2025-05-21 09:32:22.834749	401	0
8217	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:44:03.498678	2025-05-21 09:44:03.518372	200	19
8218	\N	/api/washingmachines	GET	2025-05-21 09:44:03.845613	2025-05-21 09:44:03.918534	200	72
8219	\N	/api/weather	GET	2025-05-21 09:44:03.925132	2025-05-21 09:44:03.976376	200	51
8220	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 09:44:03.989238	2025-05-21 09:44:03.991623	200	2
8223	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 10:43:03.726668	2025-05-21 10:43:03.726776	200	0
8224	\N	/api/auth/login	POST	2025-05-21 10:49:51.504921	2025-05-21 10:49:51.574909	200	69
8225	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 10:49:51.81614	2025-05-21 10:49:51.818865	200	2
8226	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 10:49:51.948161	2025-05-21 10:49:51.949759	200	1
8227	\N	/api/washingmachines	GET	2025-05-21 10:49:52.167984	2025-05-21 10:49:52.266387	200	98
8228	\N	/api/weather	GET	2025-05-21 10:49:53.062966	2025-05-21 10:49:53.123948	200	60
8229	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 10:49:53.145212	2025-05-21 10:49:53.147309	200	2
8230	\N	/api/restaurant	GET	2025-05-21 10:49:53.15062	2025-05-21 10:49:53.876704	200	726
8233	\N	/api/data/572b22e4-10a0-4a25-8c35-5483a9a64782.jpeg	GET	2025-05-21 10:49:58.297631	2025-05-21 10:49:58.297752	200	0
8234	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 10:50:00.099179	2025-05-21 10:50:00.101249	200	2
8235	yohann.chavanel@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 10:50:26.834714	2025-05-21 10:50:26.839612	200	4
8236	\N	/api/washingmachines	GET	2025-05-21 10:51:00.519591	2025-05-21 10:51:00.547233	200	27
8237	\N	/api/restaurant	GET	2025-05-21 10:51:03.818342	2025-05-21 10:51:03.818412	200	0
8239	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:11:13.712269	2025-05-21 13:11:13.800273	200	88
8240	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:11:15.104539	2025-05-21 13:11:15.191217	200	86
8241	\N	/health	GET	2025-05-21 11:13:14.897278	2025-05-21 11:13:14.897281	200	0
8242	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:33.581423	2025-05-21 11:13:33.58558	200	4
8243	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:33.877082	2025-05-21 11:13:33.879283	200	2
8244	\N	/api/washingmachines	GET	2025-05-21 11:13:33.817339	2025-05-21 11:13:33.915619	200	98
8245	\N	/api/weather	GET	2025-05-21 11:13:33.877075	2025-05-21 11:13:33.932256	200	55
8246	\N	/api/restaurant	GET	2025-05-21 11:13:33.877071	2025-05-21 11:13:33.96817	200	91
8247	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:35.32939	2025-05-21 11:13:35.332964	200	3
8248	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:38.610384	2025-05-21 11:13:38.612776	200	2
8249	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:42.695748	2025-05-21 11:13:42.698041	200	2
8250	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:45.172815	2025-05-21 11:13:45.174934	200	2
8251	\N	/api/washingmachines	GET	2025-05-21 11:13:45.317989	2025-05-21 11:13:45.357365	200	39
8252	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:45.374371	2025-05-21 11:13:45.377188	200	2
8253	\N	/api/weather	GET	2025-05-21 11:13:45.371479	2025-05-21 11:13:45.385499	200	14
8254	matis.byar@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:45.936743	2025-05-21 11:13:45.93895	200	2
8255	\N	/api/restaurant	GET	2025-05-21 11:13:45.371685	2025-05-21 11:13:46.09999	200	728
8256	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:13:50.384677	2025-05-21 11:13:50.389758	200	5
8257	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:14:49.96794	2025-05-21 11:14:49.969838	200	1
8258	\N	/api/washingmachines	GET	2025-05-21 11:14:50.086377	2025-05-21 11:14:50.122279	200	35
8259	\N	/api/weather	GET	2025-05-21 11:14:50.13252	2025-05-21 11:14:50.132538	200	0
8260	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:14:50.134877	2025-05-21 11:14:50.137635	200	2
8261	\N	/api/restaurant	GET	2025-05-21 11:14:50.134791	2025-05-21 11:14:50.134882	200	0
8262	yohann.chavanel@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:14:51.98413	2025-05-21 11:14:51.986575	200	2
8263	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:11.499152	2025-05-21 11:15:11.505087	200	5
8264	\N	/api/washingmachines	GET	2025-05-21 11:15:11.703256	2025-05-21 11:15:11.740351	200	37
8265	\N	/api/weather	GET	2025-05-21 11:15:11.823713	2025-05-21 11:15:11.823751	200	0
8266	\N	/api/restaurant	GET	2025-05-21 11:15:11.85876	2025-05-21 11:15:11.858833	200	0
8267	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:11.858749	2025-05-21 11:15:11.860829	200	2
8268	\N	/api/traq/	GET	2025-05-21 11:15:29.848512	2025-05-21 11:15:29.850253	200	1
8269	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:33.011031	2025-05-21 11:15:33.013403	200	2
8270	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:51.243116	2025-05-21 11:15:51.245138	200	2
8271	\N	/api/washingmachines	GET	2025-05-21 11:15:51.503728	2025-05-21 11:15:51.543659	200	39
8272	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:51.699608	2025-05-21 11:15:51.701218	200	1
8273	\N	/api/weather	GET	2025-05-21 11:15:51.663561	2025-05-21 11:15:51.712658	200	49
8274	\N	/api/restaurant	GET	2025-05-21 11:15:51.699671	2025-05-21 11:15:52.796376	200	1096
8275	\N	/api/washingmachines	GET	2025-05-21 11:15:54.784159	2025-05-21 11:15:54.81394	200	29
8276	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:15:58.368185	2025-05-21 11:15:58.369872	200	1
8277	\N	/api/restaurant	GET	2025-05-21 11:32:14.137369	2025-05-21 11:32:14.137532	200	0
8278	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:34:25.332048	2025-05-21 11:34:25.334182	200	2
8279	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-21 11:34:25.463252	2025-05-21 11:34:25.464483	200	1
8280	enzo.morvan@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:34:29.427091	2025-05-21 11:34:29.429493	200	2
8281	\N	/api/data/8DA6336C-6D0D-472C-B5AB-25BE4078CBBF_dea9a9d0059fd1ab.jpg	GET	2025-05-21 11:34:29.52076	2025-05-21 11:34:29.520872	200	0
8282	\N	/api/test/notifications/send-to-group	GET	2025-05-21 13:42:01.95432	2025-05-21 13:42:01.954342	500	0
8283	\N	/api/test/notifications/send-to-group	GET	2025-05-21 11:42:30.646426	2025-05-21 11:42:30.646438	500	0
8284	\N	/api/test/notifications/send-to-group	GET	2025-05-21 13:42:35.073932	2025-05-21 13:42:35.073992	401	0
8285	yohann.chavanel@imt-atlantique.net	/api/test/notifications/send-to-group	GET	2025-05-21 13:42:54.276003	2025-05-21 13:42:54.27616	500	0
8286	yohann.chavanel@imt-atlantique.net	/api/test/notifications/send-to-group	POST	2025-05-21 13:43:25.582876	2025-05-21 13:43:26.945345	200	1362
8287	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:49:52.666614	2025-05-21 11:49:52.668229	200	1
8288	\N	/api/washingmachines	GET	2025-05-21 11:49:52.884228	2025-05-21 11:49:52.98561	200	101
8289	\N	/api/restaurant	GET	2025-05-21 11:49:53.156711	2025-05-21 11:49:53.156904	200	0
8290	\N	/api/weather	GET	2025-05-21 11:49:53.12013	2025-05-21 11:49:53.169487	200	49
8292	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 11:49:53.187181	2025-05-21 11:49:53.188696	200	1
8293	yohann.chavanel@imt-atlantique.net	/api/newf/notification	GET	2025-05-21 14:00:17.292367	2025-05-21 14:00:17.292917	500	0
8294	yohann.chavanel@imt-atlantique.net	/api/newf/notification	GET	2025-05-21 14:11:26.242804	2025-05-21 14:11:26.243036	500	0
8291	\N	/api/restaurant	GET	2025-05-21 11:49:53.187557	2025-05-21 11:49:53.187655	200	0
8295	\N	/api/notifications/send-to-group	POST	2025-05-21 14:12:04.330866	2025-05-21 14:12:04.330884	500	0
8296	yohann.chavanel@imt-atlantique.net	/api/notifications/send-to-group	POST	2025-05-21 14:13:40.272133	2025-05-21 14:13:41.723762	200	1451
8297	\N	/health	GET	2025-05-21 12:15:13.024903	2025-05-21 12:15:13.024908	200	0
8298	\N	/api/auth/login	POST	2025-05-21 12:16:40.510817	2025-05-21 12:16:40.586276	200	75
8299	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:40.745464	2025-05-21 12:16:40.747426	200	1
8300	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:40.885085	2025-05-21 12:16:40.887464	200	2
8301	\N	/api/washingmachines	GET	2025-05-21 12:16:41.134343	2025-05-21 12:16:41.23459	200	100
8302	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:41.980897	2025-05-21 12:16:41.983014	200	2
8303	\N	/api/weather	GET	2025-05-21 12:16:41.937036	2025-05-21 12:16:42.022573	200	85
8304	\N	/api/restaurant	GET	2025-05-21 12:16:41.981309	2025-05-21 12:16:42.164993	200	183
8305	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:48.345175	2025-05-21 12:16:48.347988	200	2
8306	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:53.583365	2025-05-21 12:16:53.586003	200	2
8307	marina.carbone@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:16:58.894599	2025-05-21 12:16:58.897756	200	3
8308	\N	/status	GET	2025-05-21 12:17:06.15022	2025-05-21 12:17:06.150223	200	0
8309	\N	/api/statistics/global	GET	2025-05-21 12:17:06.220925	2025-05-21 12:17:06.223847	200	2
8310	\N	/api/statistics/top-users	GET	2025-05-21 12:17:06.301332	2025-05-21 12:17:06.303348	200	2
8311	\N	/status	GET	2025-05-21 12:17:16.162672	2025-05-21 12:17:16.162674	200	0
8312	\N	/status	GET	2025-05-21 12:17:26.218848	2025-05-21 12:17:26.21885	200	0
8313	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:17:30.800851	2025-05-21 12:17:30.802592	200	1
8314	\N	/api/washingmachines	GET	2025-05-21 12:17:30.992271	2025-05-21 12:17:31.024786	200	32
8315	\N	/api/weather	GET	2025-05-21 12:17:31.164588	2025-05-21 12:17:31.164613	200	0
8316	\N	/api/restaurant	GET	2025-05-21 12:17:31.197927	2025-05-21 12:17:31.198021	200	0
8317	\N	/api/restaurant	GET	2025-05-21 12:17:31.203843	2025-05-21 12:17:31.20395	200	0
8318	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 12:17:31.203864	2025-05-21 12:17:31.205183	200	1
8319	\N	/status	GET	2025-05-21 12:17:36.20563	2025-05-21 12:17:36.205632	200	0
8320	\N	/status	GET	2025-05-21 12:17:46.187748	2025-05-21 12:17:46.187749	200	0
8321	\N	/status	GET	2025-05-21 12:17:56.199569	2025-05-21 12:17:56.19957	200	0
8322	\N	/status	GET	2025-05-21 12:18:06.208354	2025-05-21 12:18:06.208356	200	0
8323	\N	/status	GET	2025-05-21 12:18:16.205928	2025-05-21 12:18:16.205933	200	0
8324	\N	/status	GET	2025-05-21 12:18:26.206021	2025-05-21 12:18:26.206022	200	0
8325	\N	/status	GET	2025-05-21 12:18:36.204768	2025-05-21 12:18:36.204769	200	0
8326	\N	/status	GET	2025-05-21 12:18:46.235045	2025-05-21 12:18:46.235047	200	0
8327	\N	/status	GET	2025-05-21 12:18:56.222663	2025-05-21 12:18:56.222665	200	0
8328	\N	/status	GET	2025-05-21 12:19:06.214959	2025-05-21 12:19:06.21496	200	0
8329	\N	/status	GET	2025-05-21 12:19:16.230302	2025-05-21 12:19:16.230303	200	0
8330	\N	/api/restaurant	GET	2025-05-21 12:21:16.979815	2025-05-21 12:21:16.979997	200	0
8331	\N	/status	GET	2025-05-21 12:21:29.204529	2025-05-21 12:21:29.204531	200	0
8332	\N	/status	GET	2025-05-21 12:21:32.777726	2025-05-21 12:21:32.777728	200	0
8333	\N	/api/statistics/global	GET	2025-05-21 12:21:32.87029	2025-05-21 12:21:32.872763	200	2
8334	\N	/api/statistics/top-users	GET	2025-05-21 12:21:32.944243	2025-05-21 12:21:32.946021	200	1
8335	\N	/status	GET	2025-05-21 12:21:58.218198	2025-05-21 12:21:58.2182	200	0
8336	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 12:58:59.453995	2025-05-21 12:58:59.454485	200	0
8337	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:04:18.53567	2025-05-21 13:04:18.537647	200	1
8338	\N	/api/washingmachines	GET	2025-05-21 13:04:18.855935	2025-05-21 13:04:18.951167	200	95
8339	\N	/api/weather	GET	2025-05-21 13:04:18.96548	2025-05-21 13:04:19.038603	200	73
8340	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:04:19.090959	2025-05-21 13:04:19.093212	200	2
8341	\N	/api/restaurant	GET	2025-05-21 13:04:19.086031	2025-05-21 13:04:19.974752	200	888
8342	\N	/api/washingmachines	GET	2025-05-21 13:04:21.106076	2025-05-21 13:04:21.142931	200	36
8343	\N	/api/washingmachines	GET	2025-05-21 13:04:29.009614	2025-05-21 13:04:29.039852	200	30
8344	\N	/api/restaurant	GET	2025-05-21 13:04:31.480677	2025-05-21 13:04:31.48076	200	0
8345	\N	/api/restaurant	GET	2025-05-21 13:04:40.219882	2025-05-21 13:04:40.220008	200	0
8346	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:06:11.355858	2025-05-21 13:06:11.357734	200	1
8347	\N	/api/washingmachines	GET	2025-05-21 13:06:11.535124	2025-05-21 13:06:11.634762	200	99
8348	\N	/api/restaurant	GET	2025-05-21 13:06:11.69381	2025-05-21 13:06:11.693907	200	0
8349	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:06:11.728257	2025-05-21 13:06:11.730407	200	2
8350	\N	/api/restaurant	GET	2025-05-21 13:06:11.736346	2025-05-21 13:06:11.736426	200	0
8351	\N	/api/weather	GET	2025-05-21 13:06:11.728181	2025-05-21 13:06:11.797697	200	69
8352	matheo.vallee@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:06:20.936805	2025-05-21 13:06:20.939763	200	2
8353	\N	/api/washingmachines	GET	2025-05-21 13:06:25.618033	2025-05-21 13:06:25.646123	200	28
8354	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 13:24:05.08485	2025-05-21 13:24:05.084958	200	0
8355	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:35.153463	2025-05-21 13:25:35.155249	200	1
8356	\N	/api/restaurant	GET	2025-05-21 13:25:35.341069	2025-05-21 13:25:35.341405	200	0
8357	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:35.340404	2025-05-21 13:25:35.34315	200	2
8358	\N	/api/washingmachines	GET	2025-05-21 13:25:35.311659	2025-05-21 13:25:35.408015	200	96
8359	\N	/api/weather	GET	2025-05-21 13:25:35.341221	2025-05-21 13:25:35.422562	200	81
8360	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:37.188117	2025-05-21 13:25:37.190064	200	1
8361	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:39.166123	2025-05-21 13:25:39.167982	200	1
8362	\N	/api/auth/login	POST	2025-05-21 13:25:56.372159	2025-05-21 13:25:56.44449	200	72
8363	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:56.629405	2025-05-21 13:25:56.631112	200	1
8364	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:56.704829	2025-05-21 13:25:56.706605	200	1
8365	\N	/api/washingmachines	GET	2025-05-21 13:25:56.793279	2025-05-21 13:25:56.822382	200	29
8366	\N	/api/weather	GET	2025-05-21 13:25:57.062918	2025-05-21 13:25:57.062941	200	0
8367	\N	/api/restaurant	GET	2025-05-21 13:25:57.076528	2025-05-21 13:25:57.076636	200	0
8368	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 13:25:57.076568	2025-05-21 13:25:57.080529	200	3
8369	\N	/api/auth/verification-code	GET	2025-05-21 13:26:32.461802	2025-05-21 13:26:32.461829	500	0
8370	\N	/favicon.ico	GET	2025-05-21 13:26:32.595466	2025-05-21 13:26:32.59547	500	0
8371	\N	/	GET	2025-05-21 14:00:29.308533	2025-05-21 14:00:29.308537	500	0
8372	\N	/api/auth/login	POST	2025-05-21 14:01:40.151885	2025-05-21 14:01:40.224787	200	72
8373	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 14:01:40.335077	2025-05-21 14:01:40.336721	200	1
8374	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 14:01:40.412792	2025-05-21 14:01:40.414323	200	1
8375	\N	/api/restaurant	GET	2025-05-21 14:01:40.64165	2025-05-21 14:01:40.641839	200	0
8376	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 14:01:40.641672	2025-05-21 14:01:40.643551	200	1
8377	\N	/api/washingmachines	GET	2025-05-21 14:01:40.5683	2025-05-21 14:01:40.656323	200	88
8378	\N	/api/weather	GET	2025-05-21 14:01:40.641655	2025-05-21 14:01:40.702612	200	60
8379	maxime.bodin@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 14:01:42.516761	2025-05-21 14:01:42.521417	200	4
8380	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 14:05:33.775363	2025-05-21 14:05:33.77546	200	0
8381	\N	/api/restaurant	GET	2025-05-21 14:05:33.833898	2025-05-21 14:05:33.834028	200	0
8382	\N	/api/weather	GET	2025-05-21 14:05:51.304802	2025-05-21 14:05:51.304832	200	0
8383	\N	/api/restaurant	GET	2025-05-21 14:05:51.304884	2025-05-21 14:05:51.304988	200	0
8384	\N	/api/washingmachines	GET	2025-05-21 14:05:51.304742	2025-05-21 14:05:51.357028	200	52
8385	\N	/api/restaurant	GET	2025-05-21 14:05:54.371115	2025-05-21 14:05:54.371189	200	0
8386	\N	/api/weather	GET	2025-05-21 14:05:54.382061	2025-05-21 14:05:54.382091	200	0
8387	\N	/api/washingmachines	GET	2025-05-21 14:05:54.358233	2025-05-21 14:05:54.387387	200	29
8388	\N	/api/restaurant	GET	2025-05-21 14:41:34.782689	2025-05-21 14:41:34.782893	200	0
8389	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:15:01.862885	2025-05-21 15:15:01.864909	200	2
8390	\N	/api/washingmachines	GET	2025-05-21 15:15:02.114124	2025-05-21 15:15:02.203313	200	89
8392	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:15:02.294805	2025-05-21 15:15:02.297412	200	2
8411	\N	/api/restaurant	GET	2025-05-21 15:24:37.884179	2025-05-21 15:24:37.88427	200	0
8420	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:37:57.770449	2025-05-21 16:37:57.772376	200	1
8421	\N	/api/washingmachines	GET	2025-05-21 16:37:57.724614	2025-05-21 16:37:57.820654	200	96
8422	\N	/api/weather	GET	2025-05-21 16:37:57.770453	2025-05-21 16:37:57.838639	200	68
8423	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:38:00.54452	2025-05-21 16:38:00.54641	200	1
8424	\N	/api/data/2F773180-2F5B-41EF-8DA2-1109D921BBE0_4d1b0b7d141e7f9e.jpg	GET	2025-05-21 16:38:00.64347	2025-05-21 16:38:00.643675	200	0
8425	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:38:05.94573	2025-05-21 16:38:05.947584	200	1
8459	\N	/status	GET	2025-05-21 16:48:21.798552	2025-05-21 16:48:21.798556	200	0
8460	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 16:49:47.588721	2025-05-21 16:49:47.588818	200	0
8461	\N	/api/restaurant	GET	2025-05-21 16:49:47.874128	2025-05-21 16:49:47.874205	200	0
8462	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:57:09.332832	2025-05-21 16:57:09.335043	200	2
8463	\N	/api/washingmachines	GET	2025-05-21 16:57:09.523127	2025-05-21 16:57:09.596666	200	73
8464	\N	/api/restaurant	GET	2025-05-21 16:57:09.60778	2025-05-21 16:57:09.607851	200	0
8465	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:57:09.613212	2025-05-21 16:57:09.614831	200	1
8466	\N	/api/weather	GET	2025-05-21 16:57:09.61318	2025-05-21 16:57:09.683545	200	70
8467	\N	/api/restaurant	GET	2025-05-21 17:03:03.488945	2025-05-21 17:03:03.489069	200	0
8468	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:12:01.454944	2025-05-21 17:12:01.46164	200	6
8469	\N	/api/washingmachines	GET	2025-05-21 17:12:01.684419	2025-05-21 17:12:01.759804	200	75
8470	\N	/api/restaurant	GET	2025-05-21 17:12:01.807592	2025-05-21 17:12:01.807684	200	0
8471	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:12:01.813622	2025-05-21 17:12:01.815469	200	1
8472	\N	/api/weather	GET	2025-05-21 17:12:01.758009	2025-05-21 17:12:01.839389	200	81
8473	\N	/api/restaurant	GET	2025-05-21 17:13:07.247362	2025-05-21 17:13:07.247484	200	0
8474	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:14:04.584211	2025-05-21 17:14:04.586318	200	2
8475	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:14:18.721281	2025-05-21 17:14:18.723343	200	2
8476	pacome.cailleteau@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:14:19.400506	2025-05-21 17:14:19.402526	200	2
8477	\N	/api/restaurant	GET	2025-05-21 17:14:27.781474	2025-05-21 17:14:27.781559	200	0
8482	\N	/api/restaurant	GET	2025-05-21 17:18:02.663416	2025-05-21 17:18:02.663529	200	0
8491	\N	/api/restaurant	GET	2025-05-21 17:18:49.978887	2025-05-21 17:18:49.979053	200	0
8495	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 17:19:12.25495	2025-05-21 17:19:12.255057	200	0
8496	\N	/api/restaurant	GET	2025-05-21 17:19:12.5218	2025-05-21 17:19:12.521892	200	0
8497	\N	/api/washingmachines	GET	2025-05-21 17:26:08.933862	2025-05-21 17:26:09.029128	200	95
8498	\N	/api/auth/login	POST	2025-05-21 18:39:10.362064	2025-05-21 18:39:10.436427	200	74
8499	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 18:39:10.823846	2025-05-21 18:39:10.825701	200	1
8500	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 18:39:11.187044	2025-05-21 18:39:11.188691	200	1
8501	\N	/api/washingmachines	GET	2025-05-21 18:39:11.780791	2025-05-21 18:39:11.873471	200	92
8502	\N	/api/weather	GET	2025-05-21 18:39:12.068977	2025-05-21 18:39:12.135107	200	66
8503	\N	/api/restaurant	GET	2025-05-21 18:39:12.241508	2025-05-21 18:39:12.241604	200	0
8508	\N	/api/restaurant	GET	2025-05-21 18:40:38.731605	2025-05-21 18:40:38.733375	200	1
8391	\N	/api/restaurant	GET	2025-05-21 15:15:02.294971	2025-05-21 15:15:02.295082	200	0
8393	\N	/api/weather	GET	2025-05-21 15:15:02.250233	2025-05-21 15:15:02.31826	200	68
8394	\N	/api/washingmachines	GET	2025-05-21 15:15:03.532169	2025-05-21 15:15:03.56101	200	28
8395	\N	/api/restaurant	GET	2025-05-21 15:15:08.514011	2025-05-21 15:15:08.514084	200	0
8396	\N	/api/restaurant	GET	2025-05-21 15:15:18.456951	2025-05-21 15:15:18.45703	200	0
8397	\N	/api/restaurant	GET	2025-05-21 15:15:20.914471	2025-05-21 15:15:20.91455	200	0
8398	\N	/api/washingmachines	GET	2025-05-21 15:15:26.314013	2025-05-21 15:15:26.342659	200	28
8399	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:15:29.72511	2025-05-21 15:15:29.727675	200	2
8400	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:15:32.197449	2025-05-21 15:15:32.201116	200	3
8401	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 15:16:54.58767	2025-05-21 15:16:54.592257	200	4
8402	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 15:17:15.803555	2025-05-21 15:17:15.808205	200	4
8403	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 15:17:20.18823	2025-05-21 15:17:20.192341	200	4
8404	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 15:17:23.200882	2025-05-21 15:17:23.205754	200	4
8405	nathaniel.guitton@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 15:17:30.757612	2025-05-21 15:17:30.774851	200	17
8406	\N	/api/traq/	GET	2025-05-21 15:18:34.060981	2025-05-21 15:18:34.062688	200	1
8407	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:19:44.17104	2025-05-21 15:19:44.173331	200	2
8408	\N	/api/restaurant	GET	2025-05-21 15:19:49.214064	2025-05-21 15:19:49.214138	200	0
8409	\N	/api/washingmachines	GET	2025-05-21 15:19:53.994314	2025-05-21 15:19:54.097498	200	103
8410	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:24:37.597882	2025-05-21 15:24:37.600517	200	2
8412	maxime.bodin@imt-atlantique.net	/api/newf/me	GET	2025-05-21 15:24:37.883911	2025-05-21 15:24:37.886035	200	2
8413	\N	/api/washingmachines	GET	2025-05-21 15:24:37.818644	2025-05-21 15:24:37.906386	200	87
8414	\N	/api/weather	GET	2025-05-21 15:24:37.884273	2025-05-21 15:24:37.944175	200	59
8415	\N	/api/restaurant	GET	2025-05-21 15:27:15.209795	2025-05-21 15:27:15.209916	200	0
8416	\N	/api/weather	GET	2025-05-21 15:27:15.209839	2025-05-21 15:27:15.267142	200	57
8417	\N	/api/washingmachines	GET	2025-05-21 15:27:15.209785	2025-05-21 15:27:15.326783	200	116
8418	aurelien.pautet@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:37:57.570137	2025-05-21 16:37:57.571849	200	1
8419	\N	/api/restaurant	GET	2025-05-21 16:37:57.77004	2025-05-21 16:37:57.770165	200	0
8426	\N	/api/data/2F773180-2F5B-41EF-8DA2-1109D921BBE0_4d1b0b7d141e7f9e.jpg	GET	2025-05-21 16:38:06.041678	2025-05-21 16:38:06.041768	200	0
8427	\N	/status	GET	2025-05-21 16:38:09.50755	2025-05-21 16:38:09.507552	200	0
8428	\N	/api/statistics/global	GET	2025-05-21 16:38:09.587563	2025-05-21 16:38:09.590972	200	3
8429	\N	/api/statistics/top-users	GET	2025-05-21 16:38:09.647348	2025-05-21 16:38:09.649728	200	2
8430	\N	/status	GET	2025-05-21 16:38:16.737526	2025-05-21 16:38:16.737528	200	0
8431	\N	/api/statistics/global	GET	2025-05-21 16:38:16.825522	2025-05-21 16:38:16.828139	200	2
8432	\N	/api/statistics/top-users	GET	2025-05-21 16:38:16.907715	2025-05-21 16:38:16.90956	200	1
8433	\N	/status	GET	2025-05-21 16:38:19.497755	2025-05-21 16:38:19.497757	200	0
8434	\N	/api/restaurant	GET	2025-05-21 16:38:32.644638	2025-05-21 16:38:32.644725	200	0
8435	\N	/api/washingmachines	GET	2025-05-21 16:38:35.957654	2025-05-21 16:38:35.992808	200	35
8436	\N	/api/traq/	GET	2025-05-21 16:38:46.53372	2025-05-21 16:38:46.535221	200	1
8437	\N	/status	GET	2025-05-21 16:38:53.337643	2025-05-21 16:38:53.337645	200	0
8438	\N	/api/statistics/global	GET	2025-05-21 16:38:53.417571	2025-05-21 16:38:53.420308	200	2
8439	\N	/api/statistics/top-users	GET	2025-05-21 16:38:53.497465	2025-05-21 16:38:53.499317	200	1
8440	\N	/status	GET	2025-05-21 16:38:56.308146	2025-05-21 16:38:56.308148	200	0
8441	\N	/api/statistics/global	GET	2025-05-21 16:38:56.39798	2025-05-21 16:38:56.400339	200	2
8442	\N	/api/statistics/top-users	GET	2025-05-21 16:38:56.457985	2025-05-21 16:38:56.459514	200	1
8443	\N	/status	GET	2025-05-21 16:39:03.429478	2025-05-21 16:39:03.42948	200	0
8444	\N	/status	GET	2025-05-21 16:39:33.54751	2025-05-21 16:39:33.547512	200	0
8445	\N	/api/auth/login	POST	2025-05-21 16:45:23.923829	2025-05-21 16:45:23.9994	401	75
8446	\N	/api/auth/login	POST	2025-05-21 16:45:27.162554	2025-05-21 16:45:27.239625	200	77
8447	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:27.314905	2025-05-21 16:45:27.316634	200	1
8448	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:27.375772	2025-05-21 16:45:27.377871	200	2
8449	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:27.624234	2025-05-21 16:45:27.625999	200	1
8450	\N	/api/washingmachines	GET	2025-05-21 16:45:27.542406	2025-05-21 16:45:27.642923	200	100
8451	\N	/api/weather	GET	2025-05-21 16:45:27.625229	2025-05-21 16:45:27.676211	200	50
8452	\N	/api/restaurant	GET	2025-05-21 16:45:27.625173	2025-05-21 16:45:28.351408	200	726
8453	\N	/api/data/467bdeed-c69f-4074-bb6e-765100a569f5_85fe4e125190b7b4.jpeg	GET	2025-05-21 16:45:32.074394	2025-05-21 16:45:32.07458	200	0
8454	test@imt-atlantique.net	/api/newf/me	PATCH	2025-05-21 16:45:46.680514	2025-05-21 16:45:46.684906	200	4
8455	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:48.877042	2025-05-21 16:45:48.879212	200	2
8456	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:49.225368	2025-05-21 16:45:49.227367	200	1
8457	test@imt-atlantique.net	/api/newf/me	GET	2025-05-21 16:45:49.589523	2025-05-21 16:45:49.593026	200	3
8458	\N	/status	GET	2025-05-21 16:48:21.798615	2025-05-21 16:48:21.798617	200	0
8478	\N	/api/restaurant	GET	2025-05-21 17:15:53.022692	2025-05-21 17:15:53.022832	200	0
8479	\N	/api/data/New_Project_943f029163e65ef3.png	GET	2025-05-21 17:15:53.022143	2025-05-21 17:15:53.022735	200	0
8480	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:18:02.004814	2025-05-21 17:18:02.009341	200	4
8481	\N	/api/washingmachines	GET	2025-05-21 17:18:02.413307	2025-05-21 17:18:02.504015	200	90
8483	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:18:02.663385	2025-05-21 17:18:02.665342	200	1
8484	\N	/api/weather	GET	2025-05-21 17:18:02.598882	2025-05-21 17:18:02.673363	200	74
8485	\N	/api/washingmachines	GET	2025-05-21 17:18:06.510953	2025-05-21 17:18:06.538449	200	27
8486	\N	/api/washingmachines	GET	2025-05-21 17:18:13.508694	2025-05-21 17:18:13.53729	200	28
8487	\N	/api/washingmachines	GET	2025-05-21 17:18:39.139141	2025-05-21 17:18:39.167274	200	28
8488	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:18:49.309936	2025-05-21 17:18:49.312594	200	2
8489	\N	/api/washingmachines	GET	2025-05-21 17:18:49.698758	2025-05-21 17:18:49.726916	200	28
8490	\N	/api/weather	GET	2025-05-21 17:18:49.92884	2025-05-21 17:18:49.928878	200	0
8492	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:18:49.978546	2025-05-21 17:18:49.980817	200	2
8493	\N	/api/washingmachines	GET	2025-05-21 17:18:51.538785	2025-05-21 17:18:51.567409	200	28
8494	nathaniel.guitton@imt-atlantique.net	/api/newf/me	GET	2025-05-21 17:19:04.118698	2025-05-21 17:19:04.120602	200	1
8504	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 18:39:12.245693	2025-05-21 18:39:12.247511	200	1
8505	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 18:40:37.318773	2025-05-21 18:40:37.322038	200	3
8506	\N	/api/washingmachines	GET	2025-05-21 18:40:38.225126	2025-05-21 18:40:38.260607	200	35
8507	\N	/api/weather	GET	2025-05-21 18:40:38.54931	2025-05-21 18:40:38.549334	200	0
8509	lucie.delestre@imt-atlantique.net	/api/newf/me	GET	2025-05-21 18:40:38.731297	2025-05-21 18:40:38.734335	200	3
8510	\N	/api/washingmachines	GET	2025-05-21 18:40:46.928085	2025-05-21 18:40:46.962766	200	34
8511	\N	/api/restaurant	GET	2025-05-21 18:40:59.120113	2025-05-21 18:40:59.120199	200	0
\.


--
-- Data for Name: reservations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.reservations (email, id_rooms, start_date, end_date) FROM stdin;
\.


--
-- Data for Name: restaurant; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.restaurant (id_restaurant, articles, updated_date, language) FROM stdin;
20	{"grilladesMidi":["Tortelloni ricotta  et épinards","Filet de Tacaud vapeur sauce vin blanc ciboulette"],"migrateurs":["Bolognaise de boeuf","Poulet aux arachides"],"cibo":["Blanquette aux haricots rouges","Autre plat végétarien disponible pôle 1 GrilladeTortelloni"],"accompMidi":["Riz bio","Tombé de chou vert","Pâtes bio"],"grilladesSoir":["Tartiflette végétarienne et légumes","Brandade de poisson basquaise"],"accompSoir":["Semoule bio","Choux fleurs"]}	2025-03-12 08:20:01.020799	1
19	{"grilladesMidi":["Brouillade d'oeufs au jambon","Steak de thon sauce chien","Pièce de boeuf  VBF / RAV supplément de 1,06€"],"migrateurs":["Colombo de poulet"],"cibo":["Gratin de christophines aux légumes"],"accompMidi":["Riz bio créole parfumé","Ecrasé de patate douce"],"grilladesSoir":["Bolognaise végétarienne et légumes","Bolognaise de boeuf basquaise"],"accompSoir":["Carottes bio et panais","Spaghetti bio"]}	2025-03-11 14:37:17.299237	1
31	{"grilladesMidi":["Zarzuela de colin  aux fruits de mer","Parmentier de dinde et champignons"],"migrateurs":["Risotto de la mer","Boulette de boeuf coriandre et sésame"],"cibo":["Dahl de lentilles  lait de coco et curry","Ramen aux oeufs brouillés"],"accompMidi":["Riz bio","Pommes de terre persillées","Carottes bio braisées","Pommes de terre persillées  bio"],"grilladesSoir":[],"accompSoir":[]}	2025-03-14 10:25:01.633984	1
32	{"grilladesMidi":["Zarzuela de colin  aux fruits de mer","Parmentier de dinde et champignons"],"migrateurs":["Boulette de boeuf coriandre et sésame",""],"cibo":["Dahl de lentilles  lait de coco et curry","Ramen aux oeufs brouillés"],"accompMidi":["Riz bio","Pommes de terre persillées","Carottes bio braisées","Pommes de terre persillées  bio"],"grilladesSoir":[],"accompSoir":[]}	2025-03-14 10:31:00.843174	1
33	{"grilladesMidi":["Gnocchi carbonara et champignons","Risotto de la mer"],"migrateurs":["Poulet yassa","Parmentier de dinde et champignons"],"cibo":["Pommes de terre au four,  fondue de lentilles, poireauxet cantal"],"accompMidi":["Riz bio","Pommes de terre persillées","Carottes bio braisées","Coeur de blé bio","Haricots plats au beurre"],"grilladesSoir":["Vol au vent de champignons et courges","Calamars à la  romaine"],"accompSoir":["Flan de carottes","Riz bio"]}	2025-03-17 08:30:00.953605	1
34	{"grilladesMidi":["Gnocchi carbonara et champignons","Risotto de la mer"],"migrateurs":["Poulet yassa","Parmentier de dinde et champignons"],"cibo":["Pommes de terre au four,  fondue de lentilles, poireauxet cantal"],"accompMidi":["Riz bio","Coeur de blé bio","Haricots plats au beurre"],"grilladesSoir":["Vol au vent de champignons et courges","Calamars à la  romaine"],"accompSoir":["Flan de carottes","Riz bio"]}	2025-03-17 09:50:00.252406	1
27	{"grilladesMidi":["Zarzuela de colin  aux fruits de mer","Rognons de boeuf sauce moutarde"],"migrateurs":["Brochette de saucisse  au pesto"],"cibo":["Blanquette aux haricots blancs","Wok brocolis et champignons"],"accompMidi":["Riz bio","Lentilles vertes bio","Haricots verts au naturel"],"grilladesSoir":["Oeuf brouillé aux champignons","Fricassée de  saucisse"],"accompSoir":["Riz bio","Brocolis gratinés bio"]}	2025-03-13 10:19:02.556529	1
28	{"grilladesMidi":["Zarzuela de colin  aux fruits de mer","Rognons de boeuf sauce moutarde"],"migrateurs":["Brochette de saucisse  au pesto"],"cibo":["Blanquette aux haricots blancs"],"accompMidi":["Riz bio","Lentilles vertes bio","Haricots verts au naturel"],"grilladesSoir":["Oeuf brouillé aux champignons","Fricassée de  saucisse"],"accompSoir":["Riz bio","Brocolis gratinés bio"]}	2025-03-13 11:10:00.387631	1
29	{"grilladesMidi":[],"migrateurs":["Zarzuela de colin  aux fruits de mer","Boulette de boeuf coriandre et sésame"],"cibo":["Dahl de lentilles  lait de coco et curry","Ramen aux oeufs brouillés"],"accompMidi":["Riz bio","Pommes de terre persillées  bio","Carottes bio braisées"],"grilladesSoir":[],"accompSoir":[]}	2025-03-14 07:30:01.101923	1
35	{"grilladesMidi":["Gnocchi carbonara et champignons","Risotto de la mer"],"migrateurs":["Poulet yassa","Parmentier de dinde et champignons"],"cibo":["Pommes de terre au four,  fondue de lentilles, poireauxet cantal"],"accompMidi":["Riz bio","Coeur de blé bio","Haricots plats au beurre"],"grilladesSoir":["Parmentier végétal de potimarron","Calamars à la  romaine"],"accompSoir":["Flan de carottes","Riz bio"]}	2025-03-17 12:50:00.961665	1
37	{"grilladesMidi":["Echine de porc braisée","Gnocchi carbonara en gratin","Risotto de poulet  et petits pois","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Couscous au merguez","Calamars à la romaine"],"cibo":["Pommes de terre au four,  fondue de lentilles, poireauxet cantal","Tourte aux blettes  et champignons"],"accompMidi":["Purée de pommes de terre bio","Epinards bio à la crème","Semoule bio"],"grilladesSoir":["Nouille de riz au wok de légumes","Endives au jambon romaine","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Brocolis","Pommes de terre vapeur"]}	2025-03-18 10:18:44.181175	1
38	{"grilladesMidi":["Echine de porc braisée","Gnocchi carbonara en gratin","Risotto de poulet  et petits pois","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Couscous au merguez"],"cibo":["Pommes de terre au four,  fondue de lentilles, poireauxet cantal","Tourte aux blettes  et champignons"],"accompMidi":["Purée de pommes de terre bio","Epinards bio à la crème","Semoule bio"],"grilladesSoir":["Nouille de riz au wok de légumes","Endives au jambon romaine","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Brocolis","Pommes de terre vapeur"]}	2025-03-18 11:20:00.98795	1
39	{"grilladesMidi":["Filet de Flet Waterzoï","Bruschetta buttenut,  chèvre et roquette","Risotto de poulet  et petits pois"],"migrateurs":["Sauté de poulet teriyaki"],"cibo":["Riz sauté aux amandes et légumes, brouillade d'oeufs","Tourte aux blettes  et champignons","Autre plat végétarien disponible pôle 1 GrilladeBruschetta"],"accompMidi":["Riz bio","Pâtes bio","Flan de carottes bio"],"grilladesSoir":["Pizza 3 fromages de légumes","Colombo de dinde romaine","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Gratin de choux fleurs","Riz bio"]}	2025-03-19 08:20:01.075389	1
40	{"grilladesMidi":["Filet de Flet Waterzoï","Bruschetta butternut,  chèvre et roquette","Risotto de poulet  et petits pois"],"migrateurs":["Sauté de poulet teriyaki"],"cibo":["Riz sauté aux amandes et légumes, brouillade d'oeufs","Tourte aux blettes  et champignons","Autre plat végétarien disponible pôle 1 GrilladeBruschetta"],"accompMidi":["Riz bio","Pâtes bio","Flan de carottes bio"],"grilladesSoir":["Pizza 3 fromages de légumes","Colombo de dinde romaine","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Gratin de choux fleurs","Riz bio"]}	2025-03-19 09:38:00.816212	1
41	{"grilladesMidi":["Filet de lieu noir  bouillon Vietnamien","Escalope de cochon sauce diable"],"migrateurs":["Sauté de dinde aux olives  et poivrons"],"cibo":["Riz sauté aux amandes et légumes, brouillade d'oeufs","Bruschetta butternut,  chèvre et roquette"],"accompMidi":["Riz bio","Gratin d'endives"],"grilladesSoir":["Quiche potiron cheddar","Haut de cuisse de poulet aux herbes"],"accompSoir":["Haricots verts","Semoule bio"]}	2025-03-20 07:15:05.414714	1
42	{"grilladesMidi":["Pizza raclette","Moules marinière"],"migrateurs":["Longe de porc laqué"],"cibo":["Quiche potiron cheddar","Gnocchi aux légumes, tomates et cheddar","Croziflette poireaux  et champignons"],"accompMidi":["Riz bio","Pommes de terre grenaille","Gratin d'endives"],"grilladesSoir":[],"accompSoir":[]}	2025-03-21 07:50:01.011732	1
43	{"grilladesMidi":["Filet de lieu noir bouillon à la Vietnamienne","Moules marinière"],"migrateurs":["Longe de porc laqué"],"cibo":["Quiche potiron cheddar","Gnocchi aux légumes, tomates et cheddar","Croziflette poireaux  et champignons"],"accompMidi":["Riz bio","Pommes de terre grenaille","Gratin d'endives","Fenouil rôti"],"grilladesSoir":[],"accompSoir":[]}	2025-03-21 09:30:01.021974	1
44	{"grilladesMidi":["Moules marinière","Pièce de boeuf  VBF / RAV supplément de 1,06€"],"migrateurs":["Longe de porc laqué"],"cibo":["Quiche potiron cheddar","Gnocchi aux légumes, tomates et cheddar","Croziflette poireaux  et champignons"],"accompMidi":["Riz bio","Pommes de terre grenaille","Gratin d'endives","Fenouil rôti"],"grilladesSoir":[],"accompSoir":[]}	2025-03-21 11:30:00.999619	1
45	{"grilladesMidi":["Viennoise de dinde, sauce Arrabiata","Boudin noir aux pommes"],"migrateurs":["Aiguillette de poulet  citronnelle"],"cibo":["Mac and cheese, brocolis cheddar"],"accompMidi":["Riz bio","Haricots verts en persillade"],"grilladesSoir":["Rôti de porc"],"accompSoir":["Gratin de choux fleurs","Riz bio"]}	2025-03-24 08:30:01.03884	1
46	{"grilladesMidi":["Viennoise de dinde, sauce Arrabiata","Boudin noir aux pommes"],"migrateurs":["Aiguillette de poulet  citronnelle"],"cibo":["Coquillettes and cheese brocolis et cheddar"],"accompMidi":["Riz bio","Haricots verts en persillade"],"grilladesSoir":["Ragoût de patates douces et courges à la cacahuètes","Rôti de porc"],"accompSoir":["Gratin de choux fleurs","Riz bio"]}	2025-03-24 09:20:00.974038	1
47	{"grilladesMidi":["Raviole aux fromages  gratinés","Encornet basquaise","Pièce de boeuf  VBF / RAV supplément de 1,06€"],"migrateurs":["Sauté de dinde  à la Marocaine"],"cibo":["Coquillettes and cheese brocolis et cheddar","Tarte poireaux, champignons  à l'emmental","Autre plat végétarien disponible pôle 1 GrilladeRavioles"],"accompMidi":["Riz bio","Pommes de terre grenaille","Haricots verts en persillade","Semoule bio","Carottes bio fondantes"],"grilladesSoir":["Tartiflette aux légumes et pommes de terre","Fricassée de poulet  aux champignons"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-03-25 08:40:00.96241	1
48	{"grilladesMidi":["Raviole aux fromages  gratinés","Encornet basquaise","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Sauté de dinde  à la Marocaine"],"cibo":["Coquillettes and cheese brocolis et cheddar","Tarte poireaux, champignons  à l'emmental","Autre plat végétarien disponible pôle 1 GrilladeRavioles"],"accompMidi":["Riz bio","Pommes de terre grenaille","Haricots verts en persillade","Semoule bio","Carottes bio fondantes"],"grilladesSoir":["Tartiflette aux légumes et pommes de terre","Fricassée de poulet  aux champignons","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-03-25 13:20:01.225558	1
49	{"grilladesMidi":["Raviole aux fromages  gratinés","Bolognaise de légumes"],"migrateurs":["Encornet basquaise","Cuisse de poulet BBC barbecue"],"cibo":["Lasagne aux blettes et épinard","Tartiflette végétarienne"],"accompMidi":["Riz bio","Fusilis bio","Gratin de poireaux"],"grilladesSoir":["Chachouka poireaux épinard et pommes de terre","Brushetta pizza aux champignons","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Patates douces rôties","Ecrasé de brocolis"]}	2025-03-26 09:16:13.845169	1
50	{"grilladesMidi":["Bolognaise de légumes","Ragout de poisson  au haricots blancs"],"migrateurs":["Encornet basquaise","Cuisse de poulet BBC barbecue"],"cibo":["Lasagne aux blettes et épinard","Tartiflette végétarienne"],"accompMidi":["Riz bio","Fusilis bio","Gratin de poireaux"],"grilladesSoir":["Chachouka poireaux épinard et pommes de terre","Brushetta pizza aux champignons","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Patates douces rôties","Ecrasé de brocolis"]}	2025-03-26 11:20:00.936773	1
51	{"grilladesMidi":["Omelette au fromage poivrons et basilic","Vol au vent de colin moutarde à l'ancienne"],"migrateurs":["Columbo de porc"],"cibo":["Lasagne aux blettes et épinard","Curry aux 2 lentilles poireaux et riz","Autre plat végétarien disponible pôle 1 GrilladeOmelette"],"accompMidi":["Riz bio","Pomme de terre vapeur","Choux fleurs rôti au paprika"],"grilladesSoir":["Galette de blé noir  aux légumes et chèvre","Blanc de poulet sauce crème"],"accompSoir":["Coeur de blé bio","Poêlée de courgettes  aux poivrons"]}	2025-03-27 08:10:00.847256	1
52	{"grilladesMidi":["Omelette au fromage poivrons et basilic","Vol au vent de colin moutarde à l'ancienne"],"migrateurs":["Columbo de porc"],"cibo":["Lasagne aux blettes et épinard","Curry aux 2 lentilles poireaux et riz","Autre plat végétarien disponible pôle 1 GrilladeOmelette"],"accompMidi":["Riz bio","Pomme de terre vapeur","Choux fleurs rôti au paprika","gratin de poireaux","Gratin de poireaux"],"grilladesSoir":["Galette de blé noir  aux légumes et chèvre","Blanc de poulet sauce crème"],"accompSoir":["Coeur de blé bio","Poêlée de courgettes  aux poivrons"]}	2025-03-27 09:10:00.975433	1
53	{"grilladesMidi":["Omelette au fromage poivrons et basilic","Vol au vent de colin moutarde à l'ancienne"],"migrateurs":["Columbo de porc"],"cibo":["Lasagne aux blettes et épinard","Autre plat végétarien disponible pôle 1 GrilladeOmelette"],"accompMidi":["Pomme de terre vapeur","Choux fleurs rôti au paprika","gratin de poireaux","Gratin de poireaux"],"grilladesSoir":["Galette de blé noir  aux légumes et chèvre","Blanc de poulet sauce crème"],"accompSoir":["Coeur de blé bio","Poêlée de courgettes  aux poivrons"]}	2025-03-27 12:20:01.028569	1
54	{"grilladesMidi":["Filet de julienne  câpres et anchois","Poulet rôti","Gratin de poisson"],"migrateurs":["Boeuf haché à la Cubaine","Sauté de porc curry et coco"],"cibo":["Crumble chèvre, panais carottes bio et courges"],"accompMidi":["Riz bio","Poêlée de légumes verts","Mélange de céréales bio"],"grilladesSoir":[],"accompSoir":[]}	2025-03-28 08:00:00.99215	1
55	{"grilladesMidi":["Filet de julienne  câpres et anchois","Poulet rôti","Gratin de poisson"],"migrateurs":["Boeuf haché à la Cubaine"],"cibo":["Crumble chèvre, panais carottes bio et courges"],"accompMidi":["Riz bio","Poêlée de légumes verts","Mélange de céréales bio"],"grilladesSoir":[],"accompSoir":[]}	2025-03-28 11:00:01.354201	1
56	{"grilladesMidi":["Filet de julienne  câpres et anchois","Encornet persillé aux artichauts","Andouillette  crème moutarde à l'ancienne"],"migrateurs":["Boeuf haché à la Cubaine","Picata de poulet  au citron"],"cibo":["Gratin de poireaux  et blé à la tomate"],"accompMidi":["Riz bio","Pâtes bio","Haricots verts bio"],"grilladesSoir":["Hachi parmentier végétal","Tajine de boulettes de boeuf"],"accompSoir":["Semoule bio","Légumes couscous aux poivrons"]}	2025-03-31 07:20:01.007638	1
57	{"grilladesMidi":["Encornet persillé aux artichauts","Andouillette  crème moutarde à l'ancienne"],"migrateurs":["Picata de poulet  au citron"],"cibo":["Gratin de poireaux  et blé à la tomate"],"accompMidi":["Riz bio","Pâtes bio","Haricots verts bio"],"grilladesSoir":["Hachi parmentier végétal","Tajine de boulettes de boeuf"],"accompSoir":["Semoule bio","Légumes couscous aux poivrons"]}	2025-03-31 10:40:01.089408	1
58	{"grilladesMidi":["Vol au vent de colin à la moutarde","Rognons à la Bordelaise","Andouillette  crème moutarde à l'ancienne"],"migrateurs":["Viennoise de poulet sauce arabiata","Encornet persillé aux artichauts"],"cibo":["Gnocchi à la ricotta  et épinards bio"],"accompMidi":["Semoule bio","Purée de pommes de terre bio","Céleri rôti et champignons en persilllade"],"grilladesSoir":["Omeltte au fromage  et fines herbes","Blanc de poulet  sauce moutarde et champignons"],"accompSoir":["Patates douces rôties","Petits pois carottes"]}	2025-04-01 06:20:01.082493	1
59	{"grilladesMidi":["Vol au vent de colin à la moutarde","Rognons à la Bordelaise","Andouillette  crème moutarde à l'ancienne"],"migrateurs":["Viennoise de poulet sauce arabiatasauce arrabbiata","Encornet persillé aux artichauts"],"cibo":["Gnocchi à la ricotta  et épinards bio"],"accompMidi":["Semoule bio","Purée de pommes de terre bio","Céleri rôti et champignons en persilllade"],"grilladesSoir":["Omeltte au fromage  et fines herbes","Blanc de poulet  sauce moutarde et champignons"],"accompSoir":["Patates douces rôties","Petits pois carottes"]}	2025-04-01 10:00:00.953748	1
60	{"grilladesMidi":["Vol au vent de colin à la moutarde","Rognons à la Bordelaise"],"migrateurs":["Encornet persillé aux artichauts"],"cibo":["Gnocchi à la ricotta  et épinards bio"],"accompMidi":["Semoule bio","Purée de pommes de terre bio","Céleri rôti et champignons en persilllade"],"grilladesSoir":["Omeltte au fromage  et fines herbes","Blanc de poulet  sauce moutarde et champignons"],"accompSoir":["Patates douces rôties","Petits pois carottes"]}	2025-04-01 10:30:00.968848	1
61	{"grilladesMidi":["Gratin de poissons","Moules au curry","Saucisse de Toulouse"],"migrateurs":["Escalope de dinde grillé mariné à l'origan"],"cibo":["Hachi parmentier végétal","Gratin de poireaux  et blé à la tomate"],"accompMidi":["RIz bio","Patates douces rôti","Brocolis, vinaigrette en persilllade","Brocolis, vinaigrette citron câpres"],"grilladesSoir":["Risotto de coquillettes,  poireaux et pesto de persil","Sauté de porc au lait  de coco et curry"],"accompSoir":["Riz bio","Gratin de blettes"]}	2025-04-02 07:00:00.995893	1
62	{"grilladesMidi":["Gratin de poisson","Moules au curry","Saucisse de Toulouse"],"migrateurs":["Escalope de dinde grillé mariné à l'origan"],"cibo":["Hachi parmentier végétal","Gratin de poireaux  et blé à la tomate"],"accompMidi":["RIz bio","Patates douces rôti","Brocolis, vinaigrette en persilllade","Brocolis, vinaigrette citron câpres"],"grilladesSoir":["Risotto de coquillettes,  poireaux et pesto de persil","Sauté de porc au lait  de coco et curry"],"accompSoir":["Riz bio","Gratin de blettes"]}	2025-04-02 07:10:00.932275	1
63	{"grilladesMidi":["Gratin de poisson","Moules au curry","Saucisse de Toulouse"],"migrateurs":["Escalope de dinde grillé mariné à l'origan"],"cibo":["Hachi parmentier végétal","Gratin de poireaux  et blé à la tomate"],"accompMidi":["RIz bio","Patates douces rôti","Brocolis  vinaigrette citron câpres"],"grilladesSoir":["Risotto de coquillettes,  poireaux et pesto de persil","Sauté de porc au lait  de coco et curry"],"accompSoir":["Riz bio","Gratin de blettes"]}	2025-04-02 08:30:00.859576	1
64	{"grilladesMidi":["Gratin de poisson","Moules au curry","Saucisse de Toulouse","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Escalope de dinde grillé mariné à l'origan"],"cibo":["Hachi parmentier végétal","Gratin de poireaux  et blé à la tomate"],"accompMidi":["RIz bio","Patates douces rôti","Brocolis  vinaigrette citron câpres"],"grilladesSoir":["Risotto de coquillettes,  poireaux et pesto de persil","Sauté de porc au lait  de coco et curry","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Riz bio","Gratin de blettes"]}	2025-04-02 09:00:01.010638	1
65	{"grilladesMidi":["Gratin de poisson","Moules au curry","Saucisse de Toulouse","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Escalope de dinde grillé mariné à l'origan"],"cibo":["Hachi parmentier végétal","Gratin de poireaux  et blé à la tomate"],"accompMidi":["RIz bio","Patates douces rôti","Brocolis  vinaigrette citron câpres"],"grilladesSoir":["Risotto de coquillettes,  poireaux et pesto de persil","Sauté de porc au lait  de coco et curry","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Riz bio","Gratin de céléri champignons"]}	2025-04-02 12:00:00.913892	1
66	{"grilladesMidi":["Escalope de dinde grillé mariné à l'origan","Moules au curry","Jambon sauce moutarde","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Poulet, miel et paprika"],"cibo":["Risotto de coquillettes,  poireaux et pesto de persil","Pastilla de légumes et lentilles vertes bio"],"accompMidi":["RIz bio","Boulgour bio","Carottes bio au cumin vinaigrette citron câpres"],"grilladesSoir":["Curry de patates douces brocolis et pois chichesnoix de cajou","Croque monsieur jambon fromage","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Penne","Brocolis"]}	2025-04-03 07:40:01.204039	1
67	{"grilladesMidi":["Escalope de dinde grillé mariné à l'origan","Moules au curry","Jambon sauce moutarde","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Poulet, miel et paprika"],"cibo":["Risotto de coquillettes,  poireaux et pesto de persil","Pastilla de légumes et lentilles vertes bio"],"accompMidi":["RIz bio","Boulgour bio","Carottes bio au cumin"],"grilladesSoir":["Curry de patates douces brocolis et pois chichesnoix de cajou","Croque monsieur jambon fromage","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Penne","Brocolis"]}	2025-04-03 08:20:01.05971	1
68	{"grilladesMidi":["Moules au curry","Jambon sauce moutarde","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Poulet, miel et paprika"],"cibo":["Risotto de coquillettes,  poireaux et pesto de persil","Pastilla de légumes et lentilles vertes bio"],"accompMidi":["RIz bio","Boulgour bio","Carottes bio au cumin"],"grilladesSoir":["Curry de patates douces brocolis et pois chichesnoix de cajou","Croque monsieur jambon fromage","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Penne","Brocolis"]}	2025-04-03 10:30:00.969373	1
69	{"grilladesMidi":["Moules au curry","Jambon sauce moutarde","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"migrateurs":["Calamars à la romaine"],"cibo":["Risotto de coquillettes,  poireaux et pesto de persil","Pastilla de légumes et lentilles vertes bio"],"accompMidi":["RIz bio","Boulgour bio","Carottes bio au cumin"],"grilladesSoir":["Curry de patates douces brocolis et pois chichesnoix de cajou","Croque monsieur jambon fromage","Plat disponible dans le frigo  rose bonne action pour les étudiants"],"accompSoir":["Penne","Brocolis"]}	2025-04-03 10:40:00.953225	1
70	{"grilladesMidi":["Filet de lingue bleue gomasio, soja et citron","Côte de porc façon Carbonade"],"migrateurs":["Pizza Napolitaine"],"cibo":["Shawarma de choux fleurs,  choux rouges mi-cuit aux pois chiches","Pastilla de légumes et lentilles vertes bio"],"accompMidi":["Pâtes bio","Blé bio","Epinards bio à la crème"],"grilladesSoir":[],"accompSoir":[]}	2025-04-04 06:30:01.005201	1
71	{"grilladesMidi":["Filet de lingue bleue gomasio, soja et citron","Côte de porc façon Carbonade"],"migrateurs":["Pizza Napolitaine"],"cibo":["Shawarma de choux fleurs,  choux rouges mi-cuit aux pois chiches","Pastilla de légumes et lentilles vertes bio","Autre plat végétarien disponible pôle 1 GrilladePizza","Autre plat végétarien disponible pôle 2 Cuistots MigrateursPizza"],"accompMidi":["Pâtes bio","Blé bio","Epinards bio à la crème"],"grilladesSoir":[],"accompSoir":[]}	2025-04-04 09:50:00.948518	1
72	{"grilladesMidi":["Pizza Napolitaines","Côte de porc façon Carbonade"],"migrateurs":["Pizza Napolitaine"],"cibo":["Shawarma de choux fleurs,  choux rouges mi-cuit aux pois chiches","Pastilla de légumes et lentilles vertes bio","Autre plat végétarien disponible pôle 1 GrilladePizza","Autre plat végétarien disponible pôle 2 Cuistots MigrateursPizza"],"accompMidi":["Pâtes bio","Blé bio","Epinards bio à la crème"],"grilladesSoir":[],"accompSoir":[]}	2025-04-04 10:40:00.995074	1
73	{"grilladesMidi":["Bruschetta, jambon cru poivrons et courgettes","Aiguillette de poulet à la Dijonnaise"],"migrateurs":["Rôti de porc à la Danoise"],"cibo":["Mac and cheese au cheddar"],"accompMidi":["Légumes marinés à l'orientale","Riz bio","Epinards bio à la crème","Pommes de terre grenaille"],"grilladesSoir":["Tourte d'aubergines boulgour bioet feta","Cordon bleu"],"accompSoir":["Pâtes bio","Haricots verts"]}	2025-04-07 07:20:00.954701	1
74	{"grilladesMidi":["Bruschetta, jambon cru poivrons et courgettes","Aiguillette de poulet à la Dijonnaise"],"migrateurs":["Rôti de porc à la Danoise"],"cibo":["Mac and cheese au cheddar"],"accompMidi":["Légumes marinés à l'orientale","Riz bio","Epinards bio à la crème","Pommes de terre grenaille"],"grilladesSoir":["Brushetta au jambon cru courgettes et poivrons","Tourte d'aubergines boulgour bioet feta","Viennoise de dinde sauce échalotes et vin rouge"],"accompSoir":["Pâtes bio","Haricots verts","Riz bio"]}	2025-04-07 14:10:00.414448	1
75	{"grilladesMidi":["Bissap de boeuf","Aiguillette de poulet à la Dijonnaise"],"migrateurs":["Penne puttanesca au poulet"],"cibo":["Mac and cheese au cheddar","Haricots rouges  et riz aux poivrons"],"accompMidi":["Pomme de terre grenaille","Fusilli bio","Chou fleur au curry"],"grilladesSoir":["Lasagne aux légumes boulgour bioet feta"],"accompSoir":[]}	2025-04-08 07:50:00.867364	1
76	{"grilladesMidi":["Bissap de boeuf","Filet d'Eiglefin  sauce vierge"],"migrateurs":["Penne puttanesca au poulet"],"cibo":["Mac and cheese au cheddar","Haricots rouges  et riz aux poivrons"],"accompMidi":["Pomme de terre grenaille","Fusilli bio","Chou fleur au curry"],"grilladesSoir":["Lasagne aux légumes boulgour bioet feta"],"accompSoir":[]}	2025-04-08 10:20:00.971802	1
77	{"grilladesMidi":["Bissap de boeuf","Cuisse de poulet  ketchup oriental","Filet de poisson  aux agrumes et fenouil","Gratin de pak chog","Champignons et courgettes"],"migrateurs":["Fricassée de chipolatas  à la tomate","Porc au curry"],"cibo":["Lasagne de légumes","Haricots rouges  et riz aux poivrons","Bagem houmous et pousses d'épinards"],"accompMidi":["Riz bio","Fusilli bio","Duo de céléri et carottes bio","Champignons et courgettes"],"grilladesSoir":["Cassoulet de courgettes  et champignons","Carbonara"],"accompSoir":["Spaghetti bio","Choux fleurs"]}	2025-04-09 07:30:00.988738	1
78	{"grilladesMidi":["Bissap de boeuf","Cuisse de poulet  ketchup oriental","Filet de poisson  aux agrumes et fenouil","Gratin de pak chog","Champignons et courgettes"],"migrateurs":["Fricassée de chipolatas  à la tomate","Porc au curry"],"cibo":["Lasagne de légumes","Haricots rouges  et riz aux poivrons","Bagem houmous et pousses d'épinards"],"accompMidi":["Riz bio","Fusilli bio","Duo de céléri et carottes bio","Champignons et courgettes"],"grilladesSoir":["Cassoulet de courgettes  et champignons","Carbonara"],"accompSoir":["Spaghetti bio","Choux fleurs"]}	2025-04-09 07:30:01.062992	1
79	{"grilladesMidi":["Bissap de boeuf","Cuisse de poulet  ketchup oriental","Filet de eiglefin aux agrumes et fenouil","Gratin de pak chog","Champignons et courgettes"],"migrateurs":["Fricassée de chipolatas  à la tomate","Porc au curry"],"cibo":["Lasagne de légumes","Haricots rouges  et riz aux poivrons","Bagel houmous et pousses d'épinards"],"accompMidi":["Riz bio","Fusilli bio","Duo de céléri et carottes bio","Champignons et courgettes"],"grilladesSoir":["Cassoulet de courgettes  et champignons","Carbonara"],"accompSoir":["Spaghetti bio","Choux fleurs"]}	2025-04-09 08:40:00.846615	1
80	{"grilladesMidi":["Bissap de boeuf","Cuisse de poulet  ketchup oriental","Filet d'eiglefin aux agrumes et fenouil","Gratin de pak chog","Champignons et courgettes"],"migrateurs":["Fricassée de chipolatas  à la tomate","Porc au curry"],"cibo":["Lasagne de légumes","Haricots rouges  et riz aux poivrons","Bagel houmous et pousses d'épinards"],"accompMidi":["Riz bio","Fusilli bio","Duo de céléri et carottes bio","Champignons et courgettes"],"grilladesSoir":["Cassoulet de courgettes  et champignons","Carbonara"],"accompSoir":["Spaghetti bio","Choux fleurs"]}	2025-04-09 08:40:50.471363	1
81	{"grilladesMidi":["Bissap de boeuf","Cuisse de poulet  ketchup oriental","Filet d'eiglefin aux agrumes et fenouil","Champignons et courgettes"],"migrateurs":["Fricassée de chipolatas  à la tomate","Porc au curry"],"cibo":["Lasagne de légumes","Haricots rouges  et riz aux poivrons","Bagel houmous et pousses d'épinards"],"accompMidi":["Riz bio","Fusilli bio","Duo de céléri et carottes bio","Gratin de choux","Champignons et courgettes"],"grilladesSoir":["Cassoulet de courgettes  et champignons","Carbonara"],"accompSoir":["Spaghetti bio","Choux fleurs"]}	2025-04-09 09:10:00.912488	1
186	{"grilladesMidi":["Cuban Ground Beef","Roast chicken"],"migrateurs":[],"cibo":["Curry with 2 lentils, leeks, spinach and rice","Lasagna with chard and spinach"],"accompMidi":["Organic steamed potatoes","Organic peas and carrots","Roasted fennel"],"grilladesSoir":[],"accompSoir":[]}	2025-05-09 15:06:01.365556	2
187	{"grilladesMidi":["Cuban Ground Beef","Roast chicken"],"migrateurs":[],"cibo":["Curry with 2 lentils, leeks, spinach and rice","Lasagna with chard and spinach"],"accompMidi":["Organic steamed potatoes","Organic peas and carrots","Roasted fennel"],"grilladesSoir":[],"accompSoir":[]}	2025-05-10 10:29:12.750909	2
195	{"grilladesMidi":["Boeuf haché à la cubaine","Poulet rôti"],"migrateurs":[],"cibo":["Curry aux 2 lentilles, poireaux, épinards et riz","Lasagne aux blettes et épinards"],"accompMidi":["Pomme de terre bio vapeur","Petits pois carottes bio","Fenouil rôti"],"grilladesSoir":[],"accompSoir":[]}	2025-05-11 07:57:36.83441	1
185	{"grilladesMidi":["Boeine","Poulet rôti"],"migrateurs":[],"cibo":["Cuux et riz","Lasas"],"accompMidi":["Poeur","Petio","Feti"],"grilladesSoir":[],"accompSoir":[]}	2025-05-09 08:40:00.921305	1
95	{"grilladesMidi":["Suprème de poulet","Porc au curry"],"migrateurs":["Saucisse de Toulouse"],"cibo":["Tarte fine aux légumes et  mozzarella","Courgettes à la marocaine"],"accompMidi":["Spaghetti bio","Fenouil au four","Duo de céléri et carottes bio","Riz bio"],"grilladesSoir":["Pizza végétarienne et champignons","Sauté de dinde basquaise"],"accompSoir":["Semoule bio","Gratin de courgettes"]}	2025-04-10 09:40:00.878096	1
96	{"grilladesMidi":["Suprème de poulet","Porc au curry"],"migrateurs":["Saucisse de Toulouse"],"cibo":["Tarte fine aux légumes et  mozzarella","Courgettes à la marocaine"],"accompMidi":["Spaghetti bio","Fenouil au four","Duo de céléri et carottes bio","Riz bio"],"grilladesSoir":["Pizza 3 fromages","Sauté de dinde basquaise"],"accompSoir":["Semoule bio","Gratin de courgettes"]}	2025-04-10 11:30:00.98713	1
115	{"grilladesMidi":["Suprème de poulet","Saucisse de Toulouse sauce moutarde et thym","Effiloché de poulet basquaise","Couscous merguez"],"migrateurs":[],"cibo":["Oeufs brouillés  légumes provençale"],"accompMidi":["Semoule bio","Pommes de terre boulangère","Brocolis","Poélée de carottes"],"grilladesSoir":[],"accompSoir":[]}	2025-04-11 12:09:13.948285	1
116	{"grilladesMidi":["至尊鸡肉","图卢兹香肠配芥末和百里香酱","巴斯克风味鸡丝","梅尔格斯库斯"],"migrateurs":[],"cibo":["普罗旺斯蔬菜炒鸡蛋"],"accompMidi":["有机粗面粉","烤土豆","西兰花","炒胡萝卜"],"grilladesSoir":[],"accompSoir":[]}	2025-04-11 23:25:51.76506	8
117	{"grilladesMidi":["Chicken Supreme","Toulouse sausage with mustard and thyme sauce","Shredded Basque-style chicken","Merguez couscous"],"migrateurs":[],"cibo":["Scrambled eggs with Provençal vegetables"],"accompMidi":["Organic semolina","Boulangère potatoes","Broccoli","Sautéed carrots"],"grilladesSoir":[],"accompSoir":[]}	2025-04-11 23:25:58.479757	2
118	{"grilladesMidi":["Pollo Supremo","Salchicha de Toulouse con salsa de mostaza y tomillo","Pollo desmenuzado al estilo vasco","Cuscús merguez"],"migrateurs":[],"cibo":["Huevos revueltos con verduras provenzales"],"accompMidi":["Sémola orgánica","Patatas boulangère","Brócoli","Zanahorias salteadas"],"grilladesSoir":[],"accompSoir":[]}	2025-04-11 23:26:04.527606	3
119	{"grilladesMidi":["Suprème de poulet","Saucisse de Toulouse sauce moutarde et thym","Effiloché de poulet basquaise","Couscous merguez"],"migrateurs":[],"cibo":["Oeufs brouillés  légumes provençale"],"accompMidi":["Semoule bio","Pommes de terre boulangère","Brocolis","Poélée de carottes"],"grilladesSoir":[],"accompSoir":[]}	2025-04-11 23:30:00.818021	1
120	{"grilladesMidi":["Suprème de poulet","Saucisse de Toulouse sauce moutarde et thym","Effiloché de poulet basquaise","Couscous merguez"],"migrateurs":[],"cibo":["Oeufs brouillés  légumes provençale"],"accompMidi":["Semoule bio","Pommes de terre boulangère","Brocolis","Poélée de carottes"],"grilladesSoir":[],"accompSoir":[]}	2025-04-12 08:55:38.339863	1
121	{"grilladesMidi":["Suprème de poulet","Saucisse de Toulouse sauce moutarde et thym","Effiloché de poulet basquaise","Couscous merguez"],"migrateurs":[],"cibo":["Oeufs brouillés  légumes provençale"],"accompMidi":["Semoule bio","Pommes de terre boulangère","Brocolis","Poélée de carottes"],"grilladesSoir":[],"accompSoir":[]}	2025-04-13 08:22:20.906546	1
122	{"grilladesMidi":["Risotto aux champignons et poireaux","Escalope de porc charcutière"],"migrateurs":["Pilon de poulet au curry curcuma et lait de coco"],"cibo":["Curry de pommes de terre choux fleurs et pois chiches","Autre plat végétarien disponible pôle 1 GrilladeRisotto"],"accompMidi":["Epinards bio","Pommes de terre boulangère"],"grilladesSoir":["Dahl de lentilles et  patates douces","Cheeseburger"],"accompSoir":["Pâtes bio","Carottes bio braisées"]}	2025-04-14 07:50:00.909323	1
123	{"grilladesMidi":["Mushroom and Leek Risotto","Charcuterie pork escalope"],"migrateurs":["Chicken drumstick with turmeric curry and coconut milk"],"cibo":["Potato, cauliflower and chickpea curry","Other vegetarian dish available pole 1 Grill Risotto"],"accompMidi":["Epinards bio","Boulangère potatoes"],"grilladesSoir":["Lentil and Sweet Potato Dahl","Cheeseburger"],"accompSoir":["Organic pasta","Braised organic carrots"]}	2025-04-14 20:36:18.435258	2
124	{"grilladesMidi":["Risotto aux champignons et poireaux","Escalope de porc charcutière"],"migrateurs":["Pilon de poulet au curry curcuma et lait de coco"],"cibo":["Curry de pommes de terre choux fleurs et pois chiches","Autre plat végétarien disponible pôle 1 Grillade Risotto"],"accompMidi":["Epinards bio","Pommes de terre boulangère"],"grilladesSoir":["Dahl de lentilles et patates douces","Cheeseburger"],"accompSoir":["Pâtes bio","Carottes bio braisées"]}	2025-04-14 20:36:19.153031	1
125	{"grilladesMidi":["Risotto aux champignons et poireaux","Escalope de porc charcutière"],"migrateurs":["Pilon de poulet au curry curcuma et lait de coco"],"cibo":["Curry de pommes de terre choux fleurs et pois chiches","Autre plat végétarien disponible pôle 1 Grillade Risotto"],"accompMidi":["Epinards bio","Pommes de terre boulangère"],"grilladesSoir":["Dahl de lentilles et patates douces","Cheeseburger"],"accompSoir":["Pâtes bio","Carottes bio braisées"]}	2025-04-15 07:09:35.881563	1
126	{"grilladesMidi":["Osso bucco de veau marengo","Frittata saveur d'orient charcutière"],"migrateurs":["Filet de dinde grillé","Légumes couscous"],"cibo":["Lasagne végétarienne"],"accompMidi":["Pâtes bio","Semoule bio","Légumes couscous","Riz bio"],"grilladesSoir":["Cannelloni végétarien patates douces","Saucisse de Toulouse/ Chipolatas grillée"],"accompSoir":["Pommes de terre grenailles","Carottes bio braisées"]}	2025-04-15 07:10:00.293137	1
127	{"grilladesMidi":["Osso bucco de veau marengo","Frittata saveur d'orient charcutière"],"migrateurs":["Filet de dinde grillé"],"cibo":["Lasagne végétarienne"],"accompMidi":["Pâtes bio","Semoule bio","Légumes couscous","Pommes de terre boulangères","Riz bio"],"grilladesSoir":["Cannelloni végétarien patates douces","Saucisse de Toulouse/ Chipolatas grillée"],"accompSoir":["Pommes de terre grenailles","Carottes bio braisées"]}	2025-04-15 09:00:00.964355	1
128	{"grilladesMidi":["Osso bucco de veau marengo","Frittata saveur d'orient"],"migrateurs":["Filet de dinde grillé"],"cibo":["Lasagne végétarienne"],"accompMidi":["Pâtes bio","Semoule bio","Légumes couscous","Pommes de terre boulangères","Riz bio"],"grilladesSoir":["Cannelloni végétarien patates douces","Saucisse de Toulouse/ Chipolatas grillée"],"accompSoir":["Pommes de terre grenailles","Carottes bio braisées"]}	2025-04-15 09:10:00.950356	1
129	{"grilladesMidi":["Osso bucco de veau marengo","Frittata saveur d'orient"],"migrateurs":["Filet de dinde grillé"],"cibo":["Lasagne végétarienne"],"accompMidi":["Pâtes bio","Semoule bio","Légumes couscous","Pommes de terre boulangères","Riz bio"],"grilladesSoir":["Chili sin carne","Saucisse de Toulouse/ Chipolatas grillée"],"accompSoir":["Pommes de terre grenailles","Carottes bio braisées"]}	2025-04-15 09:40:00.943233	1
130	{"grilladesMidi":["Osso bucco de veau marengo","Frittata saveur d'orient"],"migrateurs":["Filet de dinde grillé"],"cibo":["Lasagne végétarienne","Autre plat végétarien disponible pôle 1 Grillade Frittata"],"accompMidi":["Pâtes bio","Semoule bio","Légumes couscous","Pommes de terre boulangères","Riz bio"],"grilladesSoir":["Chili sin carne","Saucisse de Toulouse/ Chipolatas grillée"],"accompSoir":["Pommes de terre grenailles","Carottes bio braisées"]}	2025-04-15 10:00:00.971313	1
131	{"grilladesMidi":["Veal osso bucco marengo","Frittata with an oriental flavor"],"migrateurs":["Grilled turkey fillet"],"cibo":["Vegetarian lasagna","Other vegetarian dish available pole 1 Grilled Frittata"],"accompMidi":["Organic pasta","Organic semolina","Vegetable couscous","Boulangère potatoes","Organic rice"],"grilladesSoir":["Chili sin carne","Toulouse sausage/grilled chipolatas"],"accompSoir":["Baby potatoes","Braised organic carrots"]}	2025-04-15 18:56:55.623947	2
132	{"grilladesMidi":["Filet de lieu noir, concassée de tomates câpres et parmesan","Croque monsieur à provençale","Pièce de boeuf VBF / RAV supplément de 1,06€"],"migrateurs":["Aiguilette de poulet à la Sénégalaise"],"cibo":["Flan aux pommes de terre champignons et cantal","Frittata saveur d'orient"],"accompMidi":["Semoule bio","Légumes couscous","Pommes de terre boulangères"],"grilladesSoir":["Cannelloni végétarien","Nuggets sauce tartare"],"accompSoir":["Riz bio","Poêlée maraichère"]}	2025-04-16 07:10:00.966576	1
133	{"grilladesMidi":["Filet de lieu noir, concassée de tomates câpres et parmesan","Croque monsieur à provençale","Pièce de boeuf VBF / RAV supplément de 1,06€"],"migrateurs":["Aiguilette de poulet à la Sénégalaise"],"cibo":["Flan aux pommes de terre champignons et cantal","Frittata saveur d'orient"],"accompMidi":["Semoule bio","Légumes couscous","Pommes de terre boulangères","Ecrasé de patates douces","Carottes bio"],"grilladesSoir":["Cannelloni végétarien","Nuggets sauce tartare"],"accompSoir":["Riz bio","Poêlée maraichère"]}	2025-04-16 08:30:00.945445	1
134	{"grilladesMidi":["Aiguilette de poulet à la Sénégalaise","Croque monsieur à provençale"],"migrateurs":["Emincé de poulet au paprika","Hotchepot de lieu"],"cibo":["Biryani végétarien"],"accompMidi":["Riz bio","Flan de poireaux"],"grilladesSoir":["Quiche champignons ricotta","Fajitas boeuf sauce tartare"],"accompSoir":["Boulgour bio","Gratin d'épinards bio"]}	2025-04-17 07:30:00.993587	1
135	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Carotte bio","Riz bio","Pâte bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-18 08:10:00.935651	1
136	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Riz bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-18 10:30:01.014139	1
137	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Riz bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-19 13:18:47.977494	1
138	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Riz bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-20 08:09:13.382249	1
139	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Riz bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-21 00:01:35.0902	1
140	{"grilladesMidi":["Turkey escalope with basil and tomato sauce","Honey and Mustard Chicken Stir-Fry","Piece of beef","Crustless quiche with vegetables and raw ham"],"migrateurs":[],"cibo":["Buns and rolls with falafel, stewed onions and carrots"],"accompMidi":["Organic rice","Stewed green cabbage"],"grilladesSoir":[],"accompSoir":[]}	2025-04-21 13:01:50.072149	2
141	{"grilladesMidi":["Escalope de dinde au basilic sauce tomate","Sauté de poulet Miel et moutarde","Pièce de boeuf","Quiche sans pâte aux légumes et jambon cru"],"migrateurs":[],"cibo":["Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Riz bio","Etuvé de choux verts"],"grilladesSoir":[],"accompSoir":[]}	2025-04-22 06:08:25.629588	1
142	{"grilladesMidi":["Oeuf brouillé à la Portuguaise","Saucisse fumé"],"migrateurs":["Brochette de dinde citron basilic"],"cibo":["Galette de sarasin printanière de légumes et emmentale","Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Haricots verts","Riz bio","Pâte bio","Haricots vert"],"grilladesSoir":[],"accompSoir":[]}	2025-04-22 08:50:01.045765	1
143	{"grilladesMidi":["Portuguese Scrambled Eggs","Smoked sausage"],"migrateurs":["Lemon and basil turkey skewer"],"cibo":["Spring buckwheat pancake with vegetables and Emmental cheese","Buns and rolls with falafel, stewed onions and carrots"],"accompMidi":["Green beans","Organic rice","Organic pasta","Green beans"],"grilladesSoir":[],"accompSoir":[]}	2025-04-22 10:01:25.885172	2
144	{"grilladesMidi":["Oeuf brouillé à la Portuguaise","Jambon grillé"],"migrateurs":["Brochette de dinde citron basilic"],"cibo":["Galette de sarasin printanière de légumes et emmentale","Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Haricots verts","Riz bio","Pâte bio","Haricots vert"],"grilladesSoir":[],"accompSoir":[]}	2025-04-22 10:10:00.935982	1
145	{"grilladesMidi":["Portuguese Scrambled Eggs","Grilled ham"],"migrateurs":["Lemon and basil turkey skewer"],"cibo":["Spring buckwheat pancake with vegetables and Emmental cheese","Buns and rolls with falafel, stewed onions and carrots"],"accompMidi":["Green beans","Organic rice","Organic pasta","Green beans"],"grilladesSoir":[],"accompSoir":[]}	2025-04-22 11:48:27.187065	2
146	{"grilladesMidi":["Oeuf brouillé à la Portuguaise","Jambon grillé"],"migrateurs":["Brochette de dinde citron basilic"],"cibo":["Galette de sarasin printanière de légumes et emmentale","Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Haricots verts","Riz bio","Pâte bio","Haricots vert"],"grilladesSoir":["Bolognaise Végétarienne","Bolognaise au boeuf sauce tartare"],"accompSoir":["Spaghetti","Carotte bio"]}	2025-04-22 15:00:01.033942	1
147	{"grilladesMidi":["Oeuf brouillé à la Portuguaise","Jambon grillé"],"migrateurs":["Brochette de dinde citron basilic"],"cibo":["Galette de sarasin printanière de légumes et emmentale","Bun's and roll au falafel compoté d'oignons et carottes"],"accompMidi":["Haricots verts","Riz bio","Pâte bio","Haricots vert"],"grilladesSoir":["Bolognaise Végétarienne","Bolognaise au boeuf sauce tartare"],"accompSoir":["Spaghetti","Carotte bio"]}	2025-04-23 08:15:25.10519	1
148	{"grilladesMidi":["Filet d'Eglefin en croûte d'olive","Aiguillette de poulet aux poivrons"],"migrateurs":["Côte de porc curcuma et gingembre, sauce chien"],"cibo":["Galette de sarasin printanière de légumes et emmentale","Oeufs brouillés à la Portuguaise"],"accompMidi":["Gratin de pomme de terre","Riz bio","Pâte bio","Carotte bio","gratin de pomme de terre"],"grilladesSoir":[],"accompSoir":[]}	2025-04-23 08:50:00.928298	1
149	{"grilladesMidi":["Haddock fillet in olive crust","Chicken Strips with Peppers"],"migrateurs":["Pork chop with turmeric and ginger, dog sauce"],"cibo":["Spring buckwheat pancake with vegetables and Emmental cheese","Portuguese Scrambled Eggs"],"accompMidi":["Potato gratin","Organic rice","Organic pasta","Organic carrot","potato gratin"],"grilladesSoir":[],"accompSoir":[]}	2025-04-23 10:11:06.72765	2
150	{"grilladesMidi":["Filet d'Eglefin en croûte d'olive","Aiguillette de poulet aux poivrons"],"migrateurs":["Côte de porc curcuma et gingembre, sauce chien"],"cibo":["Blanquette de haricots blanc"],"accompMidi":["Gratin de pomme de terre","Riz bio","Pâte bio","Carotte bio","gratin de pomme de terre"],"grilladesSoir":[],"accompSoir":[]}	2025-04-23 10:40:00.997412	1
151	{"grilladesMidi":["Saucisse au pesto","Cari de poulet","Moules Marinières"],"migrateurs":[],"cibo":["Blanquette de haricots blanc"],"accompMidi":["Semoule bio","Purée de pois cassé bio","Gratin de poireaux"],"grilladesSoir":[],"accompSoir":[]}	2025-04-24 08:30:00.946865	1
152	{"grilladesMidi":["Saucisse au pesto","Cari de poulet","Moules Marinières"],"migrateurs":[],"cibo":["Wok de brocolis et champignons et nouilles de riz"],"accompMidi":["Semoule bio","Purée de pois cassé bio","Gratin de poireaux"],"grilladesSoir":[],"accompSoir":[]}	2025-04-24 10:20:00.95905	1
153	{"grilladesMidi":["Viennoise de colin Sauce tartare","Boulette au boeuf coriandre et sésames","Parmentier de dinde"],"migrateurs":[],"cibo":["Curry de pois chiches et courgettes","Wok de brocolis et champignons et nouilles de riz"],"accompMidi":["Haricots plats persillées","Blé bio","Purée de pois cassé bio"],"grilladesSoir":[],"accompSoir":[]}	2025-04-25 10:46:47.034531	1
154	{"grilladesMidi":["Viennoise de colin Sauce tartare","Boulette au boeuf coriandre et sésames","Parmentier de dinde"],"migrateurs":[],"cibo":["Curry de pois chiches et courgettes"],"accompMidi":["Haricots plats persillées","Blé bio","Purée de pois cassé bio"],"grilladesSoir":[],"accompSoir":[]}	2025-04-25 11:00:00.937916	1
155	{"grilladesMidi":["Viennoise de colin Sauce tartare","Boulette au boeuf coriandre et sésames","Parmentier de dinde"],"migrateurs":[],"cibo":["Curry de pois chiches et courgettes"],"accompMidi":["Haricots plats persillées","Blé bio","Purée de pois cassé bio"],"grilladesSoir":[],"accompSoir":[]}	2025-04-26 04:23:21.316007	1
156	{"grilladesMidi":["Viennoise de colin Sauce tartare","Boulette au boeuf coriandre et sésames","Parmentier de dinde"],"migrateurs":[],"cibo":["Curry de pois chiches et courgettes"],"accompMidi":["Haricots plats persillées","Blé bio","Purée de pois cassé bio"],"grilladesSoir":[],"accompSoir":[]}	2025-04-27 02:15:57.928272	1
157	{"grilladesMidi":["FIlet d'Eiglegin beurre nantais","Gnocchi à la carbonara"],"migrateurs":["Poulet au lait de coco et coriandre"],"cibo":["Curry de pois chiches et courgettes","Pomme de terre au four, fondue de lentilles et poireaux et camembert"],"accompMidi":["Blé bio","Riz bio","Haricots plats","Epinards bio à la crème"],"grilladesSoir":["Vol au vent de champignons","Calamars à la romaine sauce tartare"],"accompSoir":["Flan de carottes bio","Farfalles"]}	2025-04-28 07:20:00.962423	1
158	{"grilladesMidi":["Lasagne au saumon et épinards bio","Echine de porc braisée"],"migrateurs":["Couscous au merguez"],"cibo":["Poivron farci boulgour et ricotta","Pomme de terre au four, fondue de lentilles et poireaux et camembert"],"accompMidi":["Purée de pommes de terre bio","Riz bio","Poếlée de carottes et choux"],"grilladesSoir":["Pizza 3 fromages","Colombo de dinde"],"accompSoir":["Brocolis","Riz bio"]}	2025-04-29 07:30:00.943804	1
159	{"grilladesMidi":["Lasagne au saumon et épinards bio","Echine de porc braisée"],"migrateurs":["Couscous au merguez"],"cibo":["Poivron farci boulgour et ricotta","Pomme de terre au four, fondue de lentilles et poireaux et camembert"],"accompMidi":["Purée de pommes de terre bio","Poếlée de carottes et choux"],"grilladesSoir":["Pizza 3 fromages","Colombo de dinde"],"accompSoir":["Brocolis","Riz bio"]}	2025-04-29 08:20:00.975798	1
160	{"grilladesMidi":["Lasagne au saumon et épinards bio","Echine de porc braisée"],"migrateurs":["Couscous au merguez"],"cibo":["Pomme de terre au four, fondue de lentilles et poireaux et camembert"],"accompMidi":["Purée de pommes de terre bio","Poếlée de carottes et choux"],"grilladesSoir":["Pizza 3 fromages","Colombo de dinde"],"accompSoir":["Brocolis","Riz bio"]}	2025-04-29 09:30:00.909247	1
161	{"grilladesMidi":["Lasagne au saumon et épinards bio"],"migrateurs":["Sauté de poulet Teriayi"],"cibo":["Poivron farci boulgour et ricotta","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Haricots beurre à la tomate","Riz bio","Purée de pommes de terre bio"],"grilladesSoir":["Cuisse de poulet rôtie","Pizza 3 fromages"],"accompSoir":["Epinards bio à la crème","Grenailles persillées"]}	2025-04-30 07:20:00.990975	1
162	{"grilladesMidi":["Lasagne au saumon et épinards bio","Filet de lieu fumé Waterzoï"],"migrateurs":["Sauté de poulet Teriayi"],"cibo":["Poivron farci boulgour et ricotta","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Haricots beurre à la tomate","Riz bio","Purée de pommes de terre bio"],"grilladesSoir":["Cuisse de poulet rôtie","Pizza 3 fromages"],"accompSoir":["Epinards bio à la crème","Grenailles persillées"]}	2025-04-30 09:10:00.965117	1
163	{"grilladesMidi":["Lasagne au saumon et épinards bio","Filet de lieu fumé Waterzoï"],"migrateurs":["Sauté de poulet Teriyaki"],"cibo":["Poivron farci boulgour et ricotta","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Haricots beurre à la tomate","Riz bio","Purée de pommes de terre bio"],"grilladesSoir":["Cuisse de poulet rôtie","Pizza 3 fromages"],"accompSoir":["Epinards bio à la crème","Grenailles persillées"]}	2025-04-30 10:00:00.964809	1
164	{"grilladesMidi":["Lasagne au saumon et épinards bio","Filet de lieu fumé Waterzoï"],"migrateurs":["Sauté de poulet Teriyaki"],"cibo":["Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Haricots beurre à la tomate","Riz bio","Purée de pommes de terre bio"],"grilladesSoir":["Cuisse de poulet rôtie","Pizza 3 fromages"],"accompSoir":["Epinards bio à la crème","Grenailles persillées"]}	2025-04-30 11:10:00.934389	1
165	{"grilladesMidi":["Organic Salmon and Spinach Lasagna","Smoked fillet of waterzoï"],"migrateurs":["Teriyaki Chicken Stir-Fry"],"cibo":["Zucchini, mozzarella and arugula bruschetta"],"accompMidi":["Butter beans with tomato","Organic rice","Organic mashed potatoes"],"grilladesSoir":["Roasted chicken thigh","Pizza 3 fromages"],"accompSoir":["Organic creamed spinach","Marbled grenailles"]}	2025-05-01 22:02:34.142727	2
166	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Filet de lieu fumé Waterzoï","Escalope de porc à la diable"],"migrateurs":[],"cibo":["Tajine de légumes aux fruits secs et boulettes de pois chiches","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-02 07:10:00.987595	1
167	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Filet de lieu fumé Waterzoï"],"migrateurs":[],"cibo":["Tajine de légumes aux fruits secs et boulettes de pois chiches","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-02 09:40:00.936747	1
168	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Escalope de porc à la diable","Filet d'eiglefin bouillon à l'asiatique"],"migrateurs":[],"cibo":["Tajine de légumes aux fruits secs et boulettes de pois chiches","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-02 10:20:00.999554	1
169	{"grilladesMidi":["Turkey Sauté with Olives and Peppers","Deviled Pork Escalope","Haddock fillet in Asian broth"],"migrateurs":[],"cibo":["Vegetable tagine with dried fruit and chickpea balls","Zucchini, mozzarella and arugula bruschetta"],"accompMidi":["Granailles with onions","Organic rice","Organic carrots with cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-02 21:24:51.050172	2
170	{"grilladesMidi":["Pavo salteado con aceitunas y pimientos","Escalope de cerdo endiablado","Filete de eglefino en caldo asiático"],"migrateurs":[],"cibo":["Tajín de verduras con frutos secos y bolitas de garbanzos","Bruschetta de calabacín, mozzarella y rúcula"],"accompMidi":["Granailles con cebollas","Arroz orgánico","Zanahorias orgánicas con comino"],"grilladesSoir":[],"accompSoir":[]}	2025-05-02 21:38:27.543418	3
171	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Escalope de porc à la diable","Filet d'eiglefin bouillon à l'asiatique"],"migrateurs":[],"cibo":["Tajine de légumes aux fruits secs et boulettes de pois chiches","Bruschetta courgettes, mozzarella et roquette"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-03 00:09:34.918182	1
172	{"grilladesMidi":["Turkey Sauté with Olives and Peppers","Deviled Pork Escalope","Haddock fillet in Asian broth"],"migrateurs":[],"cibo":["Vegetable tagine with dried fruit and chickpea balls","Zucchini, mozzarella and arugula bruschetta"],"accompMidi":["Onion grenailles","Organic rice","Organic carrots with cumin"],"grilladesSoir":[],"accompSoir":[]}	2025-05-03 14:54:52.570204	2
173	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Pizza royale Waterzoï","Encornet basquaise à la diable"],"migrateurs":["Aiguilette de poulet à la citronnelle"],"cibo":["Mac and cheese brocolis cheddar"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin","Haricots verts en persillade"],"grilladesSoir":["Rôti de porc","Galette de blé noir aux légumes et chèvre"],"accompSoir":["Brocolis","Gratin de pommes de terre"]}	2025-05-05 08:00:00.940192	1
174	{"grilladesMidi":["Sauté de dinde aux olives et poivrons","Pizza royale","Encornet basquaise"],"migrateurs":["Aiguilette de poulet à la citronnelle"],"cibo":["Mac and cheese brocolis cheddar"],"accompMidi":["Grenailles aux oignons","Riz bio","Carottes bio au cumin","Haricots verts en persillade"],"grilladesSoir":["Rôti de porc","Galette de blé noir aux légumes et chèvre"],"accompSoir":["Brocolis","Gratin de pommes de terre"]}	2025-05-05 09:10:00.962567	1
175	{"grilladesMidi":["Tarte poireaux, champignons à l'emmental","Rôti de dinde, sauce champignons","Encornet basquaise"],"migrateurs":["Aiguilette de poulet à la citronnelle","Sauté de porc à la Madrilène","Rôti de porc"],"cibo":["Ravioles aux fromages gratinés","Autre plat végétarien disponible pôle 1 Grillade Tarte poireaux, champignons"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-06 07:20:00.957181	1
176	{"grilladesMidi":["Tarte poireaux, champignons à l'emmental","Rôti de dinde, sauce champignons","Encornet basquaise"],"migrateurs":["Aiguilette de poulet à la citronnelle","Sauté de porc à la Madrilène","Rôti de porc sauce champignons"],"cibo":["Ravioles aux fromages gratinés","Autre plat végétarien disponible pôle 1 Grillade Tarte poireaux, champignons"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-06 08:00:00.969826	1
177	{"grilladesMidi":["Tarte poireaux, champignons à l'emmental","Rôti de dinde, sauce champignons","Encornet basquaise"],"migrateurs":["Aiguilette de poulet à la citronnelle","Sauté de porc à la Madrilène (pois chiches et tomates)","Rôti de porc sauce champignons"],"cibo":["Ravioles aux fromages gratinés","Autre plat végétarien disponible pôle 1 Grillade Tarte poireaux, champignons"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-06 09:10:01.036297	1
178	{"grilladesMidi":["Rôti de porc pois chiche et tomates","Encornet basquaise"],"migrateurs":["Sauté de porc à la Madrilène (pois chiches et tomates)"],"cibo":["Tortelloni ricotta épinard","Autre plat végétarien disponible pôle 1 Grillade Tarte poireaux, champignons"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-06 10:50:00.975868	1
179	{"grilladesMidi":["Rôti de porc pois chiche et tomates","Encornet basquaise"],"migrateurs":["Sauté de porc à la Madrilène (pois chiches et tomates)"],"cibo":["Tortelloni ricotta épinard"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-06 11:00:01.019447	1
180	{"grilladesMidi":["Rôti de porc pois chiche et tomates","Encornet basquaise","Plat disponible dans le frigo rose bonne action pour les étudiants"],"migrateurs":["Sauté de porc à la Madrilène (pois chiches et tomates)"],"cibo":["Tortelloni ricotta épinard"],"accompMidi":["Semoule bio","Riz bio","Carottes bio fondantes","Gratin de pâtes"],"grilladesSoir":["Bruschetta pizza","Parmentier végétal","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Julienne de légumes","Pâtes bio"]}	2025-05-07 00:10:00.921103	1
181	{"grilladesMidi":["Marmite de poisson et crustacés aux haricots blancs","Sauté de porc à la Madrilène (pois chiches et tomates)","Plat disponible dans le frigo rose bonne action pour les étudiants"],"migrateurs":["Cuisse de poulet sauce barbecue"],"cibo":["Bruschetta aux légumes","Lasagne aux blettes et épinards"],"accompMidi":["Pâtes bio","Riz bio","Fenouil rôti"],"grilladesSoir":["Blanc de poulet, sauce crème","Chachouka poireaux, épinards","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Purée de patates douces","Poêlée de courgettes et poivrons"]}	2025-05-07 07:40:01.001327	1
182	{"grilladesMidi":["Marmite de poisson et crustacés aux haricots blancs","Plat disponible dans le frigo rose bonne action pour les étudiants"],"migrateurs":["Cuisse de poulet sauce barbecue"],"cibo":["Bruschetta aux légumes","Lasagne aux blettes et épinards"],"accompMidi":["Pâtes bio","Riz bio","Fenouil rôti"],"grilladesSoir":["Blanc de poulet, sauce crème","Chachouka poireaux, épinards","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Purée de patates douces","Poêlée de courgettes et poivrons"]}	2025-05-07 10:20:00.963637	1
183	{"grilladesMidi":["Marmite de poisson et crustacés aux haricots blancs","Côtes de porc sauce champignons","Plat disponible dans le frigo rose bonne action pour les étudiants"],"migrateurs":["Cuisse de poulet sauce barbecue"],"cibo":["Lasagne aux blettes et épinards"],"accompMidi":["Pâtes bio","Riz bio","Fenouil rôti"],"grilladesSoir":["Blanc de poulet, sauce crème","Chachouka poireaux, épinards","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Purée de patates douces","Poêlée de courgettes et poivrons"]}	2025-05-07 10:40:00.944493	1
184	{"grilladesMidi":["Boeine","Poulet rôti"],"migrateurs":[],"cibo":["Cuux et riz","Lasas"],"accompMidi":["Poeur","Petio","Feti"],"grilladesSoir":[],"accompSoir":[]}	2025-05-09 07:20:00.931541	1
193	{"grilladesMidi":["Boeuf haché à la cubaine","Poulet rôti"],"migrateurs":[],"cibo":["Curry aux 2 lentilles, poireaux, épinards et riz","Lasagne aux blettes et épinards"],"accompMidi":["Pomme de terre bio vapeur","Petits pois carottes bio","Fenouil rôti"],"grilladesSoir":[],"accompSoir":[]}	2025-05-10 11:41:56.475605	1
194	{"grilladesMidi":["Kubanisches Hackfleisch","Brathähnchen"],"migrateurs":[],"cibo":["Curry mit 2 Linsen, Lauch, Spinat und Reis","Lasagne mit Mangold und Spinat"],"accompMidi":["Bio-Dampfkartoffeln","Bio-Erbsen und Karotten","Gerösteter Fenchel"],"grilladesSoir":[],"accompSoir":[]}	2025-05-10 12:06:17.490361	4
196	{"grilladesMidi":["Côte de porc BBC sauce bourguignonne","Chipolatas aux herbes grillées"],"migrateurs":["Boeuf haché à la cubaine","Raviolini aux fromages sauce tomate"],"cibo":["Dahl de carottes au curcuma lentilles et pois chiches"],"accompMidi":["Grenailles rôties","Haricots beurre vapeur","Riz bio","Haricots beurres vapeur"],"grilladesSoir":["Cuisse de poulet rôti","Parmentier aux légumes","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Poêlée de légumes verts","Riz bio"]}	2025-05-12 07:11:56.457277	1
197	{"grilladesMidi":["Côte de porc BBC sauce bourguignonne","Chipolatas aux herbes grillées"],"migrateurs":["Raviolini aux fromages sauce tomate"],"cibo":["Dahl de carottes au curcuma lentilles et pois chiches","Autre plat végétarien disponible pôle 2 Cuistots Migrateurs Raviolini aux fromages"],"accompMidi":["Grenailles rôties","Haricots beurre vapeur","Riz bio","Haricots beurres vapeur"],"grilladesSoir":["Cuisse de poulet rôti","Parmentier aux légumes","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Poêlée de légumes verts","Riz bio"]}	2025-05-12 12:49:29.11574	1
198	{"grilladesMidi":["Echine de porc à la Toscane (ail, romarin)","Filet de colin, sauce vierge"],"migrateurs":["Boeuf à la chinoise (pousses de bambou, soja)"],"cibo":["Conchiglioni, bolognaise de légumes"],"accompMidi":["Pâtes bio","Fondue de poireaux","Riz bio"],"grilladesSoir":["Chili con carne","Quiche aux champignons","Plat disponible dans le frigo rose bonne action pour les étudiants"],"accompSoir":["Ecrasée de carottes bio","Pâtes bio"]}	2025-05-13 07:08:51.76474	1
199	{"grilladesMidi":["Tourte brocolis au Gorgonzola (ail, romarin)","Escalope de dinde sauce Normande"],"migrateurs":["Lieu noir à la Majorquine","Chili con carne"],"cibo":["Curry Laska","Conchiglioni, bolognaise de légumes","Autre plat végétarien disponible pôle 1 Grillade Tourte brocolis, gorgonzola"],"accompMidi":["Poêlée d'aubergines courgettes et fenouil","Riz bio"],"grilladesSoir":["Sauté de porc à la tomate","Tourte de légumes au bleu"],"accompSoir":["Petits pois aux oignons","Riz bio"]}	2025-05-14 06:58:51.818823	1
200	{"grilladesMidi":["西兰花派配戈贡佐拉奶酪（大蒜、迷迭香）","火鸡肉排佐诺曼底酱"],"migrateurs":["马略卡风味黑狭鳕","西兰花戈贡佐拉派"],"cibo":["咖喱之爱","贝壳面，蔬菜博洛尼亚酱","其他素食可供选择 1 烤西兰花派，戈贡佐拉奶酪"],"accompMidi":["煎茄子、西葫芦和茴香","有机大米"],"grilladesSoir":["番茄炒猪肉","蓝纹奶酪蔬菜派"],"accompSoir":["豌豆配洋葱","有机大米"]}	2025-05-14 13:13:54.481517	8
201	{"grilladesMidi":["Brokkoli-Pie mit Gorgonzola (Knoblauch, Rosmarin)","Putenschnitzel mit Normandiesauce"],"migrateurs":["Schwarzer Seelachs nach mallorquinischer Art","Brokkoli-Gorgonzola-Kuchen"],"cibo":["Curry-Liebe","Conchiglioni, Gemüse-Bolognese","Andere vegetarische Gerichte verfügbar Pole 1 Grill Brokkoli Pie, Gorgonzola"],"accompMidi":["Gebratene Auberginen, Zucchini und Fenchel","Bio-Reis"],"grilladesSoir":["Gebratenes Schweinefleisch mit Tomaten","Gemüsekuchen mit Blauschimmelkäse"],"accompSoir":["Erbsen mit Zwiebeln","Bio-Reis"]}	2025-05-14 13:14:21.749158	4
202	{"grilladesMidi":["Tortelloni ricotta épinards","Matelote de coin"],"migrateurs":["Sauté de poulet épicé des Iles Fidji"],"cibo":["Curry Laska","Autre plat végétarien disponible pôle 1 Grillade Tortelloni ricotta épinards"],"accompMidi":["Pomme vapeur bio","Courgettes à l'ail et sauge","Riz bio","Pommes vapeur bio"],"grilladesSoir":["Filet de lieu sauce bonne femme","Bruschetta courgettes, poivrons et olives"],"accompSoir":["Semoule bio","Salade verte Gratin de blettes"]}	2025-05-15 07:46:15.623421	1
203	{"grilladesMidi":["Tortelloni ricotta épinards","Matelote de colin"],"migrateurs":["Sauté de poulet épicé des Iles Fidji"],"cibo":["Curry Laska","Autre plat végétarien disponible pôle 1 Grillade Tortelloni ricotta épinards"],"accompMidi":["Pomme vapeur bio","Courgettes à l'ail et sauge","Riz bio","Pommes vapeur bio"],"grilladesSoir":["Filet de lieu sauce bonne femme","Bruschetta courgettes, poivrons et olives"],"accompSoir":["Semoule bio","Salade verte Gratin de blettes"]}	2025-05-15 10:00:31.121089	1
204	{"grilladesMidi":["Aiguilette de poulet aux olives","Rôti de dinde sauce échalote au vin blanc"],"migrateurs":["Fricassée de merguez aux fruits secs"],"cibo":["Parmentier aux patates douce et lentilles"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio"],"grilladesSoir":[],"accompSoir":[]}	2025-05-16 06:26:15.606073	1
205	{"grilladesMidi":["Chicken Strips with Olives","Roast turkey with shallot and white wine sauce"],"migrateurs":["Fricassee of merguez with dried fruits"],"cibo":["Sweet potato and lentil parmentier"],"accompMidi":["Organic semolina","Parsley flat beans","Organic rice"],"grilladesSoir":[],"accompSoir":[]}	2025-05-16 14:23:47.681459	2
206	{"grilladesMidi":["Aiguilette de poulet aux olives","Rôti de dinde sauce échalote au vin blanc"],"migrateurs":["Fricassée de merguez aux fruits secs"],"cibo":["Parmentier aux patates douce et lentilles"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio"],"grilladesSoir":[],"accompSoir":[]}	2025-05-17 09:57:17.760492	1
207	{"grilladesMidi":["Chicken Strips with Olives","Roast turkey with shallot and white wine sauce"],"migrateurs":["Fricassee of merguez with dried fruits"],"cibo":["Sweet potato and lentil parmentier"],"accompMidi":["Organic semolina","Parsley flat beans","Organic rice"],"grilladesSoir":[],"accompSoir":[]}	2025-05-17 11:15:21.634552	2
208	{"grilladesMidi":["Chicken Strips with Olives","Roast turkey with shallot and white wine sauce"],"migrateurs":["Fricassee of merguez with dried fruits"],"cibo":["Sweet potato and lentil parmentier"],"accompMidi":["Organic semolina","Parsley flat beans","Organic rice"],"grilladesSoir":[],"accompSoir":[]}	2025-05-18 16:06:15.588587	2
209	{"grilladesMidi":["Aiguilette de poulet aux olives","Rôti de dinde sauce échalote au vin blanc","Bruschetta façon pizza","Parmentier de boudin noir"],"migrateurs":["Escalope de dinde grillé mariné au citron et parika"],"cibo":["Parmentier aux patates douce et lentilles","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio","Patates douce rôti","Choux fleurs"],"grilladesSoir":["Sauté de porc basquaise sauce bonne femme","Lasagne végétarienne poivrons et olives"],"accompSoir":["Riz bio","Haricots verts Gratin de blettes"]}	2025-05-19 08:38:06.158613	1
210	{"grilladesMidi":["Rôti de dinde sauce échalote au vin blanc","Bruschetta façon pizza","Parmentier de boudin noir"],"migrateurs":["Escalope de dinde grillé mariné au citron et parika","Aiguilette de poulet aux olives"],"cibo":["Parmentier aux patates douce et lentilles","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio","Patates douce rôti","Choux fleurs"],"grilladesSoir":["Sauté de porc basquaise sauce bonne femme","Lasagne végétarienne poivrons et olives"],"accompSoir":["Riz bio","Haricots verts Gratin de blettes"]}	2025-05-19 08:42:47.924559	1
226	{"grilladesMidi":["Conchigliette forestier au poulet","Quiche courgettes et jambon cru","Potrine de porc rôtie"],"migrateurs":["Calamars à la Sétoise"],"cibo":["Burritos de légumes","Gnocchi tomates et mozzarella"],"accompMidi":["Pâtes bio","Carottes bio vichy"],"grilladesSoir":["Carbonara","Quiche épinards et chèvre"],"accompSoir":["Choux fleurs","Spaghetti bio"]}	2025-05-21 08:16:17.377573	1
211	{"grilladesMidi":["Roast turkey with shallot and white wine sauce","Pizza-style bruschetta","Black pudding parmentier"],"migrateurs":["Grilled turkey escalope marinated with lemon and parika","Chicken Strips with Olives"],"cibo":["Sweet potato and lentil parmentier","Eggplant risotto with tomato and pine nuts"],"accompMidi":["Organic semolina","Parsley flat beans","Organic rice","Roasted sweet potatoes","Cauliflowers"],"grilladesSoir":["Basque-style pork sauté with Bonne Femme sauce","Vegetarian lasagna with peppers and olives"],"accompSoir":["Organic rice","Green beans Swiss chard gratin"]}	2025-05-19 09:11:19.392169	2
212	{"grilladesMidi":["Bruschetta façon pizza","Parmentier de boudin noir"],"migrateurs":["Escalope de dinde grillé mariné au citron et parika","Aiguilette de poulet aux olives"],"cibo":["Parmentier aux patates douce et lentilles","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio","Patates douce rôti","Choux fleurs"],"grilladesSoir":["Sauté de porc basquaise sauce bonne femme","Lasagne végétarienne poivrons et olives"],"accompSoir":["Riz bio","Haricots verts Gratin de blettes"]}	2025-05-19 15:02:47.693092	1
213	{"grilladesMidi":["披萨式意式烤面包","黑布丁帕门蒂尔"],"migrateurs":["柠檬和帕里卡腌制的烤火鸡肉片","橄榄鸡柳"],"cibo":["红薯扁豆帕门蒂尔","茄子烩饭配番茄和松子"],"accompMidi":["有机粗面粉","欧芹扁豆","有机大米","烤红薯","花椰菜"],"grilladesSoir":["巴斯克风味猪肉炒配 Bonne Femme 酱","素食千层面配辣椒和橄榄"],"accompSoir":["土豆泥","四季豆瑞士甜菜焗烤菜"]}	2025-05-19 16:57:24.405119	8
214	{"grilladesMidi":["Bruschetta estilo pizza","Parmentier de morcela"],"migrateurs":["Escalope de peru grelhado marinado com limão e parika","Tiras de Frango com Azeitonas"],"cibo":["Parmentier de batata-doce e lentilha","Risoto de berinjela com tomate e pinhões"],"accompMidi":["Sêmola orgânica","Feijão-chato de salsa","Arroz orgânico","Batatas-doces assadas","Couve-flor"],"grilladesSoir":["Salteado de porco à moda basca com molho Bonne Femme","Lasanha vegetariana com pimentões e azeitonas"],"accompSoir":["Purê de batatas","Gratinado de feijão verde e acelga"]}	2025-05-19 17:24:30.16803	6
215	{"grilladesMidi":["Bruschetta façon pizza","Parmentier de boudin noir"],"migrateurs":["Escalope de dinde grillé mariné au citron et parika","Aiguilette de poulet aux olives"],"cibo":["Parmentier aux patates douce et lentilles","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio","Patates douce rôti","Choux fleurs"],"grilladesSoir":["Sauté de porc basquaise sauce bonne femme","Lasagne végétarienne poivrons et olives"],"accompSoir":["Purée de pomme de terre","Haricots verts Gratin de blettes"]}	2025-05-19 23:15:22.914628	1
216	{"grilladesMidi":["Bruschetta façon pizza","Parmentier de boudin noir"],"migrateurs":["Escalope de dinde grillé mariné au citron et parika","Aiguilette de poulet aux olives"],"cibo":["Parmentier aux patates douce et lentilles","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Semoule bio","Haricots plats persillés","Riz bio","Patates douce rôti","Choux fleurs"],"grilladesSoir":["Sauté de porc basquaise sauce bonne femme","Lasagne végétarienne poivrons et olives"],"accompSoir":["Purée de pomme de terre","Haricots verts Gratin de blettes"]}	2025-05-20 03:18:45.20191	1
217	{"grilladesMidi":["Pizza-style bruschetta","Black pudding parmentier"],"migrateurs":["Grilled turkey escalope marinated with lemon and parika","Chicken Strips with Olives"],"cibo":["Sweet potato and lentil parmentier","Eggplant risotto with tomato and pine nuts"],"accompMidi":["Organic semolina","Parsley flat beans","Organic rice","Roasted sweet potatoes","Cauliflowers"],"grilladesSoir":["Basque-style pork sauté with Bonne Femme sauce","Vegetarian lasagna with peppers and olives"],"accompSoir":["Mashed potatoes","Green beans Swiss chard gratin"]}	2025-05-20 07:10:36.914085	2
218	{"grilladesMidi":["Escalope de dinde grillé mariné au citron et parika","Quiche courgettes et jambon cru","Parmentier de boudin noir"],"migrateurs":["Colombo de porc"],"cibo":["Burritos de légumes","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Riz bio","Petits pois carottes"],"grilladesSoir":["Pilon de poulet mariné et rôti","Jambalaya de légumes en crumble d'avoine"],"accompSoir":["Potatoes","Gratin de courgettes"]}	2025-05-20 07:40:25.697604	1
219	{"grilladesMidi":["Grilled turkey escalope marinated with lemon and parika","Zucchini and raw ham quiche","Black pudding parmentier"],"migrateurs":["Pork Colombo"],"cibo":["Vegetable burritos","Eggplant risotto with tomato and pine nuts"],"accompMidi":["Organic rice","green peas and carrots"],"grilladesSoir":["Marinated and roasted chicken drumstick","Vegetable Jambalaya with Oat Crumble"],"accompSoir":["Potatoes","Zucchini gratin"]}	2025-05-20 07:42:34.852547	2
220	{"grilladesMidi":["Gegrilltes Putenschnitzel mariniert mit Zitrone und Parika","Zucchini-Rohschinken-Quiche","Blutwurst-Parmentier"],"migrateurs":["Schweinefleisch Colombo"],"cibo":["Gemüse-Burritos","Auberginenrisotto mit Tomaten und Pinienkernen"],"accompMidi":["Bio-Reis","grüne Erbsen und Karotten"],"grilladesSoir":["Marinierte und gebratene Hähnchenkeule","Gemüse-Jambalaya mit Hafer-Crumble"],"accompSoir":["Kartoffeln","Zucchini-Gratin"]}	2025-05-20 08:46:23.448108	4
221	{"grilladesMidi":["Quiche courgettes et jambon cru","Parmentier de boudin noir"],"migrateurs":["Colombo de porc"],"cibo":["Burritos de légumes","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Riz bio","Petits pois carottes"],"grilladesSoir":["Pilon de poulet mariné et rôti","Jambalaya de légumes en crumble d'avoine"],"accompSoir":["Potatoes","Gratin de courgettes"]}	2025-05-20 12:09:37.661504	1
222	{"grilladesMidi":["Quiche de abobrinha e presunto cru","Frango Colombo"],"migrateurs":["Porco Colombo"],"cibo":["Burritos de vegetais","Risoto de berinjela com tomate e pinhões"],"accompMidi":["Arroz orgânico","ervilhas e cenouras"],"grilladesSoir":["Coxa de frango marinada e assada","Jambalaya de vegetais com crumble de aveia"],"accompSoir":["Arroz orgânico","Gratinado de couve-flor"]}	2025-05-20 17:43:29.480171	6
223	{"grilladesMidi":["Quiche de abobrinha e presunto cru","Frango Colombo"],"migrateurs":["Porco Colombo"],"cibo":["Burritos de vegetais","Risoto de berinjela com tomate e pinhões"],"accompMidi":["Arroz orgânico","ervilhas e cenouras"],"grilladesSoir":["Coxa de frango marinada e assada","Jambalaya de vegetais com crumble de aveia"],"accompSoir":["Arroz orgânico","Gratinado de couve-flor"]}	2025-05-21 06:51:56.357408	6
224	{"grilladesMidi":["Zucchini-Rohschinken-Quiche","Hühnchen Colombo"],"migrateurs":["Schweinefleisch Colombo"],"cibo":["Gemüse-Burritos","Auberginenrisotto mit Tomaten und Pinienkernen"],"accompMidi":["Bio-Reis","grüne Erbsen und Karotten"],"grilladesSoir":["Marinierte und gebratene Hähnchenkeule","Gemüse-Jambalaya mit Hafer-Crumble"],"accompSoir":["Bio-Reis","Blumenkohlgratin"]}	2025-05-21 07:27:55.4267	4
225	{"grilladesMidi":["Quiche courgettes et jambon cru","Colombo de poulet"],"migrateurs":["Colombo de porc"],"cibo":["Burritos de légumes","Risotto d'aubergines à la tomate et pignons de pins"],"accompMidi":["Riz bio","Petits pois carottes"],"grilladesSoir":["Pilon de poulet mariné et rôti","Jambalaya de légumes en crumble d'avoine"],"accompSoir":["Riz bio","Gratin de choux fleur"]}	2025-05-21 07:32:29.414202	1
227	{"grilladesMidi":["Forest conchigliette with chicken","Zucchini and raw ham quiche","Cheese Omelette"],"migrateurs":["Calamari from Sète","Roasted pork belly"],"cibo":["Gnocchi with tomatoes and mozzarella"],"accompMidi":["Organic pasta","Vichy organic carrots","Organic rice"],"grilladesSoir":["Carbonara","Spinach and goat cheese quiche"],"accompSoir":["Cauliflowers","Organic Spaghetti"]}	2025-05-21 10:49:53.251402	2
228	{"grilladesMidi":["Conchigliette forestier au poulet","Quiche courgettes et jambon cru","Omelette aux fromages"],"migrateurs":["Calamars à la Sétoise","Potrine de porc rôtie"],"cibo":["Gnocchi tomates et mozzarella"],"accompMidi":["Pâtes bio","Carottes bio vichy","Riz bio"],"grilladesSoir":["Carbonara","Quiche épinards et chèvre"],"accompSoir":["Choux fleurs","Spaghetti bio"]}	2025-05-21 11:13:33.968065	1
229	{"grilladesMidi":["Wald-Conchigliette mit Hühnchen","Zucchini-Rohschinken-Quiche","Käseomelett"],"migrateurs":["Calamari aus Sète","Gebratener Schweinebauch"],"cibo":["Gnocchi mit Tomaten und Mozzarella"],"accompMidi":["Bio-Nudeln","Vichy Bio-Karotten","Bio-Reis"],"grilladesSoir":["Carbonara","Spinat-Ziegenkäse-Quiche"],"accompSoir":["Blumenkohl","Bio-Spaghetti"]}	2025-05-21 11:15:51.805392	4
230	{"grilladesMidi":["Conchigliette forestier au poulet","Quiche courgettes et jambon cru","Omelette aux fromages"],"migrateurs":["Calamars à la Sétoise","Potrine de porc rôtie"],"cibo":["Gnocchi tomates et mozzarella"],"accompMidi":["Pâtes bio","Carottes bio vichy","Riz bio"],"grilladesSoir":["Carbonara","Quiche épinards et comté"],"accompSoir":["Courgettes","Farfalle"]}	2025-05-21 13:29:56.194824	1
231	{"grilladesMidi":["Wald-Conchigliette mit Hühnchen","Zucchini-Rohschinken-Quiche","Käseomelett"],"migrateurs":["Calamari aus Sète","Gebratener Schweinebauch"],"cibo":["Gnocchi mit Tomaten und Mozzarella"],"accompMidi":["Bio-Nudeln","Vichy Bio-Karotten","Bio-Reis"],"grilladesSoir":["Carbonara","Spinat-Comté-Quiche"],"accompSoir":["Zucchini","Schmetterlinge"]}	2025-05-21 13:04:19.189642	4
232	{"grilladesMidi":["Forest conchigliette with chicken","Zucchini and raw ham quiche","Cheese Omelette"],"migrateurs":["Calamari from Sète","Roasted pork belly"],"cibo":["Gnocchi with tomatoes and mozzarella"],"accompMidi":["Organic pasta","Vichy organic carrots","Organic rice"],"grilladesSoir":["Carbonara","Spinach and Comté Quiche"],"accompSoir":["Courgettes","Butterflies"]}	2025-05-21 16:45:27.728703	2
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.roles (id_roles, name, description) FROM stdin;
1	ADMIN	Global administrator with all privileges
2	NEWF	Transat user
3	VALIDATED	Transat user with validated email
4	BANNED	Transat user banned from the platform
5	VERIFYING	Transat user with email not yet verified
\.


--
-- Data for Name: rooms; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rooms (id_rooms, name, description, picture, location) FROM stdin;
\.


--
-- Data for Name: services; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.services (id_services, name, type, id_category) FROM stdin;
1	RESTAURANT	RESTAURANT	\N
2	TRAQ	TRAQ	\N
3	REAL_CAMPUS	REAL_CAMPUS	\N
\.


--
-- Data for Name: support; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.support (id_support, email, subject, message, creation_date, status, response) FROM stdin;
1	yohann.chavanel@imt-atlantique.net	Test Support Request	This is a test message for my support request	2025-05-12 15:05:46.871956	pending	\N
2	yohann.chavanel@imt-atlantique.net	Test Support Request	This is a test message for my support request	2025-05-12 15:34:57.229382	pending	\N
3	yohann.chavanel@imt-atlantique.net	Test Support Request 2	This is a test message for my support request	2025-05-12 15:45:17.401713	pending	\N
4	yohann.chavanel@imt-atlantique.net	Nouveau PB depuis l&#x27;&#x27;app	Description piti a marche	2025-05-12 15:46:36.819007	pending	\N
\.


--
-- Data for Name: support_media; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.support_media (id_support_media, id_support, id_files) FROM stdin;
\.


--
-- Data for Name: traq; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.traq (id_traq, name, disabled, limited, alcohol, out_of_stock, creation_date, picture, description, price, price_half, id_traq_types) FROM stdin;
1	Kasteel Rouge	f	f	8.00	f	2025-03-01 08:15:08.434983	https://jumpseller.s3.eu-west-1.amazonaws.com/store/elvino-cl/assets/Kasteel.png	Bière belge alliant la force d’une brune et la douceur de la cerise. Gourmande et légèrement sucrée, elle offre une belle rondeur en bouche.	3.00	0.00	1
2	Kerisac	f	f	5.50	f	2025-03-01 08:27:34.340495	https://www.tempetedelouest.fr/wp-content/uploads/2023/09/logo-kerisac.png	Cidre breton artisanal, équilibré entre douceur et acidité. Ses arômes fruités et sa fine effervescence en font un incontournable.	2.50	0.00	1
3	Filou	f	f	8.50	f	2025-03-01 08:30:46.365375	https://www.vanhonsebrouck.be/wp-content/uploads/external/70ed1eeb0e074191d0646ef029fa3223-1000x0-c-default.webp	Bière blonde belge aux notes maltées et épicées. Son équilibre entre douceur et amertume séduit les amateurs de caractère.	3.00	0.00	1
6	Chips Brets	f	f	0.00	t	2025-03-01 08:36:08.923629	https://www.bretagne-allerlei.de/images/manufacturers/logo_Brets.png	Différentes saveurs sont disponnibles en fonction des stocks.	0.00	0.00	3
5	Sprite	f	t	0.00	f	2025-03-01 08:34:50.371393	https://vignette.wikia.nocookie.net/logopedia/images/4/40/Sprite_logo_2019.svg/revision/latest/scale-to-width-down/2000?cb=20180312044311	Boisson pétillante sucrée goût citroné.	0.50	0.00	2
\.


--
-- Data for Name: traq_types; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.traq_types (id_traq_types, name) FROM stdin;
2	Soft
3	Snacks
4	Food
1	Alcohol
\.


--
-- Data for Name: washing_machines; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.washing_machines (id_washing_machines, type, broken) FROM stdin;
\.


--
-- Data for Name: washing_machines_reports; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.washing_machines_reports (email, id_washing_machines, report_date) FROM stdin;
\.


--
-- Name: articles_id_articles_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.articles_id_articles_seq', 1, false);


--
-- Name: articles_types_id_articles_types_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.articles_types_id_articles_types_seq', 1, false);


--
-- Name: bassine_id_bassine_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.bassine_id_bassine_seq', 1, false);


--
-- Name: caps_id_caps_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.caps_id_caps_seq', 1, false);


--
-- Name: clubs_id_clubs_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.clubs_id_clubs_seq', 1, false);


--
-- Name: clubs_types_id_clubs_types_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.clubs_types_id_clubs_types_seq', 1, false);


--
-- Name: events_id_events_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.events_id_events_seq', 1, false);


--
-- Name: files_id_files_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.files_id_files_seq', 45, true);


--
-- Name: games_id_games_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.games_id_games_seq', 1, false);


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 10, true);


--
-- Name: languages_id_languages_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.languages_id_languages_seq', 9, true);


--
-- Name: newf_id_newf_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.newf_id_newf_seq', 31, true);


--
-- Name: realcampus_friendships_id_friendship_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.realcampus_friendships_id_friendship_seq', 3, true);


--
-- Name: realcampus_post_reactions_id_reaction_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.realcampus_post_reactions_id_reaction_seq', 1, false);


--
-- Name: realcampus_posts_id_post_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.realcampus_posts_id_post_seq', 9, true);


--
-- Name: request_statistics_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.request_statistics_id_seq', 8823, true);


--
-- Name: restaurant_id_restaurant_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.restaurant_id_restaurant_seq', 232, true);


--
-- Name: roles_id_roles_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_roles_seq', 5, true);


--
-- Name: rooms_id_rooms_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.rooms_id_rooms_seq', 1, false);


--
-- Name: services_id_services_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.services_id_services_seq', 3, true);


--
-- Name: support_id_support_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.support_id_support_seq', 4, true);


--
-- Name: support_media_id_support_media_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.support_media_id_support_media_seq', 1, false);


--
-- Name: traq_id_traq_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.traq_id_traq_seq', 7, true);


--
-- Name: traq_types_id_traq_types_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.traq_types_id_traq_types_seq', 4, true);


--
-- Name: washing_machines_id_washing_machines_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.washing_machines_id_washing_machines_seq', 1, false);


--
-- Name: articles articles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id_articles);


--
-- Name: articles_types articles_types_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles_types
    ADD CONSTRAINT articles_types_name_key UNIQUE (name);


--
-- Name: articles_types articles_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles_types
    ADD CONSTRAINT articles_types_pkey PRIMARY KEY (id_articles_types);


--
-- Name: bassine bassine_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bassine
    ADD CONSTRAINT bassine_pkey PRIMARY KEY (id_bassine);


--
-- Name: caps caps_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.caps
    ADD CONSTRAINT caps_pkey PRIMARY KEY (id_caps);


--
-- Name: clubs_members clubs_members_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_members
    ADD CONSTRAINT clubs_members_pkey PRIMARY KEY (email, id_clubs);


--
-- Name: clubs clubs_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs
    ADD CONSTRAINT clubs_name_key UNIQUE (name);


--
-- Name: clubs clubs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs
    ADD CONSTRAINT clubs_pkey PRIMARY KEY (id_clubs);


--
-- Name: clubs_respos clubs_respos_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_respos
    ADD CONSTRAINT clubs_respos_pkey PRIMARY KEY (email, id_clubs);


--
-- Name: clubs_types clubs_types_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_types
    ADD CONSTRAINT clubs_types_name_key UNIQUE (name);


--
-- Name: clubs_types clubs_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_types
    ADD CONSTRAINT clubs_types_pkey PRIMARY KEY (id_clubs_types);


--
-- Name: events_attendents events_attendents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_attendents
    ADD CONSTRAINT events_attendents_pkey PRIMARY KEY (email, id_events);


--
-- Name: events_clubs events_clubs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_clubs
    ADD CONSTRAINT events_clubs_pkey PRIMARY KEY (id_clubs, id_events);


--
-- Name: events_hosts events_hosts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_hosts
    ADD CONSTRAINT events_hosts_pkey PRIMARY KEY (email, id_events);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id_events);


--
-- Name: files files_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_pkey PRIMARY KEY (id_files);


--
-- Name: games games_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_name_key UNIQUE (name);


--
-- Name: games games_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_pkey PRIMARY KEY (id_games);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: languages languages_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_name_key UNIQUE (name);


--
-- Name: languages languages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.languages
    ADD CONSTRAINT languages_pkey PRIMARY KEY (id_languages);


--
-- Name: newf newf_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf
    ADD CONSTRAINT newf_pkey PRIMARY KEY (email);


--
-- Name: newf_roles newf_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf_roles
    ADD CONSTRAINT newf_roles_pkey PRIMARY KEY (email, id_roles);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (email, id_services);


--
-- Name: players_bassine players_bassine_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_bassine
    ADD CONSTRAINT players_bassine_pkey PRIMARY KEY (email, id_games, id_bassine);


--
-- Name: players_caps players_caps_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_caps
    ADD CONSTRAINT players_caps_pkey PRIMARY KEY (email, id_games, id_caps);


--
-- Name: realcampus_friendships realcampus_friendships_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_friendships
    ADD CONSTRAINT realcampus_friendships_pkey PRIMARY KEY (id_friendship);


--
-- Name: realcampus_friendships realcampus_friendships_user_id_friend_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_friendships
    ADD CONSTRAINT realcampus_friendships_user_id_friend_id_key UNIQUE (user_id, friend_id);


--
-- Name: realcampus_post_reactions realcampus_post_reactions_id_post_author_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions
    ADD CONSTRAINT realcampus_post_reactions_id_post_author_email_key UNIQUE (id_post, author_email);


--
-- Name: realcampus_post_reactions realcampus_post_reactions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions
    ADD CONSTRAINT realcampus_post_reactions_pkey PRIMARY KEY (id_reaction);


--
-- Name: realcampus_posts realcampus_posts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_posts
    ADD CONSTRAINT realcampus_posts_pkey PRIMARY KEY (id_post);


--
-- Name: request_statistics request_statistics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.request_statistics
    ADD CONSTRAINT request_statistics_pkey PRIMARY KEY (id);


--
-- Name: reservations reservations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_pkey PRIMARY KEY (email, id_rooms);


--
-- Name: restaurant restaurant_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restaurant
    ADD CONSTRAINT restaurant_pkey PRIMARY KEY (id_restaurant);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id_roles);


--
-- Name: rooms rooms_location_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_location_key UNIQUE (location);


--
-- Name: rooms rooms_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_name_key UNIQUE (name);


--
-- Name: rooms rooms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rooms
    ADD CONSTRAINT rooms_pkey PRIMARY KEY (id_rooms);


--
-- Name: services services_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_name_key UNIQUE (name);


--
-- Name: services services_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.services
    ADD CONSTRAINT services_pkey PRIMARY KEY (id_services);


--
-- Name: support_media support_media_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support_media
    ADD CONSTRAINT support_media_pkey PRIMARY KEY (id_support_media);


--
-- Name: support support_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support
    ADD CONSTRAINT support_pkey PRIMARY KEY (id_support);


--
-- Name: traq traq_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq
    ADD CONSTRAINT traq_pkey PRIMARY KEY (id_traq);


--
-- Name: traq_types traq_types_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq_types
    ADD CONSTRAINT traq_types_name_key UNIQUE (name);


--
-- Name: traq_types traq_types_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq_types
    ADD CONSTRAINT traq_types_pkey PRIMARY KEY (id_traq_types);


--
-- Name: washing_machines washing_machines_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.washing_machines
    ADD CONSTRAINT washing_machines_pkey PRIMARY KEY (id_washing_machines);


--
-- Name: washing_machines_reports washing_machines_reports_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.washing_machines_reports
    ADD CONSTRAINT washing_machines_reports_pkey PRIMARY KEY (email, id_washing_machines);


--
-- Name: idx_realcampus_friendships_users; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_realcampus_friendships_users ON public.realcampus_friendships USING btree (user_id, friend_id);


--
-- Name: idx_realcampus_post_reactions_post; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_realcampus_post_reactions_post ON public.realcampus_post_reactions USING btree (id_post);


--
-- Name: idx_realcampus_posts_files_1; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_realcampus_posts_files_1 ON public.realcampus_posts USING btree (id_file_1);


--
-- Name: idx_realcampus_posts_files_2; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_realcampus_posts_files_2 ON public.realcampus_posts USING btree (id_file_2);


--
-- Name: idx_request_statistics_endpoint; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_request_statistics_endpoint ON public.request_statistics USING btree (endpoint);


--
-- Name: idx_request_statistics_user_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_request_statistics_user_email ON public.request_statistics USING btree (user_email);


--
-- Name: articles articles_buyer_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_buyer_fkey FOREIGN KEY (buyer) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: articles articles_id_articles_types_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_id_articles_types_fkey FOREIGN KEY (id_articles_types) REFERENCES public.articles_types(id_articles_types) ON DELETE CASCADE;


--
-- Name: articles articles_seller_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_seller_fkey FOREIGN KEY (seller) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: bassine bassine_creator_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.bassine
    ADD CONSTRAINT bassine_creator_fkey FOREIGN KEY (creator) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: caps caps_creator_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.caps
    ADD CONSTRAINT caps_creator_fkey FOREIGN KEY (creator) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: clubs clubs_id_clubs_types_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs
    ADD CONSTRAINT clubs_id_clubs_types_fkey FOREIGN KEY (id_clubs_types) REFERENCES public.clubs_types(id_clubs_types) ON DELETE CASCADE;


--
-- Name: clubs_members clubs_members_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_members
    ADD CONSTRAINT clubs_members_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: clubs_members clubs_members_id_clubs_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_members
    ADD CONSTRAINT clubs_members_id_clubs_fkey FOREIGN KEY (id_clubs) REFERENCES public.clubs(id_clubs) ON DELETE CASCADE;


--
-- Name: clubs_respos clubs_respos_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_respos
    ADD CONSTRAINT clubs_respos_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: clubs_respos clubs_respos_id_clubs_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.clubs_respos
    ADD CONSTRAINT clubs_respos_id_clubs_fkey FOREIGN KEY (id_clubs) REFERENCES public.clubs(id_clubs) ON DELETE CASCADE;


--
-- Name: events_attendents events_attendents_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_attendents
    ADD CONSTRAINT events_attendents_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: events_attendents events_attendents_id_events_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_attendents
    ADD CONSTRAINT events_attendents_id_events_fkey FOREIGN KEY (id_events) REFERENCES public.events(id_events) ON DELETE CASCADE;


--
-- Name: events_clubs events_clubs_id_clubs_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_clubs
    ADD CONSTRAINT events_clubs_id_clubs_fkey FOREIGN KEY (id_clubs) REFERENCES public.clubs(id_clubs) ON DELETE CASCADE;


--
-- Name: events_clubs events_clubs_id_events_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_clubs
    ADD CONSTRAINT events_clubs_id_events_fkey FOREIGN KEY (id_events) REFERENCES public.events(id_events) ON DELETE CASCADE;


--
-- Name: events_hosts events_hosts_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_hosts
    ADD CONSTRAINT events_hosts_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: events_hosts events_hosts_id_events_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events_hosts
    ADD CONSTRAINT events_hosts_id_events_fkey FOREIGN KEY (id_events) REFERENCES public.events(id_events) ON DELETE CASCADE;


--
-- Name: files files_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.files
    ADD CONSTRAINT files_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: newf fk_language; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf
    ADD CONSTRAINT fk_language FOREIGN KEY (language) REFERENCES public.languages(id_languages) ON DELETE CASCADE;


--
-- Name: restaurant fk_restaurant_language; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.restaurant
    ADD CONSTRAINT fk_restaurant_language FOREIGN KEY (language) REFERENCES public.languages(id_languages) ON DELETE CASCADE;


--
-- Name: newf_roles newf_roles_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf_roles
    ADD CONSTRAINT newf_roles_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: newf_roles newf_roles_id_roles_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.newf_roles
    ADD CONSTRAINT newf_roles_id_roles_fkey FOREIGN KEY (id_roles) REFERENCES public.roles(id_roles) ON DELETE CASCADE;


--
-- Name: notifications notifications_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: notifications notifications_id_services_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_id_services_fkey FOREIGN KEY (id_services) REFERENCES public.services(id_services) ON DELETE CASCADE;


--
-- Name: players_bassine players_bassine_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_bassine
    ADD CONSTRAINT players_bassine_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: players_bassine players_bassine_id_bassine_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_bassine
    ADD CONSTRAINT players_bassine_id_bassine_fkey FOREIGN KEY (id_bassine) REFERENCES public.bassine(id_bassine) ON DELETE CASCADE;


--
-- Name: players_bassine players_bassine_id_games_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_bassine
    ADD CONSTRAINT players_bassine_id_games_fkey FOREIGN KEY (id_games) REFERENCES public.games(id_games) ON DELETE CASCADE;


--
-- Name: players_caps players_caps_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_caps
    ADD CONSTRAINT players_caps_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: players_caps players_caps_id_caps_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_caps
    ADD CONSTRAINT players_caps_id_caps_fkey FOREIGN KEY (id_caps) REFERENCES public.caps(id_caps) ON DELETE CASCADE;


--
-- Name: players_caps players_caps_id_games_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.players_caps
    ADD CONSTRAINT players_caps_id_games_fkey FOREIGN KEY (id_games) REFERENCES public.games(id_games) ON DELETE CASCADE;


--
-- Name: realcampus_friendships realcampus_friendships_friend_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_friendships
    ADD CONSTRAINT realcampus_friendships_friend_id_fkey FOREIGN KEY (friend_id) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: realcampus_friendships realcampus_friendships_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_friendships
    ADD CONSTRAINT realcampus_friendships_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: realcampus_post_reactions realcampus_post_reactions_author_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions
    ADD CONSTRAINT realcampus_post_reactions_author_email_fkey FOREIGN KEY (author_email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: realcampus_post_reactions realcampus_post_reactions_id_file_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions
    ADD CONSTRAINT realcampus_post_reactions_id_file_fkey FOREIGN KEY (id_file) REFERENCES public.files(id_files) ON DELETE CASCADE;


--
-- Name: realcampus_post_reactions realcampus_post_reactions_id_post_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_post_reactions
    ADD CONSTRAINT realcampus_post_reactions_id_post_fkey FOREIGN KEY (id_post) REFERENCES public.realcampus_posts(id_post) ON DELETE CASCADE;


--
-- Name: realcampus_posts realcampus_posts_author_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_posts
    ADD CONSTRAINT realcampus_posts_author_email_fkey FOREIGN KEY (author_email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: realcampus_posts realcampus_posts_id_file_1_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_posts
    ADD CONSTRAINT realcampus_posts_id_file_1_fkey FOREIGN KEY (id_file_1) REFERENCES public.files(id_files) ON DELETE CASCADE;


--
-- Name: realcampus_posts realcampus_posts_id_file_2_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.realcampus_posts
    ADD CONSTRAINT realcampus_posts_id_file_2_fkey FOREIGN KEY (id_file_2) REFERENCES public.files(id_files) ON DELETE CASCADE;


--
-- Name: request_statistics request_statistics_user_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.request_statistics
    ADD CONSTRAINT request_statistics_user_email_fkey FOREIGN KEY (user_email) REFERENCES public.newf(email) ON DELETE SET NULL;


--
-- Name: reservations reservations_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: reservations reservations_id_rooms_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.reservations
    ADD CONSTRAINT reservations_id_rooms_fkey FOREIGN KEY (id_rooms) REFERENCES public.rooms(id_rooms) ON DELETE CASCADE;


--
-- Name: support support_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support
    ADD CONSTRAINT support_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: support_media support_media_id_files_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support_media
    ADD CONSTRAINT support_media_id_files_fkey FOREIGN KEY (id_files) REFERENCES public.files(id_files) ON DELETE CASCADE;


--
-- Name: support_media support_media_id_support_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.support_media
    ADD CONSTRAINT support_media_id_support_fkey FOREIGN KEY (id_support) REFERENCES public.support(id_support) ON DELETE CASCADE;


--
-- Name: traq traq_id_traq_types_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.traq
    ADD CONSTRAINT traq_id_traq_types_fkey FOREIGN KEY (id_traq_types) REFERENCES public.traq_types(id_traq_types) ON DELETE CASCADE;


--
-- Name: washing_machines_reports washing_machines_reports_email_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.washing_machines_reports
    ADD CONSTRAINT washing_machines_reports_email_fkey FOREIGN KEY (email) REFERENCES public.newf(email) ON DELETE CASCADE;


--
-- Name: washing_machines_reports washing_machines_reports_id_washing_machines_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.washing_machines_reports
    ADD CONSTRAINT washing_machines_reports_id_washing_machines_fkey FOREIGN KEY (id_washing_machines) REFERENCES public.washing_machines(id_washing_machines) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

