IF EXISTS(SELECT 1
            FROM sys.objects
            WHERE name = '[ContainerBenchmark].[Project]')
    BEGIN
        DROP TABLE [ContainerBenchmark].[Project]
    END;

IF NOT EXISTS(SELECT 1
              FROM sys.objects
              WHERE object_id = OBJECT_ID('[ContainerBenchmark].[Project]')
                AND type IN (N'U'))
    BEGIN
        CREATE TABLE ContainerBenchmark.Project
        (
            gitlabId     INT          NOT NULL PRIMARY KEY,
            teamGitlabId INT          NOT NULL,
            name         VARCHAR(100) NOT NULL CONSTRAINT CHK_Project_Name_Nonempty CHECK (name <> ''),
            path         VARCHAR(100) NOT NULL UNIQUE CONSTRAINT CHK_Project_Path_Nonempty CHECK (path <> ''),
            avatarURL    VARCHAR(200) NOT NULL,
            isActive     BIT          NOT NULL,
            toSync       BIT          NOT NULL,
        )
    END;
