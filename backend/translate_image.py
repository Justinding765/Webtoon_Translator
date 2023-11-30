from flask import Flask, request, jsonify
import base64
from io import BytesIO
import requests
from PIL import Image, ImageDraw, ImageFont
from google.cloud import vision
from google.cloud import translate_v3 as translate_v3
import os
from datetime import datetime
import cv2
import numpy as np
import uuid

# Set Google Cloud credentials
os.environ['GOOGLE_APPLICATION_CREDENTIALS'] = '../credentials.json'

app = Flask(__name__)
# Function to wrap text

@app.route('/translate_image', methods=['POST'])
def translate_image():
    data = request.json
    img_url = data.get('url')

    # Initialize Google Cloud clients
    vision_client = vision.ImageAnnotatorClient()
     # Initialize the Translation client for Translation - Advanced API
    translate_client = translate_v3.TranslationServiceClient()

    # Fetch the image
    headers = {'User-Agent': 'Mozilla/5.0 ...'}
    response = requests.get(img_url, headers=headers)
    img_np = np.array(Image.open(BytesIO(response.content)))
    img = img_np.copy()

    # Detect text in the image
    image = vision.Image(content=response.content)
    response = vision_client.document_text_detection(image=image)


    # Process each block of text
    temp = ""
    for page in response.full_text_annotation.pages:
            for block in page.blocks:
                block_texts = []
                for paragraph in block.paragraphs:
                    paragraph_texts = []
                    for word in paragraph.words:
                        word_text = ''.join([symbol.text for symbol in word.symbols])
                        paragraph_texts.append(word_text)
                    block_texts.append(' '.join(paragraph_texts))

                block_text = ' '.join(block_texts)
                # Use the Advanced API for translation
                parent = "projects/omega-winter-406314/locations/global"
                response = translate_client.translate_text(
                    contents=[block_text],
                    target_language_code="en",
                    parent=parent,
                    mime_type="text/plain"  # Use "text/plain" for plain text, "text/html" for HTML content
                )
                translated_text = response.translations[0].translated_text
                # Get the bounding box for the block
                vertices = [(vertex.x, vertex.y) for vertex in block.bounding_box.vertices]
                rect_start = (vertices[0][0], vertices[0][1])
                rect_end = (vertices[2][0], vertices[2][1])
                # Create a mask for inpainting
                mask = np.zeros(img.shape[:2], np.uint8)
                cv2.rectangle(mask, rect_start, rect_end, 255, thickness=-1)

                # Inpainting to remove the text, looks at 8 pixel radius
                img = cv2.inpaint(img, mask, 8, cv2.INPAINT_TELEA)

                # Extract a small area around the text to analyze the color for opposite color
                margin = 5
                y1, y2 = max(0, rect_start[1]-margin), min(rect_end[1]+margin, img.shape[0])
                x1, x2 = max(0, rect_start[0]-margin), min(rect_end[0]+margin, img.shape[1])
                
                if (x2 > x1) and (y2 > y1):
                    color_sample_area = img[y1:y2, x1:x2]
                    average_color = np.mean(color_sample_area, axis=(0, 1))
                    if not np.isnan(average_color).any():
                        average_color = tuple(int(c) for c in average_color)
                        opposite_color = tuple(255 - c for c in average_color)  # Invert the color for translated text
                    else:
                        opposite_color = (0, 0, 0)  # Fallback opposite color (black)
                else:
                    opposite_color = (0, 0, 0)  # Fallback opposite color (black)

                def wrap_text(text, font, max_width):
                    words = text.split(' ')
                    wrapped_lines = []
                    line = ''

                    for word in words:
                        test_line = line + word + ' '
                        text_size = cv2.getTextSize(test_line, font, font_scale, 1)[0]
                        if text_size[0] <= max_width:
                            line = test_line
                        else:
                            wrapped_lines.append(line)
                            line = word + ' '
                    wrapped_lines.append(line)

                    return wrapped_lines
                
                
                # Overlay the translated text with the opposite color
                font_scale = 0.5
                font = cv2.FONT_HERSHEY_SIMPLEX
                max_width = min(rect_end[0] - rect_start[0], img.shape[1] - rect_start[0])

                wrapped_lines = wrap_text(translated_text, font, max_width)
                for line in wrapped_lines:
                    # Check if the position is within the image bounds
                    if rect_start[1] < img.shape[0]:
                        cv2.putText(img, line.strip(), rect_start, font, font_scale, opposite_color, 1, cv2.LINE_AA)
                        rect_start = (rect_start[0], rect_start[1] + int(font_scale * 30))  # Move down for the next line


    # Convert back to PIL Image and save
    img_pil = Image.fromarray(img)
    static_folder_path = '../frontend/Static/Images'  # Relative path from your script to the images folder
    unique_id = uuid.uuid4()
    filename = f"translated_image_{unique_id}.jpg"
    filepath = os.path.join(static_folder_path, filename)
    temp+= "\n\n\n\n\n"
    img_pil.save(filepath, "JPEG")

    return jsonify({"image": filepath})

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)
