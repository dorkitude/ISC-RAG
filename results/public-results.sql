--
-- PostgreSQL database dump
--

\restrict 6O5YU1eYqprrejMOgG1Bo5eLFLJSKWvgAresakQnNWDNxZGqGGmkHMoz29LsHs2

-- Dumped from database version 17.10 (Debian 17.10-1.pgdg12+1)
-- Dumped by pg_dump version 17.10 (Debian 17.10-1.pgdg12+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public_release; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA public_release;


--
-- Name: SCHEMA public_release; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON SCHEMA public_release IS 'Source-text-free aggregate ISC-RAG results for public redistribution';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: agent_evaluation_run_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.agent_evaluation_run_summary (
    id uuid,
    name text,
    question_split text,
    model text,
    capability text,
    representation text,
    representation_version text,
    embedding_model text,
    repetitions integer,
    status text,
    started_at timestamp with time zone,
    completed_at timestamp with time zone,
    attempt_count integer,
    completed_attempt_count integer,
    error_count integer
);


--
-- Name: experiment_run_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.experiment_run_summary (
    id uuid,
    kind text,
    representation text,
    representation_version text,
    embedding_model text,
    generator_model text,
    question_split text,
    corpus_split text,
    top_k integer,
    context_char_budget integer,
    context_token_budget integer,
    context_tokenizer text,
    retrieval_mode text,
    rerank_model text,
    rerank_candidates integer,
    status text,
    started_at timestamp with time zone,
    completed_at timestamp with time zone,
    answer_count integer,
    prompt_tokens bigint,
    completion_tokens bigint,
    total_tokens bigint
);


--
-- Name: judgment_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.judgment_summary (
    experiment_run_id uuid,
    judgment_key text,
    verdict text,
    answer_count integer,
    recorded_prompt_tokens bigint,
    recorded_completion_tokens bigint,
    recorded_total_tokens bigint
);


--
-- Name: pipeline_run_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.pipeline_run_summary (
    id uuid,
    kind text,
    version text,
    model text,
    status text,
    started_at timestamp with time zone,
    completed_at timestamp with time zone,
    source_unit_count integer,
    fact_count integer,
    rejection_count integer
);


--
-- Name: rephrasing_study_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.rephrasing_study_summary (
    namespace text,
    version text,
    model text,
    run_count integer,
    fact_count integer,
    first_run_at timestamp with time zone,
    last_run_at timestamp with time zone
);


--
-- Name: study_two_summary; Type: TABLE; Schema: public_release; Owner: -
--

CREATE TABLE public_release.study_two_summary (
    study_version text,
    candidate text,
    model text,
    run_count integer,
    prompt_tokens bigint,
    completion_tokens bigint,
    total_tokens bigint,
    first_run_at timestamp with time zone,
    last_run_at timestamp with time zone
);


--
-- Data for Name: agent_evaluation_run_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.agent_evaluation_run_summary (id, name, question_split, model, capability, representation, representation_version, embedding_model, repetitions, status, started_at, completed_at, attempt_count, completed_attempt_count, error_count) FROM stdin;
7b403092-9c80-4463-acdc-3e326b94fd1c	agent-evaluation-integration-check	development	accounts/fireworks/models/kimi-k2p6	named-fact-rag	facts	pilot-named-v3	nomic-ai/nomic-embed-text-v1.5	1	complete	2026-07-25 07:58:08.008463+00	2026-07-25 07:58:10.444115+00	1	1	0
40dc981a-770a-49f8-a060-800535791a30	agent-toolbox-development-20260725-expanded	agent-dev-20260725	accounts/fireworks/models/deepseek-v4-flash	fixed-chunk-rag	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	10	failed	2026-07-25 08:00:20.162661+00	2026-07-25 08:12:14.067392+00	750	750	0
6655cb4e-100b-4223-ade7-9f3dcc5f8fda	agent-toolbox-development-20260725-expanded	agent-dev-20260725	accounts/fireworks/models/deepseek-v4-flash	named-fact-rag	facts	pilot-named-v3	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 08:02:46.454817+00	2026-07-25 08:07:29.327706+00	750	750	22
29d02006-fd7c-4b3d-8cf2-54dae7adcea8	agent-toolbox-development-20260725-expanded	agent-dev-20260725	accounts/fireworks/models/kimi-k2p6	fixed-chunk-rag	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	10	failed	2026-07-25 08:07:30.486488+00	2026-07-25 08:12:14.067392+00	750	750	0
94a4fe65-6b20-4cca-82ad-34e918b265de	agent-toolbox-development-20260725-expanded	agent-dev-20260725	accounts/fireworks/models/kimi-k2p6	named-fact-rag	facts	pilot-named-v3	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 08:09:23.104419+00	2026-07-25 08:11:08.29812+00	750	750	0
67d03477-75c3-40a9-955e-51236802dc10	agent-toolbox-development-20260725-expanded-corrected-fixed	agent-dev-20260725	accounts/fireworks/models/deepseek-v4-flash	pilot-fixed-chunk-rag	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	10	failed	2026-07-25 08:12:27.723026+00	2026-07-25 22:10:45.174487+00	188	188	0
c2bfe340-2d3b-4eb3-a61f-12e2bd609b69	agent-toolbox-development-20260725-expanded-corrected-fixed-r2	agent-dev-20260725	accounts/fireworks/models/deepseek-v4-flash	pilot-fixed-chunk-rag	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 22:10:54.673487+00	2026-07-25 22:15:02.934393+00	750	750	1
48173b40-2fc8-48e6-986a-36efd8396844	agent-toolbox-development-20260725-expanded-corrected-fixed-r2	agent-dev-20260725	accounts/fireworks/models/kimi-k2p6	pilot-fixed-chunk-rag	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 22:15:03.401225+00	2026-07-25 22:17:19.181368+00	750	750	0
c77456c0-be80-4371-a876-8cfba5cfbd87	agent-toolbox-development-20260725-fact-quote	agent-dev-20260725	accounts/fireworks/models/deepseek-v4-flash	named-fact-quote-rag	facts	pilot-named-v3-quote-v1	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 22:25:23.313312+00	2026-07-25 22:29:15.332771+00	750	750	0
5c06c381-8793-4886-80ac-a0526be8fd79	agent-toolbox-development-20260725-fact-quote	agent-dev-20260725	accounts/fireworks/models/kimi-k2p6	named-fact-quote-rag	facts	pilot-named-v3-quote-v1	nomic-ai/nomic-embed-text-v1.5	10	complete	2026-07-25 22:29:16.855375+00	2026-07-25 22:31:15.215679+00	750	750	0
4016c648-a2ad-44c7-9af9-4deae00cb902	fact-development-sample-100-flash-initial	fact-dev-sample-100-20260725	accounts/fireworks/models/deepseek-v4-flash	turnwise-opinion-fact-rag	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	3	complete	2026-07-25 23:17:21.188585+00	2026-07-25 23:17:49.660896+00	60	60	0
993b131e-f7ce-4f0a-8a49-72af9cb65be2	fact-development-sample-100-flash-initial	fact-dev-sample-100-20260725	accounts/fireworks/models/deepseek-v4-flash	fact-development-fixed-rag	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	3	complete	2026-07-25 23:17:54.405062+00	2026-07-25 23:18:58.641465+00	60	60	6
ba81fbd3-8a3d-46ce-bdfa-b225737b739f	fact-development-sample-100-kimi-initial	fact-dev-sample-100-20260725	accounts/fireworks/models/kimi-k2p6	turnwise-opinion-fact-rag	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	3	complete	2026-07-25 23:19:12.660037+00	2026-07-25 23:19:31.246391+00	60	60	0
cb54fc99-81bb-4dc7-8162-f4339fd6f8c5	fact-development-sample-100-kimi-initial	fact-dev-sample-100-20260725	accounts/fireworks/models/kimi-k2p6	fact-development-fixed-rag	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	3	complete	2026-07-25 23:19:39.882166+00	2026-07-25 23:20:03.160935+00	60	60	0
\.


--
-- Data for Name: experiment_run_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.experiment_run_summary (id, kind, representation, representation_version, embedding_model, generator_model, question_split, corpus_split, top_k, context_char_budget, context_token_budget, context_tokenizer, retrieval_mode, rerank_model, rerank_candidates, status, started_at, completed_at, answer_count, prompt_tokens, completion_tokens, total_tokens) FROM stdin;
6aaaa37e-dd76-4758-9a0d-02ec7c02aac7	retrieval	facts	v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	smoke	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 06:11:24.715988+00	2026-07-24 06:11:26.283463+00	1	0	0	0
387cad02-dd9d-4dad-8f16-ef226f608be4	retrieval	fixed	v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	smoke	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 06:11:39.693263+00	2026-07-24 06:11:40.102948+00	0	0	0	0
9d469b8f-3285-43f8-bc5d-918c643160fe	retrieval	turn	v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	smoke	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 06:11:40.384853+00	2026-07-24 06:11:40.678504+00	0	0	0	0
65e01e0b-b945-49ee-933d-bb32fd539afe	retrieval	semantic	v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	smoke	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 06:11:40.965295+00	2026-07-24 06:11:41.27635+00	0	0	0	0
dd06e024-6f81-4160-8819-efccc74623f7	retrieval	fixed	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	pilot-100x100-v1	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 20:11:57.668669+00	2026-07-24 20:14:13.715843+00	0	0	0	0
1a404566-679f-4fba-b11d-511db53130e1	retrieval	turn	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	pilot-100x100-v1	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 20:12:14.969401+00	2026-07-24 20:14:26.322235+00	0	0	0	0
55d9fd67-571e-4b23-b58d-9cfb72db3483	retrieval	semantic	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	pilot-100x100-v1	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 20:12:29.834016+00	2026-07-24 20:14:41.614468+00	0	0	0	0
54c3cd9b-c86a-472b-8ad2-90db37d4f551	retrieval	facts	pilot-100x100-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	pilot-100x100-v1	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-24 20:12:45.375932+00	2026-07-24 20:14:55.732656+00	0	0	0	0
52d0973c-f697-413f-8a5b-13f29dcd7f72	retrieval	facts	pilot-coverage-v2	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	pilot-100x100-v1	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-25 00:19:34.868603+00	2026-07-25 00:21:46.974928+00	0	0	0	0
79843416-81b0-4191-9fa3-5cb311a221a2	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:27:39.731777+00	2026-07-25 23:27:41.552678+00	0	0	0	0
a9dec815-06d5-4523-b70f-8861fe538854	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:27:41.836113+00	2026-07-25 23:27:43.6746+00	0	0	0	0
f32ae639-01e0-4ca4-9479-4900c92322ea	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	1	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:05.7169+00	2026-07-25 23:28:07.815962+00	0	0	0	0
ca2ba454-10da-48e2-92c5-2592ac8c2ed0	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	1	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:08.095011+00	2026-07-25 23:28:10.319708+00	0	0	0	0
aa7f8765-4465-4a48-95e5-3887384f8e48	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	3	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:10.594749+00	2026-07-25 23:28:12.549285+00	0	0	0	0
ad023981-2f7a-4426-bde6-ec5cc2863f59	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	3	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:12.821236+00	2026-07-25 23:28:15.269709+00	0	0	0	0
a7640b9c-3d52-49fc-b6db-240c0ddba9c6	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:15.557663+00	2026-07-25 23:28:17.798964+00	0	0	0	0
6def7954-607e-4c4a-b16c-6e96f5193c75	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-25 23:28:18.082992+00	2026-07-25 23:28:20.197551+00	0	0	0	0
7e4aea46-6a7d-453a-9080-a49386073332	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	1	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:20.892921+00	2026-07-26 00:33:22.680971+00	0	0	0	0
e5581de8-e412-491d-b25e-66effd17ff0c	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	1	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:22.958273+00	2026-07-26 00:33:24.934721+00	0	0	0	0
6a481ef7-91f1-4a75-9129-ddf6bf65434a	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	3	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:25.21138+00	2026-07-26 00:33:26.878083+00	0	0	0	0
b97cc8a1-0b29-4b21-a68d-0a405789f2bd	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	3	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:27.155355+00	2026-07-26 00:33:28.869412+00	0	0	0	0
a618d1bf-7906-477b-b509-4ddde1a78631	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:29.144722+00	2026-07-26 00:33:30.957922+00	0	0	0	0
417a4347-91d8-4647-a8bf-d45099a99cf6	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/gpt-oss-120b	fact-dev-sample-100-expansion-20260725	test	5	\N	\N	\N	\N	\N	\N	complete	2026-07-26 00:33:31.23561+00	2026-07-26 00:33:33.062614+00	0	0	0	0
d998bf41-0cb9-4291-9806-52e8f2f8aa93	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	2048	\N	\N	\N	\N	\N	complete	2026-07-26 00:39:53.927042+00	2026-07-26 00:44:08.00763+00	100	59747	20711	80458
92a23195-e90c-4dc0-8dc7-b7c25d20f77a	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	2048	\N	\N	\N	\N	\N	complete	2026-07-26 00:44:43.269029+00	2026-07-26 00:48:26.249405+00	100	51485	20014	71499
6b0106e8-8ec8-473e-9989-b225bfc5ff0b	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	512	\N	\N	\N	\N	\N	complete	2026-07-26 00:56:02.965422+00	2026-07-26 00:59:37.19288+00	100	18063	17429	35492
76e177e5-b41c-4817-9ec2-9f7efc295922	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	256	cl100k_base	\N	\N	\N	complete	2026-07-26 01:54:46.677793+00	2026-07-26 01:56:15.04589+00	100	36343	2641	38984
28d47a4d-69f5-45ea-82ed-fe41429cc861	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	512	cl100k_base	\N	\N	\N	complete	2026-07-26 02:03:22.030855+00	2026-07-26 02:05:33.904427+00	100	62861	2578	65439
dfa0dd13-d115-4b0e-b453-8ffdd354a66b	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	1024	cl100k_base	\N	\N	\N	complete	2026-07-26 02:05:34.186537+00	2026-07-26 02:07:07.739787+00	100	113194	2685	115879
12225f72-ed38-422b-a163-dd5e6599308f	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 02:07:08.030178+00	2026-07-26 02:08:54.51709+00	100	138186	2562	140748
7371831f-e91c-4455-ae73-e1a049f444f2	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	256	cl100k_base	\N	\N	\N	failed	2026-07-26 02:08:55.446998+00	2026-07-26 02:09:03.012579+00	7	2760	260	3020
8fc06526-982a-4c7b-9518-1ea213f5fa05	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	256	cl100k_base	\N	\N	\N	complete	2026-07-26 02:09:31.989531+00	2026-07-26 02:11:06.908604+00	100	40163	1762	41925
5cee9df7-76e8-4ebd-943f-22356adf1e3f	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	512	cl100k_base	\N	\N	\N	complete	2026-07-26 02:11:07.533264+00	2026-07-26 02:12:43.88792+00	100	64342	2015	66357
bc241efa-dd77-41cd-8d47-f250b0507c10	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	1024	cl100k_base	\N	\N	\N	complete	2026-07-26 02:12:44.181085+00	2026-07-26 02:14:20.142489+00	100	88359	2248	90607
af62b994-b63b-428c-a386-e884f2c93a52	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 02:14:20.439324+00	2026-07-26 02:16:04.363067+00	100	177913	2215	180128
a0e9f6c9-7031-4995-a709-5c3db18ac137	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	256	cl100k_base	\N	\N	\N	complete	2026-07-26 02:16:04.648969+00	2026-07-26 02:18:25.89542+00	100	40186	1841	42027
24e5e888-259b-4f04-ab4c-b8a0023cc0f5	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	512	cl100k_base	\N	\N	\N	complete	2026-07-26 02:18:26.290514+00	2026-07-26 02:20:19.098915+00	100	64349	2038	66387
14bc49ae-9929-4814-88dd-382813bf2b69	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	1024	cl100k_base	\N	\N	\N	complete	2026-07-26 02:20:19.524732+00	2026-07-26 02:22:53.002981+00	100	89810	2299	92109
2796a54b-241a-4dc3-94c6-d7f77d0b670b	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 02:22:53.418309+00	2026-07-26 02:24:45.78934+00	100	178359	2286	180645
d189aa2b-04be-4270-8b61-5fe9eb8bc86b	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	256	cl100k_base	\N	\N	\N	complete	2026-07-26 02:24:46.080539+00	2026-07-26 02:26:16.781199+00	100	38702	1783	40485
f38ca0d9-6d21-4481-b8ed-56c3311af5a9	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	512	cl100k_base	\N	\N	\N	complete	2026-07-26 02:26:17.064567+00	2026-07-26 02:28:03.022964+00	100	61077	2040	63117
f6c48244-fccf-49f4-9c1e-c3794d6cf648	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	1024	cl100k_base	\N	\N	\N	complete	2026-07-26 02:28:03.308515+00	2026-07-26 02:29:34.280905+00	100	84636	2377	87013
4ef714a6-728a-4763-a99b-22b31ab1961f	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 02:29:35.257003+00	2026-07-26 02:31:32.976883+00	100	187855	2251	190106
4ee5a60b-0f88-4d20-b1e8-33c3751efa46	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 02:49:36.946841+00	2026-07-26 02:50:12.691533+00	5	1065	81	1146
42379df3-e4ec-4e16-a28a-7c573e59f61b	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 02:50:12.984836+00	2026-07-26 02:50:17.481029+00	5	1684	87	1771
5c327145-46f1-463b-9829-61d0544c8b91	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 02:50:17.755687+00	2026-07-26 02:50:23.529856+00	5	2292	108	2400
6fd7aad9-0500-45d0-8593-9bf1ddeb10c7	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 02:50:23.80016+00	2026-07-26 02:50:28.754135+00	5	3755	118	3873
e51224df-d184-4d7c-bb7a-e1673de7e0f7	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 02:50:29.043576+00	2026-07-26 02:50:34.33499+00	5	6960	116	7076
c1717381-4259-4695-b030-ca2b62ce6379	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 02:50:34.625163+00	2026-07-26 02:50:40.66313+00	5	9484	114	9598
1429dfa6-ad4b-4047-ab9e-a6d1de743e9d	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 02:50:40.947888+00	2026-07-26 02:50:45.047332+00	5	4233	92	4325
cfc6c85e-7644-451a-9eea-6e013a06526e	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 02:50:45.321189+00	2026-07-26 02:50:50.261271+00	5	10555	102	10657
187225ed-364c-4fa8-842d-0bd9b55504ff	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 02:50:50.54137+00	2026-07-26 02:50:59.452228+00	5	15990	102	16092
b17a11c2-92c1-4a1c-8b98-d48774f6d1fc	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 02:50:59.738914+00	2026-07-26 02:51:10.61784+00	5	34125	106	34231
e5315e7e-6182-4023-b51d-f5fb7239d278	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 02:51:10.903573+00	2026-07-26 02:51:35.937655+00	5	70233	104	70337
42564cb1-a8ad-42d3-a08c-86942f08f35c	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 02:51:36.243264+00	2026-07-26 02:51:53.206864+00	5	105778	97	105875
9bb5148e-9243-4498-8378-b5c1e6d67f68	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 02:51:53.493157+00	2026-07-26 02:51:56.277962+00	5	3358	92	3450
3088ef78-36c5-4d28-af1f-73cf5cd619c4	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 02:51:56.565338+00	2026-07-26 02:52:02.159803+00	5	9817	103	9920
76715c27-806e-4bf7-a202-e24de7f596c0	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 02:52:02.450944+00	2026-07-26 02:52:07.247416+00	5	16375	106	16481
d1009eff-138b-4dff-96fb-0471bac724ae	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 02:52:07.530495+00	2026-07-26 02:52:12.38738+00	5	34444	104	34548
e1f2ed06-6d61-4071-b73a-7bf7ae8fedd7	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 02:52:12.675229+00	2026-07-26 02:52:50.537447+00	5	68501	104	68605
7ba0eaa3-c3d8-44e8-8aae-c8c779034c1e	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 02:52:50.816819+00	2026-07-26 02:53:09.074193+00	5	104781	101	104882
c8edfbad-3e27-4e1c-9a49-05183af4c6e4	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 02:53:09.354134+00	2026-07-26 02:53:12.487838+00	5	3061	92	3153
cab8c518-1f86-40be-b36d-a3c45fc328aa	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 02:53:12.773197+00	2026-07-26 02:53:21.01702+00	5	8518	92	8610
04964dd6-a00c-4fa7-8325-2676e5366ecb	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 02:53:21.293577+00	2026-07-26 02:53:29.755415+00	5	14116	106	14222
5f5eab39-ac46-4f6f-8ae9-f08da1c1c07a	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 02:53:30.039707+00	2026-07-26 02:53:35.32069+00	5	25843	106	25949
d40d723b-dda1-49fd-a72b-1fa36c058171	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 02:53:35.595899+00	2026-07-26 02:53:51.993156+00	5	50258	104	50362
ab0104a4-1b1f-4282-9cd0-2736e0704504	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 02:53:52.277077+00	2026-07-26 02:54:09.66777+00	5	77481	112	77593
c09be366-1289-4585-83b1-1436bba02f5e	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 03:04:36.033666+00	2026-07-26 03:15:05.986508+00	100	21387	1896	23283
2da5b00e-8f98-428e-be22-1a5bb7173205	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 03:15:06.291646+00	2026-07-26 03:18:28.352347+00	100	34516	2037	36553
68e765eb-4053-433a-8d63-b8f238f2b5e0	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 03:18:29.335041+00	2026-07-26 03:20:56.814234+00	100	47389	2160	49549
474a41c9-2a76-49c1-89e6-126b22b32db1	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 03:20:57.773137+00	2026-07-26 03:23:50.074175+00	100	79564	2180	81744
356d25e4-d974-48ab-8a52-50bde9468f16	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 03:23:51.025082+00	2026-07-26 03:25:57.48185+00	100	139079	2304	141383
4dd13da9-5773-451f-bd3f-28777aba9cb5	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 03:25:58.450814+00	2026-07-26 03:28:44.474836+00	100	186881	2339	189220
937f1cbd-c964-4995-9f06-58ea7e44841b	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 03:28:45.47245+00	2026-07-26 03:30:53.277523+00	100	86630	2254	88884
591e9c90-5a06-42b3-9b4b-ffb5a29f886d	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 03:29:27.49364+00	2026-07-26 03:32:56.662315+00	100	88001	2349	90350
5442ef5f-dabb-46c8-af8e-e5e19ed1d210	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 03:29:27.495446+00	2026-07-26 03:31:16.426159+00	100	70470	2456	72926
bc19fd4a-ff74-4e27-8fe3-4d4043591e6d	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 03:30:54.268106+00	2026-07-26 03:33:17.860513+00	100	228873	2236	231109
2efb2ea8-1dbd-4814-8975-37d7a26cfdc0	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 03:31:16.646074+00	2026-07-26 03:33:33.443141+00	100	187014	2230	189244
1083a416-9948-40de-b24b-7c58789d3d1a	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	complete	2026-07-26 03:32:56.923737+00	2026-07-26 03:34:44.145852+00	100	232463	2347	234810
263a1af9-649d-4550-b77f-d60ad6d14882	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 03:33:18.814046+00	2026-07-26 03:35:26.571633+00	100	371149	2268	373417
f8b6874a-adf4-4384-9436-2059c3f95d59	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 03:33:33.476308+00	2026-07-26 03:35:33.684352+00	100	306196	2212	308408
e9842c68-6937-4768-b013-5f7e87517d18	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	5	0	0		\N	\N	\N	complete	2026-07-26 03:34:44.186026+00	2026-07-26 03:36:46.424156+00	100	373385	2269	375654
53cacbbf-fd3a-4846-af03-96c90f4bdf9f	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 03:35:27.528082+00	2026-07-26 03:37:34.837959+00	100	732002	2229	734231
4310f88e-bb1a-48ea-9209-c9ca00bbdc68	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 03:35:33.925988+00	2026-07-26 03:37:35.912262+00	100	601284	2313	603597
babb185f-a8bd-4b0b-81b9-b5a726fc0b00	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	10	0	0		\N	\N	\N	complete	2026-07-26 03:36:46.470349+00	2026-07-26 03:38:57.490701+00	100	734657	2215	736872
c12c67c1-923b-48c6-b672-c871a8b2d01f	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 03:37:35.73426+00	2026-07-26 03:41:59.400104+00	100	1444219	2300	1446519
897e0b71-7232-4872-af0e-b10b228dc6ca	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 03:37:35.946279+00	2026-07-26 03:44:12.027532+00	100	1167385	2359	1169744
062d7c2b-b15f-4773-be22-bf04d122ee8d	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	20	0	0		\N	\N	\N	complete	2026-07-26 03:38:57.736278+00	2026-07-26 03:44:12.644499+00	100	1428203	2358	1430561
a5aee811-f912-475c-9f86-b1811fc36164	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 03:42:00.378124+00	2026-07-26 03:48:23.962666+00	100	2157709	2150	2159859
dab29081-c149-4abb-816d-5bd97ecfe202	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 03:44:12.302883+00	2026-07-26 03:48:53.674843+00	100	1718785	2429	1721214
0ee42343-f27d-4665-867c-ed9d0e146e3c	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	30	0	0		\N	\N	\N	complete	2026-07-26 03:44:12.703569+00	2026-07-26 03:49:31.503639+00	100	2130481	2343	2132824
bb02383f-926c-4513-af90-23d83d6af31f	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	1	0	0		\N	\N	\N	complete	2026-07-26 03:48:24.972202+00	2026-07-26 03:49:42.882815+00	100	88001	2288	90289
c3d4b70c-8ff2-4932-ba82-751b418511c8	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	3	0	0		\N	\N	\N	failed	2026-07-26 03:49:43.201024+00	2026-07-26 03:50:04.137258+00	9	18731	159	18890
ce011eea-b5af-4c37-b260-73c7792333ec	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 04:32:26.27962+00	2026-07-26 04:34:23.642444+00	100	191926	2257	194183
76d3ad12-40fd-4292-823c-f2f4a5a54325	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	4096	cl100k_base	\N	\N	\N	complete	2026-07-26 04:34:24.129986+00	2026-07-26 04:37:30.034373+00	100	201203	2319	203522
354a3b3a-4258-488f-a362-0bcfcbf0b275	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	8192	cl100k_base	\N	\N	\N	complete	2026-07-26 04:37:31.041389+00	2026-07-26 04:41:47.279915+00	100	201203	2276	203479
b84228f0-176e-439a-9146-e4c7f8a480bb	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 04:41:48.312637+00	2026-07-26 04:47:37.50606+00	100	201203	2264	203467
737b9df9-158c-45c8-bf16-6ca2dbc3bdaf	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 04:47:38.514823+00	2026-07-26 04:49:19.25689+00	100	177913	2204	180117
85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	4096	cl100k_base	\N	\N	\N	complete	2026-07-26 04:49:19.57255+00	2026-07-26 04:51:12.84429+00	100	392425	2251	394676
57834359-fbc2-4b7e-a776-f74df6b53df6	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	8192	cl100k_base	\N	\N	\N	complete	2026-07-26 04:51:13.176953+00	2026-07-26 04:53:59.050739+00	100	798012	2210	800222
b13cc9e8-1615-40ce-a13d-00021ee64647	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 04:53:59.373306+00	2026-07-26 04:59:28.703326+00	100	1622026	2313	1624339
f92ee367-79d6-4a88-b891-4f2305eb9d22	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 04:59:29.722725+00	2026-07-26 05:01:54.374728+00	100	178359	2132	180491
ada7f9b2-cc08-4524-a280-fcdfd2a37e0f	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	4096	cl100k_base	\N	\N	\N	complete	2026-07-26 05:01:55.491972+00	2026-07-26 05:04:22.360691+00	100	385162	2107	387269
99b978b9-966b-44a9-96b2-ca3fe4857eac	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	8192	cl100k_base	\N	\N	\N	complete	2026-07-26 05:04:23.331251+00	2026-07-26 05:07:17.217956+00	100	798217	2169	800386
a72590e0-ecc7-40ab-b47b-de381f8b948f	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 05:07:18.274955+00	2026-07-26 05:13:36.969543+00	100	1624723	2399	1627122
e42a506f-01a7-4253-8db8-f287fc3f504c	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 05:13:37.968509+00	2026-07-26 05:15:49.354873+00	100	187855	2260	190115
0192d427-8800-45ed-8e22-94caecef6e03	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	4096	cl100k_base	\N	\N	\N	complete	2026-07-26 05:15:49.641639+00	2026-07-26 05:18:15.50353+00	100	394152	2247	396399
4e1fe17f-deb1-4132-87d5-bdd056aa0d5d	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	8192	cl100k_base	\N	\N	\N	complete	2026-07-26 05:18:16.519044+00	2026-07-26 05:21:19.398788+00	100	802723	2271	804994
f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd	retrieval	semantic	fact-dev-sample-100-semantic-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 05:21:19.721568+00	2026-07-26 05:27:41.88578+00	100	1633777	2358	1636135
a00fb789-a14d-4d9a-9285-61d4af91f005	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 05:22:12.226554+00	2026-07-26 05:24:12.871008+00	100	221009	2316	223325
519b1c35-d314-4c3b-82a4-f60e7cbaa088	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	4096	cl100k_base	\N	\N	\N	complete	2026-07-26 05:24:13.158519+00	2026-07-26 05:27:39.447581+00	100	431051	2212	433263
5cf07cfc-a135-4080-b7e3-92c19c58d992	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	8192	cl100k_base	\N	\N	\N	complete	2026-07-26 05:27:39.765165+00	2026-07-26 05:32:36.024111+00	100	842276	2310	844586
adbfce87-cac8-4aed-80e8-d3262a094b41	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 05:32:37.022828+00	2026-07-26 05:37:21.064015+00	100	1307200	2210	1309410
a3abb2f5-623d-4e3a-ac28-1b9fb66b4143	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:01:32.079603+00	2026-07-26 06:03:04.749051+00	100	215514	3637	219151
92a84632-23f1-4a96-8634-51a94124997b	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:03:05.802801+00	2026-07-26 06:05:25.055476+00	100	1278724	3455	1282179
3f4d3129-63bd-43f4-aecf-925ec8c9e8ea	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:05:26.094767+00	2026-07-26 06:07:03.904752+00	100	174303	3637	177940
27606ee0-a78d-4f8e-9932-2823459cd04a	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:07:04.901644+00	2026-07-26 06:09:44.134045+00	100	1596495	3551	1600046
bbb15eb7-e886-48ce-a081-f7175f32533e	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:09:45.215475+00	2026-07-26 06:11:30.09655+00	100	174693	3585	178278
437b3aa8-04e9-47b7-a7c7-93340c94d854	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/kimi-k2p6	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:11:30.435885+00	2026-07-26 06:14:09.22053+00	100	1599591	3696	1603287
374b1d89-9d27-4c5c-8fe2-1b1b45134662	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:14:09.529449+00	2026-07-26 06:17:58.241678+00	100	218054	3620	221674
13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:17:58.56896+00	2026-07-26 06:22:08.952544+00	100	1299729	3705	1303434
35cbf6f5-a71a-4241-8936-30fcecea8cbf	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:22:09.9698+00	2026-07-26 06:25:47.682059+00	100	175536	3716	179252
1ce8371e-4d9c-41cf-bd42-b68507c04911	retrieval	fixed	fact-dev-sample-100-fixed-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:25:48.01802+00	2026-07-26 06:30:18.789651+00	100	1614197	4178	1618375
0334d09b-8c21-4d89-958b-0da9258ae68e	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	\N	\N	\N	complete	2026-07-26 06:30:19.793268+00	2026-07-26 06:33:59.86041+00	100	175921	3450	179371
3522c9f5-7734-454d-948f-89245c630aeb	retrieval	turn	fact-dev-sample-100-turn-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/glm-5p2	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	\N	\N	\N	complete	2026-07-26 06:34:00.18387+00	2026-07-26 06:38:12.989419+00	100	1616241	3994	1620235
ed4e1cc1-5a33-4f54-a91b-8c50bc45942b	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-26 08:48:53.248465+00	2026-07-26 08:50:22.399171+00	100	182134	2376	184510
77b616ce-fc29-46db-9339-ac313885f878	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	hybrid		100	complete	2026-07-26 08:50:22.750928+00	2026-07-26 08:55:50.774997+00	100	182021	2381	184402
95b8f4e3-3718-4c70-8f1c-3de34e9aab93	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	hybrid	accounts/fireworks/models/qwen3-reranker-8b	50	complete	2026-07-26 08:55:51.087684+00	2026-07-26 09:09:49.26143+00	100	182225	2573	184798
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	vector		100	complete	2026-07-26 09:09:50.370281+00	2026-07-26 09:13:56.875419+00	100	1623122	2387	1625509
e58993e6-1f2e-4e8c-af89-d6139259e281	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	hybrid		100	complete	2026-07-26 09:13:57.188396+00	2026-07-26 09:19:02.250653+00	100	1623107	2372	1625479
453ddf0c-0605-48b3-a6b0-1f60674f3dcf	retrieval	fixed	fact-dev-sample-100-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	16384	cl100k_base	hybrid	accounts/fireworks/models/qwen3-reranker-8b	50	complete	2026-07-26 09:19:03.28076+00	2026-07-26 09:35:38.776391+00	100	1613999	2373	1616372
b2e31edc-f490-4eb2-b834-e080fba87dd5	retrieval	facts	turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	fact-dev-sample-100-expansion-20260725	test	300	0	2048	cl100k_base	vector	accounts/fireworks/models/qwen3-reranker-8b	100	complete	2026-07-26 09:35:39.796023+00	2026-07-26 09:39:15.680719+00	100	221425	2265	223690
04440a4f-75cc-441c-b749-280208e0847b	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	256	cl100k_base	vector		100	complete	2026-07-26 19:17:47.49303+00	2026-07-26 19:28:02.770127+00	499	186582	10170	196752
fcd33788-0881-4897-b1e1-29aabe71ef1a	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	256	cl100k_base	vector		100	complete	2026-07-26 19:17:47.493436+00	2026-07-26 19:30:38.760359+00	499	197166	9503	206669
264380aa-479b-4110-8176-7a3b6f2f8294	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	256	cl100k_base	vector		100	complete	2026-07-26 19:17:47.493639+00	2026-07-26 19:27:43.526887+00	499	198145	9248	207393
1d7d20ab-4048-45b8-904e-f3a042b5ad8d	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	256	cl100k_base	vector		100	complete	2026-07-26 19:17:47.493721+00	2026-07-26 19:35:16.718528+00	499	200194	9346	209540
de00b901-f068-4553-98a6-5de1c92f9815	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	complete	2026-07-26 19:35:17.813343+00	2026-07-26 19:46:16.241857+00	499	314590	10650	325240
19fd0831-b224-443f-a13f-01746c6e9ea8	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	complete	2026-07-26 19:46:17.347422+00	2026-07-26 19:55:44.067157+00	499	317292	10515	327807
f95e0ade-83c7-4108-8131-6c6dfde7388c	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	complete	2026-07-26 19:55:45.053712+00	2026-07-26 20:06:36.705149+00	499	312943	10455	323398
dbada771-c92d-4199-b044-a1fe4f2c25cb	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	failed	2026-07-26 20:06:37.934122+00	2026-07-26 20:16:39.595813+00	496	306165	10526	316691
30611993-b157-49a7-b451-9a4311dedb23	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	1024	cl100k_base	vector		100	complete	2026-07-26 20:16:40.658969+00	2026-07-26 20:25:13.945496+00	499	544180	10714	554894
0cf02bbf-d09b-4ce8-99ea-58eb5ab14ee1	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	1024	cl100k_base	vector		100	complete	2026-07-26 20:25:14.956269+00	2026-07-26 20:34:20.043659+00	499	446508	11091	457599
3959efad-f8ec-4f17-ba67-bcccc677352b	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	1024	cl100k_base	vector		100	failed	2026-07-26 20:34:21.073247+00	2026-07-26 20:45:07.100279+00	471	424639	10435	435074
34a958ff-2e73-4510-941a-05daeba85d50	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	1024	cl100k_base	vector		100	complete	2026-07-26 20:45:08.206346+00	2026-07-26 20:55:20.21668+00	499	440086	11268	451354
010d3fd4-3ff1-49e4-8c71-aa3b04bdb49c	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	2048	cl100k_base	vector		100	complete	2026-07-26 20:55:21.315642+00	2026-07-26 21:06:42.700822+00	499	633673	10879	644552
3742139e-8fee-49e0-9e6a-74264863f56d	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:06:43.789643+00	2026-07-26 21:16:45.153938+00	499	888237	11315	899552
d997c177-8802-498b-8220-bc036950d847	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:16:46.342246+00	2026-07-26 21:25:23.360681+00	499	899540	11116	910656
c839fc4c-37ac-44fb-b59c-bdd71fc597c0	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:25:24.443842+00	2026-07-26 21:34:09.112263+00	499	932690	11325	944015
2f9634c4-726e-45e2-8e90-a59cea90dd3d	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:34:10.213004+00	2026-07-26 21:43:02.467934+00	499	1103258	10955	1114213
caf8170f-d920-4ad3-a03a-d49b87b38350	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:43:03.509131+00	2026-07-26 21:51:27.861075+00	499	888237	11234	899471
ef1954a1-a60a-43b8-914c-ee01fab35bfa	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-26 21:51:28.966726+00	2026-07-26 22:01:10.772322+00	499	899540	11304	910844
c6b10ae2-e9b7-4ecc-ab49-d76eb92834a0	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-26 22:01:11.898927+00	2026-07-26 22:12:31.539578+00	499	932690	11312	944002
69de4988-5157-4cec-8f2f-869381b47432	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	4096	cl100k_base	vector		100	complete	2026-07-26 22:12:32.705985+00	2026-07-26 22:22:12.590558+00	499	2106999	10926	2117925
f292d9ca-5e2f-44b6-a175-f322073c7ad4	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	4096	cl100k_base	vector		100	complete	2026-07-26 22:22:13.625023+00	2026-07-26 22:30:54.903742+00	499	1936868	11111	1947979
7b587744-cd46-4764-aefb-cb7639456fa8	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	4096	cl100k_base	vector		100	complete	2026-07-26 22:30:55.965209+00	2026-07-26 22:40:18.551399+00	499	1935365	11328	1946693
711eed2f-8550-4a29-88a7-03b8ebc430f3	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	4096	cl100k_base	vector		100	complete	2026-07-26 22:40:19.57732+00	2026-07-26 22:49:13.543424+00	499	1962607	11028	1973635
61d548d6-0fe7-4a4d-9a57-0400564eb73f	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	8192	cl100k_base	vector		100	running	2026-07-26 22:49:14.588409+00	\N	259	1889867	5606	1895473
cd2d33c4-7125-42be-a7a3-5176ba984c64	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	8192	cl100k_base	vector		100	complete	2026-07-26 23:01:47.584999+00	2026-07-26 23:19:01.940941+00	499	3648031	10886	3658917
70e5b79f-263a-406b-bec9-8570ad56dc4b	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	8192	cl100k_base	vector		100	complete	2026-07-26 23:01:47.585013+00	2026-07-26 23:16:55.079186+00	499	3975051	11097	3986148
8155633d-4c03-4da7-9bbb-66a344361bfc	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	1024	cl100k_base	vector		100	complete	2026-07-26 23:01:47.585145+00	2026-07-26 23:11:10.618321+00	499	450985	11113	462098
be729536-d824-41e9-93b3-864608a80889	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	failed	2026-07-26 23:01:47.585181+00	2026-07-26 23:13:24.913605+00	496	306165	10550	316715
59ab526c-3301-4930-aa93-7f9233ce2453	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	8192	cl100k_base	vector		100	complete	2026-07-26 23:19:02.971459+00	2026-07-26 23:34:03.448472+00	499	3991173	10962	4002135
1c56273b-105a-4b3b-b627-57fbca3af8c3	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	failed	2026-07-26 23:21:32.294039+00	2026-07-26 23:30:47.248299+00	496	306165	10614	316779
655b8a6c-0ff5-4d0e-9c39-e99c1a56e1f6	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	20	0	512	cl100k_base	vector		100	complete	2026-07-26 23:31:45.963496+00	2026-07-26 23:43:27.226195+00	499	308315	10660	318975
3e5c61fc-6124-4a87-8eef-a8bf819fda9d	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	8192	cl100k_base	vector		100	complete	2026-07-26 23:34:04.950055+00	2026-07-26 23:50:38.743179+00	499	4013479	11357	4024836
e1540d49-375a-4ebd-8385-98f2cb976156	retrieval	facts	heldout500-turnwise-opinion-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	vector		100	complete	2026-07-26 23:50:39.641636+00	2026-07-27 00:13:02.923758+00	499	4231187	10859	4242046
d728cf0d-8e04-4f80-8487-a1876ee25aa4	retrieval	fixed	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	vector		100	complete	2026-07-27 00:13:04.589121+00	2026-07-27 00:45:20.787914+00	499	8095784	11776	8107560
38b99945-45d0-4156-aafa-8a8f0073be48	retrieval	turn	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	vector		100	complete	2026-07-27 00:45:21.813149+00	2026-07-27 01:10:50.68086+00	499	8097916	12000	8109916
57ce84ee-cf8c-4863-a7cd-463ace593dff	retrieval	semantic	heldout500-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	vector		100	complete	2026-07-27 01:10:51.794469+00	2026-07-27 01:38:41.203324+00	499	8125393	11978	8137371
63667afc-f06f-4006-ae00-a23627ad85a4	retrieval	fixed	heldout500-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	hybrid	accounts/fireworks/models/qwen3-reranker-8b	50	running	2026-07-27 07:44:07.398831+00	\N	349	50805	7991	58796
d43882dc-7339-44ac-9cd0-75f67dfc669f	retrieval	fixed	heldout500-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	vector		100	complete	2026-07-27 07:53:47.830052+00	2026-07-27 08:05:26.031635+00	499	909897	11632	921529
c916bb0a-81d1-4ffd-b555-0fed0b547cf4	retrieval	fixed	heldout500-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	hybrid		100	complete	2026-07-27 08:05:27.223155+00	2026-07-27 08:23:22.680703+00	499	910524	11868	922392
bfebee8c-6fa7-4cf2-8557-9995a9c7cc51	retrieval	fixed	heldout500-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	2048	cl100k_base	hybrid	accounts/fireworks/models/qwen3-reranker-8b	50	complete	2026-07-27 08:23:23.783007+00	2026-07-27 09:36:48.47979+00	499	909970	11972	921942
3fb83e2c-0f72-4909-8ea8-011c2255de55	retrieval	fixed	heldout500-ctx-v1	nomic-ai/nomic-embed-text-v1.5	accounts/fireworks/models/deepseek-v4-flash	heldout-500-20260726	test	300	0	16384	cl100k_base	hybrid	accounts/fireworks/models/qwen3-reranker-8b	50	complete	2026-07-27 09:36:49.65042+00	2026-07-27 11:12:54.003834+00	499	8073307	12242	8085549
\.


--
-- Data for Name: judgment_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.judgment_summary (experiment_run_id, judgment_key, verdict, answer_count, recorded_prompt_tokens, recorded_completion_tokens, recorded_total_tokens) FROM stdin;
010d3fd4-3ff1-49e4-8c71-aa3b04bdb49c	judge_kimi_heldout500_v1	correct	393	0	0	0
010d3fd4-3ff1-49e4-8c71-aa3b04bdb49c	judge_kimi_heldout500_v1	incorrect	106	0	0	0
010d3fd4-3ff1-49e4-8c71-aa3b04bdb49c	source_support_kimi_heldout500_v1	supported	377	0	0	0
010d3fd4-3ff1-49e4-8c71-aa3b04bdb49c	source_support_kimi_heldout500_v1	unsupported	122	0	0	0
0192d427-8800-45ed-8e22-94caecef6e03	judge_kimi_frontier_xl_v1	correct	61	0	0	0
0192d427-8800-45ed-8e22-94caecef6e03	judge_kimi_frontier_xl_v1	incorrect	39	0	0	0
0192d427-8800-45ed-8e22-94caecef6e03	source_support_kimi_frontier_xl_v1	supported	64	0	0	0
0192d427-8800-45ed-8e22-94caecef6e03	source_support_kimi_frontier_xl_v1	unsupported	36	0	0	0
0334d09b-8c21-4d89-958b-0da9258ae68e	judge_kimi_frontier_xl_v1	correct	72	0	0	0
0334d09b-8c21-4d89-958b-0da9258ae68e	judge_kimi_frontier_xl_v1	incorrect	28	0	0	0
0334d09b-8c21-4d89-958b-0da9258ae68e	source_support_kimi_frontier_xl_v1	supported	72	0	0	0
0334d09b-8c21-4d89-958b-0da9258ae68e	source_support_kimi_frontier_xl_v1	unsupported	28	0	0	0
04440a4f-75cc-441c-b749-280208e0847b	judge_kimi_heldout500_v1	correct	347	0	0	0
04440a4f-75cc-441c-b749-280208e0847b	judge_kimi_heldout500_v1	incorrect	152	0	0	0
04440a4f-75cc-441c-b749-280208e0847b	source_support_kimi_heldout500_v1	supported	341	0	0	0
04440a4f-75cc-441c-b749-280208e0847b	source_support_kimi_heldout500_v1	unsupported	158	0	0	0
04964dd6-a00c-4fa7-8325-2676e5366ecb	judge_kimi_depth_v1	correct	4	639	196	835
04964dd6-a00c-4fa7-8325-2676e5366ecb	judge_kimi_depth_v1	incorrect	1	148	32	180
04964dd6-a00c-4fa7-8325-2676e5366ecb	source_support_kimi_depth_v1	supported	3	1916	187	2103
04964dd6-a00c-4fa7-8325-2676e5366ecb	source_support_kimi_depth_v1	unsupported	2	777	66	843
062d7c2b-b15f-4773-be22-bf04d122ee8d	judge_flash_depth100_v1	correct	69	13161	2801	15962
062d7c2b-b15f-4773-be22-bf04d122ee8d	judge_flash_depth100_v1	incorrect	31	5557	1402	6959
062d7c2b-b15f-4773-be22-bf04d122ee8d	source_support_flash_depth100_v1	supported	66	63602	2775	66377
062d7c2b-b15f-4773-be22-bf04d122ee8d	source_support_flash_depth100_v1	unsupported	34	6320	308	6628
0cf02bbf-d09b-4ce8-99ea-58eb5ab14ee1	judge_kimi_heldout500_v1	correct	290	0	0	0
0cf02bbf-d09b-4ce8-99ea-58eb5ab14ee1	judge_kimi_heldout500_v1	incorrect	209	0	0	0
0cf02bbf-d09b-4ce8-99ea-58eb5ab14ee1	source_support_kimi_heldout500_v1	supported	294	0	0	0
0cf02bbf-d09b-4ce8-99ea-58eb5ab14ee1	source_support_kimi_heldout500_v1	unsupported	205	0	0	0
0ee42343-f27d-4665-867c-ed9d0e146e3c	judge_flash_depth100_v1	correct	75	14077	3033	17110
0ee42343-f27d-4665-867c-ed9d0e146e3c	judge_flash_depth100_v1	incorrect	25	4627	1154	5781
0ee42343-f27d-4665-867c-ed9d0e146e3c	source_support_flash_depth100_v1	supported	72	68316	2899	71215
0ee42343-f27d-4665-867c-ed9d0e146e3c	source_support_flash_depth100_v1	unsupported	28	6388	283	6671
1083a416-9948-40de-b24b-7c58789d3d1a	judge_flash_depth100_v1	correct	59	11126	2375	13501
1083a416-9948-40de-b24b-7c58789d3d1a	judge_flash_depth100_v1	incorrect	41	7581	1810	9391
1083a416-9948-40de-b24b-7c58789d3d1a	source_support_flash_depth100_v1	supported	55	51368	2155	53523
1083a416-9948-40de-b24b-7c58789d3d1a	source_support_flash_depth100_v1	unsupported	45	7386	352	7738
12225f72-ed38-422b-a163-dd5e6599308f	judge_kimi_v1	correct	86	14556	4545	19101
12225f72-ed38-422b-a163-dd5e6599308f	judge_kimi_v1	incorrect	14	2672	860	3532
12225f72-ed38-422b-a163-dd5e6599308f	source_support_kimi_v1	supported	87	24114	4697	28811
12225f72-ed38-422b-a163-dd5e6599308f	source_support_kimi_v1	unsupported	13	1566	418	1984
13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2	judge_kimi_frontier_xl_v1	correct	99	0	0	0
13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2	judge_kimi_frontier_xl_v1	incorrect	1	0	0	0
13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2	source_support_kimi_frontier_xl_v1	supported	90	0	0	0
13fe5abc-f7db-4e42-9e3b-fa62a1cf13a2	source_support_kimi_frontier_xl_v1	unsupported	10	0	0	0
1429dfa6-ad4b-4047-ab9e-a6d1de743e9d	judge_kimi_depth_v1	correct	4	624	191	815
1429dfa6-ad4b-4047-ab9e-a6d1de743e9d	judge_kimi_depth_v1	incorrect	1	148	32	180
1429dfa6-ad4b-4047-ab9e-a6d1de743e9d	source_support_kimi_depth_v1	supported	3	2677	156	2833
1429dfa6-ad4b-4047-ab9e-a6d1de743e9d	source_support_kimi_depth_v1	unsupported	2	716	86	802
14bc49ae-9929-4814-88dd-382813bf2b69	judge_kimi_v1	correct	56	9744	3023	12767
14bc49ae-9929-4814-88dd-382813bf2b69	judge_kimi_v1	incorrect	44	7254	2073	9327
14bc49ae-9929-4814-88dd-382813bf2b69	source_support_kimi_v1	supported	57	51302	3004	54306
14bc49ae-9929-4814-88dd-382813bf2b69	source_support_kimi_v1	unsupported	43	3842	367	4209
187225ed-364c-4fa8-842d-0bd9b55504ff	judge_kimi_depth_v1	correct	4	630	225	855
187225ed-364c-4fa8-842d-0bd9b55504ff	judge_kimi_depth_v1	incorrect	1	148	39	187
187225ed-364c-4fa8-842d-0bd9b55504ff	source_support_kimi_depth_v1	supported	3	2668	155	2823
187225ed-364c-4fa8-842d-0bd9b55504ff	source_support_kimi_depth_v1	unsupported	2	731	58	789
19fd0831-b224-443f-a13f-01746c6e9ea8	judge_kimi_heldout500_v1	correct	225	0	0	0
19fd0831-b224-443f-a13f-01746c6e9ea8	judge_kimi_heldout500_v1	incorrect	274	0	0	0
19fd0831-b224-443f-a13f-01746c6e9ea8	source_support_kimi_heldout500_v1	supported	231	0	0	0
19fd0831-b224-443f-a13f-01746c6e9ea8	source_support_kimi_heldout500_v1	unsupported	268	0	0	0
1ce8371e-4d9c-41cf-bd42-b68507c04911	judge_kimi_frontier_xl_v1	correct	91	0	0	0
1ce8371e-4d9c-41cf-bd42-b68507c04911	judge_kimi_frontier_xl_v1	incorrect	9	0	0	0
1ce8371e-4d9c-41cf-bd42-b68507c04911	source_support_kimi_frontier_xl_v1	supported	93	0	0	0
1ce8371e-4d9c-41cf-bd42-b68507c04911	source_support_kimi_frontier_xl_v1	unsupported	7	0	0	0
1d7d20ab-4048-45b8-904e-f3a042b5ad8d	judge_kimi_heldout500_v1	correct	126	0	0	0
1d7d20ab-4048-45b8-904e-f3a042b5ad8d	judge_kimi_heldout500_v1	incorrect	373	0	0	0
1d7d20ab-4048-45b8-904e-f3a042b5ad8d	source_support_kimi_heldout500_v1	supported	134	0	0	0
1d7d20ab-4048-45b8-904e-f3a042b5ad8d	source_support_kimi_heldout500_v1	unsupported	365	0	0	0
24e5e888-259b-4f04-ab4c-b8a0023cc0f5	judge_kimi_v1	correct	45	7560	2406	9966
24e5e888-259b-4f04-ab4c-b8a0023cc0f5	judge_kimi_v1	incorrect	55	9158	2592	11750
24e5e888-259b-4f04-ab4c-b8a0023cc0f5	source_support_kimi_v1	supported	46	40178	2447	42625
24e5e888-259b-4f04-ab4c-b8a0023cc0f5	source_support_kimi_v1	unsupported	54	3648	332	3980
263a1af9-649d-4550-b77f-d60ad6d14882	judge_flash_depth100_v1	correct	72	13573	2897	16470
263a1af9-649d-4550-b77f-d60ad6d14882	judge_flash_depth100_v1	incorrect	28	5050	1276	6326
263a1af9-649d-4550-b77f-d60ad6d14882	source_support_flash_depth100_v1	supported	64	61466	2657	64123
263a1af9-649d-4550-b77f-d60ad6d14882	source_support_flash_depth100_v1	unsupported	36	12062	535	12597
264380aa-479b-4110-8176-7a3b6f2f8294	judge_kimi_heldout500_v1	correct	135	0	0	0
264380aa-479b-4110-8176-7a3b6f2f8294	judge_kimi_heldout500_v1	incorrect	364	0	0	0
264380aa-479b-4110-8176-7a3b6f2f8294	source_support_kimi_heldout500_v1	supported	141	0	0	0
264380aa-479b-4110-8176-7a3b6f2f8294	source_support_kimi_heldout500_v1	unsupported	358	0	0	0
27606ee0-a78d-4f8e-9932-2823459cd04a	judge_glm_frontier_xl_v1	correct	87	0	0	0
27606ee0-a78d-4f8e-9932-2823459cd04a	judge_glm_frontier_xl_v1	incorrect	13	0	0	0
27606ee0-a78d-4f8e-9932-2823459cd04a	source_support_glm_frontier_xl_v1	supported	89	0	0	0
27606ee0-a78d-4f8e-9932-2823459cd04a	source_support_glm_frontier_xl_v1	unsupported	11	0	0	0
2796a54b-241a-4dc3-94c6-d7f77d0b670b	judge_kimi_v1	correct	59	10122	3173	13295
2796a54b-241a-4dc3-94c6-d7f77d0b670b	judge_kimi_v1	incorrect	41	6874	1943	8817
2796a54b-241a-4dc3-94c6-d7f77d0b670b	source_support_kimi_v1	supported	62	56058	3152	59210
2796a54b-241a-4dc3-94c6-d7f77d0b670b	source_support_kimi_v1	unsupported	38	2788	229	3017
28d47a4d-69f5-45ea-82ed-fe41429cc861	judge_kimi_v1	correct	81	13633	4313	17946
28d47a4d-69f5-45ea-82ed-fe41429cc861	judge_kimi_v1	incorrect	19	3641	1079	4720
28d47a4d-69f5-45ea-82ed-fe41429cc861	source_support_kimi_v1	supported	83	22308	4486	26794
28d47a4d-69f5-45ea-82ed-fe41429cc861	source_support_kimi_v1	unsupported	17	1492	408	1900
2da5b00e-8f98-428e-be22-1a5bb7173205	judge_flash_depth100_v1	correct	71	12854	2878	15732
2da5b00e-8f98-428e-be22-1a5bb7173205	judge_flash_depth100_v1	incorrect	29	5504	1396	6900
2da5b00e-8f98-428e-be22-1a5bb7173205	source_support_flash_depth100_v1	supported	73	19650	3086	22736
2da5b00e-8f98-428e-be22-1a5bb7173205	source_support_flash_depth100_v1	unsupported	27	723	133	856
2efb2ea8-1dbd-4814-8975-37d7a26cfdc0	judge_flash_depth100_v1	correct	62	11761	2485	14246
2efb2ea8-1dbd-4814-8975-37d7a26cfdc0	judge_flash_depth100_v1	incorrect	38	6834	1686	8520
2efb2ea8-1dbd-4814-8975-37d7a26cfdc0	source_support_flash_depth100_v1	supported	55	44191	2284	46475
2efb2ea8-1dbd-4814-8975-37d7a26cfdc0	source_support_flash_depth100_v1	unsupported	45	8893	547	9440
2f9634c4-726e-45e2-8e90-a59cea90dd3d	judge_kimi_heldout500_v1	correct	425	0	0	0
2f9634c4-726e-45e2-8e90-a59cea90dd3d	judge_kimi_heldout500_v1	incorrect	74	0	0	0
2f9634c4-726e-45e2-8e90-a59cea90dd3d	source_support_kimi_heldout500_v1	supported	414	0	0	0
2f9634c4-726e-45e2-8e90-a59cea90dd3d	source_support_kimi_heldout500_v1	unsupported	85	0	0	0
30611993-b157-49a7-b451-9a4311dedb23	judge_kimi_heldout500_v1	correct	388	0	0	0
30611993-b157-49a7-b451-9a4311dedb23	judge_kimi_heldout500_v1	incorrect	111	0	0	0
30611993-b157-49a7-b451-9a4311dedb23	source_support_kimi_heldout500_v1	supported	376	0	0	0
30611993-b157-49a7-b451-9a4311dedb23	source_support_kimi_heldout500_v1	unsupported	123	0	0	0
3088ef78-36c5-4d28-af1f-73cf5cd619c4	judge_kimi_depth_v1	correct	4	637	202	839
3088ef78-36c5-4d28-af1f-73cf5cd619c4	judge_kimi_depth_v1	incorrect	1	148	40	188
3088ef78-36c5-4d28-af1f-73cf5cd619c4	source_support_kimi_depth_v1	supported	3	2445	140	2585
3088ef78-36c5-4d28-af1f-73cf5cd619c4	source_support_kimi_depth_v1	unsupported	2	640	43	683
34a958ff-2e73-4510-941a-05daeba85d50	judge_kimi_heldout500_v1	correct	294	0	0	0
34a958ff-2e73-4510-941a-05daeba85d50	judge_kimi_heldout500_v1	incorrect	205	0	0	0
34a958ff-2e73-4510-941a-05daeba85d50	source_support_kimi_heldout500_v1	supported	293	0	0	0
34a958ff-2e73-4510-941a-05daeba85d50	source_support_kimi_heldout500_v1	unsupported	206	0	0	0
3522c9f5-7734-454d-948f-89245c630aeb	judge_kimi_frontier_xl_v1	correct	91	0	0	0
3522c9f5-7734-454d-948f-89245c630aeb	judge_kimi_frontier_xl_v1	incorrect	9	0	0	0
3522c9f5-7734-454d-948f-89245c630aeb	source_support_kimi_frontier_xl_v1	supported	92	0	0	0
3522c9f5-7734-454d-948f-89245c630aeb	source_support_kimi_frontier_xl_v1	unsupported	8	0	0	0
356d25e4-d974-48ab-8a52-50bde9468f16	judge_flash_depth100_v1	correct	83	15454	3320	18774
356d25e4-d974-48ab-8a52-50bde9468f16	judge_flash_depth100_v1	incorrect	17	3096	788	3884
356d25e4-d974-48ab-8a52-50bde9468f16	source_support_flash_depth100_v1	supported	85	25327	3567	28894
356d25e4-d974-48ab-8a52-50bde9468f16	source_support_flash_depth100_v1	unsupported	15	1112	213	1325
35cbf6f5-a71a-4241-8936-30fcecea8cbf	judge_kimi_frontier_xl_v1	correct	75	0	0	0
35cbf6f5-a71a-4241-8936-30fcecea8cbf	judge_kimi_frontier_xl_v1	incorrect	25	0	0	0
35cbf6f5-a71a-4241-8936-30fcecea8cbf	source_support_kimi_frontier_xl_v1	supported	76	0	0	0
35cbf6f5-a71a-4241-8936-30fcecea8cbf	source_support_kimi_frontier_xl_v1	unsupported	24	0	0	0
3742139e-8fee-49e0-9e6a-74264863f56d	judge_kimi_heldout500_v1	correct	308	0	0	0
3742139e-8fee-49e0-9e6a-74264863f56d	judge_kimi_heldout500_v1	incorrect	191	0	0	0
3742139e-8fee-49e0-9e6a-74264863f56d	source_support_kimi_heldout500_v1	supported	311	0	0	0
3742139e-8fee-49e0-9e6a-74264863f56d	source_support_kimi_heldout500_v1	unsupported	188	0	0	0
374b1d89-9d27-4c5c-8fe2-1b1b45134662	judge_kimi_frontier_xl_v1	correct	98	0	0	0
374b1d89-9d27-4c5c-8fe2-1b1b45134662	judge_kimi_frontier_xl_v1	incorrect	2	0	0	0
374b1d89-9d27-4c5c-8fe2-1b1b45134662	source_support_kimi_frontier_xl_v1	supported	93	0	0	0
374b1d89-9d27-4c5c-8fe2-1b1b45134662	source_support_kimi_frontier_xl_v1	unsupported	7	0	0	0
38b99945-45d0-4156-aafa-8a8f0073be48	judge_kimi_heldout500_v1	correct	357	0	0	0
38b99945-45d0-4156-aafa-8a8f0073be48	judge_kimi_heldout500_v1	incorrect	142	0	0	0
38b99945-45d0-4156-aafa-8a8f0073be48	source_support_kimi_heldout500_v1	supported	357	0	0	0
38b99945-45d0-4156-aafa-8a8f0073be48	source_support_kimi_heldout500_v1	unsupported	142	0	0	0
3e5c61fc-6124-4a87-8eef-a8bf819fda9d	judge_kimi_heldout500_v1	correct	344	0	0	0
3e5c61fc-6124-4a87-8eef-a8bf819fda9d	judge_kimi_heldout500_v1	incorrect	155	0	0	0
3e5c61fc-6124-4a87-8eef-a8bf819fda9d	source_support_kimi_heldout500_v1	supported	344	0	0	0
3e5c61fc-6124-4a87-8eef-a8bf819fda9d	source_support_kimi_heldout500_v1	unsupported	155	0	0	0
3f4d3129-63bd-43f4-aecf-925ec8c9e8ea	judge_glm_frontier_xl_v1	correct	74	0	0	0
3f4d3129-63bd-43f4-aecf-925ec8c9e8ea	judge_glm_frontier_xl_v1	incorrect	26	0	0	0
3f4d3129-63bd-43f4-aecf-925ec8c9e8ea	source_support_glm_frontier_xl_v1	supported	77	0	0	0
3f4d3129-63bd-43f4-aecf-925ec8c9e8ea	source_support_glm_frontier_xl_v1	unsupported	23	0	0	0
3fb83e2c-0f72-4909-8ea8-011c2255de55	judge_kimi_heldout500_v1	correct	436	0	0	0
3fb83e2c-0f72-4909-8ea8-011c2255de55	judge_kimi_heldout500_v1	incorrect	63	0	0	0
3fb83e2c-0f72-4909-8ea8-011c2255de55	source_support_kimi_heldout500_v1	supported	435	0	0	0
3fb83e2c-0f72-4909-8ea8-011c2255de55	source_support_kimi_heldout500_v1	unsupported	64	0	0	0
42379df3-e4ec-4e16-a28a-7c573e59f61b	judge_kimi_depth_v1	correct	3	473	165	638
42379df3-e4ec-4e16-a28a-7c573e59f61b	judge_kimi_depth_v1	incorrect	2	292	87	379
42379df3-e4ec-4e16-a28a-7c573e59f61b	source_support_kimi_depth_v1	supported	3	679	150	829
42379df3-e4ec-4e16-a28a-7c573e59f61b	source_support_kimi_depth_v1	unsupported	2	0	0	0
42564cb1-a8ad-42d3-a08c-86942f08f35c	judge_kimi_depth_v1	correct	5	775	287	1062
42564cb1-a8ad-42d3-a08c-86942f08f35c	source_support_kimi_depth_v1	supported	4	3399	195	3594
42564cb1-a8ad-42d3-a08c-86942f08f35c	source_support_kimi_depth_v1	unsupported	1	731	58	789
4310f88e-bb1a-48ea-9209-c9ca00bbdc68	judge_flash_depth100_v1	correct	68	12965	2818	15783
4310f88e-bb1a-48ea-9209-c9ca00bbdc68	judge_flash_depth100_v1	incorrect	32	5717	1448	7165
4310f88e-bb1a-48ea-9209-c9ca00bbdc68	source_support_flash_depth100_v1	supported	60	48391	2466	50857
4310f88e-bb1a-48ea-9209-c9ca00bbdc68	source_support_flash_depth100_v1	unsupported	40	9300	558	9858
437b3aa8-04e9-47b7-a7c7-93340c94d854	judge_glm_frontier_xl_v1	correct	87	0	0	0
437b3aa8-04e9-47b7-a7c7-93340c94d854	judge_glm_frontier_xl_v1	incorrect	13	0	0	0
437b3aa8-04e9-47b7-a7c7-93340c94d854	source_support_glm_frontier_xl_v1	supported	86	0	0	0
437b3aa8-04e9-47b7-a7c7-93340c94d854	source_support_glm_frontier_xl_v1	unsupported	14	0	0	0
453ddf0c-0605-48b3-a6b0-1f60674f3dcf	judge_kimi_frontier_xl_v1	correct	87	0	0	0
453ddf0c-0605-48b3-a6b0-1f60674f3dcf	judge_kimi_frontier_xl_v1	incorrect	13	0	0	0
453ddf0c-0605-48b3-a6b0-1f60674f3dcf	source_support_kimi_frontier_xl_v1	supported	86	0	0	0
453ddf0c-0605-48b3-a6b0-1f60674f3dcf	source_support_kimi_frontier_xl_v1	unsupported	14	0	0	0
474a41c9-2a76-49c1-89e6-126b22b32db1	judge_flash_depth100_v1	correct	80	14724	3241	17965
474a41c9-2a76-49c1-89e6-126b22b32db1	judge_flash_depth100_v1	incorrect	20	3737	961	4698
474a41c9-2a76-49c1-89e6-126b22b32db1	source_support_flash_depth100_v1	supported	82	23551	3435	26986
474a41c9-2a76-49c1-89e6-126b22b32db1	source_support_flash_depth100_v1	unsupported	18	993	182	1175
4dd13da9-5773-451f-bd3f-28777aba9cb5	judge_flash_depth100_v1	correct	84	15590	3357	18947
4dd13da9-5773-451f-bd3f-28777aba9cb5	judge_flash_depth100_v1	incorrect	16	2991	750	3741
4dd13da9-5773-451f-bd3f-28777aba9cb5	source_support_flash_depth100_v1	supported	85	25597	3608	29205
4dd13da9-5773-451f-bd3f-28777aba9cb5	source_support_flash_depth100_v1	unsupported	15	1344	241	1585
4e1fe17f-deb1-4132-87d5-bdd056aa0d5d	judge_kimi_frontier_xl_v1	correct	65	0	0	0
4e1fe17f-deb1-4132-87d5-bdd056aa0d5d	judge_kimi_frontier_xl_v1	incorrect	35	0	0	0
4e1fe17f-deb1-4132-87d5-bdd056aa0d5d	source_support_kimi_frontier_xl_v1	supported	68	0	0	0
4e1fe17f-deb1-4132-87d5-bdd056aa0d5d	source_support_kimi_frontier_xl_v1	unsupported	32	0	0	0
4ee5a60b-0f88-4d20-b1e8-33c3751efa46	judge_kimi_depth_v1	correct	2	308	112	420
4ee5a60b-0f88-4d20-b1e8-33c3751efa46	judge_kimi_depth_v1	incorrect	3	449	115	564
4ee5a60b-0f88-4d20-b1e8-33c3751efa46	source_support_kimi_depth_v1	supported	2	450	95	545
4ee5a60b-0f88-4d20-b1e8-33c3751efa46	source_support_kimi_depth_v1	unsupported	3	0	0	0
4ef714a6-728a-4763-a99b-22b31ab1961f	judge_kimi_v1	correct	63	10936	3345	14281
4ef714a6-728a-4763-a99b-22b31ab1961f	judge_kimi_v1	incorrect	37	6019	1641	7660
4ef714a6-728a-4763-a99b-22b31ab1961f	source_support_kimi_v1	supported	63	48782	3494	52276
4ef714a6-728a-4763-a99b-22b31ab1961f	source_support_kimi_v1	unsupported	37	3121	244	3365
519b1c35-d314-4c3b-82a4-f60e7cbaa088	judge_kimi_frontier_xl_v1	correct	87	0	0	0
519b1c35-d314-4c3b-82a4-f60e7cbaa088	judge_kimi_frontier_xl_v1	incorrect	13	0	0	0
519b1c35-d314-4c3b-82a4-f60e7cbaa088	source_support_kimi_frontier_xl_v1	supported	85	0	0	0
519b1c35-d314-4c3b-82a4-f60e7cbaa088	source_support_kimi_frontier_xl_v1	unsupported	15	0	0	0
53cacbbf-fd3a-4846-af03-96c90f4bdf9f	judge_flash_depth100_v1	correct	69	12986	2785	15771
53cacbbf-fd3a-4846-af03-96c90f4bdf9f	judge_flash_depth100_v1	incorrect	31	5599	1410	7009
53cacbbf-fd3a-4846-af03-96c90f4bdf9f	source_support_flash_depth100_v1	supported	64	61913	2635	64548
53cacbbf-fd3a-4846-af03-96c90f4bdf9f	source_support_flash_depth100_v1	unsupported	36	8818	378	9196
5442ef5f-dabb-46c8-af8e-e5e19ed1d210	judge_flash_depth100_v1	correct	57	11132	2406	13538
5442ef5f-dabb-46c8-af8e-e5e19ed1d210	judge_flash_depth100_v1	incorrect	43	7690	1961	9651
5442ef5f-dabb-46c8-af8e-e5e19ed1d210	source_support_flash_depth100_v1	supported	51	41100	2056	43156
5442ef5f-dabb-46c8-af8e-e5e19ed1d210	source_support_flash_depth100_v1	unsupported	49	9818	559	10377
57834359-fbc2-4b7e-a776-f74df6b53df6	judge_kimi_frontier_xl_v1	correct	68	0	0	0
57834359-fbc2-4b7e-a776-f74df6b53df6	judge_kimi_frontier_xl_v1	incorrect	32	0	0	0
57834359-fbc2-4b7e-a776-f74df6b53df6	source_support_kimi_frontier_xl_v1	supported	68	0	0	0
57834359-fbc2-4b7e-a776-f74df6b53df6	source_support_kimi_frontier_xl_v1	unsupported	32	0	0	0
57ce84ee-cf8c-4863-a7cd-463ace593dff	judge_kimi_heldout500_v1	correct	362	0	0	0
57ce84ee-cf8c-4863-a7cd-463ace593dff	judge_kimi_heldout500_v1	incorrect	137	0	0	0
57ce84ee-cf8c-4863-a7cd-463ace593dff	source_support_kimi_heldout500_v1	supported	361	0	0	0
57ce84ee-cf8c-4863-a7cd-463ace593dff	source_support_kimi_heldout500_v1	unsupported	138	0	0	0
59ab526c-3301-4930-aa93-7f9233ce2453	judge_kimi_heldout500_v1	correct	340	0	0	0
59ab526c-3301-4930-aa93-7f9233ce2453	judge_kimi_heldout500_v1	incorrect	159	0	0	0
59ab526c-3301-4930-aa93-7f9233ce2453	source_support_kimi_heldout500_v1	supported	336	0	0	0
59ab526c-3301-4930-aa93-7f9233ce2453	source_support_kimi_heldout500_v1	unsupported	163	0	0	0
5c327145-46f1-463b-9829-61d0544c8b91	judge_kimi_depth_v1	correct	3	488	184	672
5c327145-46f1-463b-9829-61d0544c8b91	judge_kimi_depth_v1	incorrect	2	292	79	371
5c327145-46f1-463b-9829-61d0544c8b91	source_support_kimi_depth_v1	supported	3	919	153	1072
5c327145-46f1-463b-9829-61d0544c8b91	source_support_kimi_depth_v1	unsupported	2	0	0	0
5cee9df7-76e8-4ebd-943f-22356adf1e3f	judge_kimi_v1	correct	46	7805	2513	10318
5cee9df7-76e8-4ebd-943f-22356adf1e3f	judge_kimi_v1	incorrect	54	8884	2496	11380
5cee9df7-76e8-4ebd-943f-22356adf1e3f	source_support_kimi_v1	supported	44	38350	2422	40772
5cee9df7-76e8-4ebd-943f-22356adf1e3f	source_support_kimi_v1	unsupported	56	4675	350	5025
5cf07cfc-a135-4080-b7e3-92c19c58d992	judge_kimi_frontier_xl_v1	correct	88	0	0	0
5cf07cfc-a135-4080-b7e3-92c19c58d992	judge_kimi_frontier_xl_v1	incorrect	12	0	0	0
5cf07cfc-a135-4080-b7e3-92c19c58d992	source_support_kimi_frontier_xl_v1	supported	83	0	0	0
5cf07cfc-a135-4080-b7e3-92c19c58d992	source_support_kimi_frontier_xl_v1	unsupported	17	0	0	0
5f5eab39-ac46-4f6f-8ae9-f08da1c1c07a	judge_kimi_depth_v1	correct	4	639	224	863
5f5eab39-ac46-4f6f-8ae9-f08da1c1c07a	judge_kimi_depth_v1	incorrect	1	148	39	187
5f5eab39-ac46-4f6f-8ae9-f08da1c1c07a	source_support_kimi_depth_v1	supported	3	1916	206	2122
5f5eab39-ac46-4f6f-8ae9-f08da1c1c07a	source_support_kimi_depth_v1	unsupported	2	777	40	817
655b8a6c-0ff5-4d0e-9c39-e99c1a56e1f6	judge_kimi_heldout500_v1	correct	233	0	0	0
655b8a6c-0ff5-4d0e-9c39-e99c1a56e1f6	judge_kimi_heldout500_v1	incorrect	266	0	0	0
655b8a6c-0ff5-4d0e-9c39-e99c1a56e1f6	source_support_kimi_heldout500_v1	supported	235	0	0	0
655b8a6c-0ff5-4d0e-9c39-e99c1a56e1f6	source_support_kimi_heldout500_v1	unsupported	264	0	0	0
68e765eb-4053-433a-8d63-b8f238f2b5e0	judge_flash_depth100_v1	correct	75	13828	3131	16959
68e765eb-4053-433a-8d63-b8f238f2b5e0	judge_flash_depth100_v1	incorrect	25	4633	1201	5834
68e765eb-4053-433a-8d63-b8f238f2b5e0	source_support_flash_depth100_v1	supported	77	21499	3231	24730
68e765eb-4053-433a-8d63-b8f238f2b5e0	source_support_flash_depth100_v1	unsupported	23	733	134	867
69de4988-5157-4cec-8f2f-869381b47432	judge_kimi_heldout500_v1	correct	421	0	0	0
69de4988-5157-4cec-8f2f-869381b47432	judge_kimi_heldout500_v1	incorrect	78	0	0	0
69de4988-5157-4cec-8f2f-869381b47432	source_support_kimi_heldout500_v1	supported	409	0	0	0
69de4988-5157-4cec-8f2f-869381b47432	source_support_kimi_heldout500_v1	unsupported	90	0	0	0
6b0106e8-8ec8-473e-9989-b225bfc5ff0b	judge	correct	64	13202	2767	15969
6b0106e8-8ec8-473e-9989-b225bfc5ff0b	judge	incorrect	36	8605	1903	10508
6fd7aad9-0500-45d0-8593-9bf1ddeb10c7	judge_kimi_depth_v1	correct	3	455	150	605
6fd7aad9-0500-45d0-8593-9bf1ddeb10c7	judge_kimi_depth_v1	incorrect	2	331	67	398
6fd7aad9-0500-45d0-8593-9bf1ddeb10c7	source_support_kimi_depth_v1	supported	4	1352	202	1554
6fd7aad9-0500-45d0-8593-9bf1ddeb10c7	source_support_kimi_depth_v1	unsupported	1	0	0	0
70e5b79f-263a-406b-bec9-8570ad56dc4b	judge_kimi_heldout500_v1	correct	320	0	0	0
70e5b79f-263a-406b-bec9-8570ad56dc4b	judge_kimi_heldout500_v1	incorrect	179	0	0	0
70e5b79f-263a-406b-bec9-8570ad56dc4b	source_support_kimi_heldout500_v1	supported	316	0	0	0
70e5b79f-263a-406b-bec9-8570ad56dc4b	source_support_kimi_heldout500_v1	unsupported	183	0	0	0
711eed2f-8550-4a29-88a7-03b8ebc430f3	judge_kimi_heldout500_v1	correct	326	0	0	0
711eed2f-8550-4a29-88a7-03b8ebc430f3	judge_kimi_heldout500_v1	incorrect	173	0	0	0
711eed2f-8550-4a29-88a7-03b8ebc430f3	source_support_kimi_heldout500_v1	supported	322	0	0	0
711eed2f-8550-4a29-88a7-03b8ebc430f3	source_support_kimi_heldout500_v1	unsupported	177	0	0	0
737b9df9-158c-45c8-bf16-6ca2dbc3bdaf	judge_kimi_frontier_xl_v1	correct	62	0	0	0
737b9df9-158c-45c8-bf16-6ca2dbc3bdaf	judge_kimi_frontier_xl_v1	incorrect	38	0	0	0
737b9df9-158c-45c8-bf16-6ca2dbc3bdaf	source_support_kimi_frontier_xl_v1	supported	63	0	0	0
737b9df9-158c-45c8-bf16-6ca2dbc3bdaf	source_support_kimi_frontier_xl_v1	unsupported	37	0	0	0
76715c27-806e-4bf7-a202-e24de7f596c0	judge_kimi_depth_v1	correct	4	639	243	882
76715c27-806e-4bf7-a202-e24de7f596c0	judge_kimi_depth_v1	incorrect	1	148	39	187
76715c27-806e-4bf7-a202-e24de7f596c0	source_support_kimi_depth_v1	supported	3	2447	142	2589
76715c27-806e-4bf7-a202-e24de7f596c0	source_support_kimi_depth_v1	unsupported	2	640	43	683
76e177e5-b41c-4817-9ec2-9f7efc295922	judge_flash_v2	correct	73	13366	2977	16343
76e177e5-b41c-4817-9ec2-9f7efc295922	judge_flash_v2	incorrect	27	5586	1474	7060
76e177e5-b41c-4817-9ec2-9f7efc295922	judge_kimi_probe	correct	1	178	65	243
76e177e5-b41c-4817-9ec2-9f7efc295922	judge_kimi_v1	correct	74	12324	4032	16356
76e177e5-b41c-4817-9ec2-9f7efc295922	judge_kimi_v1	incorrect	26	5045	1509	6554
76e177e5-b41c-4817-9ec2-9f7efc295922	source_support_flash_v1	supported	73	19896	3044	22940
76e177e5-b41c-4817-9ec2-9f7efc295922	source_support_flash_v1	unsupported	27	3445	549	3994
76e177e5-b41c-4817-9ec2-9f7efc295922	source_support_kimi_v1	supported	77	19658	4019	23677
76e177e5-b41c-4817-9ec2-9f7efc295922	source_support_kimi_v1	unsupported	23	2068	450	2518
77b616ce-fc29-46db-9339-ac313885f878	judge_kimi_frontier_xl_v1	correct	79	0	0	0
77b616ce-fc29-46db-9339-ac313885f878	judge_kimi_frontier_xl_v1	incorrect	21	0	0	0
77b616ce-fc29-46db-9339-ac313885f878	source_support_kimi_frontier_xl_v1	supported	81	0	0	0
77b616ce-fc29-46db-9339-ac313885f878	source_support_kimi_frontier_xl_v1	unsupported	19	0	0	0
7b587744-cd46-4764-aefb-cb7639456fa8	judge_kimi_heldout500_v1	correct	332	0	0	0
7b587744-cd46-4764-aefb-cb7639456fa8	judge_kimi_heldout500_v1	incorrect	167	0	0	0
7b587744-cd46-4764-aefb-cb7639456fa8	source_support_kimi_heldout500_v1	supported	329	0	0	0
7b587744-cd46-4764-aefb-cb7639456fa8	source_support_kimi_heldout500_v1	unsupported	170	0	0	0
7ba0eaa3-c3d8-44e8-8aae-c8c779034c1e	judge_kimi_depth_v1	correct	5	784	242	1026
7ba0eaa3-c3d8-44e8-8aae-c8c779034c1e	source_support_kimi_depth_v1	supported	4	3178	172	3350
7ba0eaa3-c3d8-44e8-8aae-c8c779034c1e	source_support_kimi_depth_v1	unsupported	1	640	43	683
8155633d-4c03-4da7-9bbb-66a344361bfc	judge_kimi_heldout500_v1	correct	296	0	0	0
8155633d-4c03-4da7-9bbb-66a344361bfc	judge_kimi_heldout500_v1	incorrect	203	0	0	0
8155633d-4c03-4da7-9bbb-66a344361bfc	source_support_kimi_heldout500_v1	supported	303	0	0	0
8155633d-4c03-4da7-9bbb-66a344361bfc	source_support_kimi_heldout500_v1	unsupported	196	0	0	0
85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865	judge_kimi_frontier_xl_v1	correct	70	0	0	0
85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865	judge_kimi_frontier_xl_v1	incorrect	30	0	0	0
85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865	source_support_kimi_frontier_xl_v1	supported	71	0	0	0
85a4d5ef-7d2b-4a8a-9f05-0bd4b283c865	source_support_kimi_frontier_xl_v1	unsupported	29	0	0	0
897e0b71-7232-4872-af0e-b10b228dc6ca	judge_flash_depth100_v1	correct	68	12981	2745	15726
897e0b71-7232-4872-af0e-b10b228dc6ca	judge_flash_depth100_v1	incorrect	32	5730	1465	7195
897e0b71-7232-4872-af0e-b10b228dc6ca	source_support_flash_depth100_v1	supported	62	52220	2601	54821
897e0b71-7232-4872-af0e-b10b228dc6ca	source_support_flash_depth100_v1	unsupported	38	8291	505	8796
8fc06526-982a-4c7b-9518-1ea213f5fa05	judge_kimi_v1	correct	23	3642	1075	4717
8fc06526-982a-4c7b-9518-1ea213f5fa05	judge_kimi_v1	incorrect	77	12756	3515	16271
8fc06526-982a-4c7b-9518-1ea213f5fa05	source_support_kimi_v1	supported	24	19609	1310	20919
8fc06526-982a-4c7b-9518-1ea213f5fa05	source_support_kimi_v1	unsupported	76	3201	226	3427
92a23195-e90c-4dc0-8dc7-b7c25d20f77a	judge	correct	45	9578	2020	11598
92a23195-e90c-4dc0-8dc7-b7c25d20f77a	judge	incorrect	55	13176	2889	16065
92a84632-23f1-4a96-8634-51a94124997b	judge_glm_frontier_xl_v1	correct	91	0	0	0
92a84632-23f1-4a96-8634-51a94124997b	judge_glm_frontier_xl_v1	incorrect	9	0	0	0
92a84632-23f1-4a96-8634-51a94124997b	source_support_glm_frontier_xl_v1	supported	90	0	0	0
92a84632-23f1-4a96-8634-51a94124997b	source_support_glm_frontier_xl_v1	unsupported	10	0	0	0
937f1cbd-c964-4995-9f06-58ea7e44841b	judge_flash_depth100_v1	correct	58	10996	2377	13373
937f1cbd-c964-4995-9f06-58ea7e44841b	judge_flash_depth100_v1	incorrect	42	7622	1920	9542
937f1cbd-c964-4995-9f06-58ea7e44841b	source_support_flash_depth100_v1	supported	51	47383	2085	49468
937f1cbd-c964-4995-9f06-58ea7e44841b	source_support_flash_depth100_v1	unsupported	49	10634	561	11195
95b8f4e3-3718-4c70-8f1c-3de34e9aab93	judge_kimi_frontier_xl_v1	correct	88	0	0	0
95b8f4e3-3718-4c70-8f1c-3de34e9aab93	judge_kimi_frontier_xl_v1	incorrect	12	0	0	0
95b8f4e3-3718-4c70-8f1c-3de34e9aab93	source_support_kimi_frontier_xl_v1	supported	89	0	0	0
95b8f4e3-3718-4c70-8f1c-3de34e9aab93	source_support_kimi_frontier_xl_v1	unsupported	11	0	0	0
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0	judge_kimi_frontier_xl_v1	correct	78	0	0	0
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0	judge_kimi_frontier_xl_v1	incorrect	22	0	0	0
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0	source_support_kimi_frontier_xl_v1	supported	80	0	0	0
9694cde6-94a4-4cd7-88f7-e6e9903cd8a0	source_support_kimi_frontier_xl_v1	unsupported	20	0	0	0
99b978b9-966b-44a9-96b2-ca3fe4857eac	judge_kimi_frontier_xl_v1	correct	64	0	0	0
99b978b9-966b-44a9-96b2-ca3fe4857eac	judge_kimi_frontier_xl_v1	incorrect	36	0	0	0
99b978b9-966b-44a9-96b2-ca3fe4857eac	source_support_kimi_frontier_xl_v1	supported	63	0	0	0
99b978b9-966b-44a9-96b2-ca3fe4857eac	source_support_kimi_frontier_xl_v1	unsupported	37	0	0	0
9bb5148e-9243-4498-8378-b5c1e6d67f68	judge_kimi_depth_v1	correct	4	624	221	845
9bb5148e-9243-4498-8378-b5c1e6d67f68	judge_kimi_depth_v1	incorrect	1	148	39	187
9bb5148e-9243-4498-8378-b5c1e6d67f68	source_support_kimi_depth_v1	supported	3	2447	141	2588
9bb5148e-9243-4498-8378-b5c1e6d67f68	source_support_kimi_depth_v1	unsupported	2	625	67	692
a00fb789-a14d-4d9a-9285-61d4af91f005	judge_kimi_frontier_xl_v1	correct	92	0	0	0
a00fb789-a14d-4d9a-9285-61d4af91f005	judge_kimi_frontier_xl_v1	incorrect	8	0	0	0
a00fb789-a14d-4d9a-9285-61d4af91f005	source_support_kimi_frontier_xl_v1	supported	88	0	0	0
a00fb789-a14d-4d9a-9285-61d4af91f005	source_support_kimi_frontier_xl_v1	unsupported	12	0	0	0
a0e9f6c9-7031-4995-a709-5c3db18ac137	judge_kimi_v1	correct	21	3376	968	4344
a0e9f6c9-7031-4995-a709-5c3db18ac137	judge_kimi_v1	incorrect	79	13102	3567	16669
a0e9f6c9-7031-4995-a709-5c3db18ac137	source_support_kimi_v1	supported	23	18428	1234	19662
a0e9f6c9-7031-4995-a709-5c3db18ac137	source_support_kimi_v1	unsupported	77	2138	170	2308
a3abb2f5-623d-4e3a-ac28-1b9fb66b4143	judge_glm_frontier_xl_v1	correct	90	0	0	0
a3abb2f5-623d-4e3a-ac28-1b9fb66b4143	judge_glm_frontier_xl_v1	incorrect	10	0	0	0
a3abb2f5-623d-4e3a-ac28-1b9fb66b4143	source_support_glm_frontier_xl_v1	supported	91	0	0	0
a3abb2f5-623d-4e3a-ac28-1b9fb66b4143	source_support_glm_frontier_xl_v1	unsupported	9	0	0	0
a5aee811-f912-475c-9f86-b1811fc36164	judge_flash_depth100_v1	correct	69	12873	2737	15610
a5aee811-f912-475c-9f86-b1811fc36164	judge_flash_depth100_v1	incorrect	31	5641	1434	7075
a5aee811-f912-475c-9f86-b1811fc36164	source_support_flash_depth100_v1	supported	63	59295	2582	61877
a5aee811-f912-475c-9f86-b1811fc36164	source_support_flash_depth100_v1	unsupported	37	9240	443	9683
a72590e0-ecc7-40ab-b47b-de381f8b948f	judge_kimi_frontier_xl_v1	correct	75	0	0	0
a72590e0-ecc7-40ab-b47b-de381f8b948f	judge_kimi_frontier_xl_v1	incorrect	25	0	0	0
a72590e0-ecc7-40ab-b47b-de381f8b948f	source_support_kimi_frontier_xl_v1	supported	74	0	0	0
a72590e0-ecc7-40ab-b47b-de381f8b948f	source_support_kimi_frontier_xl_v1	unsupported	26	0	0	0
ab0104a4-1b1f-4282-9cd0-2736e0704504	judge_kimi_depth_v1	correct	5	793	279	1072
ab0104a4-1b1f-4282-9cd0-2736e0704504	source_support_kimi_depth_v1	supported	4	2647	215	2862
ab0104a4-1b1f-4282-9cd0-2736e0704504	source_support_kimi_depth_v1	unsupported	1	777	57	834
ada7f9b2-cc08-4524-a280-fcdfd2a37e0f	judge_kimi_frontier_xl_v1	correct	59	0	0	0
ada7f9b2-cc08-4524-a280-fcdfd2a37e0f	judge_kimi_frontier_xl_v1	incorrect	41	0	0	0
ada7f9b2-cc08-4524-a280-fcdfd2a37e0f	source_support_kimi_frontier_xl_v1	supported	59	0	0	0
ada7f9b2-cc08-4524-a280-fcdfd2a37e0f	source_support_kimi_frontier_xl_v1	unsupported	41	0	0	0
adbfce87-cac8-4aed-80e8-d3262a094b41	judge_kimi_frontier_xl_v1	correct	86	0	0	0
adbfce87-cac8-4aed-80e8-d3262a094b41	judge_kimi_frontier_xl_v1	incorrect	14	0	0	0
adbfce87-cac8-4aed-80e8-d3262a094b41	source_support_kimi_frontier_xl_v1	supported	83	0	0	0
adbfce87-cac8-4aed-80e8-d3262a094b41	source_support_kimi_frontier_xl_v1	unsupported	17	0	0	0
af62b994-b63b-428c-a386-e884f2c93a52	judge_kimi_v1	correct	64	11153	3397	14550
af62b994-b63b-428c-a386-e884f2c93a52	judge_kimi_v1	incorrect	36	5767	1645	7412
af62b994-b63b-428c-a386-e884f2c93a52	source_support_kimi_v1	supported	64	58993	3529	62522
af62b994-b63b-428c-a386-e884f2c93a52	source_support_kimi_v1	unsupported	36	3458	240	3698
b13cc9e8-1615-40ce-a13d-00021ee64647	judge_kimi_frontier_xl_v1	correct	75	0	0	0
b13cc9e8-1615-40ce-a13d-00021ee64647	judge_kimi_frontier_xl_v1	incorrect	25	0	0	0
b13cc9e8-1615-40ce-a13d-00021ee64647	source_support_kimi_frontier_xl_v1	supported	76	0	0	0
b13cc9e8-1615-40ce-a13d-00021ee64647	source_support_kimi_frontier_xl_v1	unsupported	24	0	0	0
b17a11c2-92c1-4a1c-8b98-d48774f6d1fc	judge_kimi_depth_v1	correct	4	639	248	887
b17a11c2-92c1-4a1c-8b98-d48774f6d1fc	judge_kimi_depth_v1	incorrect	1	148	40	188
b17a11c2-92c1-4a1c-8b98-d48774f6d1fc	source_support_kimi_depth_v1	supported	3	2677	152	2829
b17a11c2-92c1-4a1c-8b98-d48774f6d1fc	source_support_kimi_depth_v1	unsupported	2	731	43	774
b2e31edc-f490-4eb2-b834-e080fba87dd5	judge_kimi_frontier_xl_v1	correct	92	0	0	0
b2e31edc-f490-4eb2-b834-e080fba87dd5	judge_kimi_frontier_xl_v1	incorrect	8	0	0	0
b2e31edc-f490-4eb2-b834-e080fba87dd5	source_support_kimi_frontier_xl_v1	supported	86	0	0	0
b2e31edc-f490-4eb2-b834-e080fba87dd5	source_support_kimi_frontier_xl_v1	unsupported	14	0	0	0
babb185f-a8bd-4b0b-81b9-b5a726fc0b00	judge_flash_depth100_v1	correct	69	12905	2707	15612
babb185f-a8bd-4b0b-81b9-b5a726fc0b00	judge_flash_depth100_v1	incorrect	31	5675	1392	7067
babb185f-a8bd-4b0b-81b9-b5a726fc0b00	source_support_flash_depth100_v1	supported	65	61297	2702	63999
babb185f-a8bd-4b0b-81b9-b5a726fc0b00	source_support_flash_depth100_v1	unsupported	35	6048	308	6356
bb02383f-926c-4513-af90-23d83d6af31f	judge_flash_depth100_v1	correct	55	10495	2272	12767
bb02383f-926c-4513-af90-23d83d6af31f	judge_flash_depth100_v1	incorrect	45	8160	2066	10226
bb02383f-926c-4513-af90-23d83d6af31f	source_support_flash_depth100_v1	supported	52	48705	2135	50840
bb02383f-926c-4513-af90-23d83d6af31f	source_support_flash_depth100_v1	unsupported	48	8277	451	8728
bbb15eb7-e886-48ce-a081-f7175f32533e	judge_glm_frontier_xl_v1	correct	70	0	0	0
bbb15eb7-e886-48ce-a081-f7175f32533e	judge_glm_frontier_xl_v1	incorrect	30	0	0	0
bbb15eb7-e886-48ce-a081-f7175f32533e	source_support_glm_frontier_xl_v1	supported	73	0	0	0
bbb15eb7-e886-48ce-a081-f7175f32533e	source_support_glm_frontier_xl_v1	unsupported	27	0	0	0
bc19fd4a-ff74-4e27-8fe3-4d4043591e6d	judge_flash_depth100_v1	correct	67	12662	2721	15383
bc19fd4a-ff74-4e27-8fe3-4d4043591e6d	judge_flash_depth100_v1	incorrect	33	5923	1478	7401
bc19fd4a-ff74-4e27-8fe3-4d4043591e6d	source_support_flash_depth100_v1	supported	62	58166	2573	60739
bc19fd4a-ff74-4e27-8fe3-4d4043591e6d	source_support_flash_depth100_v1	unsupported	38	7675	366	8041
bc241efa-dd77-41cd-8d47-f250b0507c10	judge_kimi_v1	correct	60	10396	3287	13683
bc241efa-dd77-41cd-8d47-f250b0507c10	judge_kimi_v1	incorrect	40	6555	1873	8428
bc241efa-dd77-41cd-8d47-f250b0507c10	source_support_kimi_v1	supported	59	53290	3092	56382
bc241efa-dd77-41cd-8d47-f250b0507c10	source_support_kimi_v1	unsupported	41	3917	368	4285
bfebee8c-6fa7-4cf2-8557-9995a9c7cc51	judge_kimi_heldout500_v1	correct	439	0	0	0
bfebee8c-6fa7-4cf2-8557-9995a9c7cc51	judge_kimi_heldout500_v1	incorrect	60	0	0	0
bfebee8c-6fa7-4cf2-8557-9995a9c7cc51	source_support_kimi_heldout500_v1	supported	434	0	0	0
bfebee8c-6fa7-4cf2-8557-9995a9c7cc51	source_support_kimi_heldout500_v1	unsupported	65	0	0	0
c09be366-1289-4585-83b1-1436bba02f5e	judge_flash_depth100_v1	correct	56	9972	2257	12229
c09be366-1289-4585-83b1-1436bba02f5e	judge_flash_depth100_v1	incorrect	44	8293	2201	10494
c09be366-1289-4585-83b1-1436bba02f5e	source_support_flash_depth100_v1	supported	51	12581	2110	14691
c09be366-1289-4585-83b1-1436bba02f5e	source_support_flash_depth100_v1	unsupported	49	3903	785	4688
c12c67c1-923b-48c6-b672-c871a8b2d01f	judge_flash_depth100_v1	correct	72	13569	2913	16482
c12c67c1-923b-48c6-b672-c871a8b2d01f	judge_flash_depth100_v1	incorrect	28	5092	1267	6359
c12c67c1-923b-48c6-b672-c871a8b2d01f	source_support_flash_depth100_v1	supported	67	64515	2824	67339
c12c67c1-923b-48c6-b672-c871a8b2d01f	source_support_flash_depth100_v1	unsupported	33	6571	308	6879
c1717381-4259-4695-b030-ca2b62ce6379	judge_kimi_depth_v1	correct	3	455	134	589
c1717381-4259-4695-b030-ca2b62ce6379	judge_kimi_depth_v1	incorrect	2	329	71	400
c1717381-4259-4695-b030-ca2b62ce6379	source_support_kimi_depth_v1	supported	4	1271	232	1503
c1717381-4259-4695-b030-ca2b62ce6379	source_support_kimi_depth_v1	unsupported	1	0	0	0
c6b10ae2-e9b7-4ecc-ab49-d76eb92834a0	judge_kimi_heldout500_v1	correct	323	0	0	0
c6b10ae2-e9b7-4ecc-ab49-d76eb92834a0	judge_kimi_heldout500_v1	incorrect	176	0	0	0
c6b10ae2-e9b7-4ecc-ab49-d76eb92834a0	source_support_kimi_heldout500_v1	supported	328	0	0	0
c6b10ae2-e9b7-4ecc-ab49-d76eb92834a0	source_support_kimi_heldout500_v1	unsupported	171	0	0	0
c839fc4c-37ac-44fb-b59c-bdd71fc597c0	judge_kimi_heldout500_v1	correct	323	0	0	0
c839fc4c-37ac-44fb-b59c-bdd71fc597c0	judge_kimi_heldout500_v1	incorrect	176	0	0	0
c839fc4c-37ac-44fb-b59c-bdd71fc597c0	source_support_kimi_heldout500_v1	supported	325	0	0	0
c839fc4c-37ac-44fb-b59c-bdd71fc597c0	source_support_kimi_heldout500_v1	unsupported	174	0	0	0
c8edfbad-3e27-4e1c-9a49-05183af4c6e4	judge_kimi_depth_v1	correct	3	475	162	637
c8edfbad-3e27-4e1c-9a49-05183af4c6e4	judge_kimi_depth_v1	incorrect	2	296	82	378
c8edfbad-3e27-4e1c-9a49-05183af4c6e4	source_support_kimi_depth_v1	supported	2	1240	89	1329
c8edfbad-3e27-4e1c-9a49-05183af4c6e4	source_support_kimi_depth_v1	unsupported	3	762	70	832
c916bb0a-81d1-4ffd-b555-0fed0b547cf4	judge_kimi_heldout500_v1	correct	377	0	0	0
c916bb0a-81d1-4ffd-b555-0fed0b547cf4	judge_kimi_heldout500_v1	incorrect	122	0	0	0
c916bb0a-81d1-4ffd-b555-0fed0b547cf4	source_support_kimi_heldout500_v1	supported	374	0	0	0
c916bb0a-81d1-4ffd-b555-0fed0b547cf4	source_support_kimi_heldout500_v1	unsupported	125	0	0	0
cab8c518-1f86-40be-b36d-a3c45fc328aa	judge_kimi_depth_v1	correct	4	624	219	843
cab8c518-1f86-40be-b36d-a3c45fc328aa	judge_kimi_depth_v1	incorrect	1	148	40	188
cab8c518-1f86-40be-b36d-a3c45fc328aa	source_support_kimi_depth_v1	supported	3	1916	167	2083
cab8c518-1f86-40be-b36d-a3c45fc328aa	source_support_kimi_depth_v1	unsupported	2	762	93	855
caf8170f-d920-4ad3-a03a-d49b87b38350	judge_kimi_heldout500_v1	correct	310	0	0	0
caf8170f-d920-4ad3-a03a-d49b87b38350	judge_kimi_heldout500_v1	incorrect	189	0	0	0
caf8170f-d920-4ad3-a03a-d49b87b38350	source_support_kimi_heldout500_v1	supported	316	0	0	0
caf8170f-d920-4ad3-a03a-d49b87b38350	source_support_kimi_heldout500_v1	unsupported	183	0	0	0
cd2d33c4-7125-42be-a7a3-5176ba984c64	judge_kimi_heldout500_v1	correct	415	0	0	0
cd2d33c4-7125-42be-a7a3-5176ba984c64	judge_kimi_heldout500_v1	incorrect	84	0	0	0
cd2d33c4-7125-42be-a7a3-5176ba984c64	source_support_kimi_heldout500_v1	supported	400	0	0	0
cd2d33c4-7125-42be-a7a3-5176ba984c64	source_support_kimi_heldout500_v1	unsupported	99	0	0	0
cfc6c85e-7644-451a-9eea-6e013a06526e	judge_kimi_depth_v1	correct	4	630	222	852
cfc6c85e-7644-451a-9eea-6e013a06526e	judge_kimi_depth_v1	incorrect	1	148	39	187
cfc6c85e-7644-451a-9eea-6e013a06526e	source_support_kimi_depth_v1	supported	3	2668	151	2819
cfc6c85e-7644-451a-9eea-6e013a06526e	source_support_kimi_depth_v1	unsupported	2	731	43	774
d1009eff-138b-4dff-96fb-0471bac724ae	judge_kimi_depth_v1	correct	5	786	286	1072
d1009eff-138b-4dff-96fb-0471bac724ae	source_support_kimi_depth_v1	supported	4	3180	177	3357
d1009eff-138b-4dff-96fb-0471bac724ae	source_support_kimi_depth_v1	unsupported	1	640	43	683
d189aa2b-04be-4270-8b61-5fe9eb8bc86b	judge_kimi_v1	correct	28	4615	1340	5955
d189aa2b-04be-4270-8b61-5fe9eb8bc86b	judge_kimi_v1	incorrect	72	11814	3232	15046
d189aa2b-04be-4270-8b61-5fe9eb8bc86b	source_support_kimi_v1	supported	29	21434	1503	22937
d189aa2b-04be-4270-8b61-5fe9eb8bc86b	source_support_kimi_v1	unsupported	71	1515	210	1725
d40d723b-dda1-49fd-a72b-1fa36c058171	judge_kimi_depth_v1	correct	5	786	275	1061
d40d723b-dda1-49fd-a72b-1fa36c058171	source_support_kimi_depth_v1	supported	4	2640	204	2844
d40d723b-dda1-49fd-a72b-1fa36c058171	source_support_kimi_depth_v1	unsupported	1	777	67	844
d43882dc-7339-44ac-9cd0-75f67dfc669f	judge_kimi_heldout500_v1	correct	369	0	0	0
d43882dc-7339-44ac-9cd0-75f67dfc669f	judge_kimi_heldout500_v1	incorrect	130	0	0	0
d43882dc-7339-44ac-9cd0-75f67dfc669f	source_support_kimi_heldout500_v1	supported	365	0	0	0
d43882dc-7339-44ac-9cd0-75f67dfc669f	source_support_kimi_heldout500_v1	unsupported	134	0	0	0
d728cf0d-8e04-4f80-8487-a1876ee25aa4	judge_kimi_heldout500_v1	correct	344	0	0	0
d728cf0d-8e04-4f80-8487-a1876ee25aa4	judge_kimi_heldout500_v1	incorrect	155	0	0	0
d728cf0d-8e04-4f80-8487-a1876ee25aa4	source_support_kimi_heldout500_v1	supported	345	0	0	0
d728cf0d-8e04-4f80-8487-a1876ee25aa4	source_support_kimi_heldout500_v1	unsupported	154	0	0	0
d997c177-8802-498b-8220-bc036950d847	judge_kimi_heldout500_v1	correct	313	0	0	0
d997c177-8802-498b-8220-bc036950d847	judge_kimi_heldout500_v1	incorrect	186	0	0	0
d997c177-8802-498b-8220-bc036950d847	source_support_kimi_heldout500_v1	supported	324	0	0	0
d997c177-8802-498b-8220-bc036950d847	source_support_kimi_heldout500_v1	unsupported	175	0	0	0
d998bf41-0cb9-4291-9806-52e8f2f8aa93	judge	correct	85	18511	3640	22151
d998bf41-0cb9-4291-9806-52e8f2f8aa93	judge	incorrect	15	4277	790	5067
dab29081-c149-4abb-816d-5bd97ecfe202	judge_flash_depth100_v1	correct	72	13738	2896	16634
dab29081-c149-4abb-816d-5bd97ecfe202	judge_flash_depth100_v1	incorrect	28	5055	1281	6336
dab29081-c149-4abb-816d-5bd97ecfe202	source_support_flash_depth100_v1	supported	67	56760	2773	59533
dab29081-c149-4abb-816d-5bd97ecfe202	source_support_flash_depth100_v1	unsupported	33	6819	407	7226
de00b901-f068-4553-98a6-5de1c92f9815	judge_kimi_heldout500_v1	correct	371	0	0	0
de00b901-f068-4553-98a6-5de1c92f9815	judge_kimi_heldout500_v1	incorrect	128	0	0	0
de00b901-f068-4553-98a6-5de1c92f9815	source_support_kimi_heldout500_v1	supported	365	0	0	0
de00b901-f068-4553-98a6-5de1c92f9815	source_support_kimi_heldout500_v1	unsupported	134	0	0	0
dfa0dd13-d115-4b0e-b453-8ffdd354a66b	judge_kimi_v1	correct	87	14831	4672	19503
dfa0dd13-d115-4b0e-b453-8ffdd354a66b	judge_kimi_v1	incorrect	13	2521	795	3316
dfa0dd13-d115-4b0e-b453-8ffdd354a66b	source_support_kimi_v1	supported	86	23766	4669	28435
dfa0dd13-d115-4b0e-b453-8ffdd354a66b	source_support_kimi_v1	unsupported	14	1836	526	2362
e1540d49-375a-4ebd-8385-98f2cb976156	judge_kimi_heldout500_v1	correct	417	0	0	0
e1540d49-375a-4ebd-8385-98f2cb976156	judge_kimi_heldout500_v1	incorrect	82	0	0	0
e1540d49-375a-4ebd-8385-98f2cb976156	source_support_kimi_heldout500_v1	supported	404	0	0	0
e1540d49-375a-4ebd-8385-98f2cb976156	source_support_kimi_heldout500_v1	unsupported	95	0	0	0
e1f2ed06-6d61-4071-b73a-7bf7ae8fedd7	judge_kimi_depth_v1	correct	5	786	294	1080
e1f2ed06-6d61-4071-b73a-7bf7ae8fedd7	source_support_kimi_depth_v1	supported	4	3180	201	3381
e1f2ed06-6d61-4071-b73a-7bf7ae8fedd7	source_support_kimi_depth_v1	unsupported	1	640	44	684
e42a506f-01a7-4253-8db8-f287fc3f504c	judge_kimi_frontier_xl_v1	correct	63	0	0	0
e42a506f-01a7-4253-8db8-f287fc3f504c	judge_kimi_frontier_xl_v1	incorrect	37	0	0	0
e42a506f-01a7-4253-8db8-f287fc3f504c	source_support_kimi_frontier_xl_v1	supported	67	0	0	0
e42a506f-01a7-4253-8db8-f287fc3f504c	source_support_kimi_frontier_xl_v1	unsupported	33	0	0	0
e51224df-d184-4d7c-bb7a-e1673de7e0f7	judge_kimi_depth_v1	correct	4	640	201	841
e51224df-d184-4d7c-bb7a-e1673de7e0f7	judge_kimi_depth_v1	incorrect	1	144	34	178
e51224df-d184-4d7c-bb7a-e1673de7e0f7	source_support_kimi_depth_v1	supported	4	1350	237	1587
e51224df-d184-4d7c-bb7a-e1673de7e0f7	source_support_kimi_depth_v1	unsupported	1	0	0	0
e5315e7e-6182-4023-b51d-f5fb7239d278	judge_kimi_depth_v1	correct	5	786	283	1069
e5315e7e-6182-4023-b51d-f5fb7239d278	source_support_kimi_depth_v1	supported	4	3410	186	3596
e5315e7e-6182-4023-b51d-f5fb7239d278	source_support_kimi_depth_v1	unsupported	1	731	71	802
e58993e6-1f2e-4e8c-af89-d6139259e281	judge_kimi_frontier_xl_v1	correct	80	0	0	0
e58993e6-1f2e-4e8c-af89-d6139259e281	judge_kimi_frontier_xl_v1	incorrect	20	0	0	0
e58993e6-1f2e-4e8c-af89-d6139259e281	source_support_kimi_frontier_xl_v1	supported	81	0	0	0
e58993e6-1f2e-4e8c-af89-d6139259e281	source_support_kimi_frontier_xl_v1	unsupported	19	0	0	0
e9842c68-6937-4768-b013-5f7e87517d18	judge_flash_depth100_v1	correct	62	11729	2504	14233
e9842c68-6937-4768-b013-5f7e87517d18	judge_flash_depth100_v1	incorrect	38	6908	1734	8642
e9842c68-6937-4768-b013-5f7e87517d18	source_support_flash_depth100_v1	supported	56	52627	2260	54887
e9842c68-6937-4768-b013-5f7e87517d18	source_support_flash_depth100_v1	unsupported	44	8316	467	8783
ed4e1cc1-5a33-4f54-a91b-8c50bc45942b	judge_kimi_frontier_xl_v1	correct	76	0	0	0
ed4e1cc1-5a33-4f54-a91b-8c50bc45942b	judge_kimi_frontier_xl_v1	incorrect	24	0	0	0
ed4e1cc1-5a33-4f54-a91b-8c50bc45942b	source_support_kimi_frontier_xl_v1	supported	78	0	0	0
ed4e1cc1-5a33-4f54-a91b-8c50bc45942b	source_support_kimi_frontier_xl_v1	unsupported	22	0	0	0
ef1954a1-a60a-43b8-914c-ee01fab35bfa	judge_kimi_heldout500_v1	correct	313	0	0	0
ef1954a1-a60a-43b8-914c-ee01fab35bfa	judge_kimi_heldout500_v1	incorrect	186	0	0	0
ef1954a1-a60a-43b8-914c-ee01fab35bfa	source_support_kimi_heldout500_v1	supported	323	0	0	0
ef1954a1-a60a-43b8-914c-ee01fab35bfa	source_support_kimi_heldout500_v1	unsupported	176	0	0	0
f292d9ca-5e2f-44b6-a175-f322073c7ad4	judge_kimi_heldout500_v1	correct	316	0	0	0
f292d9ca-5e2f-44b6-a175-f322073c7ad4	judge_kimi_heldout500_v1	incorrect	183	0	0	0
f292d9ca-5e2f-44b6-a175-f322073c7ad4	source_support_kimi_heldout500_v1	supported	311	0	0	0
f292d9ca-5e2f-44b6-a175-f322073c7ad4	source_support_kimi_heldout500_v1	unsupported	188	0	0	0
f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd	judge_kimi_frontier_xl_v1	correct	70	0	0	0
f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd	judge_kimi_frontier_xl_v1	incorrect	30	0	0	0
f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd	source_support_kimi_frontier_xl_v1	supported	72	0	0	0
f2cbaf1f-afcd-43c9-8aa1-2e7ece16dcbd	source_support_kimi_frontier_xl_v1	unsupported	28	0	0	0
f38ca0d9-6d21-4481-b8ed-56c3311af5a9	judge_kimi_v1	correct	45	7634	2356	9990
f38ca0d9-6d21-4481-b8ed-56c3311af5a9	judge_kimi_v1	incorrect	55	9082	2538	11620
f38ca0d9-6d21-4481-b8ed-56c3311af5a9	source_support_kimi_v1	supported	44	32571	2395	34966
f38ca0d9-6d21-4481-b8ed-56c3311af5a9	source_support_kimi_v1	unsupported	56	3778	438	4216
f6c48244-fccf-49f4-9c1e-c3794d6cf648	judge_kimi_v1	correct	61	10843	3339	14182
f6c48244-fccf-49f4-9c1e-c3794d6cf648	judge_kimi_v1	incorrect	39	6231	1767	7998
f6c48244-fccf-49f4-9c1e-c3794d6cf648	source_support_kimi_v1	supported	61	47687	3456	51143
f6c48244-fccf-49f4-9c1e-c3794d6cf648	source_support_kimi_v1	unsupported	39	2504	267	2771
f8b6874a-adf4-4384-9436-2059c3f95d59	judge_flash_depth100_v1	correct	67	12641	2707	15348
f8b6874a-adf4-4384-9436-2059c3f95d59	judge_flash_depth100_v1	incorrect	33	5938	1440	7378
f8b6874a-adf4-4384-9436-2059c3f95d59	source_support_flash_depth100_v1	supported	59	47397	2386	49783
f8b6874a-adf4-4384-9436-2059c3f95d59	source_support_flash_depth100_v1	unsupported	41	9390	505	9895
f92ee367-79d6-4a88-b891-4f2305eb9d22	judge_kimi_frontier_xl_v1	correct	59	0	0	0
f92ee367-79d6-4a88-b891-4f2305eb9d22	judge_kimi_frontier_xl_v1	incorrect	41	0	0	0
f92ee367-79d6-4a88-b891-4f2305eb9d22	source_support_kimi_frontier_xl_v1	supported	60	0	0	0
f92ee367-79d6-4a88-b891-4f2305eb9d22	source_support_kimi_frontier_xl_v1	unsupported	40	0	0	0
f95e0ade-83c7-4108-8131-6c6dfde7388c	judge_kimi_heldout500_v1	correct	231	0	0	0
f95e0ade-83c7-4108-8131-6c6dfde7388c	judge_kimi_heldout500_v1	incorrect	268	0	0	0
f95e0ade-83c7-4108-8131-6c6dfde7388c	source_support_kimi_heldout500_v1	supported	237	0	0	0
f95e0ade-83c7-4108-8131-6c6dfde7388c	source_support_kimi_heldout500_v1	unsupported	262	0	0	0
fcd33788-0881-4897-b1e1-29aabe71ef1a	judge_kimi_heldout500_v1	correct	140	0	0	0
fcd33788-0881-4897-b1e1-29aabe71ef1a	judge_kimi_heldout500_v1	incorrect	359	0	0	0
fcd33788-0881-4897-b1e1-29aabe71ef1a	source_support_kimi_heldout500_v1	supported	145	0	0	0
fcd33788-0881-4897-b1e1-29aabe71ef1a	source_support_kimi_heldout500_v1	unsupported	354	0	0	0
\.


--
-- Data for Name: pipeline_run_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.pipeline_run_summary (id, kind, version, model, status, started_at, completed_at, source_unit_count, fact_count, rejection_count) FROM stdin;
20527931-6fea-4604-acee-5f9982552d9e	fact-extraction	v1	accounts/fireworks/models/llama4-maverick-instruct	failed	2026-07-24 06:09:04.898203+00	2026-07-24 06:09:05.04328+00	0	0	0
f442ed12-456c-4567-9fdd-5a7cbef47f6e	fact-extraction	v1	accounts/fireworks/models/deepseek-v4-pro	failed	2026-07-24 06:09:45.781486+00	2026-07-24 06:09:47.749757+00	0	0	0
d9da8329-1f28-4e5c-9488-f50eceeaea03	fact-extraction	v1	accounts/fireworks/models/deepseek-v4-pro	failed	2026-07-24 06:10:16.052028+00	2026-07-24 06:10:18.819608+00	0	0	0
7039cbb7-15c0-4d07-a5e4-4159648af53c	fact-extraction	v1	accounts/fireworks/models/gpt-oss-120b	failed	2026-07-24 06:10:27.679773+00	2026-07-24 06:10:28.035468+00	0	0	0
25300579-b3bb-4a38-b9e6-277d91a2fb16	fact-extraction	v1	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-24 06:10:42.38661+00	2026-07-24 06:10:46.752694+00	8	8	0
4627238c-925f-4f2c-85d4-811765f82776	fact-extraction	pilot-100x100-v1	accounts/fireworks/models/gpt-oss-120b	failed	2026-07-24 19:52:32.089621+00	2026-07-24 19:52:41.937962+00	0	0	0
39da2be5-f01d-42ee-ac5d-8aca27be5259	fact-extraction	pilot-100x100-v1	accounts/fireworks/models/deepseek-v4-pro	failed	2026-07-24 19:52:51.554669+00	2026-07-24 19:54:32.099026+00	0	0	0
417e1eec-ed7a-45b6-bb81-9a49e8ac5191	fact-extraction	pilot-100x100-v1	accounts/fireworks/models/gpt-oss-120b	failed	2026-07-24 19:54:49.580072+00	2026-07-24 19:54:49.822662+00	0	0	0
4a362604-a345-49dd-99ba-bb9959c2ba02	fact-extraction	pilot-100x100-v1	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-24 19:55:04.370075+00	2026-07-24 19:55:09.536432+00	0	0	0
fa494c29-6485-4352-9c79-413c7275b625	fact-extraction	pilot-100x100-v1	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-24 19:55:16.016983+00	2026-07-24 20:11:10.344816+00	629	0	0
d8588c56-dcea-473f-92d2-838d3ef29cb4	fact-extraction	pilot-coverage-v2	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:06:19.469744+00	2026-07-25 00:17:42.785828+00	607	0	0
46ddd163-ae5c-4e2e-9a9f-ddf3fa0bc103	fact-extraction	pilot-coverage-v2	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:06:33.614578+00	2026-07-25 00:16:23.249906+00	534	0	0
dfcc26ce-6188-44ed-96e2-4be5a0f3444c	fact-extraction	pilot-coverage-v2	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:06:48.186008+00	2026-07-25 00:17:07.05513+00	597	0	0
d3143ecb-e088-4273-9f04-8b221adaa8c0	fact-extraction	pilot-coverage-v2	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:07:04.652755+00	2026-07-25 00:18:44.847219+00	615	0	0
8081d307-c1bb-4095-8653-676e20d7e6fe	fact-extraction	pilot-named-v3	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:30:41.344347+00	2026-07-25 00:31:10.885092+00	0	0	0
3749bd3a-0302-447f-a40c-5c2d13877106	fact-extraction	pilot-named-v3	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:31:22.992374+00	2026-07-25 00:45:48.195172+00	1260	630	0
4dfd88b4-47d5-49b2-857e-cbe062c0e27b	fact-extraction	pilot-named-v3	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:31:37.233843+00	2026-07-25 00:43:40.106466+00	1010	505	0
9d8b87f9-cf2c-4233-92a9-67ba5d4ed5c0	fact-extraction	pilot-named-v3	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:31:53.693514+00	2026-07-25 00:45:40.820493+00	1236	618	0
4e549945-8aa7-4af7-94ac-82b5ef526806	fact-extraction	pilot-named-v3	accounts/fireworks/models/gpt-oss-120b	complete	2026-07-25 00:32:08.0628+00	2026-07-25 00:46:35.723511+00	1304	624	0
9756ec6f-33ad-4da6-afe4-ea1b633510ba	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/deepseek-v4-flash	failed	2026-07-25 22:42:18.518105+00	2026-07-25 22:42:48.553308+00	0	0	0
2702d050-7440-4af0-9568-c9c29e000886	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/deepseek-v4-flash	failed	2026-07-25 22:42:58.77842+00	2026-07-25 22:43:47.842314+00	0	0	0
4bb4b8f7-7beb-471e-9bbd-a641012f2b17	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/gpt-oss-120b	failed	2026-07-25 22:44:19.723991+00	2026-07-25 22:47:02.39406+00	0	0	0
be851568-def8-4c34-8cd0-c1adca922049	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	failed	2026-07-25 22:47:03.321838+00	2026-07-25 22:47:57.13027+00	0	0	0
b3a6e44e-3452-4818-a464-e8c52f064c34	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	failed	2026-07-25 22:48:34.168662+00	2026-07-25 22:53:28.432797+00	0	0	0
32bef7e9-90a7-42c9-a666-364c0a49d266	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	complete	2026-07-25 22:49:20.652515+00	2026-07-25 22:50:15.570235+00	0	0	0
d3135c62-e796-4e15-b070-7976c4590b33	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	failed	2026-07-25 22:51:01.4269+00	2026-07-25 22:53:28.432797+00	0	0	0
d017d157-84bc-4937-973e-f2676e5a86b6	fact-extraction	turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	complete	2026-07-25 22:55:19.691721+00	2026-07-25 23:14:06.030317+00	13064	13064	0
d38055a0-04c7-4d72-bf30-e497c6e0c34e	fact-extraction	heldout500-turnwise-opinion-v1	accounts/fireworks/models/gpt-oss-120b	failed	2026-07-26 17:25:53.242072+00	2026-07-26 17:25:53.512992+00	0	0	0
a9e543c1-9edd-472c-82bf-df2866f2e943	fact-extraction	heldout500-turnwise-opinion-v1	accounts/fireworks/models/deepseek-v4-flash	failed	2026-07-26 17:27:01.480431+00	2026-07-26 17:29:48.7129+00	0	0	0
b04a982a-4086-473e-8e63-594972c97156	fact-extraction	heldout500-turnwise-opinion-v1	accounts/fireworks/models/kimi-k2p6	complete	2026-07-26 17:31:37.027621+00	2026-07-26 18:55:17.232655+00	69746	69746	0
\.


--
-- Data for Name: rephrasing_study_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.rephrasing_study_summary (namespace, version, model, run_count, fact_count, first_run_at, last_run_at) FROM stdin;
deepseek-fact-rephrasing	v1	accounts/fireworks/models/deepseek-v4-pro	10	259	2026-07-24 09:42:27.840563+00	2026-07-24 10:17:36.279918+00
deepseek-manual-prompt-diagnostic	v1	accounts/fireworks/models/deepseek-v4-pro	4	90	2026-07-24 10:43:05.652407+00	2026-07-24 10:44:42.242992+00
fact-prompt-study-6-dual-field	v1	accounts/fireworks/models/deepseek-v4-pro	7	81	2026-07-24 18:59:07.811952+00	2026-07-24 19:03:31.206863+00
\.


--
-- Data for Name: study_two_summary; Type: TABLE DATA; Schema: public_release; Owner: -
--

COPY public_release.study_two_summary (study_version, candidate, model, run_count, prompt_tokens, completion_tokens, total_tokens, first_run_at, last_run_at) FROM stdin;
study5-dialogue-units	dialogue-act-pairs	accounts/fireworks/models/deepseek-v4-pro	1	1164	1151	2315	2026-07-24 19:00:28.62742+00	2026-07-24 19:00:28.62742+00
study5-dialogue-units	dialogue-act-predicate-units	accounts/fireworks/models/deepseek-v4-pro	1	1943	1942	3885	2026-07-24 19:01:39.078771+00	2026-07-24 19:01:39.078771+00
study5-dialogue-units	direct	accounts/fireworks/models/deepseek-v4-pro	1	788	1259	2047	2026-07-24 19:00:48.484024+00	2026-07-24 19:00:48.484024+00
study5-dialogue-units	prefilter-nli	accounts/fireworks/models/deepseek-v4-pro	1	3105	1692	4797	2026-07-24 19:01:19.213713+00	2026-07-24 19:01:19.213713+00
study5-dialogue-units	quote-capsules-exact	accounts/fireworks/models/deepseek-v4-pro	1	1013	526	1539	2026-07-24 19:00:55.37524+00	2026-07-24 19:00:55.37524+00
study5-dialogue-units	quote-repair	accounts/fireworks/models/deepseek-v4-pro	1	3014	1770	4784	2026-07-24 19:02:53.499822+00	2026-07-24 19:02:53.499822+00
study5-dialogue-units	substantive-render	accounts/fireworks/models/deepseek-v4-pro	1	2630	739	3369	2026-07-24 19:02:27.132757+00	2026-07-24 19:02:27.132757+00
v2	atomic-attributed	accounts/fireworks/models/deepseek-v4-pro	3	6555	2515	9070	2026-07-24 11:45:48.965566+00	2026-07-24 11:46:07.662748+00
v2	attributed-packed	accounts/fireworks/models/deepseek-v4-pro	3	6171	3462	9633	2026-07-24 11:41:49.942086+00	2026-07-24 11:42:15.780977+00
v2	bundle	accounts/fireworks/models/deepseek-v4-pro	3	7782	4336	12118	2026-07-24 11:29:59.593609+00	2026-07-24 11:31:46.297705+00
v2	carryover	accounts/fireworks/models/deepseek-v4-pro	3	8881	5453	14334	2026-07-24 11:20:41.449459+00	2026-07-24 11:21:45.697897+00
v2	clause-anchored	accounts/fireworks/models/deepseek-v4-pro	3	47729	5996	53725	2026-07-24 12:09:14.049933+00	2026-07-24 12:09:35.018836+00
v2	decomposed-nli	accounts/fireworks/models/deepseek-v4-pro	3	54193	8708	62901	2026-07-24 12:53:23.97307+00	2026-07-24 12:54:48.958648+00
v2	direct	accounts/fireworks/models/deepseek-v4-pro	3	5883	3364	9247	2026-07-24 11:17:13.614419+00	2026-07-24 11:17:44.722463+00
v2	expanded-bundle	accounts/fireworks/models/deepseek-v4-pro	3	6351	2018	8369	2026-07-24 11:35:29.732116+00	2026-07-24 11:35:54.288533+00
v2	expanded-packed	accounts/fireworks/models/deepseek-v4-pro	3	6822	4035	10857	2026-07-24 11:36:38.298604+00	2026-07-24 11:37:15.352264+00
v2	framed-render	accounts/fireworks/models/deepseek-v4-pro	3	35791	11423	47214	2026-07-24 13:04:19.914883+00	2026-07-24 13:05:25.333714+00
v2	highrecall-nli	accounts/fireworks/models/deepseek-v4-pro	3	18097	4264	22361	2026-07-24 12:36:14.09784+00	2026-07-24 12:36:41.733793+00
v2	highrecall-verified	accounts/fireworks/models/deepseek-v4-pro	3	18231	6801	25032	2026-07-24 12:30:20.479591+00	2026-07-24 12:31:23.814815+00
v2	predicate-repair	accounts/fireworks/models/deepseek-v4-pro	3	41647	10331	51978	2026-07-24 13:24:20.359951+00	2026-07-24 13:26:19.659116+00
v2	prefilter-nli	accounts/fireworks/models/deepseek-v4-pro	3	15905	4868	20773	2026-07-24 12:44:31.976639+00	2026-07-24 12:45:34.673641+00
v2	proof-obligation	accounts/fireworks/models/deepseek-v4-pro	3	59130	6640	65770	2026-07-24 12:22:53.781543+00	2026-07-24 12:24:17.295253+00
v2	quote-repair	accounts/fireworks/models/deepseek-v4-pro	3	15720	5751	21471	2026-07-24 13:13:28.978608+00	2026-07-24 13:14:14.547687+00
v2	speaker-packed	accounts/fireworks/models/deepseek-v4-pro	3	8710	5336	14046	2026-07-24 11:18:17.038036+00	2026-07-24 11:19:27.680572+00
v2	substantive-render	accounts/fireworks/models/deepseek-v4-pro	3	33061	7966	41027	2026-07-24 13:35:44.843146+00	2026-07-24 13:36:46.251291+00
v2	turnwise-local-staged	accounts/fireworks/models/deepseek-v4-pro	3	99715	11735	111450	2026-07-24 11:57:12.366557+00	2026-07-24 11:58:05.087542+00
v2	turnwise-staged	accounts/fireworks/models/deepseek-v4-pro	3	52569	9794	62363	2026-07-24 11:50:39.70789+00	2026-07-24 11:51:39.353726+00
v2-c4-sbar	fable-c4-sbar	accounts/fireworks/models/deepseek-v4-pro	3	20368	7915	28283	2026-07-24 13:50:38.329826+00	2026-07-24 13:52:57.999107+00
v2-dialogue-act-pairs	dialogue-act-pairs	accounts/fireworks/models/deepseek-v4-pro	3	9274	7073	16347	2026-07-24 14:10:08.506879+00	2026-07-24 14:11:58.868232+00
v2-dialogue-act-predicate-units	dialogue-act-predicate-units	accounts/fireworks/models/deepseek-v4-pro	3	19748	12134	31882	2026-07-24 14:27:20.68455+00	2026-07-24 14:29:35.603695+00
v2-glm52	clause-anchored	accounts/fireworks/models/glm-5p2	3	43222	4567	47789	2026-07-24 12:12:32.467072+00	2026-07-24 12:12:59.02586+00
v2-quote-capsules	quote-capsules	accounts/fireworks/models/deepseek-v4-pro	1	3407	1476	4883	2026-07-24 14:48:17.731845+00	2026-07-24 14:48:17.731845+00
v2-quote-capsules-exact	quote-capsules-exact	accounts/fireworks/models/deepseek-v4-pro	3	8821	4111	12932	2026-07-24 14:51:43.760146+00	2026-07-24 14:52:20.258286+00
\.


--
-- PostgreSQL database dump complete
--

\unrestrict 6O5YU1eYqprrejMOgG1Bo5eLFLJSKWvgAresakQnNWDNxZGqGGmkHMoz29LsHs2
