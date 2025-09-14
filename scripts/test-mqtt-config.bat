@echo off
REM Test MQTT Configuration Script for Windows

echo 🧪 Testing MQTT Configuration...

REM Test mosquitto.conf syntax
echo 📝 Testing mosquitto.conf syntax...
docker run --rm -v "%cd%\mosquitto.conf:/mosquitto/config/mosquitto.conf" eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

if %errorlevel% equ 0 (
    echo ✅ mosquitto.conf syntax is valid
) else (
    echo ❌ mosquitto.conf has syntax errors
    exit /b 1
)

REM Test mosquitto.dev.conf syntax
echo 📝 Testing mosquitto.dev.conf syntax...
docker run --rm -v "%cd%\mosquitto.dev.conf:/mosquitto/config/mosquitto.conf" eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

if %errorlevel% equ 0 (
    echo ✅ mosquitto.dev.conf syntax is valid
) else (
    echo ❌ mosquitto.dev.conf has syntax errors
    exit /b 1
)

REM Test mosquitto.prod.conf syntax (if password file exists)
if exist "mosquitto_config\passwd" (
    echo 📝 Testing mosquitto.prod.conf syntax...
    docker run --rm -v "%cd%\mosquitto.prod.conf:/mosquitto/config/mosquitto.conf" -v "%cd%\mosquitto_config\passwd:/mosquitto/config/passwd" -v "%cd%\mosquitto_config\acl:/mosquitto/config/acl" eclipse-mosquitto:2.0 mosquitto -c /mosquitto/config/mosquitto.conf -t

    if %errorlevel% equ 0 (
        echo ✅ mosquitto.prod.conf syntax is valid
    ) else (
        echo ❌ mosquitto.prod.conf has syntax errors
        exit /b 1
    )
) else (
    echo ⚠️ Skipping mosquitto.prod.conf test (no password file found)
)

echo 🎉 All MQTT configurations are valid!