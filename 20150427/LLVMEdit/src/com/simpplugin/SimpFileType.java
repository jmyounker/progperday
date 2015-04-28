package com.simpplugin;

import com.intellij.openapi.fileTypes.LanguageFileType;
import com.intellij.openapi.util.IconLoader;
import com.intellij.openapi.vfs.CharsetToolkit;
import com.intellij.openapi.vfs.VirtualFile;
import com.intellij.openapi.vfs.encoding.EncodingManager;
import org.jetbrains.annotations.NotNull;

import javax.swing.*;
import java.nio.charset.Charset;

/**
 * Created by jeff on 4/24/15.
 */
public class SimpFileType extends LanguageFileType {
    public static final Icon FILE_ICON = IconLoader.getIcon("/fileTypes/properties.png");
    public static final LanguageFileType INSTANCE = new SimpFileType();

    private SimpFileType() {
        super(SimpLanguage.INSTANCE);
    }

    @NotNull
    public String getName() {
        return "Simp";
    }

    @NotNull
    public String getDescription() {
        return "Simple mathematical expressions";
    }

    @NotNull
    public String getDefaultExtension() {
        return "simp";
    }

    public Icon getIcon() {
        return FILE_ICON;
    }

    public String getCharset(@NotNull VirtualFile file, final byte[] content) {
        Charset charset = EncodingManager.getInstance().getDefaultCharsetForPropertiesFiles(file);
        String defaultCharsetName = charset == null ? CharsetToolkit.getDefaultSystemCharset().name() : charset.name();
        return defaultCharsetName;
    }
}
