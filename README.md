# IMP


## About
Service that evaluates player Impact on the final result (IMP). Based on:
<ul>
    <li>Minutes Played</li>
    <li>Plus-Minus</li>
    <li>Final Game Score</li>
</ul>

IMP provides deep analysis of player performance considering actual playing time volume.

## Supported Leagues

<ul>
    <li><a href="https://www.nba.com/">NBA</a></li>
    <li><a href="https://tambov.ilovebasket.ru/competitions/89960">Tambov MLBL</a></li>
</ul>

## CLI

<ul>
    <li>
        <b>cron</b>
        <p>
            Every day at 12AM service saves yesterday played games for each league.
        </p>
    </li>
    <li>
        <b>seed:list</b>
        <p>
            Prints all available database seeders.
        </p>
    </li>
    <li>
        <b>seed:db</b>
        <p>
            Seed database with provided seeder.
        </p>
    </li>
    <li>
        <b>save-game</b>
        <p>
            Save game from provided league with provided game id.
        </p>
    </li>
    <li>
        <b>save-games-by-date</b>
        <p>
            Save already played games for provided league and date (in format dd-mm-yyyy). If date is not provided, saves today games.
        </p>
    </li>
    <li>
        <b>save-game-by-team</b>
        <p>
            Save already played games for provided league and team id
        </p>
    </li>
</ul>

## Help
<ul>
    <li>
        <b>generate mock</b>
        <p>
            mockgen -source={path/to/file}/{file_name}.go -destination={path/to/file}/mocks/mock_{file_name}.go -package={same package as file}
        </p>
    </li>
    <li>
        <b>test with coverage</b>
        <p>
            There is make command in Makefile called 'test-with-coverage'
            <br>
            After running this command you can see coverage report in ./coverage.out.
            <br>
            In IDEA you can open 'Coverage' tab and see coverage and import this file.
        </p>
    </li>
</ul>