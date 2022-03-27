package server

import dev.inmo.tgbotapi.extensions.api.files.downloadFile
import dev.inmo.tgbotapi.extensions.api.send.media.sendDocument
import dev.inmo.tgbotapi.extensions.api.send.reply
import dev.inmo.tgbotapi.extensions.behaviour_builder.telegramBotWithBehaviourAndLongPolling
import dev.inmo.tgbotapi.extensions.behaviour_builder.triggers_handling.onCommand
import dev.inmo.tgbotapi.extensions.behaviour_builder.triggers_handling.onMedia
import dev.inmo.tgbotapi.requests.abstracts.asMultipartFile
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers

import org.apache.tika.metadata.Metadata
import org.apache.tika.parser.AutoDetectParser
import org.apache.tika.sax.BodyContentHandler

import java.io.ByteArrayOutputStream
import java.io.FileOutputStream
import java.nio.file.Files
import kotlin.system.exitProcess


suspend fun main() {
    val botToken = System.getenv("TELEGRAM_TOKEN")

    if (botToken.isNullOrBlank()) {
        println("Token not found!")
        exitProcess(1)
    }

    telegramBotWithBehaviourAndLongPolling(botToken, CoroutineScope(Dispatchers.IO)) {
        onCommand("start") {
            val welcomeMessage = """
                |Welcome to Content Extractor Bot
                |This bot will help you extract text from many kinds of formats, including images!.
                |
                |Available Command
                |/start - Getting started
                |/help  - Help
                |
                |""".trimMargin()
            reply(it, welcomeMessage)
        }

        onCommand("help") {
            reply(it, """
                |This bot will help you extract text from many kinds of formats (including Images).
                |Send a file to this chat room, and the bot will start working.

                |The list of supported formats can be found in https://tika.apache.org/1.28.1/formats.html
            """.trimMargin())
        }

        onMedia(initialFilter = null) {
            reply(it, "Please wait, your file is being extracted!")
            val file = bot.downloadFile(it.content.media).inputStream()
            val parser = AutoDetectParser()

            val outputStream = ByteArrayOutputStream()
            val handler = BodyContentHandler(outputStream)

            parser.parse(file, handler, Metadata())

            val content = outputStream.toString()

            if(outputStream.size() > 4096) {
                reply(it, """Text length exceeds telegram message limit :(
                    |The text will be sent as a text file
                """.trimMargin())

                val tempFile = Files.createTempFile("",".txt").toFile()
                val fos = FileOutputStream(tempFile)
                outputStream.writeTo(fos)
                fos.close()

                sendDocument(it.chat, tempFile.asMultipartFile())
                tempFile.delete()
            } else if(outputStream.size() == 0) {
                reply(it, "No content detected!")
            } else {
                reply(it, content)
            }

        }
    }.second.join()
}