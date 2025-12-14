import {
  StyleSheet,
  View,
  Text,
  TextInput,
  Pressable,
  KeyboardAvoidingView,
  ScrollView,
  Platform,
} from 'react-native';
import { Image } from 'expo-image';

type Props = {
  photoUri: string;
  description: string;
  isLoading?: boolean;
  onDescriptionChange: (text: string) => void;
  onRetake: () => void;
  onRegister: () => void;
};

export function PhotoPreview({
  photoUri,
  description,
  isLoading,
  onDescriptionChange,
  onRetake,
  onRegister,
}: Props) {
  return (
    <KeyboardAvoidingView
      style={styles.flex1}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      keyboardVerticalOffset={Platform.OS === 'ios' ? 100 : 0}
    >
      <ScrollView
        contentContainerStyle={styles.scrollContainer}
        keyboardShouldPersistTaps="handled"
        showsVerticalScrollIndicator={false}
      >
        <Text style={styles.title}>この写真で登録しますか？</Text>
        <View style={styles.previewContainer}>
          <Image
            source={{ uri: photoUri }}
            style={styles.preview}
            contentFit="cover"
            contentPosition="center"
          />
        </View>

        <View style={styles.inputContainer}>
          <Text style={styles.inputLabel}>
            これから生まれるモンスターの名前
          </Text>
          <TextInput
            style={styles.descriptionInput}
            placeholder="例：とらくえ"
            placeholderTextColor="#999"
            value={description}
            onChangeText={onDescriptionChange}
          />
        </View>

        <View style={styles.buttonRow}>
          <Pressable
            style={[styles.button, styles.retakeButton]}
            onPress={onRetake}
          >
            <Text style={styles.buttonText}>撮り直す</Text>
          </Pressable>
          <Pressable
            style={[styles.button, styles.registerButton]}
            disabled={isLoading}
            aria-busy={isLoading}
            onPress={onRegister}
          >
            <Text style={styles.buttonText}>
              {isLoading ? 'アップロード中...' : '登録する'}
            </Text>
          </Pressable>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  flex1: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  scrollContainer: {
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#333',
  },
  previewContainer: {
    width: '100%',
    aspectRatio: 1,
    borderRadius: 16,
    overflow: 'hidden',
    marginTop: 16,
  },
  preview: {
    flex: 1,
  },
  inputContainer: {
    width: '100%',
    marginTop: 16,
  },
  inputLabel: {
    fontSize: 16,
    fontWeight: '600',
    color: '#333',
    marginBottom: 6,
  },
  descriptionInput: {
    width: '100%',
    padding: 12,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 12,
    backgroundColor: '#fff',
    fontSize: 16,
    color: '#333',
  },
  buttonRow: {
    flexDirection: 'row',
    gap: 16,
    marginTop: 30,
  },
  button: {
    paddingVertical: 14,
    paddingHorizontal: 24,
    borderRadius: 12,
    backgroundColor: '#007AFF',
  },
  retakeButton: {
    backgroundColor: '#8E8E93',
  },
  registerButton: {
    backgroundColor: '#34C759',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
