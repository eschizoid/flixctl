-- CREATE VIEW  library AS
    SELECT metadata_items.title,
           metadata_items.title_sort,
           round(cast(media_parts.size as float) / 1024 / 1024 / 1024, 2) || ' Gb'             AS file_size,
           printf("%.0f", round(cast(media_items.width as float) / 1.77777778, 0)) || 'p'      AS video_quality,
           CASE
             WHEN media_items.audio_codec = 'aac'
                     THEN 'AAC'
             WHEN media_items.audio_codec = 'ac3'
                     THEN 'Dolby Digital'
             WHEN media_items.audio_codec = 'dca'
                     THEN 'DTS'
             WHEN media_items.audio_codec = 'eac3'
                     THEN 'Dolby Digital Plus'
             WHEN media_items.audio_codec = 'flac'
                     THEN 'FLAC'
             WHEN media_items.audio_codec = 'truehd'
                     THEN 'Dolby TrueHD'
             ELSE 'Other Audio' END                                                            AS [audio_codecs],
           media_items.audio_channels,
           printf("%.2f", round(cast(media_streams.bitrate as float) / 1000000, 2)) || ' Mbps' AS birate,
           metadata_items.year,
           (cast(media_items.duration as numeric) / 60000) || ' mins'                          AS duration_mins,
           printf("%.2f", cast(media_items.duration as numeric) / 60000 / 60 +
                          cast(media_items.duration as numeric) / 60000 % 60 / 100.0) || ' hr' AS duration_hours,
           metadata_items.content_rating,
           strftime('%d-%m-%Y', metadata_items.added_at)                                       AS added_date,
           media_parts.file                                                                    AS file_location
    FROM metadata_items
           JOIN media_parts
           JOIN media_items
           JOIN media_streams
    WHERE metadata_items.library_section_id = 3
      AND media_streams.codec = 'h264'
      AND metadata_items.id = media_items.metadata_item_id
      AND media_items.id = media_parts.media_item_id
      AND media_items.id = media_streams.media_item_id
    ORDER BY metadata_items.title_sort,
             metadata_items.year;

-- -- DROP VIEW
-- DROP VIEW LIBRARY;
--
-- -- SELECT VIEW
-- SELECT * FROM LIBRARY;
